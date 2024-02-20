#!/bin/bash

mkdir -p ./specs

# Pull images
for image in $(cat ./images.txt); do
    echo "Pulling ${image}"
    docker pull $image

    # Does it have gpu in the name?
    echo $image | grep gpu
    retval=$?
    if [[ "$retval" -eq 0 ]]; then
        echo "Image has GPU"
        hasGpu=yes
    else
        echo "Image does not has GPU"
        hasGpu=no
    fi
    tag=$(python -c "print('$image'.split(':')[-1])")
    cmd=". /etc/profile && /tmp/data/install-and-generate.sh $hasGpu /tmp/data/specs/compspec-$tag.json"
    echo "docker run --entrypoint /bin/bash -v $PWD:/tmp/data $image -c "$cmd""
    docker run --entrypoint /bin/bash -v $PWD:/tmp/data $image -c "$cmd"

    # This generates ./specs/compspec-intel-mpi-rocky-9-amd64.json, let's push to a registry with oras
    oras push $image-compspec --artifact-type application/org.supercontainers.compspec ./specs/compspec-$tag.json:application/org.supercontainers.compspec
done