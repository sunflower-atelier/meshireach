#!/bin/sh
ssh ubuntu@160.16.220.69 << EOC
  cd ~/meshireach
  docker-compose down
  git fetch origin deploy
  git reset --hard origin/deploy
  docker image prune -f
  docker container prune -f
  docker-compose build --no-cache
  nohup docker-compose up -d
EOC
