#!/bin/bash

# выбрал локальный путь для хранения данных 
docker run -d -p 9000:9000 -p 9001:9001 \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin \
  -v ~/minio_data:/data \
  --name testMinio \
  myminio:latest