# Palm Recognition Attendance System - User Manual

Welcome to the User Manual! This document provides a high-level explanation of how the system works, its core workflows, and the data model behind it. This guide is designed to be easily understood by both administrators and standard users.

---

## 1. What Does This Project Do?

This project is a modern **Biometric Attendance System** that uses palm recognition technology. Instead of using keycards or passwords, employees simply place their hand over a scanner to check in and out of work. 

The system consists of three main parts:
1. **Web Admin Dashboard:** Used by HR or Administrators to manage users, devices, and view attendance reports.
2. **Physical Scanners (Hardware Devices):** The actual machines stationed at entrances that scan palms.
3. **Mobile App / Employee Portal:** Used by employees to securely pair their palm with their account and view their own attendance history.

---

## 2. Core Workflows (How to Use It)

### Phase 1: Initial Setup (Admin)
Before employees can start checking in, the system needs to be set up:
1. **Create Employees:** The Admin logs into the Web Dashboard and creates profiles for all employees (assigning them employee codes and roles).
2. **Register Devices:** The Admin registers the physical palm scanners in the system (e.g., "Main Entrance Scanner", "Back Door Scanner") so the server knows which machines are authorized.

### Phase 2: Palm Enrollment (The Pairing Flow)
An employee cannot check in until their palm is securely linked to their account. This is a one-time process:
1. **Start:** The administrator or employee clicks "Enroll New Palm" on the physical scanner's touchscreen.
2. **QR Code:** The scanner displays a unique QR code on its screen.
3. **Approve:** An Administrator scans the QR code using the admin app/dashboard and enters the employee's ID code. This temporarily links the scanner session to that specific employee.
4. **Scan Palm:** The scanner prompts the employee to place their hand over the sensor. The machine reads the palm, encrypts the data, and saves it to the employee's database profile.

### Phase 3: Daily Attendance (Checking In/Out)
Once enrolled, daily usage is completely frictionless:
1. **Walk Up & Scan:** The employee walks up to any registered scanner and places their hand over it.
2. **Automatic Logging:** The scanner instantly recognizes the palm and marks them as "Present" or "Late" for the day. If they scan again later, it updates their "Check Out" time.

### Phase 4: Reporting & Dashboard
1. **Live Monitoring:** The Admin can view a live dashboard showing who is present, absent, or late today.
2. **Generate Reports:** At the end of the month, the Admin can filter by Month and Department to generate automated attendance summaries for payroll.

---

## 3. Entity-Relationship (ER) Model

Below is the database structure that powers the system, showing how Users, Devices, Palm Templates, and Attendance Logs are connected.

```mermaid
erDiagram
    USERS {
        uuid id PK
        string employee_code "UNIQUE"
        string full_name
        string email "UNIQUE"
        string phone "UNIQUE"
        string password_hash
        string role "ADMIN or EMPLOYEE"
        string department
        string status
        timestamp created_at
    }

    DEVICES {
        uuid id PK
        string device_code "UNIQUE"
        string name
        string location
        string status
        timestamp last_seen_at
    }

    DEVICE_PAIRING_SESSIONS {
        uuid id PK
        uuid device_id FK
        uuid user_id FK
        string session_token "UNIQUE"
        string purpose "enrollment or verification"
        string status "pending, approved, completed"
        timestamp expires_at
    }

    PALM_TEMPLATES {
        uuid id PK
        uuid user_id FK
        uuid registered_device_id FK
        string hand_side
        bytea template_encrypted
        bytea template_nonce
        int embedding_dim
        string status
    }

    ATTENDANCE_LOGS {
        uuid id PK
        uuid user_id FK
        uuid device_id FK
        date attendance_date "UNIQUE per user"
        timestamp check_in_time
        timestamp check_out_time
        string status "present, late, incomplete, absent"
    }

    %% Relationships
    DEVICES ||--o{ DEVICE_PAIRING_SESSIONS : "creates"
    USERS ||--o{ DEVICE_PAIRING_SESSIONS : "approves"
    
    USERS ||--o{ PALM_TEMPLATES : "owns"
    DEVICES |o--o{ PALM_TEMPLATES : "registered via"
    
    USERS ||--o{ ATTENDANCE_LOGS : "logs"
    DEVICES |o--o{ ATTENDANCE_LOGS : "scanned by"
```

