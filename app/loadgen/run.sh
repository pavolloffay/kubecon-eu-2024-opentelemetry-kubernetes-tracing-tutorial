#!/bin/bash

SLOWDOWN_IN_SECONDS=${SLOWDOWN_IN_SECONDS:-3}
URL=${URL:-"http://frontend:4000/"}

while true;
do
MAX=$(($(($RANDOM%30))+10))

for i in `seq 1 ${MAX}`
do
PLAYERS=(`echo -e "Pavol\nBenedikt\nYuri\nKristina\nSeverin" | shuf`)
timeout 5 curl -s "${URL}?player1=${PLAYERS[1]}&player2=${PLAYERS[2]}" &
sleep 1
done

wait
 
sleep ${SLOWDOWN_IN_SECONDS}
done
