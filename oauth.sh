#!/bin/bash
go build main.go
cp main bin/
rm main 
git add bin/main
git commit -m "main"
