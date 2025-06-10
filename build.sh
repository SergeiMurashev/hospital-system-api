#!/bin/bash

# Copy proto directory to each service
cp -r proto hospital-service/
cp -r proto timetable-service/
cp -r proto document-service/

# Build services
docker-compose build 