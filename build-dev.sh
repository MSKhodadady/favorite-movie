#!/usr/bin/bash

# folders
frontend_folder="nex-favorite-movie"
backend_folder="go-favorite-movie"
build_output="build-output"

# building the frontend files
echo    "----------<< build frontend >>----------"
## going to frontend folder
cd      $frontend_folder
## building
NEXT_PUBLIC_SERVER_ADDRESS="/api" npm run build

# buiding backend files
echo    "----------<< build backend >>----------"
## goid to backend folder
cd      ../$backend_folder
CGO_ENABLED="0" GOOS=linux GOARCH=amd64 go build -o go-bin

# packing files
echo    "----------<< pack files >>----------"
## going to home folder
cd      ..
## delete previous build output
rm -r   $build_output
## creating new build output
mkdir   $build_output
## going to that folder
cd      $build_output
## cp env file
cp      ../env.production.json      ./env.json
## mv frontend build output files
mv      ../$frontend_folder/out     ./frontend
## mv backend build output files
mv      ../$backend_folder/go-bin   ./
