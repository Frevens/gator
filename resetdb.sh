#!/usr/bin/env bash
set -e

DB_URL="postgres://postgres:postgres@localhost:5432/gator"

goose -dir sql/schema postgres "$DB_URL" down

sleep 3

goose -dir sql/schema postgres "$DB_URL" up
