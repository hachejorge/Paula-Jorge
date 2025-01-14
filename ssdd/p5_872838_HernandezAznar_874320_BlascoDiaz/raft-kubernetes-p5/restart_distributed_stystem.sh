#!/bin/bash

echo "Eliminando cluster..."
docker stop kind-worker
docker stop kind-worker{2,3,4}
docker stop kind-control-plane
docker stop kind-registry

docker rm kind-worker
docker rm kind-worker{2,3,4}
docker rm kind-control-plane
docker rm kind-registry

echo "Creando cluster..."
./kind-with-registry.sh

echo "\nRecompilando los ejecutables del cliente y servidor..."
rm Dockerfiles/cliente/cltraft
rm Dockerfiles/servidor/srvraft
CGO_ENABLED=0 go build -o Dockerfiles/servidor/srvraft cmd/srvraft/main.go
CGO_ENABLED=0 go build -o Dockerfiles/cliente/cltraft pkg/cltraft/cltraft.go

echo "Creando las imagenes en Docker..."
docker build Dockerfiles/servidor/. -t localhost:5001/servidor:latest
docker push localhost:5001/servidor:latest
docker build Dockerfiles/cliente/. -t localhost:5001/cliente:latest
docker push localhost:5001/cliente:latest

echo "Ejecutando Kubernetes..."
kubectl delete statefulset raft
kubectl delete pod client
kubectl delete service raft-service
kubectl create -f statefulset_go.yaml
