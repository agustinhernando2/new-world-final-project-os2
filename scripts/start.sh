#!/bin/sh
echo "Starting swagger documentation..."

swag init -g cmd/server/main.go -d . --parseDependency --parseInternal
swag fmt


echo "-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*"
echo "Starting application..."
air -c .air.toml
