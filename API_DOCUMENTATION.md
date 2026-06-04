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
* **Response (201 Created):**
  ```json
  {
    "success": true,
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "full_name": "Jane Doe",
      "role": "employee"
    },
    "access_token": "eyJhbGciOi...",
    "refresh_token": "eyJhbGciOi..."
  }
  ```

### 1.2 Login
**POST** `/auth/login`
* **Body:**
  ```json
  {
    "phone": "0812345678",
    "password": "securepassword"
  }
  ```
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "full_name": "Jane Doe",
      "role": "employee"
    },
    "access_token": "eyJhbGciOi...",
    "refresh_token": "eyJhbGciOi..."
  }
  ```

### 1.3 Refresh Token
**POST** `/auth/refresh`
* **Body:**
  ```json
  {
    "refresh_token": "eyJhbGciOi..."
  }
  ```
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "access_token": "new_eyJhbGciOi...",
    "refresh_token": "new_eyJhbGciOi..."
  }
  ```

---

## 2. Mobile App (User) APIs

*Requires Header: `Authorization: Bearer <access_token>`*

### 2.1 Get My Profile
**GET** `/me`
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "employee_code": "EMP-001",
      "full_name": "Jane Doe",
      "email": "jane@example.com",
      "phone": "0812345678",
      "role": "employee",
      "department": "IT",
      "status": "active",
      "created_at": "2026-06-03T10:00:00Z"
    }
  }
  ```

### 2.2 Change Password
**PATCH** `/me/password`
* **Body:** 
  ```json
  {
    "old_password": "currentpassword",
    "new_password": "newpassword123"
  }
  ```
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Password changed successfully"
  }
  ```

### 2.3 View My Attendance
**GET** `/me/attendance?page=1&limit=30`
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "data": [
      {
         "id": "att-123",
         "attendance_date": "2026-06-03T00:00:00Z",
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
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "data": [
      {
        "id": "tpl-123",
        "hand_side": "right",
        "status": "active",
        "created_at": "2026-06-01T12:00:00Z"
      }
    ]
  }
  ```

**DELETE** `/me/palm-templates/:id`
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Palm template revoked"
  }
  ```

---

## 3. Pairing Flow (Mobile App <-> Device)

*Requires Header: `Authorization: Bearer <access_token>`*

### 3.1 Scan QR Code
**POST** `/pairing/scan`
* **Body:** 
  ```json
  {
    "session_token": "a1b2c3d4e5f6..."
  }
  ```
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "session_id": "session-uuid",
    "device_id": "device-uuid",
    "device_name": "Front Door Scanner",
    "purpose": "enrollment"
  }
  ```

### 3.2 Approve Pairing
**POST** `/pairing/approve`
* **Body:** 
  ```json
  {
    "session_id": "session-uuid",
    "purpose": "enrollment"
  }
  ```
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Pairing approved successfully"
  }
  ```

---

## 4. Hardware Device APIs (Raspberry Pi)

*Authenticates using `device_code` in the JSON body. No JWT required.*

### 4.1 Device Heartbeat
**POST** `/devices/heartbeat`
* **Body:** `{"device_code": "DEV-001"}`
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Heartbeat updated"
  }
  ```

### 4.2 Create Pairing Session (Generate QR)
**POST** `/devices/pairing-sessions`
* **Body:** `{"device_code": "DEV-001", "purpose": "enrollment"}`
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "session_id": "session-uuid",
    "session_token": "a1b2c3d4e5f6...",
    "expires_at": "2026-06-04T10:05:00Z"
  }
  ```

