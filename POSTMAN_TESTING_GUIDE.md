# Complete API Testing Guide (Postman)

This guide provides a step-by-step workflow for testing the entire Palm Recognition Attendance System from start to finish. It covers creating users, registering devices, the full QR enrollment flow, and processing attendance.

**Prerequisites:** 
- Open Postman and select the `Palm Recognition API` collection.
- Make sure your server is running (or you are connected via cloudflared).
- Set your `base_url` variable in Postman (e.g., `http://localhost:8080/api/v1` or your public URL).

---

## Phase 1: Setup Admin & Hardware Device

### Step 1: Register an Admin User
First, we need to create the main administrator account.
* **Endpoint:** `1. Authentication` -> `Register`
* **Body (Raw JSON):**
```json
{
  "phone": "02099999999",
  "password": "adminpassword",
  "full_name": "Admin User",
  "email": "admin@example.com",
  "role": "admin"
}
```
* **Action:** Send the request. Note down the returned `access_token`.
* **Postman Variable:** Go to your collection variables and set the `access_token` variable to the token you just received.

### Step 2: Register a Physical Device
The admin needs to register a physical palm scanner into the system so it has a valid `device_code`.
* **Endpoint:** `5. Web Admin` -> `Create Device`
* **Body (Raw JSON):**
```json
{
  "device_code": "DEV-001",
  "name": "Main Entrance Scanner",
  "location": "Lobby"
}
```
* **Action:** Send the request. The system now recognizes `DEV-001`.

---

## Phase 2: Create an Employee

### Step 3: Register an Employee
Let's pretend a new employee is registering via the mobile app.
* **Endpoint:** `1. Authentication` -> `Register`
* **Body (Raw JSON):**
```json
{
  "phone": "02011111111",
  "password": "userpassword",
  "full_name": "John Doe",
  "email": "john@example.com",
  "role": "employee"
}
```
* **Action:** Send the request. Copy the `access_token` returned.
* **Important:** In Postman, change your collection variable `access_token` to THIS employee's token for Phase 3!

---

## Phase 3: The Palm Enrollment Flow (QR Pairing)

This is the most complex flow, simulating the interaction between the Hardware Device and the Employee's Mobile App.

### Step 4: Device Requests a Pairing Session
The physical scanner generates a QR code on its screen.
* **Endpoint:** `4. Hardware Device` -> `Create Pairing Session (QR)`
* **Body (Raw JSON):**
```json
{
  "device_code": "DEV-001",
  "purpose": "enrollment"
}
```
* **Action:** Send request. Note down the `session_id` and the `qr_code_data` from the response.

### Step 5: Employee Scans the QR Code
The employee scans the QR code using their mobile app (using their employee access token).
* **Endpoint:** `3. Pairing Flow` -> `Scan QR Code`
* **Body (Raw JSON):**
```json
{
  "session_token": "<paste_the_qr_code_data_here>"
}
```
* **Action:** Send request. The server will respond with the device info ("Main Entrance Scanner").

### Step 6: Employee Approves the Pairing
The employee clicks "Approve" on their phone to link their account to the scanner.
* **Endpoint:** `3. Pairing Flow` -> `Approve Pairing`
* **Body (Raw JSON):**
```json
{
  "session_id": "<paste_the_session_id_here>",
  "purpose": "enrollment"
}
```
* **Action:** Send request. 

### Step 7: Device Checks Session Status
The scanner checks if the user approved the session.
* **Endpoint:** `4. Hardware Device` -> `Check Pairing Status`
* *(Replace `session-uuid` in the URL with your `session_id`)*
* **Action:** Send request. The `status` should now be `"approved"`, and it will return a `session_token`. Note this token!

### Step 8: Device Enrolls the Palm
The scanner reads the user's palm and sends the data to the server using the approved session.
* **Endpoint:** `4. Hardware Device` -> `Enroll Palm`
* **Body (Raw JSON):**
```json
{
  "device_code": "DEV-001",
  "session_token": "<paste_the_approved_session_token_here>",
  "hand_side": "right",
  "model_version": "v1.0",
  "embedding_dim": 2,
  "embeddings": [[0.15, 0.88]],
  "liveness_passed": true,
  "quality_score": 0.95,
  "thermal_min": 34.1,
  "thermal_max": 36.5,
  "thermal_avg": 35.3
}
```
* **Action:** Send request. The palm template is now saved to the employee!

---

## Phase 4: Taking Attendance

### Step 9: Process Attendance
The employee arrives the next morning and scans their palm. The device sends the palm vectors to the server to find a match and record attendance.
* **Endpoint:** `4. Hardware Device` -> `Process Attendance`
* **Body (Raw JSON):**
```json
{
  "device_code": "DEV-001",
  "model_version": "v1.0",
  "embedding_dim": 2,
  "embeddings": [[0.15, 0.88]], 
  "liveness_passed": true,
  "quality_score": 0.98,
  "thermal_min": 34.0,
  "thermal_max": 36.6,
  "thermal_avg": 35.5
}
```
* **Action:** Send request. The server will match the embeddings `[0.15, 0.88]` to John Doe, create an attendance record, and return success!

---

## Phase 5: Verifying the Data

### Step 10: Employee Checks Their Attendance
* **Endpoint:** `2. Mobile App (User)` -> `View My Attendance`
* **Action:** Send request (still using employee's access token). You should see the attendance log you just created.

### Step 11: Admin Checks Everything
* Change your Postman `access_token` variable back to the **Admin token** from Step 1.
* **Endpoint:** `5. Web Admin` -> `Get Users` (See that John Doe has `is_palm_registered: true`).
* **Endpoint:** `5. Web Admin` -> `Get Global Attendance` (See the company-wide attendance log).
