#!/bin/bash

# Mapa de nodos y puertos
declare -A nodes=(
    ["192.168.3.8"]="31137"
    ["192.168.3.7"]="31138"
    ["192.168.3.6"]="31139"
)

# Directorios
PROJECT_DIR="/misc/alumnos/sd/sd2425/a872838/practica3/cmd/srvraft"

# Crear un array con las claves (nodos) y ordenarlas
sorted_nodes=($(for node in "${!nodes[@]}"; do echo $node; done | sort))

# Construir lista de endpoints
endpoints=()
for node in "${sorted_nodes[@]}"; do
    endpoints+=("${node}:${nodes[$node]}")
done

# Iniciar servidores
index=0
for node in "${sorted_nodes[@]}"; do
    port=${nodes[$node]}

    echo "Iniciando servidor en $node:$port..."
    ssh "$node" "
        cd $PROJECT_DIR &&
        nohup ./main $index ${endpoints[*]} > /dev/null 2>&1 &
        #nohup go run $BIN_FILE $index ${endpoints[*]} > /dev/null 2>&1 &
    "
    ((index++))
done
