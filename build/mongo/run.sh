#!/bin/bash

NAME="testMongo"
PORT="27017"
DATABASE="labyrinth"
VERSION="latest"

docker run -d \
  --name "$NAME" \
  -p "$PORT:$PORT" \
  -e MONGO_INITDB_DATABASE="$DATABASE" \
  mymongodb:"$VERSION" \
  mongod --replSet rs0 --bind_ip_all --noauth

echo "Waiting for MongoDB to start..."
sleep 5

# Инициализируем replica set
echo "Initializing replica set..."
docker exec -it "$NAME" mongosh --eval "rs.initiate({
  _id: 'rs0',
  members: [
    { _id: 0, host: 'localhost:27017' }
  ]
})"