### 4.3 Check Pairing Status
**GET** `/devices/pairing-sessions/:session_id/status`
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "status": "approved",
    "user_id": "user-uuid"
  }
  ```

### 4.4 Enroll Palm
**POST** `/devices/palm/enroll`
* **Body:**
  ```json
  {
    "device_code": "DEV-001",
    "session_token": "a1b2c3d4e5f6...",
    "hand_side": "right",
    "model_version": "v1.0",
    "embedding_dim": 128,
    "embeddings": [[0.12, 0.44]],
    "liveness_passed": true,
    "quality_score": 0.98,
    "thermal_min": 33.5,
    "thermal_max": 36.2,
    "thermal_avg": 35.1
  }
  ```
* **Response (201 Created):**
  ```json
  {
    "success": true,
    "message": "Palm enrolled successfully"
  }
  ```

### 4.5 Identify Palm (No Attendance)
**POST** `/devices/palm/identify`
* **Body:** Requires `device_code`, `embedding`, and thermal/quality metrics.
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "user": {
      "id": "user-uuid",
      "full_name": "Jane Doe"
    },
    "score": 0.95
  }
  ```

### 4.6 Process Attendance (Check In/Out)
**POST** `/devices/attendance/palm`
* **Body:** Same as 4.5.
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "action": "check_in",
    "user": {
      "id": "user-uuid",
      "full_name": "Jane Doe"
    },
    "message": "Check-in successful"
  }
  ```

---

## 5. Web Admin APIs

*Requires Header: `Authorization: Bearer <access_token>`*
*Requires Role: `admin`*

### 5.1 User Management
**GET** `/admin/users?page=1&limit=20`
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "data": [
      {
        "id": "user-uuid",
        "employee_code": "EMP-001",
        "full_name": "Jane Doe",
        "email": "jane@example.com",
        "role": "employee",
        "status": "active"
      }
    ],
    "pagination": { "page": 1, "limit": 20, "total": 100 }
  }
  ```

**GET** `/admin/users/search?q=John`
* **Response (200 OK):** (Same array of users as above, without pagination).

**GET** `/admin/users/:id`
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "data": {
        "id": "user-uuid",
        "employee_code": "EMP-001",
        "full_name": "Jane Doe",
        "status": "active"
    }
  }
  ```

**POST** `/admin/users`
* **Body:** `{"phone": "0811111111", "full_name": "John", "password": "pass", "role": "employee"}`
* **Response (201 Created):**
  ```json
  {
    "success": true,
    "message": "User created successfully",
    "data": { "id": "user-uuid" }
  }
  ```

**PATCH** `/admin/users/:id`
* **Body:** `{"full_name": "Jane", "department": "HR", "status": "active"}`
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "User updated"
  }
  ```

**DELETE** `/admin/users/:id`
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "User deleted"
  }
  ```

### 5.2 Device Management
**GET** `/admin/devices?page=1&limit=20`
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "data": [
      {
        "id": "device-uuid",
        "device_code": "DEV-001",
        "device_name": "Main Entrance",
        "location_name": "Lobby",
        "status": "active",
        "last_seen_at": "2026-06-04T10:00:00Z"
      }
    ]
  }
  ```

**POST** `/admin/devices`
* **Body:** `{"device_code": "DEV-002", "name": "Entrance B", "location": "Lobby"}`
* **Response (201 Created):**
  ```json
  {
    "success": true,
    "message": "Device registered"
  }
  ```

**PATCH** `/admin/devices/:id`
* **Body:** `{"name": "New Name", "status": "inactive"}`
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Device updated"
  }
  ```

### 5.3 Global Attendance Monitoring
**GET** `/admin/attendance?page=1&limit=50`
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "data": [
      {
         "id": "att-123",
         "user_id": "user-uuid",
         "device_id": "device-uuid",
         "attendance_date": "2026-06-03T00:00:00Z",
         "check_in_time": "2026-06-03T08:00:00Z",
         "check_out_time": "2026-06-03T17:00:00Z",
         "status": "present"
      }
    ],
    "pagination": { "page": 1, "limit": 50, "total": 1000 }
  }
  ```

**GET** `/admin/attendance/users/:user_id/history?page=1&limit=30`
* **Response (200 OK):** (Same format as above, filtered for the specific user ID).
