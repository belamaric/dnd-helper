#!/bin/bash

docker run -d  \
    -e GIT_SYNC_REPO=git@github.com:johnbelamaric/dnd-helper.git  \
    -e GIT_SYNC_DEST=/git/dnd-helper  \
    -e GIT_SYNC_BRANCH=master \
    -e GIT_SYNC_REV=FETCH_HEAD \
    -e GIT_SYNC_WAIT=10 \
    -v /home/jbelamaric:/git \
    -v /home/jbelamaric/.ssh:/root/.ssh  openweb/git-sync:0.0.1

while [ ! -d /home/jbelamaric/dnd-helper/html ]; do
  echo Waiting for git-sync to pull the repo...
  sleep 2
done

echo Building the tool...
docker run -v /home/jbelamaric/dnd-helper/statblock5e:/go/src/github.com/johnbelamaric/dnd-helper/statblock5e infoblox/buildtool sh -c "cd /go/src/github.com/johnbelamaric/dnd-helper/statblock5e && go get && go build"

ID=$(docker run -d -v /home/jbelamaric/dnd-helper/Caddyfile:/etc/Caddyfile -v /home/jbelamaric/dnd-helper/html:/var/site/dmtool.info/html -v /home/jbelamaric/certs:/root/.caddy -p 80:80 -p 443:443 -v /home/jbelamaric/logs:/var/site/dmtool.info/logs abiosoft/caddy)

docker run -d  -v /home/jbelamaric/DnDAppFiles/Bestiary:/root/data -v /home/jbelamaric/dnd-helper/statblock5e/statblock5e:/root/statblock5e --net=container:$ID ubuntu sh -c "/root/statblock5e -s :8000 -d /root"