### 3.1 Chen's E-R Notation

Here is the database structure of the Palm Recognition Attendance System modeled in Chen's E-R notation style:

```mermaid
flowchart TD
    %% Styling
    classDef entity fill:#fff,stroke:#000,stroke-width:1px;
    classDef attribute fill:#fff,stroke:#000,stroke-width:1px,rx:20,ry:20;
    classDef relationship fill:#fff,stroke:#000,stroke-width:1px,shape:diamond;
    classDef primaryKey fill:#fff,stroke:#000,stroke-width:1px,rx:20,ry:20,text-decoration:underline;

    %% Entities
    User[User]:::entity
    PalmTemplate[Palm_Template]:::entity
    AttendanceLog[Attendance_Log]:::entity
    Device[Device]:::entity
    DevicePairingSession[Device_Pairing_Session]:::entity

    %% Relationships
    Has{Has}:::relationship
    Logs{Logs}:::relationship
    ScannedBy{Scanned_By}:::relationship
    Creates{Creates}:::relationship
    Approves{Approves}:::relationship

    %% Connections for relationships
    User ---|1| Has ---|N| PalmTemplate
    User ---|1| Logs ---|N| AttendanceLog
    Device ---|1| ScannedBy ---|N| AttendanceLog
    Device ---|1| Creates ---|N| DevicePairingSession
    User ---|1| Approves ---|N| DevicePairingSession

    %% Attributes for User
    U_ID([user_id]):::primaryKey
    U_EmpCode([employee_code]):::attribute
    U_Name([full_name]):::attribute
    U_Phone([phone]):::attribute
    U_Role([role]):::attribute
    U_Pass([password_hash]):::attribute

    User --- U_ID
    User --- U_EmpCode
    User --- U_Name
    User --- U_Phone
    User --- U_Role
    User --- U_Pass

    %% Attributes for PalmTemplate
    PT_ID([template_id]):::primaryKey
    PT_Hand([hand_side]):::attribute
    PT_Encrypted([template_encrypted blob]):::attribute
    PT_Status([status]):::attribute

    PalmTemplate --- PT_ID
    PalmTemplate --- PT_Hand
    PalmTemplate --- PT_Encrypted
    PalmTemplate --- PT_Status

    %% Attributes for AttendanceLog
    AL_ID([log_id]):::primaryKey
    AL_Date([attendance_date]):::attribute
    AL_InTime([check_in_time]):::attribute
    AL_Status([status]):::attribute
    AL_Score([confidence_score float]):::attribute

    AttendanceLog --- AL_ID
    AttendanceLog --- AL_Date
    AttendanceLog --- AL_InTime
    AttendanceLog --- AL_Status
    AttendanceLog --- AL_Score

    %% Attributes for Device
    D_ID([device_id]):::primaryKey
    D_Code([device_code]):::attribute
    D_Name([name]):::attribute

    Device --- D_ID
    Device --- D_Code
    Device --- D_Name

    %% Attributes for DevicePairingSession
    DPS_ID([session_id]):::primaryKey
    DPS_Token([session_token]):::attribute
    DPS_Purpose([purpose]):::attribute
    DPS_Status([status]):::attribute

    DevicePairingSession --- DPS_ID
    DevicePairingSession --- DPS_Token
    DevicePairingSession --- DPS_Purpose
    DevicePairingSession --- DPS_Status
```

