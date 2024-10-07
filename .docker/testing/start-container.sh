#!/bin/bash


go mod vendor
go mod tidy


CompileDaemon  --build="go build -o cmd/app/app cmd/app/main.go" --command=./cmd/app/app -color=true -include=app.yaml -include=database.yaml
