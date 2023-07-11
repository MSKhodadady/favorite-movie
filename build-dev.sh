#!/usr/bin/bash

next_project="nex-favorite-movie"
go_projcet="go-favorite-movie"

echo "--- build frontend ---"
cd $next_project
NEXT_PUBLIC_SERVER_ADDRESS="/api" npm run build

cd ..

echo "--- build backend ---"
cd $go_projcet
CGO_ENABLED="0" GOOS=linux GOARCH=amd64 go build -o go-bin

cd ..

echo "--- pack files ---"
rm -r build-output
mkdir build-output
cd build-output

mv ../$next_project/out ./frontend

mv ../$go_projcet/go-bin ./
cp ../$go_projcet/env.json ./
cp ../$go_projcet/localhost.pem ./
cp ../$go_projcet/localhost-key.pem ./