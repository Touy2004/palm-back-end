-- 001_init.sql
-- Initial schema for Palm Recognition Attendance System

-- Enable UUID extension if not already enabled (postgres specific)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 9.1 Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_code VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(20) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'EMPLOYEE',
    department VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 9.2 Devices Table
CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(100),
    status VARCHAR(20) DEFAULT 'active',
    last_seen_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 9.3 Device Pairing Sessions Table
CREATE TABLE device_pairing_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    session_token TEXT UNIQUE NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    purpose VARCHAR(50) NOT NULL,
    status VARCHAR(30) DEFAULT 'pending',
    expires_at TIMESTAMP NOT NULL,
    scanned_at TIMESTAMP,
    approved_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 9.4 Palm Templates Table
CREATE TABLE palm_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    hand_side VARCHAR(10) NOT NULL,
    template_encrypted BYTEA NOT NULL,
    template_nonce BYTEA NOT NULL,
    embedding_dim INT NOT NULL DEFAULT 128,
    model_version VARCHAR(100) NOT NULL,
    threshold NUMERIC(5, 4) DEFAULT 0.8200,
    status VARCHAR(30) DEFAULT 'active',
    registered_device_id UUID REFERENCES devices(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    revoked_at TIMESTAMP
);

-- 9.5 Attendance Logs Table
CREATE TABLE attendance_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id UUID REFERENCES devices(id) ON DELETE SET NULL,
    attendance_date DATE NOT NULL,
    check_in_time TIMESTAMP,
    check_out_time TIMESTAMP,
    check_in_score NUMERIC(6, 5),
    check_out_score NUMERIC(6, 5),
    check_in_liveness BOOLEAN DEFAULT false,
    check_out_liveness BOOLEAN DEFAULT false,
    status VARCHAR(20) DEFAULT 'present',
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, attendance_date)
);

-- 9.6 Palm Authentication Attempts (Audit Log)
CREATE TABLE palm_auth_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    device_id UUID REFERENCES devices(id) ON DELETE SET NULL,
    template_id UUID REFERENCES palm_templates(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL,
    score NUMERIC(6, 5),
    threshold NUMERIC(5, 4),
    liveness_passed BOOLEAN DEFAULT false,
    quality_score NUMERIC(6, 5),
    thermal_min NUMERIC(6, 2),
    thermal_max NUMERIC(6, 2),
    thermal_avg NUMERIC(6, 2),
    result VARCHAR(30) NOT NULL,
    failure_reason TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
