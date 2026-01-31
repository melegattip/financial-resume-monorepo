#!/bin/sh

# Script de inicio para Cloud Run
echo "🚀 Iniciando nginx en puerto $PORT"

# Usar el puerto proporcionado por Cloud Run o 8080 por defecto
export PORT=${PORT:-8080}

# Sustituir la variable PORT en la configuración de nginx
envsubst '${PORT}' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf

# Mostrar la configuración generada para debug
echo "📋 Configuración nginx generada:"
cat /etc/nginx/conf.d/default.conf

# Iniciar nginx
nginx -g "daemon off;" 