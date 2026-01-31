#!/bin/sh
# Script para iniciar serve en puerto 8080 para Cloud Run
echo "Starting serve on port 8080"
serve -s build -l 8080 