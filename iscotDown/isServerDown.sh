#! /bin/sh

response=$(curl -s localhost:3000/status)
if [ "$response" != "running" ] && [ find . -cmin -480 | grep -q "sentStatus.txt" ]; then
        date -u > ./sentStatus.txt
        ./main
fi
