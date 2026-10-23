#!/bin/bash
set -Ceu

act \
  -W .github/workflows/lamda-syncstorageimage-deploy.yml \
  --secret-file act_secretfile \
  --container-architecture linux/amd64
