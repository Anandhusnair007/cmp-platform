#!/bin/bash

# Setup local PostgreSQL database for CMP Platform
# This script sets up the database without Docker

set -e

echo "🔧 Setting up PostgreSQL database for CMP Platform..."

# Database configuration from local.env
DB_USER="cmp_user"
DB_PASSWORD="cmp_pass"
DB_NAME="cmp_db"

# Check if PostgreSQL is running
if ! pg_isready -q; then
    echo "❌ PostgreSQL is not running. Please start it first:"
    echo "   sudo systemctl start postgresql"
    exit 1
fi

echo "✅ PostgreSQL is running"

# Create database user and database
echo "📝 Creating database user and database..."
sudo -u postgres psql <<EOF
-- Drop existing if needed (for fresh setup)
-- DROP DATABASE IF EXISTS ${DB_NAME};
-- DROP USER IF EXISTS ${DB_USER};

-- Create user if doesn't exist
DO \$\$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = '${DB_USER}') THEN
        CREATE USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}' SUPERUSER;
        RAISE NOTICE 'User ${DB_USER} created';
    ELSE
        ALTER USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}' SUPERUSER;
        RAISE NOTICE 'User ${DB_USER} already exists, password updated';
    END IF;
END
\$\$;

-- Create database if doesn't exist
SELECT 'CREATE DATABASE ${DB_NAME} OWNER ${DB_USER}'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '${DB_NAME}')\gexec

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};

EOF

echo "✅ Database setup complete!"
echo ""
echo "📊 Database Info:"
echo "   Host: localhost"
echo "   Port: 5432"
echo "   User: ${DB_USER}"
echo "   Password: ${DB_PASSWORD}"
echo "   Database: ${DB_NAME}"
echo ""
echo "🎯 Next steps:"
echo "   1. Run migrations: make migrate-up"
echo "   2. Start backend services"
echo "   3. Start frontend"
