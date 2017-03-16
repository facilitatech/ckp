#!/usr/bin/env bash

# Colors
GREEN='\033[0;32m'
BLACK='\033[0;30m'
DARK_GRAY='\033[1;30m'
RED='\033[0;31m'
LIGHT_RED='\033[1;31m'
GREEN='\033[0;32m'
LIGHT_GREEN='\033[1;32m'
ORANGE='\033[0;33m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
LIGHT_BLUE='\033[1;34m'
PURPLE='\033[0;35m'
LIGHT_PURPLE='\033[1;35m'
CYAN='\033[0;36m'
LIGHT_CYAN='\033[1;36m'
LIGHT_GRAY='\033[0;37m'
WHITE='\033[1;37m'
NC='\033[0m'

clear

# Cabeçalho
echo ' '
printf "${GREEN}https://github.com/totalbr/dependency-check-php for the canonical source repository \n"
printf "Copyright (c) facilita.tech - 2016-2017\n"
printf "(http://facilita.tech) ${NC}"
echo ' '

if [ $(uname) == "Darwin" ]; then
    ENVIRONMENT='MAC'
else
    ENVIRONMENT='LINUX'
fi
echo ' '

if [ $ENVIRONMENT == 'LINUX' ]; then

    if which figlet > /dev/null; then
        printf "${GREEN}"
        figlet facilita.tech
    else
        apt-get install -y figlet
        printf "${GREEN}"
        figlet facilita.tech
    fi
    echo ' '
    printf "${NC}"
else
	if which figlet > /dev/null; then
		printf "${GREEN}"
		figlet facilita.tech
		printf "${GREEN}dependency-check-php \n${NC}"
	fi
	printf "${NC}"
echo ''
fi

# Docker
if which docker > /dev/null; then
    printf "${ORANGE}DOCKER${NC}\n"
    printf "${LIGHT_PURPLE}Generate new containers ?${NC} ${WHITE}[ ${PURPLE}1 ${WHITE}]${NC} \n${LIGHT_PURPLE}Delete all containers ?${NC} ${WHITE}  [ ${PURPLE}2 ${WHITE}]${NC} \n${LIGHT_PURPLE}Start new build ?${NC} ${WHITE}        [ ${PURPLE}3 ${WHITE}]${NC}\n"
    read gerar

    if [ -n "$gerar" ]; then
        if [ $gerar == '1' ]; then
            printf "${ORANGE}Generating new containers ... ${NC}\n"
            docker-compose ps
            docker-compose up -d
            docker-compose ps
        fi
        if [ $gerar == '2' ]; then
            printf "${ORANGE}Removing all containers ... ${NC}\n"
            docker-compose kill
            docker-compose rm
        fi
        if [ $gerar == '3' ]; then
        	printf "${LIGHT_PURPLE}Would you like to start a new compilation with cache?${NC} ${WHITE} [ ${PURPLE}yes ${WHITE}]: ${NC} "
        	read cache

        	printf "${ORANGE}Starting a new build process ... ${NC}\n"
        	if [ -n "$cache" ]; then
				if [ $cache == 'no' ]; then
					docker-compose build --no-cache
				fi
				if [ $cache == 'yes' ]; then
					docker-compose build
				fi
        	else
        	    docker-compose build
        	fi
        fi
    fi
    echo ' '
else
    printf "${BLUE}Installation of docker not found${NC}\n"
fi