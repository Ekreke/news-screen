#!/bin/sh

. ./utils.sh 

# TODO: progress bar
cat << "EOF"
       _             _        
      | |           | |       
   ___| | ___ __ ___| | _____ 
  / _ \ |/ / '__/ _ \ |/ / _ \
 |  __/   <| | |  __/   <  __/
  \___|_|\_\_|  \___|_|\_\___|
                              
                              
EOF

print_notice "run news-screen with shell..."

# Run Docker check
print_notice "Checking Docker environment..."
if ! ./check_docker.sh; then
    print_error "Docker installed check failed. Please ensure Docker is properly installed."
    exit 1
fi

# Deploy DailyHot container
print_notice "Deploying DailyHot API..."
if ! ./third_party/deploy-dailyhot.sh; then
    print_error "Failed to deploy DailyHot API"
    exit 1
fi

# Deploy news-screen service
print_notice "Building and deploying news-screen service..."
# Build Docker image
if ! docker_build_output=$(docker build -t news-screen . 2>&1); then
    print_error "Failed to build news-screen image"
    print_error "Docker Build Error: $docker_build_output"
    exit 1
fi

# Run the container
print_notice "Starting news-screen container..."
if ! docker_output=$(docker run --rm -d \
    -p 8000:8000 \
    -p 9000:9000 \
    -v "$(pwd)/configs:/data/conf" \
    news-screen 2>&1); then
    print_error "Failed to start news-screen container"
    print_error "Docker Error: $docker_output"
    exit 1
fi

print_success "news-screen service is now running"