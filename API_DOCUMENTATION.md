# Palm Recognition System - API Documentation

This documentation is designed to help Frontend and Mobile Developers connect to the Palm Recognition Backend API. 

**Base URL:** `http://localhost:8080/api/v1` (Update port/host based on your `.env` configuration)

---

## 1. Authentication & Security Workflow

All protected routes require an `Authorization` header containing the JWT access token.
**Format:** `Authorization: Bearer <access_token>`

### 1.1 Login
**POST** `/auth/login`
Logs in the user and returns an access token (for requests) and a refresh token (to get a new access token without logging in again).
* **Body:**
  ```json
  {
    "phone": "0812345678",
    "password": "mysecretpassword"
  }
  ```
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "user": { "id": "uuid", "full_name": "John Doe", "role": "employee" },
    "access_token": "eyJhb...",
    "refresh_token": "eyJhb..."
  }
  ```

### 1.2 Refresh Token
**POST** `/auth/refresh`
When the access token expires, use this to get a new pair.
* **Body:**
  ```json
  {
    "refresh_token": "eyJhb..."
  }
  ```
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "access_token": "new_eyJhb...",
    "refresh_token": "new_eyJhb..."
  }
  ```

---

## 2. Mobile App: User Profile & Management Workflow

*Note: All endpoints below require the `Authorization: Bearer <access_token>` header.*

### 2.1 Get My Profile
**GET** `/me`
Fetches the currently authenticated user's details.
* **Response (200 OK):** Returns the user JSON object.

### 2.2 Change Password
**PATCH** `/me/password`
Allows the user to securely change their password.
* **Body:**
  ```json
  {
    "old_password": "currentPassword123",
    "new_password": "newSecurePassword456"
  }
  ```

### 2.3 View My Attendance
**GET** `/me/attendance?page=1&limit=30`
Fetches the user's daily attendance logs.
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "data": [
      {
         "attendance_date": "2026-06-03",
         "check_in_time": "2026-06-03T08:00:00Z",
         "check_out_time": "2026-06-03T17:00:00Z",
         "status": "present"
      }
    ],
    "pagination": { "page": 1, "limit": 30, "total": 1 }
  }
  ```

### 2.4 Manage Enrolled Palms
**GET** `/me/palm-templates`
Returns all palm templates currently enrolled for the user.

**DELETE** `/me/palm-templates/:id`
Deletes a specific palm template by its UUID.

---

## 3. The QR Pairing Flow (Mobile App + Hardware Device)

This is the flow used when a user wants to enroll a new palm using the physical scanner device.

### Step 1: Device requests a QR Code
**POST** `/devices/pairing-sessions` (Called by Raspberry Pi)
* **Body:** `{"device_code": "DEV-001", "purpose": "enrollment"}`
* **Response:** Returns a `session_id` and a `session_token` (The device converts `session_token` into a visual QR code on its screen).

### Step 2: Mobile App Scans the QR Code
**POST** `/pairing/scan` (Called by Mobile App, requires Auth)
* **Body:**
  ```json
  {
    "session_token": "hex-string-from-qr-code"
  }
  ```
* **Response (200 OK):** Returns the `session_id` and device details so the mobile app can ask the user "Do you want to connect to Scanner DEV-001?"

### Step 3: Mobile App Approves the Pairing
**POST** `/pairing/approve` (Called by Mobile App, requires Auth)
* **Body:**
  ```json
  {
    "session_id": "uuid-from-step-2",
    "purpose": "enrollment"
  }
  ```
* **Result:** The backend links the user's ID to the device session. The device screen will now update and ask the user to place their hand on the scanner.

---

## 4. Hardware Device Operations (Raspberry Pi)

*Note: Device endpoints do not use JWT. They authenticate using the physical `device_code`.*

### 4.1 Device Heartbeat (Status check)
**POST** `/devices/heartbeat`
* **Body:** `{"device_code": "DEV-001"}`

### 4.2 Enroll a Palm (After QR Pairing)
**POST** `/devices/palm/enroll`
Saves the physical palm scan into the database (encrypted).
* **Body:**
  ```json
  {
    "device_code": "DEV-001",
    "session_token": "hex-string-from-qr-code",
    "hand_side": "right",
    "model_version": "v1.0",
    "embedding_dim": 128,
    "embeddings": [[0.12, 0.44, 0.55, ...]], 
    "liveness_passed": true,
    "quality_score": 0.98,
    "thermal_min": 33.5,
    "thermal_max": 36.2,
    "thermal_avg": 35.1
  }
  ```

### 4.3 Check-in / Check-out (Daily Attendance)
**POST** `/devices/attendance/palm`
The user places their hand on the scanner. The backend figures out who they are using Cosine Similarity and automatically records a Check-in or Check-out.
* **Body:**
  ```json
  {
    "device_code": "DEV-001",
    "model_version": "v1.0",
    "embedding_dim": 128,
    "embedding": [0.12, 0.44, 0.55, ...], 
    "liveness_passed": true,
    "quality_score": 0.98,
    "thermal_min": 33.5,
    "thermal_max": 36.2,
    "thermal_avg": 35.1
  }
  ```
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "action": "check_in",
    "user": { "id": "uuid", "full_name": "John Doe" },
    "message": "Check-in success"
  }
  ```

### 4.4 Pure Identification (No Attendance Logging)
**POST** `/devices/palm/identify`
Same body as `4.3`, but this endpoint simply returns who the user is without saving anything to the attendance logs (useful for unlocking a door).
