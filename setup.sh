#!/bin/bash
docker network create --subnet=172.18.0.0/24 stalinnet
#docker run -d --name="n1" --ip 172.18.0.2 --network stalinnet goin

for i in `seq 3 10`;
do
	docker run -d --name="n$i" --ip 172.18.0.$i --network stalinnet --env NETNODE="172.18.0.2" --env NETINT="172.18.0.$i" goin
done
