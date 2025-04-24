#!/bin/bash

# переменные 
DOCKER_IMAGE_NAME="myminio"
VERSION="latest"
DOCKER_BUILD_CONTEXT="."

docker build -t "$DOCKER_IMAGE_NAME:$VERSION" "$DOCKER_BUILD_CONTEXT"