#!/bin/bash
go build src/main.go
sudo chmod +x main
pm2 start ./main --name worker