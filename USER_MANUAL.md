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

**Note on Revocation:** If an employee needs to re-register their palm (e.g., they injured their hand), they cannot delete their own palm template from the mobile app for security reasons. Only an Administrator can remove or revoke a palm template using the Web Admin Dashboard.

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
flowchart LR
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
    U_ID(["user_id (uuid)"]):::primaryKey
    U_EmpCode(["employee_code (string)"]):::attribute
    U_Name(["full_name (string)"]):::attribute
    U_Phone(["phone (string)"]):::attribute
    U_Role(["role (string)"]):::attribute
    U_Pass(["password_hash (string)"]):::attribute

    User --- U_ID
    User --- U_EmpCode
    User --- U_Name
    User --- U_Phone
    User --- U_Role
    User --- U_Pass

    %% Attributes for PalmTemplate
    PT_ID(["template_id (uuid)"]):::primaryKey
    PT_Hand(["hand_side (string)"]):::attribute
    PT_Encrypted(["template_encrypted (blob)"]):::attribute
    PT_Status(["status (string)"]):::attribute

    PalmTemplate --- PT_ID
    PalmTemplate --- PT_Hand
    PalmTemplate --- PT_Encrypted
    PalmTemplate --- PT_Status

    %% Attributes for AttendanceLog
    AL_ID(["log_id (uuid)"]):::primaryKey
    AL_Date(["attendance_date (date)"]):::attribute
    AL_InTime(["check_in_time (timestamp)"]):::attribute
    AL_Status(["status (string)"]):::attribute
    AL_Score(["confidence_score (float)"]):::attribute

    AttendanceLog --- AL_ID
    AttendanceLog --- AL_Date
    AttendanceLog --- AL_InTime
    AttendanceLog --- AL_Status
    AttendanceLog --- AL_Score

    %% Attributes for Device
    D_ID(["device_id (uuid)"]):::primaryKey
    D_Code(["device_code (string)"]):::attribute
    D_Name(["name (string)"]):::attribute

    Device --- D_ID
    Device --- D_Code
    Device --- D_Name

    %% Attributes for DevicePairingSession
    DPS_ID(["session_id (uuid)"]):::primaryKey
    DPS_Token(["session_token (string)"]):::attribute
    DPS_Purpose(["purpose (string)"]):::attribute
    DPS_Status(["status (string)"]):::attribute

    DevicePairingSession --- DPS_ID
    DevicePairingSession --- DPS_Token
    DevicePairingSession --- DPS_Purpose
    DevicePairingSession --- DPS_Status
