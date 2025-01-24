#!/bin/bash

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check Docker on Windows (Git Bash)
check_docker_windows() {
    if command_exists docker.exe; then
        echo "Docker is installed on Windows."
        docker.exe --version
        return 0
    else
        echo "Docker is not installed on Windows."
        return 1
    fi
}

# Function to check Docker on Unix-like systems (Linux/macOS)
check_docker_unix() {
    if command_exists docker; then
        echo "Docker is installed."
        docker --version
        return 0
    else
        echo "Docker is not installed."
        return 1
    fi
}

# Main script
case "$(uname -s)" in
    Linux*)
        check_docker_unix
        ;;
    Darwin*)
        check_docker_unix
        ;;
    MINGW*|CYGWIN*|MSYS*)
        check_docker_windows
        ;;
    *)
        echo "Unsupported operating system"
        exit 1
        ;;
esac