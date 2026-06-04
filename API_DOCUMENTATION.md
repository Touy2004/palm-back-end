# Palm Recognition System - API Documentation

This documentation covers **every** endpoint in the system, organized by application roles: Public/Auth, Mobile App (User), Hardware Device (Raspberry Pi), and Web Admin (Admin).

**Base URL:** `http://localhost:8080/api/v1`

---

## 1. Authentication & Public Workflow

### 1.1 Register a new account
**POST** `/auth/register`
* **Body:**
  ```json
  {
    "phone": "0812345678",
    "password": "securepassword",
    "full_name": "Jane Doe",
    "email": "jane@example.com"
  }
  ```
* **Response (201 Created):** Returns user details and tokens.

### 1.2 Login
**POST** `/auth/login`
* **Body:**
  ```json
  {
    "phone": "0812345678",
    "password": "securepassword"
  }
  ```
* **Response (200 OK):** Returns `user`, `access_token`, and `refresh_token`.

### 1.3 Refresh Token
**POST** `/auth/refresh`
* **Body:**
  ```json
  {
    "refresh_token": "eyJhb..."
  }
  ```
* **Response (200 OK):** Returns a new `access_token` and `refresh_token`.

---

## 2. Mobile App (User) APIs

*Requires Header: `Authorization: Bearer <access_token>`*

### 2.1 Get My Profile
**GET** `/me`
* **Response:** Returns the authenticated user's JSON details.

### 2.2 Change Password
**PATCH** `/me/password`
* **Body:** `{"old_password": "...", "new_password": "..."}`

### 2.3 View My Attendance
**GET** `/me/attendance?page=1&limit=30`
* **Response:** Returns paginated attendance logs for the logged-in user.

### 2.4 Manage Enrolled Palms
**GET** `/me/palm-templates`
* **Response:** Returns list of active palm templates registered to this user.

**DELETE** `/me/palm-templates/:id`
* **Response:** Revokes a specific palm template.

---

## 3. Pairing Flow (Mobile App <-> Device)

*Requires Header: `Authorization: Bearer <access_token>`*

### 3.1 Scan QR Code
**POST** `/pairing/scan`
* **Body:** `{"session_token": "hex-string-from-qr"}`
* **Response:** Returns pairing session details.

### 3.2 Approve Pairing
**POST** `/pairing/approve`
* **Body:** `{"session_id": "uuid", "purpose": "enrollment"}`
* **Response:** Approves the session, linking the user to the physical device.

---

## 4. Hardware Device APIs (Raspberry Pi)

*Authenticates using `device_code` in the JSON body. No JWT required.*

### 4.1 Device Heartbeat
**POST** `/devices/heartbeat`
* **Body:** `{"device_code": "DEV-001"}`
* **Response:** Updates `last_seen_at`.

### 4.2 Create Pairing Session (Generate QR)
**POST** `/devices/pairing-sessions`
* **Body:** `{"device_code": "DEV-001", "purpose": "enrollment"}`
* **Response:** Returns `session_id` and `session_token` (used to display QR).

### 4.3 Check Pairing Status
**GET** `/devices/pairing-sessions/:session_id/status`
* **Response:** Returns status (`pending`, `scanned`, `approved`, `completed`).

### 4.4 Enroll Palm
**POST** `/devices/palm/enroll`
* **Body:**
  ```json
  {
    "device_code": "DEV-001",
    "session_token": "hex-string",
    "hand_side": "right",
    "model_version": "v1.0",
    "embedding_dim": 128,
    "embeddings": [[0.12, 0.44...]],
    "liveness_passed": true,
    "quality_score": 0.98,
    "thermal_min": 33.5,
    "thermal_max": 36.2,
    "thermal_avg": 35.1
  }
  ```

### 4.5 Identify Palm (No Attendance)
**POST** `/devices/palm/identify`
* **Body:** Requires `device_code`, `embedding`, and thermal/quality metrics.
* **Response:** Returns the identified user.

### 4.6 Process Attendance (Check In/Out)
**POST** `/devices/attendance/palm`
* **Body:** Same as 4.5.
* **Response:** Returns user details and action (`check_in` or `check_out`).

---

## 5. Web Admin APIs

*Requires Header: `Authorization: Bearer <access_token>`*
*Requires Role: `admin`*

### 5.1 User Management
**GET** `/admin/users?page=1&limit=20`
* **Response:** Paginated list of all users.

**GET** `/admin/users/search?q=John`
* **Response:** Searches users by name, email, phone, or employee_code.

**GET** `/admin/users/:id`
* **Response:** Get details of a specific user.

**POST** `/admin/users`
* **Body:** `{"phone": "...", "full_name": "...", "password": "...", "role": "employee", "department": "IT"}`
* **Response:** Creates a new user (Admin portal creation).

**PATCH** `/admin/users/:id`
* **Body:** `{"full_name": "Jane", "department": "HR", "status": "active"}`
* **Response:** Updates user info.

**DELETE** `/admin/users/:id`
* **Response:** Deletes a user completely.

### 5.2 Device Management
**GET** `/admin/devices?page=1&limit=20`
* **Response:** Paginated list of registered devices.

**POST** `/admin/devices`
* **Body:** `{"device_code": "DEV-002", "name": "Entrance B", "location": "Lobby"}`
* **Response:** Registers a new physical scanner device.

**PATCH** `/admin/devices/:id`
* **Body:** `{"name": "New Name", "status": "inactive"}`
* **Response:** Updates device info.

### 5.3 Global Attendance Monitoring
**GET** `/admin/attendance?page=1&limit=50`
* **Response:** View global attendance history across the whole company.

**GET** `/admin/attendance/users/:user_id/history?page=1&limit=30`
* **Response:** View the attendance history of a specific user.
