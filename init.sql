-- Script of initialize container MySQL

-- Create database dafoncvdb if not exists
CREATE DATABASE IF NOT EXISTS dafoncvdb;

-- Create user dafoncv if not exists
CREATE USER IF NOT EXISTS 'dafoncv'@'%' IDENTIFIED BY 'dafoncv';

-- Grant privileges to user dafoncv on database dafoncvdb
GRANT ALL PRIVILEGES ON dafoncvdb.* TO 'dafoncv'@'%';

-- Apply changes
FLUSH PRIVILEGES; 