#!/usr/bin/env bash

ENVBR=$1
ITEM="gopa"

IMAGE_NAME="registry-vpc.cn-shanghai.aliyuncs.com/graviti-devops/gopa"
BRANCH=$(git rev-parse --abbrev-ref HEAD)
COMMIT_ID=`git log --oneline -1 | awk '{print $1}'`
DATE_TIME=`date +%Y%m%d`
NETWORK_NAME=GlabelNet
HOST=$(curl http://100.100.100.200/latest/meta-data/private-ipv4)

declare -A APOLLO_ADDR
APOLLO_ADDR=([dev]="172.19.50.234:8080" [fat]="172.19.100.92:8080" [uat]="172.19.40.72:8080" [pro]="172.19.10.140:8080")

if [[ $ENVBR == "dev" ]] || [[ $ENVBR == "fat" ]] || [[ $ENVBR == "uat" ]] || [[ $ENVBR == "pro" ]];then
    
    IMG_TAG=${ENVBR}-${BRANCH}-${COMMIT_ID}-${DATE_TIME}
    echo "$IMAGE_NAME:$IMG_TAG"

    cd ../ && docker build -f deploy/Dockerfile -t ${IMAGE_NAME}:${IMG_TAG} --no-cache . && cd deploy/
    docker push ${IMAGE_NAME}:${IMG_TAG}

    # create docker network
    net_flag=`docker network ls|grep ${NETWORK_NAME}`
    if [ $? -ne 0 ];then
      docker network create --subnet 10.200.0.0/16 --ip-range 10.200.200.0/24 ${NETWORK_NAME}
    fi

    # adapt config
    sed -i "/image: \S*${ITEM}\:*/c\    image: ${IMAGE_NAME}:${IMG_TAG}" docker-compose.yaml
    sed -i "s/APP_ENV=\S*/APP_ENV=online/g" docker-compose.yaml
    sed -i "s/ENVIRONMENT=\S*/ENVIRONMENT=${ENVBR}/g" docker-compose.yaml
    sed -i "s/APOLLO_ADDR=\S*/APOLLO_ADDR=${APOLLO_ADDR["${ENVBR}"]}/g" docker-compose.yaml
    sed -i "s/HOST_IP=\S*/HOST_IP=${HOST_IP}/g" docker-compose.yaml

    # setup container
    docker-compose up -d
else
    echo "ERROR! Checkout ENV Config Please!"
    exit 1
fi
