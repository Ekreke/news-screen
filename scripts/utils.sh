#!/bin/bash

print_success() {
    # shellcheck disable=SC2039
    echo -e "\033[32m✔️ $1\033[0m"
}

print_error() {
    # shellcheck disable=SC2039
    echo -e "\033[31m❌ $1\033[0m"
}

print_prompt() {
    # shellcheck disable=SC2039
    echo -e "\033[31m $1\033[0m"
}

print_warning() {
    # shellcheck disable=SC2039
    echo -e "\033[33m⚠️  $1\033[0m"
}

print_notice() {
    # shellcheck disable=SC2039
    echo -e "\033[35m📢 $1\033[0m"
}