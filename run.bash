#!/bin/bash -ex

tag=$1

if [[ -z "$tag" ]]; then
	echo 'No tag given. Changing command to ./run.bash example.'
	tag=example
fi

docker build -t benlubar/it_was_inevitable:"$tag" \
	--build-arg "tag=$tag" .
if [[ "$tag" == "example" ]]; then
	docker run --rm -ti \
		--security-opt seccomp="$(pwd)/seccomp.json" \
		benlubar/it_was_inevitable:"$tag"
else
	docker stop it_was_inevitable_"$tag" || :
	docker rm -v it_was_inevitable_"$tag" || :
	docker run -d --name it_was_inevitable_"$tag" \
		--restart unless-stopped \
		--security-opt seccomp="$(pwd)/seccomp.json" \
		benlubar/it_was_inevitable:"$tag"
fi
