#!/usr/bin/env just --justfile

update:
  go get -u
  go mod tidy -v


dev:
    air & \
    tailwindcss -i ui/assets/input.css -o ui/static/css/main.css --watch & \
    browser-sync start \
      --files 'ui/html/**/*.tmpl, ui/static/**/*.css' \
      --port 4001 \
      --proxy 'localhost:4000' \
      --middleware 'function(req, res, next) { \
        res.setHeader("Cache-Control", "no-cache, no-store, must-revalidate"); \
        return next(); \
      }'



prod:
    tailwindcss -i ui/assets/input.css -o ui/static/css/main.css --minify