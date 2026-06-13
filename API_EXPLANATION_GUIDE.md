# Understanding the API: A Conceptual Guide

If you want to understand how the Palm Recognition Attendance System works behind the scenes, this guide explains the logic and purpose of every single API route in plain English, along with the underlying SQL commands executed in the database.

Instead of showing raw code and JSON structures, this document tells the story of how data moves between the Mobile App, the Hardware Scanner, the Web Dashboard, and the Server.

---

## 1. Authentication (The Front Door)
*These routes are public and act as the gateway into the system.*

* **`POST /auth/register`**: Allows a new employee to sign up for the system using their phone number and a password. It creates their profile in the database but notes that their "palm is not yet registered."
  ```sql
  INSERT INTO users (id, phone, password_hash, full_name, role, department, employee_code, created_at) 
  VALUES ($1, $2, $3, $4, $5, $6, $7, NOW());
  ```

* **`POST /auth/login`**: The standard login route. When a user or admin logs in, the server verifies their credentials and gives them an Access Token (lasts 15 minutes) and a Refresh Token (lasts 7 days).
  ```sql
  SELECT id, phone, password_hash, role FROM users WHERE phone = $1 LIMIT 1;
  ```

* **`POST /auth/refresh`**: Because the Access Token expires quickly for security, the mobile app or web dashboard silently calls this route in the background to get a new Access Token without forcing the user to log in again. *(Typically relies on JWT signature verification rather than a direct database query).*

---

## 2. Mobile App Routes (The Employee Portal)
*These routes are used by the employee's mobile app. The server knows who is calling these routes based on their Access Token.*

* **`GET /me/profile`**: Fetches the user's personal details (name, employee code, department) to display on their app's home screen.
  ```sql
  SELECT id, phone, email, full_name, role, department, employee_code 
  FROM users WHERE id = $1;
  ```

* **`PATCH /me/password`**: Allows the user to securely change their password.
  ```sql
  SELECT password_hash FROM users WHERE id = $1;
  UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1;
  ```

* **`GET /me/attendance`**: Fetches a paginated list of the user's daily check-in and check-out logs so they can track their own attendance history.
  ```sql
  SELECT * FROM attendance_logs 
  WHERE user_id = $1 
  ORDER BY attendance_date DESC 
  LIMIT 20 OFFSET 0;
  ```

* **`GET /me/palms`**: Lists all the biometric palm templates the user has registered (e.g., Left Hand, Right Hand), so they know if they are ready to use the scanners.
  ```sql
  SELECT id, model_version, hand_side, status, created_at 
  FROM palm_templates WHERE user_id = $1;
  ```

---

## 3. The Pairing Flow (Connecting Phone to Scanner)
*These routes handle the secure "handshake" between an employee's phone and a physical scanner when they are registering their palm for the first time.*

* **`POST /pairing/scan`**: When the employee points their phone camera at the scanner's QR code, the app sends the QR data to this route. The server checks if the QR code is valid and updates the session status to "scanned".
  ```sql
  SELECT * FROM device_pairing_sessions WHERE session_token = $1;
  UPDATE device_pairing_sessions SET status = 'scanned', scanned_at = NOW() WHERE id = $2;
  ```

* **`POST /pairing/approve`**: The employee clicks "Approve" on their phone and selects which hand they want to enroll (Left or Right). This route tells the server: "I am at this scanner, and I authorize it to scan my selected palm and link it to my account."
  ```sql
  UPDATE device_pairing_sessions 
  SET user_id = $1, hand_side = $2, status = 'approved', approved_at = NOW() 
  WHERE session_token = $3;
  ```

---

## 4. Hardware Scanner Routes (The Machine's Brain)
*These routes are used exclusively by the physical Raspberry Pi scanners stationed at doors. The scanners have their own special Admin tokens.*

### System Health
* **`POST /devices/heartbeat`**: The scanner calls this every few minutes to tell the server "I am online and working."
  ```sql
  UPDATE devices SET status = 'active' WHERE device_code = $1;
  ```

### The Enrollment Process (Creating the QR Code)
* **`POST /devices/pairing-session`**: The scanner calls this to say "Generate a random QR code for me." The server saves the code in the database and gives it to the scanner to display on its screen.
  ```sql
  INSERT INTO device_pairing_sessions (id, device_id, session_token, purpose, status, expires_at) 
  VALUES ($1, $2, $3, 'enrollment', 'pending', $4);
  ```