```
<p align="center"><b>ຮູບທີ 3.20 ສະແດງຮູບ E-R Model</b></p>

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
    System[ລະບົບຢືນຢັນຕົວຕົນດ້ວຍລາຍຝາມື]:::root
    
    %% Level 1 Modules
    System --> UM[1. ການຈັດການຂໍ້ມູນຜູ້ໃຊ້]:::module
    System --> DM[2. ການຈັດການອຸປະກອນ]:::module
    System --> PE[3. ການລົງທະບຽນລາຍຝາມື]:::module
    System --> AP[4. ປະມວນຜົນການລົງເວລາ]:::module
    System --> RA[5. ການລາຍງານ ແລະ ວິເຄາະ]:::module
    
    %% Level 2 Processes
    UM --> UM1[1.1 ການຢືນຢັນຕົວຕົນຜູ້ໃຊ້]:::process
    UM --> UM2[1.2 ຈັດການຂໍ້ມູນສ່ວນຕົວ]:::process
    UM --> UM3[1.3 ຈັດການສະຖານະການເຂົ້າໃຊ້]:::process
    UM --> UM4[1.4 ຍົກເລີກຂໍ້ມູນລາຍຝາມື]:::process
    
    DM --> DM1[2.1 ລົງທະບຽນເຄື່ອງສະແກນ]:::process
    DM --> DM2[2.2 ອັບເດດການຕັ້ງຄ່າ]:::process
    DM --> DM3[2.3 ຕິດຕາມສະຖານະອຸປະກອນ]:::process
    
    PE --> PE1[3.1 ສ້າງ QR Code ສໍາລັບເຊດຊັນ]:::process
    PE --> PE2[3.2 ກວດສອບ QR ໃນແອັບ]:::process
    PE --> PE3[3.3 ການອະນຸມັດຈາກຜູ້ດູແລ]:::process
    PE --> PE4[3.4 ບັນທຶກ ແລະ ເຂົ້າລະຫັດລາຍຝາມື]:::process
    
    AP --> AP1[4.1 ສະກັດເວັກເຕີລາຍຝາມື]:::process
    AP --> AP2[4.2 ປຽບທຽບຂໍ້ມູນທາງຊີວະມິຕິ]:::process
    AP --> AP3[4.3 ບັນທຶກສະຖານະການລົງເວລາ]:::process
    
    RA --> RA1[5.1 ກະດານຄວບຄຸມຜູ້ດູແລ]:::process
    RA --> RA2[5.2 ສ້າງລາຍງານ]:::process
    RA --> RA3[5.3 ເບິ່ງປະຫວັດຂອງອົງກອນ]:::process
    RA --> RA4[5.4 ເບິ່ງປະຫວັດສ່ວນຕົວ]:::process

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
flowchart TD
    %% External Entities
    Admin["ຜູ້ດູແລ"]:::entity
    Employee["ພະນັກງານ"]:::entity
    Device["ອຸປະກອນສະແກນ"]:::entity

    %% Central System
    System["0.0<hr/>ລະບົບຢືນຢັນຕົວຕົນດ້ວຍລາຍຝາມື"]:::process

    %% Flows - Admin
    Admin -- "- ຂໍ້ມູນຜູ້ໃຊ້ ແລະ ອຸປະກອນ<br/>- ອະນຸມັດເຊດຊັນ" --> System
    System -- "- ລາຍງານກະດານຄວບຄຸມ<br/>- ປະຫວັດການລົງເວລາ" --> Admin

    %% Flows - Employee
    Employee -- "- ສະແກນລາຍຝາມື<br/>- ເບິ່ງປະຫວັດ" --> System
    System -- "- ຜົນການລົງເວລາ<br/>- ປະຫວັດສ່ວນຕົວ" --> Employee

    %% Flows - Device
    Device -- "- ຂໍ້ມູນເວັກເຕີລາຍຝາມື<br/>- ສະຖານະອຸປະກອນ (Heartbeat)" --> System
    System -- "- ເຊດຊັນ QR<br/>- ຜົນການຢືນຢັນ" --> Device

    %% Styling
    classDef entity fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
    classDef process fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
    classDef datastore fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
```

### 5.2 DFD Level 1 (Main Processes)
Level 1 breaks down the main system into its primary sub-processes and shows how they interact with the database stores.

```mermaid
flowchart TD
    %% External Entities
    Admin["ຜູ້ດູແລ"]:::entity
    Employee["ພະນັກງານ"]:::entity
    Device["ອຸປະກອນສະແກນ"]:::entity

    %% Processes
    P1["1.0<hr/>ການຈັດການຂໍ້ມູນຜູ້ໃຊ້"]:::process
    P2["2.0<hr/>ການຈັດການອຸປະກອນ"]:::process
    P3["3.0<hr/>ການລົງທະບຽນລາຍຝາມື"]:::process
    P4["4.0<hr/>ປະມວນຜົນການລົງເວລາ"]:::process
    P5["5.0<hr/>ການລາຍງານ ແລະ ວິເຄາະ"]:::process

    %% Data Stores
    D1["D1 &nbsp;|&nbsp; ແຟ້ມຂໍ້ມູນຜູ້ໃຊ້"]:::datastore
    D2["D2 &nbsp;|&nbsp; ແຟ້ມຂໍ້ມູນອຸປະກອນ"]:::datastore
    D3["D3 &nbsp;|&nbsp; ແຟ້ມຂໍ້ມູນຕົ້ນແບບລາຍຝາມື"]:::datastore
    D4["D4 &nbsp;|&nbsp; ແຟ້ມຂໍ້ມູນປະຫວັດການເຂົ້າ-ອອກ"]:::datastore
    D5["D5 &nbsp;|&nbsp; ແຟ້ມຂໍ້ມູນເຊດຊັນ"]:::datastore

    %% Flows
    Admin --> P1
    P1 --> D1

    Admin --> P2
    Device --> P2
    P2 --> D2

    Device --> P3
    Employee --> P3
    Admin --> P3
    P3 --> D1
    P3 --> D3
    P3 --> D5

    Device --> P4
    P4 --> D3
    P4 --> D4
    P4 --> Device

    Admin --> P5
    Employee --> P5
    D4 --> P5
    P5 --> Admin
    P5 --> Employee

    %% Styling
    classDef entity fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
    classDef process fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
    classDef datastore fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
```

