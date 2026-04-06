#!/bin/sh
set -e

DATA_DIR="/app/data"
ENV_FILE="$DATA_DIR/.env"

# Ensure data directory exists
mkdir -p "$DATA_DIR/uploads"

# Generate persistent secrets on first run
if [ ! -f "$ENV_FILE" ]; then
    echo "==> First run detected. Generating secrets..."

    # Generate random JWT secret (64 hex chars)
    JWT_SECRET=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c 64)

    # Generate random AES key (32 hex chars for AES-256)
    AES_KEY=$(cat /dev/urandom | tr -dc 'a-f0-9' | head -c 32)

    # Save to persistent volume
    cat > "$ENV_FILE" <<EOF
# Auto-generated secrets - stored in persistent volume
# These survive container rebuilds/redeployments
JWT_SECRET=$JWT_SECRET
AES_KEY=$AES_KEY
EOF

    echo "==> Secrets generated and saved to $ENV_FILE"
    echo "==> JWT_SECRET and AES_KEY will persist across redeployments"
else
    echo "==> Loading existing secrets from $ENV_FILE"
fi

# Load saved secrets (only if not already set via environment)
while IFS='=' read -r key value; do
    # Skip comments and empty lines
    case "$key" in
        \#*|"") continue ;;
    esac
    # Trim whitespace from key and value
    key=$(echo "$key" | tr -d '[:space:]')
    value=$(echo "$value" | tr -d '[:space:]')
    [ -z "$key" ] && continue

    # Only set if not already provided via environment variable
    current_val=""
    current_val=$(printenv "$key" 2>/dev/null) || true
    if [ -z "$current_val" ]; then
        export "${key}=${value}"
        echo "==> Loaded $key from saved config"
    else
        echo "==> Using provided $key from environment"
    fi
done < "$ENV_FILE"

# Set defaults for other variables if not set
export DB_PATH="${DB_PATH:-/app/data/cvbuilder.db}"
export GIN_MODE="${GIN_MODE:-release}"
export TZ="${TZ:-Asia/Baghdad}"

echo "==> Starting CV Builder..."
echo "    Port: ${PORT:-8080}"
echo "    Database: $DB_PATH"
echo "    Mode: $GIN_MODE"

# Start the server
exec ./server