### Understanding the Relationships:
* **Users & Palm Templates (1-to-Many):** One user can have multiple palm templates (e.g., left hand, right hand).
* **Users & Attendance (1-to-Many):** One user has many attendance logs (one per day).
* **Devices & Sessions (1-to-Many):** A device generates unique pairing sessions (QR codes) for users to scan.
* **Devices & Attendance (1-to-Many):** A device processes many attendance check-ins.

---

## 4. Process Hierarchy Chart

The Process Hierarchy Chart (or Functional Decomposition Diagram) breaks down the entire system into its core functional modules and their sub-processes.

```mermaid
flowchart TD
    %% Main System
    System[Palm Recognition Attendance System]:::root
    
    %% Level 1 Modules
    System --> UM[1. User Management]:::module
    System --> DM[2. Device Management]:::module
    System --> PE[3. Palm Enrollment]:::module
    System --> AP[4. Attendance Processing]:::module
    System --> RA[5. Reporting & Analytics]:::module
    
    %% Level 2 Processes
    UM --> UM1[1.1 Add Employee]:::process
    UM --> UM2[1.2 Edit Profile]:::process
    UM --> UM3[1.3 Manage Status]:::process
    
    DM --> DM1[2.1 Register Device]:::process
    DM --> DM2[2.2 Monitor Status]:::process
    
    PE --> PE1[3.1 Generate QR Code]:::process
    PE --> PE2[3.2 Admin Approval]:::process
    PE --> PE3[3.3 Capture Palm]:::process
    
    AP --> AP1[4.1 Scan Palm]:::process
    AP --> AP2[4.2 Verify Identity]:::process
    AP --> AP3[4.3 Log Attendance]:::process
    
    RA --> RA1[5.1 Live Dashboard]:::process
    RA --> RA2[5.2 Generate Reports]:::process
    RA --> RA3[5.3 Personal History]:::process

    %% Styling
    classDef root fill:#1e293b,stroke:#000,stroke-width:2px,color:#fff,font-weight:bold
    classDef module fill:#3b82f6,stroke:#1e40af,stroke-width:1px,color:#fff,font-weight:bold
    classDef process fill:#f8fafc,stroke:#94a3b8,stroke-width:1px,color:#334155
```

---

## 5. Data Flow Diagrams (DFD)

Data Flow Diagrams map out how information flows through the system, from external entities (users/devices) into processes and data stores.

### 5.1 DFD Level 0 (Context Diagram)
The Context Diagram provides a bird's-eye view of the entire system as a single process, showing its interactions with external entities.

```mermaid
flowchart LR
    %% External Entities
    Admin[Admin]:::entity
    Employee[Employee]:::entity
    Device[Hardware Scanner]:::entity

    %% Central System
    System((0.0<br/>Palm Recognition<br/>Attendance System)):::process

    %% Flows - Admin
    Admin -- "User & Device Data" --> System
    System -- "Reports & Dashboards" --> Admin

    %% Flows - Employee
    Employee -- "Login" --> System
    System -- "Attendance History & Status" --> Employee

    %% Flows - Device
    Device -- "Scanned Palm Vectors" --> System
    System -- "Pairing Sessions & Auth Results" --> Device

    %% Styling
    classDef entity fill:#f1f5f9,stroke:#334155,stroke-width:2px,shape:rect
    classDef process fill:#eff6ff,stroke:#2563eb,stroke-width:2px,shape:circle
```

### 5.2 DFD Level 1 (Main Processes)
Level 1 breaks down the main system into its primary sub-processes and shows how they interact with the database stores.

