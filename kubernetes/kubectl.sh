#!/usr/bin/env bash

kubectl apply -f vaccinator-networkpolicy.yaml,api-persistentvolumeclaim.yaml,database-postgres-persistentvolumeclaim.yaml,pgadmin-deployment.yaml,postgres-deployment.yaml,vaccinator-deployment.yaml,pgadmin-service.yaml,postgres-service.yaml,vaccinator-service.yaml