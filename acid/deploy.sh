#!/bin/bash

# Check if both environment and IMAGE_TAG arguments are provided
if [ $# -ne 2 ]; then
    echo "Error: Both environment and IMAGE_TAG arguments are required."
    echo "Usage: $0 <env> <image_tag>"
    echo "Accepted values for <env>: prd, dev"
    exit 1
fi

# Get environment and IMAGE_TAG from command line arguments
ENV=$1
IMAGE_TAG=$2

# Validate environment
if [ "$ENV" != "prd" ] && [ "$ENV" != "dev" ]; then
    echo "Error: Invalid environment. Accepted values are 'prd' or 'dev'."
    exit 1
fi

# Set folder based on environment
FOLDER=$ENV
NAMESPACE=ns-coffee-order-demo-$FOLDER

# Delete the previous deployment
kubectl delete -f inventory-service/$FOLDER/deployment.yaml -n $NAMESPACE

# Use sed to replace IMAGE_TAG in the deployment.yaml file
sed -i "s|IMAGE_TAG|$IMAGE_TAG|g" inventory-service/$FOLDER/deployment.yaml

kubectl apply -f inventory-service/$FOLDER/deployment.yaml -n $NAMESPACE
kubectl apply -f inventory-service/$FOLDER/service.yaml -n $NAMESPACE

# Revert the change in deployment.yaml to avoid committing the modified file
sed -i "s|$IMAGE_TAG|IMAGE_TAG|g" inventory-service/$FOLDER/deployment.yaml