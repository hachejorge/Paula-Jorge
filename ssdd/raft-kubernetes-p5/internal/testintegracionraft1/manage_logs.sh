#!/bin/bash

# Definir los nombres de los archivos de entrada
file1="../../logs_raft/192.168.3.6_31139.txt"
file2="../../logs_raft/192.168.3.7_31138.txt"
file3="../../logs_raft/192.168.3.8_31137.txt"

# Definir los nombres de los archivos de salida
join_file="join.txt"
sorted_file="sorted_join.txt"

# Verificar que los archivos de entrada existen
if [[ ! -f "$file1" || ! -f "$file2" || ! -f "$file3" ]]; then
    echo "Error: Uno o más archivos de entrada no existen."
    exit 1
fi

# Extraer las primeras 100 líneas de cada archivo y combinarlas en join.txt
head -n 300 "$file1" > "$join_file"
head -n 300 "$file2" >> "$join_file"
head -n 300 "$file3" >> "$join_file"

# Ordenar las líneas del archivo combinado por el tercer atributo (campo), separados por espacios,
# y guardar el resultado en sorted_join.txt
sort -k3,3 "$join_file" > "$sorted_file"

# Confirmación de creación de archivos
echo "Archivo combinado creado: $join_file"
echo "Archivo ordenado creado: $sorted_file"