### 5.3 DFD Level 2 (Process Decomposition)
Below are the decomposed Data Flow Diagrams for each of the main processes (1.0 to 4.0), showing the specific sub-processes involved.

#### Process 1: User Management
```mermaid
flowchart LR
    Admin["ຜູ້ດູແລ"]:::entity
    Employee["ພະນັກງານ"]:::entity
    
    P1_1["1.1<hr/>ການຢືນຢັນຕົວຕົນຜູ້ໃຊ້"]:::process
    P1_2["1.2<hr/>ຈັດການຂໍ້ມູນສ່ວນຕົວ"]:::process
    P1_3["1.3<hr/>ຈັດການສະຖານະການເຂົ້າໃຊ້"]:::process
    P1_4["1.4<hr/>ຍົກເລີກຂໍ້ມູນລາຍຝາມື"]:::process
    
    D1["D1 &nbsp;|&nbsp; ແຟ້ມຂໍ້ມູນຜູ້ໃຊ້"]:::datastore
    D3["D3 &nbsp;|&nbsp; ແຟ້ມຂໍ້ມູນຕົ້ນແບບລາຍຝາມື"]:::datastore
    
    Admin --> P1_1
    Employee --> P1_1
    P1_1 --> D1
    
    Admin --> P1_2
    P1_2 --> D1
    
    Admin --> P1_3
    P1_3 --> D1
    
    Admin --> P1_4
    P1_4 --> D3

    classDef entity fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
    classDef process fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
    classDef datastore fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
```

#### Process 2: Device Management
```mermaid
flowchart LR
    Admin["ຜູ້ດູແລ"]:::entity
    Device["ອຸປະກອນສະແກນ"]:::entity
    
    P2_1["2.1<hr/>ລົງທະບຽນເຄື່ອງສະແກນ"]:::process
    P2_2["2.2<hr/>ອັບເດດການຕັ້ງຄ່າ"]:::process
    P2_3["2.3<hr/>ຕິດຕາມສະຖານະອຸປະກອນ"]:::process
    
    D2["D2 &nbsp;|&nbsp; ແຟ້ມຂໍ້ມູນອຸປະກອນ"]:::datastore
    
    Admin --> P2_1
    P2_1 --> D2
    
    Admin --> P2_2
    P2_2 --> D2
    
    Device --> P2_3
    P2_3 --> D2

    classDef entity fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
    classDef process fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
    classDef datastore fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
```

#### Process 3: Palm Enrollment
```mermaid
flowchart TD
    Device["ອຸປະກອນສະແກນ"]:::entity
    MobileApp["ແອັບມືຖື"]:::entity
    Admin["ຜູ້ດູແລ"]:::entity
    
    P3_1["3.1<hr/>ສ້າງ QR Code ສໍາລັບເຊດຊັນ"]:::process
    P3_2["3.2<hr/>ກວດສອບ QR ໃນແອັບ"]:::process
    P3_3["3.3<hr/>ການອະນຸມັດຈາກຜູ້ດູແລ"]:::process
    P3_4["3.4<hr/>ບັນທຶກ ແລະ ເຂົ້າລະຫັດລາຍຝາມື"]:::process
    
    D1["D1 &nbsp;|&nbsp; ແຟ້ມຂໍ້ມູນຜູ້ໃຊ້"]:::datastore
    D3["D3 &nbsp;|&nbsp; ແຟ້ມຂໍ້ມູນຕົ້ນແບບລາຍຝາມື"]:::datastore
    D5["D5 &nbsp;|&nbsp; ແຟ້ມຂໍ້ມູນເຊດຊັນ"]:::datastore
    
    Device --> P3_1
    P3_1 --> D5
    P3_1 --> MobileApp
    MobileApp --> P3_2
    P3_2 --> D5
    P3_2 --> Admin
    
    Admin --> P3_3
    P3_3 --> D5
    P3_3 --> D1
    P3_3 --> Device
    
    Device --> P3_4
    P3_4 --> D3

    classDef entity fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
    classDef process fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
    classDef datastore fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
```

