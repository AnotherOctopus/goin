#!/bin/bash
docker network create --subnet=172.18.0.0/24 stalinnet
docker run -d --ip 172.18.0.2 --network stalinnet goin
for i in `seq 3 10`;
do
        docker run -d --ip 172.18.0.$i --network stalinnet --env NETNODE='172.18.0.2' goin
done
