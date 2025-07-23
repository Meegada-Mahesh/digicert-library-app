#!/usr/bin/env bash

set -e

echo "Setting up Digicert Library App for local development..."

# Check for Docker
if ! command -v docker &> /dev/null; then
    echo "Docker not found. Please install Docker before running this script."
    exit 1
fi

# Check if user is in docker group
if ! groups $USER | grep -q docker; then
    echo "⚠️  Adding user to docker group..."
    sudo usermod -aG docker $USER
    echo "✅ User added to docker group. Please log out and log back in, then run this script again."
    exit 0
fi

# sudo aa-remove-unknown

# Check for Docker Compose (support both docker-compose and docker compose)
if command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker-compose"
elif docker compose version &> /dev/null; then
    COMPOSE_CMD="docker compose"
else
    echo "Docker Compose not found. Please install Docker Compose before running this script."
    exit 1
fi


# Check for Go
if ! command -v go &> /dev/null; then
    echo "Go not found. Please install Go before running this script."
    exit 1
fi

# Copy .env.example to .env if .env does not exist
if [ ! -f .env ]; then
    if [ -f .env.example ]; then
        cp .env.example .env
        echo ".env file created from .env.example. Please update credentials as needed."
    else
        echo "No .env or .env.example file found. Please create your .env file with DB credentials."
        exit 1
    fi
fi

echo "Tidying Go modules..."
go mod tidy

echo "Building and starting containers..."
$COMPOSE_CMD up --build -d

echo "Waiting for services to start..."
sleep 10

# Test if the app is responding
echo "Testing app connectivity..."
sleep 5
if curl -f http://localhost:8080 &>/dev/null; then
    echo "✅ App is running successfully at http://localhost:8080"
else
    echo "⚠️  App may still be starting. Check status with: $COMPOSE_CMD ps"
    echo "📋 View logs with: $COMPOSE_CMD logs -f"
fi

echo "🎉 Setup complete!"
echo "📱 App URL: http://localhost:8080"
echo "🗄️  MySQL: localhost:3307"