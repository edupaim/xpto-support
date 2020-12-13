#!/bin/bash

while true
do
	if (( $RANDOM % 2 )); then
		if (( $RANDOM % 2 )); then
		    echo "Request to route \"/negatives\" query customerDocument = 51537476467"
		    curl -H "x-auth-key: vgtd61gBEpw6HNWTovzDPuQkXTDS6H0P" localhost:81/negatives?customerDocument=51537476467
	        else
		    echo "Request to route \"/negatives\" query customerDocument = 26658236674"
		    curl -H "x-auth-key: vgtd61gBEpw6HNWTovzDPuQkXTDS6H0P" localhost:81/negatives?customerDocument=26658236674
		fi

	else
	      echo "Request to route \"/legacy/integrate\""
	      curl -H "x-auth-key: vgtd61gBEpw6HNWTovzDPuQkXTDS6H0P" -X POST localhost:81/legacy/integrate
	fi

	sleep 0.5
done