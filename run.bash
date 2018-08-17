#!/bin/bash -ex

docker build -t benlubar/it_was_inevitable .
docker stop it_was_inevitable || :
docker rm -v it_was_inevitable || :
docker run -d --name it_was_inevitable \
	--restart unless-stopped \
	--security-opt seccomp="`pwd`/seccomp.json" \
	benlubar/it_was_inevitable