```mermaid
flowchart TD
    %% External Entities
    Admin[Admin]:::entity
    Employee[Employee]:::entity
    Device[Hardware Scanner]:::entity

    %% Processes
    P1((1.0<br/>Manage Setup)):::process
    P2((2.0<br/>Enroll Palm)):::process
    P3((3.0<br/>Process Attendance)):::process
    P4((4.0<br/>Generate Reports)):::process

    %% Data Stores
    D1[(D1: Users DB)]:::datastore
    D2[(D2: Devices DB)]:::datastore
    D3[(D3: Templates DB)]:::datastore
    D4[(D4: Attendance DB)]:::datastore

    %% Admin Setup Flows
    Admin -- "User/Device Info" --> P1
    P1 -- "Save Data" --> D1
    P1 -- "Save Data" --> D2

    %% Palm Enrollment Flows
    Device -- "Request Session" --> P2
    Employee -- "Scan Palm" --> P2
    Admin -- "Approve Session (Enter Emp Code)" --> P2
    P2 -- "Palm Template" --> D3
    P2 -- "Save Template" --> D3
    Device -- "Send Palm Data" --> P2

    %% Attendance Processing Flows
    Device -- "Send Scan Data" --> P3
    P3 -- "Fetch Template" --> D3
    P3 -- "Save Check-in" --> D4
    P3 -- "Return Result" --> Device

    %% Reporting Flows
    Admin -- "Request Reports" --> P4
    Employee -- "View My Logs" --> P4
    D4 -- "Fetch Logs" --> P4
    P4 -- "Report Data" --> Admin
    P4 -- "History Data" --> Employee

    %% Styling
    classDef entity fill:#f1f5f9,stroke:#334155,stroke-width:2px,shape:rect
    classDef process fill:#eff6ff,stroke:#2563eb,stroke-width:2px,shape:circle
    classDef datastore fill:#fcfdfd,stroke:#0f172a,stroke-width:2px,shape:cylinder
```

### 5.3 DFD Level 2 (Process Decomposition)
Below are the decomposed Data Flow Diagrams for each of the main processes (1.0 to 4.0), showing the specific sub-processes involved.

#### Process 1: Manage Setup
```mermaid
flowchart LR
    Admin[Admin]:::entity
    
    P1_1((1.1<br/>Create/Edit User)):::process
    P1_2((1.2<br/>Register Device)):::process
    
    D1[(D1: Users DB)]:::datastore
    D2[(D2: Devices DB)]:::datastore
    
    Admin -- "User Data" --> P1_1
    Admin -- "Device Data" --> P1_2
    
    P1_1 -- "Write Profile" --> D1
    P1_2 -- "Write Device details" --> D2

    classDef entity fill:#f1f5f9,stroke:#334155,stroke-width:2px,shape:rect
    classDef process fill:#eff6ff,stroke:#2563eb,stroke-width:2px,shape:circle
    classDef datastore fill:#fcfdfd,stroke:#0f172a,stroke-width:2px,shape:cylinder
```

#### Process 2: Enroll Palm
```mermaid
flowchart TD
    Device[Hardware Scanner]:::entity
    Employee[Employee]:::entity
    Admin[Admin]:::entity
    
    P2_1((2.1<br/>Generate QR Session)):::process
    P2_2((2.2<br/>Admin Approves Session)):::process
    P2_3((2.3<br/>Scan & Extract Palm Vectors)):::process
    
    D1[(D1: Users DB)]:::datastore
    D3[(D3: Templates DB)]:::datastore
    
    Device -- "Display Session QR" --> P2_1
    Admin -- "Scan QR & Enter Emp Code" --> P2_2
    P2_2 -- "Session Approved" --> P2_3
    P2_1 -- "Show QR" --> Device
    
    P2_2 -- "Validate User" --> D1
    P2_2 -- "Session Approved" --> P2_3
    
    Device -- "Send Palm Data" --> P2_3
    P2_3 -- "Save Secure Vector" --> D3

    classDef entity fill:#f1f5f9,stroke:#334155,stroke-width:2px,shape:rect
    classDef process fill:#eff6ff,stroke:#2563eb,stroke-width:2px,shape:circle
    classDef datastore fill:#fcfdfd,stroke:#0f172a,stroke-width:2px,shape:cylinder
```

