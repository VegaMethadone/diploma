#!/bin/bash

# Переменные
NAME="testMongo"
PORT="27017"
ROOT_USER="root"
ROOT_PASSWORD="0000"
DATABASE="labyrinth"
VERSION="latest"

# Запуск контейнера MongoDB
docker run -d \
  --name "$NAME" \
  -p "$PORT:$PORT" \
  -e MONGO_INITDB_ROOT_USERNAME="$ROOT_USER" \
  -e MONGO_INITDB_ROOT_PASSWORD="$ROOT_PASSWORD" \
  -e MONGO_INITDB_DATABASE="$DATABASE" \
  mymongodb:"$VERSION"