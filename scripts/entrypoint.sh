#!/bin/sh

export DB_DSN="postgres://${DATABASE_USER}:${DATABASE_PASSWORD}@postgres:5432/${DATABASE_NAME}?sslmode=disable"

goose -dir ./migrations postgres "$DB_DSN" up

./service