#### Process 4: Attendance Processing
```mermaid
flowchart TD
    Device["ອຸປະກອນສະແກນ"]:::entity
    
    P4_1["4.1<hr/>ສະກັດເວັກເຕີລາຍຝາມື"]:::process
    P4_2["4.2<hr/>ປຽບທຽບຂໍ້ມູນທາງຊີວະມິຕິ"]:::process
    P4_3["4.3<hr/>ບັນທຶກສະຖານະການລົງເວລາ"]:::process
    
    D3["D3 &nbsp;|&nbsp; ແຟ້ມຂໍ້ມູນຕົ້ນແບບລາຍຝາມື"]:::datastore
    D4["D4 &nbsp;|&nbsp; ແຟ້ມຂໍ້ມູນປະຫວັດການເຂົ້າ-ອອກ"]:::datastore
    
    Device --> P4_1
    P4_1 --> P4_2
    
    P4_2 --> D3
    P4_2 --> P4_3
    
    P4_3 --> D4
    P4_3 --> Device

    classDef entity fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
    classDef process fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
    classDef datastore fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
```

#### Process 5: Reporting & Analytics
```mermaid
flowchart LR
    Admin["ຜູ້ດູແລ"]:::entity
    Employee["ພະນັກງານ"]:::entity
    
    P5_1["5.1<hr/>ກະດານຄວບຄຸມຜູ້ດູແລ"]:::process
    P5_2["5.2<hr/>ສ້າງລາຍງານ"]:::process
    P5_3["5.3<hr/>ເບິ່ງປະຫວັດຂອງອົງກອນ"]:::process
    P5_4["5.4<hr/>ເບິ່ງປະຫວັດສ່ວນຕົວ"]:::process
    
    D1["D1 &nbsp;|&nbsp; ແຟ້ມຂໍ້ມູນຜູ້ໃຊ້"]:::datastore
    D4["D4 &nbsp;|&nbsp; ແຟ້ມຂໍ້ມູນປະຫວັດການເຂົ້າ-ອອກ"]:::datastore
    
    Admin --> P5_1
    D4 --> P5_1
    
    Admin --> P5_2
    D4 --> P5_2
    D1 --> P5_2
    
    Admin --> P5_3
    D4 --> P5_3
    
    Employee --> P5_4
    D4 --> P5_4

    classDef entity fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
    classDef process fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
    classDef datastore fill:#fff,stroke:#000,stroke-width:1px,shape:rect,color:#000
```

---

### 6. Data Dictionary (ພົດຈະນານຸກົມຂໍ້ມູນ)

ຕາຕະລາງລຸ່ມນີ້ອະທິບາຍໂຄງສ້າງຂອງຖານຂໍ້ມູນທັງ 5 ຕາຕະລາງທີ່ໃຊ້ໃນລະບົບ, ເຊິ່ງປະກອບມີລາຍລະອຽດຂອງຟິວ, ປະເພດຂໍ້ມູນ, ຂໍ້ກໍານົດຕ່າງໆ, ແລະ ຄວາມສຳພັນ.

### 6.1 `users`
ເກັບຂໍ້ມູນໂປຣໄຟລ໌ຂອງພະນັກງານ ແລະ ຜູ້ດູແລລະບົບທັງໝົດ.

| ຊື່ຟິວ (Column) | ປະເພດຂໍ້ມູນ (Type) | ຂໍ້ກຳນົດ (Constraints) | ຄຳອະທິບາຍ (Description) |
|---|---|---|---|
| `id` | UUID | Primary Key | ລະຫັດສະເພາະສຳລັບຜູ້ໃຊ້ |
| `employee_code` | VARCHAR(50) | Unique, Not Null | ລະຫັດພະນັກງານ |
| `full_name` | VARCHAR(100) | Not Null | ຊື່ ແລະ ນາມສະກຸນຂອງຜູ້ໃຊ້ |
| `email` | VARCHAR(100) | Unique | ທີ່ຢູ່ອີເມວຂອງຜູ້ໃຊ້ |
| `phone` | VARCHAR(20) | Unique | ເບີໂທລະສັບຂອງຜູ້ໃຊ້ |
| `password_hash` | VARCHAR(255) | Not Null | ລະຫັດຜ່ານທີ່ຖືກເຂົ້າລະຫັດແບບ Bcrypt |
| `role` | VARCHAR(20) | Default 'EMPLOYEE' | ກຳນົດສິດເປັນ 'ADMIN' ຫຼື 'EMPLOYEE' |
| `department` | VARCHAR(50) | | ພະແນກທີ່ຜູ້ໃຊ້ສັງກັດຢູ່ |
| `status` | VARCHAR(20) | Default 'active' | ສະຖານະບັນຊີຜູ້ໃຊ້ (active/inactive) |
| `created_at` | TIMESTAMP | Default NOW() | ເວລາທີ່ສ້າງໂປຣໄຟລ໌ |
| `updated_at` | TIMESTAMP | Default NOW() | ເວລາທີ່ມີການອັບເດດໂປຣໄຟລ໌ລ່າສຸດ |