* **`GET /devices/pairing-status`**: While the QR code is on the screen, the scanner constantly polls this route asking, "Did an employee approve this QR code on their phone yet?"
  ```sql
  SELECT status, user_id, hand_side FROM device_pairing_sessions WHERE session_token = $1;
  ```

* **`POST /devices/enroll-palm`**: The scanner reads the physical palm, encrypts the biometric vector data, and sends it to this route. The server saves this highly secure template into the database and permanently links it to the user.
  ```sql
  INSERT INTO palm_templates (id, user_id, registered_device_id, template_encrypted, template_nonce, hand_side, status) 
  VALUES ($1, $2, $3, $4, $5, $6, 'active');
  
  UPDATE device_pairing_sessions SET status = 'completed' WHERE session_token = $7;
  ```

### The Daily Attendance Process
* **`POST /devices/identify`**: If the scanner just wants to know *who* a palm belongs to (without logging attendance), it sends the palm vector here. The server compares it against all templates and returns the matched User ID.
  ```sql
  SELECT id, user_id, template_encrypted, template_nonce FROM palm_templates WHERE status = 'active';
  -- (Comparison of biometric vectors happens in backend memory or via a vector database)
  ```

* **`POST /devices/process-attendance`**: The main daily route! The scanner reads a palm and sends the data here.
  1. Identifies the user (same as above).
  2. Checks if they have already checked in today.
  3. If no, it creates a "Check In" timestamp.
  4. If yes, it updates their "Check Out" timestamp.
  ```sql
  -- Step 2: Check for existing attendance today
  SELECT * FROM attendance_logs WHERE user_id = $1 AND attendance_date = CURRENT_DATE LIMIT 1;

  -- Step 3: If no check-in exists yet today
  INSERT INTO attendance_logs (id, user_id, device_id, attendance_date, check_in_time, status) 
  VALUES ($1, $2, $3, CURRENT_DATE, NOW(), 'present');

  -- Step 4: If checking out
  UPDATE attendance_logs SET check_out_time = NOW(), status = 'present' WHERE id = $1;
  ```

---

## 5. Web Admin Routes (The Command Center)
*These routes are used by HR and Administrators using the Web Dashboard to monitor the whole system.*

### Monitoring & Reporting
* **`GET /admin/dashboard`**: Fetches the high-level numbers for today (e.g., Total Devices Online, Check-ins Today, Total Users).
  ```sql
  SELECT COUNT(*) FROM devices WHERE status = 'active';
  SELECT COUNT(*) FROM attendance_logs WHERE attendance_date = CURRENT_DATE;
  SELECT COUNT(*) FROM users;
  ```

* **`GET /admin/reports`**: Used to generate the Monthly Payroll Reports.
  ```sql
  SELECT user_id, 
         COUNT(CASE WHEN status = 'present' THEN 1 END) as present_days,
         COUNT(CASE WHEN status = 'late' THEN 1 END) as late_days
  FROM attendance_logs 
  WHERE attendance_date >= $1 AND attendance_date <= $2 
  GROUP BY user_id;
  ```

* **`GET /admin/attendance`**: Fetches a raw, global list of every single check-in across the entire company for auditing.
  ```sql
  SELECT * FROM attendance_logs ORDER BY created_at DESC LIMIT $1 OFFSET $2;
  ```

* **`GET /admin/attendance/users/:id/history`**: Fetches the raw attendance logs for one specific employee.
  ```sql
  SELECT * FROM attendance_logs WHERE user_id = $1 ORDER BY attendance_date DESC;
  ```

### Management
* **`GET, POST, PATCH, DELETE /admin/users`**: Standard CRUD routes to manage employees.
  ```sql
  SELECT * FROM users ORDER BY created_at DESC;
  INSERT INTO users (...) VALUES (...);
  UPDATE users SET department = $1, role = $2 WHERE id = $3;
  DELETE FROM users WHERE id = $1;
  ```

* **`GET, POST, PATCH /admin/devices`**: Standard CRUD routes to manage hardware scanners.
  ```sql
  SELECT * FROM devices ORDER BY created_at DESC;
  INSERT INTO devices (device_code, name, location, status) VALUES (...);
  UPDATE devices SET name = $1, status = $2 WHERE id = $3;
  ```
