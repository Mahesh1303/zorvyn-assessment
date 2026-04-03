-- Drop tables first (respect FK dependency)
DROP TABLE IF EXISTS financial_records;
DROP TABLE IF EXISTS users;

-- Drop ENUM types
DROP TYPE IF EXISTS record_type;
DROP TYPE IF EXISTS user_role;
DROP EXTENSION IF EXISTS "pgcrypto";