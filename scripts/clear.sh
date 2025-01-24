#!/bin/bash

# Source utility functions
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
. "${SCRIPT_DIR}/utils.sh"

print_warning "This is a shell scripts for test"
print_notice "Starting cleanup of containers..."

# Clean up DailyHot API container
print_notice "Checking for DailyHot API container..."
if docker ps -a | grep -q "imsyy/dailyhot-api"; then
    print_warning "Found DailyHot API container. Removing it..."
    if docker rm -f $(docker ps -a | grep "imsyy/dailyhot-api" | awk '{print $1}'); then
        print_success "Successfully removed DailyHot API container"
    else
        print_error "Failed to remove DailyHot API container"
    fi
else
    print_notice "No DailyHot API container found"
fi

# Clean up news-screen container
print_notice "Checking for news-screen container..."
if docker ps -a | grep -q "news-screen"; then
    print_warning "Found news-screen container. Removing it..."
    if docker rm -f $(docker ps -a | grep "news-screen" | awk '{print $1}'); then
        print_success "Successfully removed news-screen container"
    else
        print_error "Failed to remove news-screen container"
    fi
else
    print_notice "No news-screen container found"
fi

print_success "Cleanup completed"