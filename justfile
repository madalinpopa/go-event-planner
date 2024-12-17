#!/usr/bin/env just --justfile

update:
  go get -u
  go mod tidy -v


dev:
    tailwindcss -i ui/assets/input.css -o ui/static/css/main.css --watch

prod:
    tailwindcss -i ui/assets/input.css -o ui/static/css/main.css --minify