#!/bin/bash

# Restart Docker Compose services
docker-compose -f  down
docker-compose -f  up -d
