#!/bin/bash

print_success() {
    printf "\033[32m✔️ %s\033[0m\n" "$1"
}

print_error() {
    printf "\033[31m❌ %s\033[0m\n" "$1"
}

print_prompt() {
    printf "\033[31m %s\033[0m\n" "$1"
}

print_warning() {
    printf "\033[33m⚠️ %s\033[0m\n" "$1"
}

print_notice() {
    printf "\033[35m📢 %s\033[0m\n" "$1"
}