### 6.2 `devices`
ເກັບຂໍ້ມູນອຸປະກອນສະແກນຮາດແວທີ່ລົງທະບຽນໃນລະບົບ.

| ຊື່ຟິວ (Column) | ປະເພດຂໍ້ມູນ (Type) | ຂໍ້ກຳນົດ (Constraints) | ຄຳອະທິບາຍ (Description) |
|---|---|---|---|
| `id` | UUID | Primary Key | ລະຫັດສະເພາະສຳລັບອຸປະກອນ |
| `device_code` | VARCHAR(50) | Unique, Not Null | ໝາຍເລກເຄື່ອງ ຫຼື ລະຫັດຮາດແວ |
| `name` | VARCHAR(100) | Not Null | ຊື່ອຸປະກອນ (ເຊັ່ນ: "ເຄື່ອງສະແກນໜ້າປະຕູ") |
| `location` | VARCHAR(100) | | ສະຖານທີ່ຕັ້ງຂອງອຸປະກອນ |
| `status` | VARCHAR(20) | Default 'active' | ສະຖານະການເຮັດວຽກຂອງອຸປະກອນ |
| `last_seen_at` | TIMESTAMP | | ເວລາລ່າສຸດທີ່ອຸປະກອນເຊື່ອມຕໍ່ກັບເຊີບເວີ |
| `created_at` | TIMESTAMP | Default NOW() | ເວລາທີ່ລົງທະບຽນ |

### 6.3 `device_pairing_sessions`
ເຊດຊັນຊົ່ວຄາວທີ່ສ້າງຂຶ້ນເມື່ອເຄື່ອງສະແກນສ້າງ QR code ສຳລັບລົງທະບຽນລາຍຝາມື.

| ຊື່ຟິວ (Column) | ປະເພດຂໍ້ມູນ (Type) | ຂໍ້ກຳນົດ (Constraints) | ຄຳອະທິບາຍ (Description) |
|---|---|---|---|
| `id` | UUID | Primary Key | ລະຫັດເຊດຊັນ |
| `device_id` | UUID | Foreign Key | ເຄື່ອງສະແກນທີ່ສ້າງເຊດຊັນນີ້ |
| `session_token` | TEXT | Unique, Not Null | ໂທເຄັນທີ່ຝັງຢູ່ໃນ QR code |
| `user_id` | UUID | Foreign Key | ຜູ້ໃຊ້ທີ່ອະນຸມັດເຊດຊັນນີ້ |
| `purpose` | VARCHAR(50) | Not Null | ຈຸດປະສົງ (ປົກກະຕິແມ່ນ 'enrollment') |
| `status` | VARCHAR(30) | Default 'pending' | ສະຖານະ (pending, approved, completed) |
| `expires_at` | TIMESTAMP | Not Null | ເວລາທີ່ QR code ໝົດອາຍຸ |
| `scanned_at` | TIMESTAMP | | ເວລາທີ່ຜູ້ໃຊ້ສະແກນ QR |
| `approved_at` | TIMESTAMP | | ເວລາທີ່ຜູ້ໃຊ້ກົດ "ອະນຸມັດ" |
| `completed_at` | TIMESTAMP | | ເວລາທີ່ຂໍ້ມູນເວັກເຕີລາຍຝາມືຖືກບັນທຶກ |
| `created_at` | TIMESTAMP | Default NOW() | ເວລາທີ່ສ້າງເຊດຊັນ |

