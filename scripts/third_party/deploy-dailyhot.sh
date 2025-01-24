#!/bin/bash

# Source utility functions
source ../utils.sh

print_notice "Starting DailyHot API deployment..."

# Pull the latest image
print_notice "Pulling the latest DailyHot API image..."
if docker pull imsyy/dailyhot-api:latest; then
    print_success "Successfully pulled the latest image"
else
    print_error "Failed to pull the image"
    exit 1
fi

# Check if container already exists and remove it
if docker ps -a | grep -q "imsyy/dailyhot-api"; then
    print_warning "Found existing DailyHot API container. Removing it..."
    docker rm -f $(docker ps -a | grep "imsyy/dailyhot-api" | awk '{print $1}')
fi

# Run the container
# Run the container
print_notice "Starting DailyHot API container..."
if ! docker_run_output=$(docker run --restart always -p 6688:6688 -d imsyy/dailyhot-api:latest 2>&1); then
    print_error "Failed to start the container"
    print_error "Docker Error: $docker_run_output"
    exit 1
fi

print_success "DailyHot API container is now running"
print_notice "The API is accessible at http://localhost:6688"