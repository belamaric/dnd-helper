#!/bin/bash

ID=$(docker run -d -v /home/jbelamaric/dnd-helper/Caddyfile:/etc/Caddyfile -v /home/jbelamaric/dnd-helper/html:/var/site/dmtool.info/html -v /home/jbelamaric/certs:/root/.caddy -p 80:2015 -v /home/jbelamaric/logs:/var/site/dmtool.info/logs abiosoft/caddy)

docker run -d  -v /home/jbelamaric/DnDAppFiles/Bestiary:/root/data -v /home/jbelamaric/dnd-helper/statblock5e/statblock5e:/root/statblock5e --net=container:$ID ubuntu sh -c "/root/statblock5e -s :8000 -d /root"
