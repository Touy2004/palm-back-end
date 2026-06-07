# Understanding the API: A Conceptual Guide

If you want to understand how the Palm Recognition Attendance System works behind the scenes, this guide explains the logic and purpose of every single API route in plain English. 

Instead of showing raw code and JSON structures, this document tells the story of how data moves between the Mobile App, the Hardware Scanner, the Web Dashboard, and the Server.

---

## 1. Authentication (The Front Door)
*These routes are public and act as the gateway into the system.*

* **`POST /auth/register`**: Allows a new employee to sign up for the system using their phone number and a password. It creates their profile in the database but notes that their "palm is not yet registered."
* **`POST /auth/login`**: The standard login route. When a user or admin logs in, the server gives them two "keys": an Access Token (lasts 15 minutes) and a Refresh Token (lasts 7 days).
* **`POST /auth/refresh`**: Because the Access Token expires quickly for security, the mobile app or web dashboard silently calls this route in the background to get a new Access Token without forcing the user to log in again.

---

## 2. Mobile App Routes (The Employee Portal)
*These routes are used by the employee's mobile app. The server knows who is calling these routes based on their Access Token.*

* **`GET /me/profile`**: Fetches the user's personal details (name, employee code, department) to display on their app's home screen.
* **`PATCH /me/password`**: Allows the user to securely change their password.
* **`GET /me/attendance`**: Fetches a paginated list of the user's daily check-in and check-out logs so they can track their own attendance history.
* **`GET /me/palms`**: Lists all the biometric palm templates the user has registered (e.g., Left Hand, Right Hand), so they know if they are ready to use the scanners.

---

## 3. The Pairing Flow (Connecting Phone to Scanner)
*These routes handle the secure "handshake" between an employee's phone and a physical scanner when they are registering their palm for the first time.*

* **`POST /pairing/scan`**: When the employee points their phone camera at the scanner's QR code, the app sends the QR data to this route. The server checks if the QR code is valid and tells the app which scanner it belongs to (e.g., "Main Entrance Scanner").
* **`POST /pairing/approve`**: The employee clicks "Approve" on their phone. This route tells the server: "I am at this scanner, and I authorize it to scan my palm and link it to my account."

---

## 4. Hardware Scanner Routes (The Machine's Brain)
*These routes are used exclusively by the physical Raspberry Pi scanners stationed at doors. The scanners have their own special Admin tokens.*

### System Health
* **`POST /devices/heartbeat`**: The scanner calls this every few minutes to tell the server "I am online and working." The server updates its `last_seen_at` timestamp in the database.

### The Enrollment Process (Creating the QR Code)
* **`POST /devices/pairing-session`**: The scanner calls this to say "Generate a random QR code for me." The server saves the code in the database and gives it to the scanner to display on its screen.
* **`GET /devices/pairing-status`**: While the QR code is on the screen, the scanner constantly polls this route asking, "Did an employee approve this QR code on their phone yet?" Once the server says "Yes, John Doe approved it!", the scanner proceeds.
* **`POST /devices/enroll-palm`**: The scanner reads the physical palm, encrypts the biometric vector data, and sends it to this route. The server saves this highly secure template into the database and permanently links it to John Doe.

### The Daily Attendance Process
* **`POST /devices/identify`**: If the scanner just wants to know *who* a palm belongs to (without logging attendance), it sends the palm vector here. The server compares it against all templates and returns the matched User ID.
* **`POST /devices/process-attendance`**: The main daily route! The scanner reads a palm and sends the data here. The server:
  1. Identifies the user.
  2. Checks if they have already checked in today.
  3. If no, it creates a "Check In" timestamp.
  4. If yes, it updates their "Check Out" timestamp.
  5. Replies to the scanner with the user's name so the scanner screen can say "Welcome, John Doe!"

---

## 5. Web Admin Routes (The Command Center)
*These routes are used by HR and Administrators using the Web Dashboard to monitor the whole system.*

### Monitoring & Reporting
* **`GET /admin/dashboard`**: Fetches the high-level numbers for today (e.g., Total Devices Online, Check-ins Today, Total Users).
* **`GET /admin/reports`**: Used to generate the Monthly Payroll Reports. The admin can filter by Month and Department to get aggregated stats (Total Present days, Total Late days) for every employee.
* **`GET /admin/attendance`**: Fetches a raw, global list of every single check-in across the entire company for auditing.
* **`GET /admin/attendance/users/:id/history`**: Fetches the raw attendance logs for one specific employee.

### Management
* **`GET, POST, PATCH, DELETE /admin/users`**: Standard routes to view all employees, create new accounts manually, edit their departments, or fire/deactivate them.
* **`GET, POST, PATCH /admin/devices`**: Standard routes to register new hardware scanners in the system, rename them, or check if they are currently online.
