#!/bin/bash
docker build . -t goin
docker network create --subnet=172.18.0.0/24 stalinnet
docker run -t --name="n1" --ip 172.18.0.2 --network stalinnet --env NETINT='172.18.0.2' goin &
sleep 5
docker run -t --name="n2" --ip 172.18.0.3 --network stalinnet --env NETNODE='172.18.0.2' --env NETINT='172.18.0.3' goin &
sleep 5
