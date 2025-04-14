#!/bin/bash

# Переменные
DOCKER_IMAGE_NAME="mymongodb"
VERSION="latest"
DOCKER_BUILD_CONTEXT="."

docker build -t "$DOCKER_IMAGE_NAME:$VERSION" "$DOCKER_BUILD_CONTEXT"
