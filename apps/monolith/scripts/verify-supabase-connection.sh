#!/bin/bash

# Script de verificación de conexión a Supabase
# Verifica que podamos conectarnos a la nueva base de datos

# Colores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================="
echo "Verificación de Conexión a Supabase"
echo "========================================="
echo ""

# Verificar que DATABASE_URL esté configurada
if [ -z "$DATABASE_URL" ]; then
    echo -e "${RED}❌ ERROR: DATABASE_URL no está configurada${NC}"
    echo "Por favor, ejecuta: export DATABASE_URL='tu-connection-string'"
    exit 1
fi

echo -e "${YELLOW}🔍 Verificando conexión...${NC}"
echo ""

# Intentar conectar y ejecutar una query simple
psql "$DATABASE_URL" -c "SELECT current_database(), version();" 2>&1

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ Conexión exitosa a Supabase${NC}"
    echo ""
    
    # Verificar que la base está vacía (nueva)
    echo -e "${YELLOW}🔍 Verificando tablas existentes...${NC}"
    TABLE_COUNT=$(psql "$DATABASE_URL" -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';")
    
    if [ "$TABLE_COUNT" -eq 0 ]; then
        echo -e "${GREEN}✅ Base de datos nueva (sin tablas)${NC}"
    else
        echo -e "${YELLOW}⚠️  Hay $TABLE_COUNT tablas en la base de datos${NC}"
    fi
else
    echo ""
    echo -e "${RED}❌ Error al conectar a Supabase${NC}"
    echo ""
    echo "Posibles causas:"
    echo "1. Contraseña incorrecta (verifica URL encoding)"
    echo "2. Región incorrecta en la connection string"
    echo "3. Firewall bloqueando la conexión"
    exit 1
fi

echo ""
echo "========================================="
