#!/bin/bash

# Lista de IPs de las máquinas remotas
MACHINES=("192.168.3.6" "192.168.3.7" "192.168.3.8")

# Usuario en las máquinas remotas
REMOTE_USER="a872838"

# Comando para buscar y matar procesos en cada máquina
REMOTE_COMMAND='
    echo "Buscando procesos main o go para el usuario $USER...";
    PIDS=$(ps -lu $USER | grep -E "main" | awk '\''{print $4}'\'');
    if [ -z "$PIDS" ]; then
        echo "No se encontraron procesos main o go.";
    else
        echo "Procesos encontrados: $PIDS";
        for PID in $PIDS; do
            kill -9 $PID 2>/dev/null && echo "Proceso $PID eliminado." || echo "Error al eliminar proceso $PID.";
        done
    fi;
'

# Verificar si se pasó un argumento
if [ $# -eq 1 ]; then
    # Validar si el argumento es un número entre 0 y 2
    if [[ $1 =~ ^[0-2]$ ]]; then
        MACHINE=${MACHINES[$1]}
        echo "Conectando a la máquina $MACHINE para detener procesos..."
        ssh "$REMOTE_USER@$MACHINE" "$REMOTE_COMMAND"
        echo "Finalizado en la máquina $MACHINE."
        exit 0
    else
        echo "Error: el parámetro debe ser un número entre 0 y 2."
        exit 1
    fi
fi

# Si no hay argumentos, ejecutar el script para todas las máquinas
for MACHINE in "${MACHINES[@]}"; do
    echo "Conectando a $MACHINE..."
    ssh "$REMOTE_USER@$MACHINE" "$REMOTE_COMMAND"
    echo "Finalizado en $MACHINE."

done

echo "Finalizado en todas las máquinas."