#### Process 3: Process Attendance
```mermaid
flowchart TD
    Device[Hardware Scanner]:::entity
    
    P3_1((3.1<br/>Read Palm Vector)):::process
    P3_2((3.2<br/>Match Identity)):::process
    P3_3((3.3<br/>Record Attendance Log)):::process
    
    D3[(D3: Templates DB)]:::datastore
    D4[(D4: Attendance DB)]:::datastore
    
    Device -- "Scanned Palm Data" --> P3_1
    P3_1 -- "Forward Vector" --> P3_2
    
    P3_2 -- "Fetch Templates" --> D3
    P3_2 -- "Matched User ID" --> P3_3
    
    P3_3 -- "Log check-in/out" --> D4
    P3_3 -- "Success/Fail Result" --> Device

    classDef entity fill:#f1f5f9,stroke:#334155,stroke-width:2px,shape:rect
    classDef process fill:#eff6ff,stroke:#2563eb,stroke-width:2px,shape:circle
    classDef datastore fill:#fcfdfd,stroke:#0f172a,stroke-width:2px,shape:cylinder
```

#### Process 4: Generate Reports
```mermaid
flowchart LR
    Admin[Admin]:::entity
    Employee[Employee]:::entity
    
    P4_1((4.1<br/>Query Daily History)):::process
    P4_2((4.2<br/>Compute Monthly Summary)):::process
    P4_3((4.3<br/>Export Formatted Data)):::process
    
    D1[(D1: Users DB)]:::datastore
    D4[(D4: Attendance DB)]:::datastore
    
    Employee -- "Request History" --> P4_1
    Admin -- "Request Summary" --> P4_2
    
    D4 -- "Raw Logs" --> P4_1
    D4 -- "Raw Logs" --> P4_2
    D1 -- "User Depts/Roles" --> P4_2
    
    P4_1 -- "Process Data" --> P4_3
    P4_2 -- "Aggregate Data" --> P4_3
    
    P4_3 -- "Dashboard & PDF/Excel" --> Admin
    P4_3 -- "List View" --> Employee

    classDef entity fill:#f1f5f9,stroke:#334155,stroke-width:2px,shape:rect
    classDef process fill:#eff6ff,stroke:#2563eb,stroke-width:2px,shape:circle
    classDef datastore fill:#fcfdfd,stroke:#0f172a,stroke-width:2px,shape:cylinder
```

---

## 6. Data Dictionary

The following tables define the structure, data types, and constraints for each table in the database.

### 6.1 `users`
Stores all employee and administrator profiles.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | Primary Key | Unique identifier for the user |
| `employee_code` | VARCHAR(50) | Unique, Not Null | Public ID/Code for the employee |
| `full_name` | VARCHAR(100) | Not Null | User's full name |
| `email` | VARCHAR(100) | Unique | User's email address |
| `phone` | VARCHAR(20) | Unique | User's phone number |
| `password_hash` | VARCHAR(255) | Not Null | Bcrypt hashed password |
| `role` | VARCHAR(20) | Default 'EMPLOYEE' | 'ADMIN' or 'EMPLOYEE' |
| `department` | VARCHAR(50) | | Department the user belongs to |
| `status` | VARCHAR(20) | Default 'active' | User account status (active/inactive) |
| `created_at` | TIMESTAMP | Default NOW() | Profile creation time |
| `updated_at` | TIMESTAMP | Default NOW() | Profile last update time |

### 6.2 `devices`
Stores physical hardware scanners registered to the system.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | Primary Key | Unique identifier for the device |
| `device_code` | VARCHAR(50) | Unique, Not Null | Serial number or hardware ID |
| `name` | VARCHAR(100) | Not Null | Human-readable name (e.g. "Lobby Scanner") |
| `location` | VARCHAR(100) | | Physical location of the device |
| `status` | VARCHAR(20) | Default 'active' | Device operational status |
| `last_seen_at` | TIMESTAMP | | Last time the device pinged the server |
| `created_at` | TIMESTAMP | Default NOW() | Registration time |

