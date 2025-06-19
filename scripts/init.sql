-- Create database if not exists
-- This script will be executed when PostgreSQL container starts

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- You can add any initial data or additional setup here
-- For example:
-- INSERT INTO users (id, email, name) VALUES (uuid_generate_v4(), 'admin@example.com', 'Admin User'); 