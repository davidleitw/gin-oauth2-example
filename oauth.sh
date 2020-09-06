#!/bin/bash
go build main.go
cp main bin/
rm main 
git add . 
git commit -m "main"
git push Heroku