#!/bin/bash
nvm use 16.18.1
go build src/main.go
sudo chmod +x main
pm2 start ./main --name worker