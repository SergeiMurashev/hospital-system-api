#!/bin/bash

# Copy proto directory
cp -r ../proto .

# Build and run with docker-compose
docker-compose up --build 