#!/bin/bash

commands=("docker exec -it $1 mkdir /app"
	  "docker exec -it $1 apk add --no-cache bash"
	  "docker cp $2 $1:/app"
	  "docker commit --change='ENTRYPOINT [\"/app/$4\"]' $1 $1/image:$3"
	  "docker save $1/image:$3 -o /images/$1.tar.gz")

for command in "${commands[@]}"
do
  echo "Executing command \"$command\""
  OUTPUT=$($command 2>&1)

  if [ $? -ne 0 ]; then
    printf "$OUTPUT\n"
    exit 1
  fi
done