### 6.3 `device_pairing_sessions`
Temporary sessions created when a scanner generates a QR code for palm enrollment.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | Primary Key | Session identifier |
| `device_id` | UUID | Foreign Key | The scanner generating the session |
| `session_token` | TEXT | Unique, Not Null | The token embedded in the QR code |
| `user_id` | UUID | Foreign Key | The user who approved the session |
| `purpose` | VARCHAR(50) | Not Null | Usually 'enrollment' |
| `status` | VARCHAR(30) | Default 'pending' | pending, approved, completed |
| `expires_at` | TIMESTAMP | Not Null | When the QR code expires |
| `scanned_at` | TIMESTAMP | | When the user scanned the QR |
| `approved_at` | TIMESTAMP | | When the user clicked "Approve" |
| `completed_at` | TIMESTAMP | | When the palm vector was saved |
| `created_at` | TIMESTAMP | Default NOW() | Session creation time |

### 6.4 `palm_templates`
Securely stores the encrypted biometric vectors.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | Primary Key | Template identifier |
| `user_id` | UUID | Foreign Key, Not Null | The owner of the template |
| `hand_side` | VARCHAR(10) | Not Null | 'left' or 'right' |
| `template_encrypted` | BYTEA | Not Null | AES-256 encrypted vector data |
| `template_nonce` | BYTEA | Not Null | Nonce used for decryption |
| `embedding_dim` | INT | Default 128 | Vector dimensions |
| `model_version` | VARCHAR(100) | Not Null | AI model used for extraction |
| `threshold` | NUMERIC(5,4) | Default 0.8200 | Match threshold for this template |
| `status` | VARCHAR(30) | Default 'active' | Template status |
| `registered_device_id`| UUID | Foreign Key | Device used to capture template |
| `created_at` | TIMESTAMP | Default NOW() | Enrollment time |
| `updated_at` | TIMESTAMP | Default NOW() | Last update time |
| `revoked_at` | TIMESTAMP | | If the template was revoked |

### 6.5 `attendance_logs`
Records daily check-ins and check-outs for all users.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | Primary Key | Log identifier |
| `user_id` | UUID | Foreign Key, Not Null | Employee who checked in |
| `device_id` | UUID | Foreign Key | Scanner used for attendance |
| `attendance_date` | DATE | Not Null | The calendar date of the log |
| `check_in_time` | TIMESTAMP | | Exact time of first scan |
| `check_out_time` | TIMESTAMP | | Exact time of latest scan |
| `check_in_score` | NUMERIC(6,5) | | Biometric confidence score (0-1) |
| `check_out_score`| NUMERIC(6,5) | | Biometric confidence score (0-1) |
| `check_in_liveness`| BOOLEAN | Default false | Passed hardware anti-spoofing |
| `check_out_liveness`| BOOLEAN | Default false | Passed hardware anti-spoofing |
| `status` | VARCHAR(20) | Default 'present' | present, late, incomplete, absent |
| `created_at` | TIMESTAMP | Default NOW() | Time the log was created |

*(Note: `user_id` + `attendance_date` is a UNIQUE compound key to prevent duplicate daily records).*

---

## 7. Step-by-Step Flowcharts

These flowcharts break down the exact step-by-step logic and decision paths for the system as a whole, and for each specific role (Admin, Employee, and Hardware Device).

### 7.1 System Overview Flowchart
This shows the high-level life cycle from system setup to daily usage.

```mermaid
flowchart TD
    Start([Start]) --> Setup[Admin Registers Devices & Employees]
    Setup --> App[Employee Logs into Mobile App]
    App --> Enroll[Palm Enrollment Phase]
    Enroll --> QR[Hardware shows QR Code]
    QR --> Scan[Admin Scans QR & Approves]
    Scan --> Capture[Hardware Captures Palm Vector]
    Capture --> Daily[Daily Attendance Phase]
    Daily --> ScanPalm[Employee Places Hand on Scanner]
    ScanPalm --> Verify{Match Found?}
    Verify -- Yes --> Log[Log Attendance as Present/Late]
    Verify -- No --> Reject[Show Access Denied]
    Log --> Dashboard[Admin Views Live Dashboard & Reports]
    Reject --> Daily
    Dashboard --> End([End])
```