### 6.4 `palm_templates`
ເກັບຮັກສາຂໍ້ມູນເວັກເຕີຊີວະມິຕິທີ່ຖືກເຂົ້າລະຫັດຢ່າງປອດໄພ.

| ຊື່ຟິວ (Column) | ປະເພດຂໍ້ມູນ (Type) | ຂໍ້ກຳນົດ (Constraints) | ຄຳອະທິບາຍ (Description) |
|---|---|---|---|
| `id` | UUID | Primary Key | ລະຫັດຕົ້ນແບບ |
| `user_id` | UUID | Foreign Key, Not Null | ເຈົ້າຂອງຕົ້ນແບບນີ້ |
| `hand_side` | VARCHAR(10) | Not Null | ມືຊ້າຍ ຫຼື ມືຂວາ ('left' ຫຼື 'right') |
| `template_encrypted` | BYTEA | Not Null | ຂໍ້ມູນເວັກເຕີທີ່ເຂົ້າລະຫັດດ້ວຍ AES-256 |
| `template_nonce` | BYTEA | Not Null | Nonce ທີ່ໃຊ້ສຳລັບຖອດລະຫັດ |
| `embedding_dim` | INT | Default 128 | ຂະໜາດຂອງເວັກເຕີ (ມິຕິ) |
| `model_version` | VARCHAR(100) | Not Null | ໂມເດວ AI ທີ່ໃຊ້ສະກັດຂໍ້ມູນ |
| `threshold` | NUMERIC(5,4) | Default 0.8200 | ຄ່າຄວາມແມ່ນຍຳ (Threshold) ສຳລັບຕົ້ນແບບນີ້ |
| `status` | VARCHAR(30) | Default 'active' | ສະຖານະຂອງຕົ້ນແບບ |
| `registered_device_id`| UUID | Foreign Key | ອຸປະກອນທີ່ໃຊ້ບັນທຶກຕົ້ນແບບ |
| `created_at` | TIMESTAMP | Default NOW() | ເວລາທີ່ລົງທະບຽນ |
| `updated_at` | TIMESTAMP | Default NOW() | ເວລາທີ່ມີການອັບເດດລ່າສຸດ |
| `revoked_at` | TIMESTAMP | | ເວລາທີ່ຕົ້ນແບບຖືກຍົກເລີກ (ຖ້າມີ) |

### 6.5 `attendance_logs`
ບັນທຶກປະຫວັດການເຂົ້າ-ອອກວຽກປະຈຳວັນຂອງຜູ້ໃຊ້ທັງໝົດ.

| ຊື່ຟິວ (Column) | ປະເພດຂໍ້ມູນ (Type) | ຂໍ້ກຳນົດ (Constraints) | ຄຳອະທິບາຍ (Description) |
|---|---|---|---|
| `id` | UUID | Primary Key | ລະຫັດປະຫວັດ |
| `user_id` | UUID | Foreign Key, Not Null | ພະນັກງານທີ່ລົງເວລາ |
| `device_id` | UUID | Foreign Key | ເຄື່ອງສະແກນທີ່ໃຊ້ລົງເວລາ |
| `attendance_date` | DATE | Not Null | ວັນທີຂອງການລົງເວລາ |
| `check_in_time` | TIMESTAMP | | ເວລາທີ່ສະແກນຄັ້ງທຳອິດ (ເຂົ້າວຽກ) |
| `check_out_time` | TIMESTAMP | | ເວລາທີ່ສະແກນຄັ້ງລ່າສຸດ (ອອກວຽກ) |
| `check_in_score` | NUMERIC(6,5) | | ຄະແນນຄວາມໝັ້ນໃຈທາງຊີວະມິຕິ (0-1) |
| `check_out_score`| NUMERIC(6,5) | | ຄະແນນຄວາມໝັ້ນໃຈທາງຊີວະມິຕິ (0-1) |
| `check_in_liveness`| BOOLEAN | Default false | ຜ່ານການກວດສອບການປອມແປງຈາກຮາດແວ |
| `check_out_liveness`| BOOLEAN | Default false | ຜ່ານການກວດສອບການປອມແປງຈາກຮາດແວ |
| `status` | VARCHAR(20) | Default 'present' | ສະຖານະ (present, late, incomplete, absent) |
| `created_at` | TIMESTAMP | Default NOW() | ເວລາທີ່ສ້າງປະຫວັດນີ້ |

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
