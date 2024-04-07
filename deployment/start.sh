#!/bin/bash

# Check if the Docker images exist locally
if ! docker images | grep -q 'blockservice' ||  ! docker images | grep -q 'peer'; then
    echo "Docker images not found locally. Loading images from image archive..."
    # Load images from image archive
    ./load-images.sh
fi

# Start Docker Compose services
docker-compose up -d
