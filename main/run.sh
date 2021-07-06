#!/bin/bash

if ! command -v reflex &> /dev/null
then
  echo "Not found reflex tool"
  echo "Downloading"
  go get github.com/cespare/reflex
fi

cd ..
sudo docker-compose stop
sudo docker-compose up -d --build
cd main

reflex -s -R vendor. -r \.go$ -- bash start.sh
