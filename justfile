#!/usr/bin/env just --justfile

# Variables
db_path := "./database/events.db"
migrations_dir := "database/migrations"
css_input := "ui/assets/input.css"
css_output := "ui/static/css/main.css"
dev_port := "4000"
browser_sync_port := "4001"

# Update Go dependencies
update:
    go get -u
    go mod tidy -v

# Run development server with live reload
dev:
    air & \
    tailwindcss -i {{css_input}} -o {{css_output}} --watch & \
    browser-sync start \
      --files 'ui/html/**/*.tmpl, ui/static/**/*.css' \
      --port {{browser_sync_port}} \
      --proxy 'localhost:{{dev_port}}' \
      --middleware 'function(req, res, next) { \
        res.setHeader("Cache-Control", "no-cache, no-store, must-revalidate"); \
        return next(); \
      }'

# Build production CSS
build:
    tailwindcss -i {{css_input}} -o {{css_output}} --minify

# Run database migrations
migrate command="up":
    goose sqlite3 {{db_path}} {{command}} --dir={{migrations_dir}}

# Create new migration
makemigrations name:
    goose sqlite3 {{db_path}} create {{name}} sql --dir={{migrations_dir}}