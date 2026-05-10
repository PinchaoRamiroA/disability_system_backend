#!/bin/bash
set -e

DB_NAME="disability_system_db"
DB_USER="root"
DB_PASSWORD="password"
DB_HOST="localhost"
DB_PORT="5432"

echo "=== Database Migration Script ==="
echo "Database: $DB_NAME"
echo ""

# Verificar estado actual
echo "Verificando estado actual..."
TABLE_COUNT=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE';" 2>/dev/null | tr -d ' ')
echo "Tablas existentes: $TABLE_COUNT"

# Si ya hay tablas, preguntar si desea resetear
if [ "$TABLE_COUNT" -gt 0 ]; then
    echo ""
    echo "⚠️  La base de datos ya contiene tablas."
    read -p "¿Desea resetear la base de datos? (yes/no): " -r
    if [[ "$REPLY" == "yes" ]]; then
        echo "Reseteando base de datos..."
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" 2>/dev/null
        echo "Base de datos reseteada."
    else
        echo "Ejecutando solo seed data..."
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f migrations/000002_seed_data.up.sql 2>/dev/null
        echo ""
        echo "✅ Seed data aplicado exitosamente!"
        exit 0
    fi
fi

# Si migrate está disponible, usarlo
if command -v migrate &> /dev/null; then
    echo "Usando golang-migrate..."
    DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
    migrate -path migrations -database "$DB_URL" up
else
    echo "Ejecutando migraciones SQL directamente..."
    
    # Ejecutar schema
    echo "  - Aplicando schema..."
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f migrations/000001_init_schema.up.sql 2>/dev/null
    
    # Esperar a que se cree la tabla de migraciones
    sleep 1
    
    # Insertar seed data
    echo "  - Insertando datos iniciales..."
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f migrations/000002_seed_data.up.sql 2>/dev/null
fi

echo ""
echo "=== Resumen ==="
echo "Tablas creadas: $(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c '\dt' 2>/dev/null | grep -c 'table' || echo '0')"

echo ""
echo "Verificando datos seed..."
echo "  Roles: $(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c 'SELECT COUNT(*) FROM rol;' 2>/dev/null | tr -d ' ') registros"
echo "  Estados incapacidad: $(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c 'SELECT COUNT(*) FROM estado_incapacidad;' 2>/dev/null | tr -d ' ') registros"
echo "  Tipos entidad: $(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c 'SELECT COUNT(*) FROM tipo_entidad;' 2>/dev/null | tr -d ' ') registros"
echo "  Entidades: $(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c 'SELECT COUNT(*) FROM entidad;' 2>/dev/null | tr -d ' ') registros"
echo "  Permisos: $(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c 'SELECT COUNT(*) FROM permisos;' 2>/dev/null | tr -d ' ') registros"

echo ""
echo "✅ Migración completada exitosamente!"