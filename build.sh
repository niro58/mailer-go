#!/bin/bash
set -e

kubectl config use-context personal

echo "=== Building niro-mailer ==="
docker build --no-cache -f ./Dockerfile -t ghcr.io/niro58/mailer-go:latest .

echo "=== Pushing niro-mailer ==="
docker push ghcr.io/niro58/mailer-go:latest

echo "=== Deploying niro-mailer ==="
kubectl apply -f .k8s/deployment.yaml
kubectl apply -f .k8s/service.yaml
kubectl apply -f .k8s/ingress.yaml
kubectl rollout restart deployment/niro-mailer -n personal

echo "=== Done ==="
