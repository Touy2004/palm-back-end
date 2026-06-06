# Testing Failed Authentication Attempts

This guide provides examples on how to trigger different types of **Failed Attempts** in the system. You can test these by sending a `POST` request to the device attendance endpoint.

> **Note:** These are **Hardware APIs**, so they do **NOT** require an `Authorization: Bearer` token. However, you must include the `Content-Type: application/json` header and a valid `device_code` in the body.

Ensure that you have an active device registered in your system and replace `"DEV-001"` with your actual `device_code`.

---

### 1. Liveness Check Failed
This happens when the hardware sensor detects a fake hand (e.g., a photograph or silicone mold).

**cURL Command:**
```bash
curl -X POST https://api.phoudthasone.com/api/v1/api/v1/devices/attendance/palm \
  -H "Content-Type: application/json" \
  -d '{
    "device_code": "DEV-001",
    "model_version": "v1.0",
    "embedding_dim": 128,
    "embeddings": [[0.12, 0.44]],
    "liveness_passed": false,
    "quality_score": 0.98,
    "thermal_min": 33.5,
    "thermal_max": 36.2,
    "thermal_avg": 35.1
  }'
```
* **Expected Result:** The server returns `400 Bad Request` and logs a failed attempt in the database with the reason `"Liveness check failed"`.

---

### 2. Palm Not Recognized
This happens when a real hand is scanned (liveness passed), but the embedding features do not closely match any enrolled templates in the database (similarity score is too low).

**cURL Command:**
```bash
curl -X POST https://api.phoudthasone.com/api/v1/api/v1/devices/attendance/palm \
  -H "Content-Type: application/json" \
  -d '{
    "device_code": "DEV-001",
    "model_version": "v1.0",
    "embedding_dim": 128,
    "embeddings": [[0.0, 0.0, 0.0]], 
    "liveness_passed": true,
    "quality_score": 0.85,
    "thermal_min": 34.0,
    "thermal_max": 36.5,
    "thermal_avg": 35.8
  }'
```
* **Expected Result:** The server returns `400 Bad Request` and logs a failed attempt with the reason `"Palm not recognized"`. 
*(Note: Sending an array of `0.0` guarantees a low similarity score against real templates).*

---

### 3. User Inactive
This happens when a valid palm is scanned and successfully recognized, but the employee's account has been deactivated or suspended by the Admin.

**How to test:**
1. Enroll a palm for a test user successfully.
2. Go to the Web Admin dashboard and change that user's status to `inactive`.
3. Send a valid attendance payload containing that user's actual embeddings.

**cURL Command:** *(Replace embeddings with the actual enrolled embeddings)*
```bash
curl -X POST https://api.phoudthasone.com/api/v1/api/v1/devices/attendance/palm \
  -H "Content-Type: application/json" \
  -d '{
    "device_code": "DEV-001",
    "model_version": "v1.0",
    "embedding_dim": 128,
    "embeddings": [[0.12, 0.44]], 
    "liveness_passed": true,
    "quality_score": 0.95,
    "thermal_min": 34.0,
    "thermal_max": 36.5,
    "thermal_avg": 35.8
  }'
```
* **Expected Result:** The server returns `400 Bad Request` and logs a failed attempt with the reason `"User inactive or not found"`.

---

### 4. Poor Thermal Quality (Hardware Diagnostics)
While poor thermal quality won't automatically fail the authentication *if* the embeddings match perfectly, you can simulate an invalid thermal read (e.g., a cold object) to see how it looks in the Admin Dashboard logs.

**cURL Command:**
```bash
curl -X POST https://api.phoudthasone.com/api/v1/api/v1/devices/attendance/palm \
  -H "Content-Type: application/json" \
  -d '{
    "device_code": "DEV-001",
    "model_version": "v1.0",
    "embedding_dim": 128,
    "embeddings": [[0.0, 0.0]], 
    "liveness_passed": false,
    "quality_score": 0.30,
    "thermal_min": 15.0,
    "thermal_max": 18.0,
    "thermal_avg": 16.5
  }'
```
* **Dashboard View:** When you check the "Failed Attempts" table in the web admin, you will clearly see the thermal average is `16.5°C` (abnormally cold) and the quality score is `0.30`, giving you diagnostic data on why the device is failing!
