-- FinanceTracker Database Initialization Script
-- This script only sets up essential PostgreSQL extensions and functions.
-- Tables and schema will be managed by Atlas migrations.

-- Enable UUID extension (required for UUID primary keys)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable pgcrypto extension (for additional UUID functions)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create trigger function for updated_at timestamp
-- This function will be used by Atlas migrations to create triggers
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Note: Tables, indexes, and constraints will be created by Atlas migrations