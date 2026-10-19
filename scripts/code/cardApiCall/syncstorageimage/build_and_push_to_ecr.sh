#!/bin/bash
set -Ceu

# 環境変数を読み込む
set -a
source .env
set +a

AWS_ECR_REPO_NAME=${ECR_REPOSITORY_NAME_SYNCSIMAGE:-"ocg-sync-storage-image"}

# ログイン
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com

# ビルド
docker build -t $AWS_ECR_REPO_NAME -f ./code/cardApiCall/cmd/aws_lambda/syncStorageImageWithDB/Dockerfile ./code/cardApiCall

# タグ付け
docker tag $AWS_ECR_REPO_NAME:latest $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$AWS_ECR_REPO_NAME:latest

# プッシュ
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$AWS_ECR_REPO_NAME:latest
