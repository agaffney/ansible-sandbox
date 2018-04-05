#!/bin/bash

BASE_IMAGE=alpine:latest
IMAGE_NAME=ansible-sandbox
ANSIBLE_VERSION=2.5.0

set -x

docker run -ti --name ${IMAGE_NAME}-build $BASE_IMAGE sh -c "apk update; apk add py3-cryptography py3-pynacl py3-bcrypt; pip3 install ansible==${ANSIBLE_VERSION}; ansible-playbook --version"
docker commit ${IMAGE_NAME}-build ${IMAGE_NAME}:${ANSIBLE_VERSION}
docker tag ${IMAGE_NAME}:${ANSIBLE_VERSION} ${IMAGE_NAME}:latest
docker rm ${IMAGE_NAME}-build
