#!/usr/bin/env bash

set -e

# Docker Compose Install
curl -L https://github.com/docker/compose/releases/download/${COMPOSE_VERSION}/docker-compose-`uname -s`-`uname -m` > docker-compose
chmod +x docker-compose
sudo mv docker-compose /usr/local/bin

# Download and install Docker libs
curl -L https://github.com/Ortus-Solutions/docker-buildfiles/archive/master.zip > docker.zip
unzip docker.zip -d workbench
mv workbench/docker-buildfiles-master workbench/docker

# CommandBox Keys
sudo apt-key adv --keyserver keys.gnupg.net --recv 6DA70622
sudo echo "deb http://downloads.ortussolutions.com/debs/noarch /" | sudo tee -a /etc/apt/sources.list.d/commandbox.list

# Core testing install
sudo apt-get update && sudo apt-get --assume-yes install commandbox
box install
box server start