### 7.2 Admin Flowchart
The step-by-step flow for an Administrator using the Web Dashboard.

```mermaid
flowchart TD
    Start([Admin Login]) --> Dashboard[View Live Dashboard]
    Dashboard --> Choice{Choose Action}
    
    Choice -- Manage Users --> UserMenu[User Management]
    UserMenu --> AddUser[Add/Edit Employee Details]
    AddUser --> SaveUser[(Save to DB)]
    SaveUser --> Choice
    
    Choice -- Manage Devices --> DeviceMenu[Device Management]
    DeviceMenu --> AddDevice[Register New Scanner]
    AddDevice --> SaveDevice[(Save to DB)]
    SaveDevice --> Choice
    
    Choice -- View Reports --> ReportsMenu[Monthly Reports]
    ReportsMenu --> Filter[Select Month & Department]
    Filter --> View[View Aggregated Data]
    View --> Export[Export PDF / Excel]
    Export --> Choice
    
    Choice -- Logout --> End([End Session])
```

### 7.3 Employee Flowchart
The step-by-step flow for an Employee using the Mobile App.

```mermaid
flowchart TD
    Start([Open App]) --> Login[Login with Employee Credentials]
    Login --> Home[View Home Dashboard]
    Home --> Choice{Choose Action}
    
    Choice -- View History --> History[Check Personal Logs]
    History --> List[View Daily Check-in/out Times]
    List --> Choice
    
    Choice -- Enroll Palm --> ScannerQR[Walk to Hardware Scanner]
    ScannerQR --> Cam[Open App Camera]
    Cam --> Scan[Scan Scanner's QR Code]
    Scan -- Valid? --> Valid{Is Session Valid?}
    Valid -- No --> ScanErr[Show Error]
    Valid -- Yes --> Approve[Admin Approves & Enters Emp Code]
    Approve --> Wait[Employee Places Hand on Scanner]
    Wait --> Success[Enrollment Complete]
    Success --> Choice
    
    Choice -- Logout --> End([End Session])
```

### 7.4 Hardware Device (Raspberry Pi) Flowchart
The step-by-step operational loop of the physical scanning device.

```mermaid
flowchart TD
    Boot([Power On & Boot OS]) --> Network[Connect to Network & API]
    Network --> Idle((IDLE STATE))
    
    Idle --> Detect{Detect Event}
    
    Detect -- Admin Button Pressed --> GenQR[Generate Secure Session Token]
    GenQR --> ShowQR[Display QR Code on Screen]
    ShowQR --> WaitApprove{Wait for App Approval}
    WaitApprove -- Timeout --> TimeoutErr[Show Timeout Error] --> Idle
    WaitApprove -- Approved --> ScanEnroll[Prompt: 'Place Hand to Enroll']
    ScanEnroll --> CapEnroll[Capture & Encrypt Palm Vector]
    CapEnroll --> SendEnroll[Send Vector to API]
    SendEnroll --> Success1[Show 'Enrollment Success'] --> Idle
    
    Detect -- Proximity Sensor / Hand Detected --> ScanAuth[Prompt: 'Scanning...']
    ScanAuth --> CapAuth[Capture Palm Vector]
    CapAuth --> SendAuth[Send Vector to API for Verification]
    SendAuth --> ApiRes{API Response}
    
    ApiRes -- Match Found --> ShowOK[Show Name & 'Check-in Successful']
    ShowOK --> GreenLED[Flash Green LED / Beep] --> Idle
    
    ApiRes -- No Match --> ShowErr[Show 'User Not Found / Try Again']
    ShowErr --> RedLED[Flash Red LED / Error Beep] --> Idle
```
