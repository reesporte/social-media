# __socialmedia__ docs

__socialmedia__ is composed of three main services: 
1. [`proxy`](./services/proxy.md) : NGINX reverse proxy that handles all incoming requests and enforces https.
2. [`social_media`](./services/social_media.md) : Go backend and HTML/CSS/JS frontend.
3. [`db`](./services/db.md) : PostgreSQL database for storing users and posts.

## Pre-requisites
1. [docker](https://docs.docker.com/get-docker/)
2. [docker-compose](https://docs.docker.com/compose/install/)

## Build and run locally 
To build and run locally:
1. Ensure port 80 is open.
2. `docker-compose -f docker-compose.dev.yml up --build -d`

That's it! 

## Build and run on a server
1. Ensure ports 80 and 443 are open.
2. Replace `proxy/cert.crt` with your ssl certificate. 
3. Replace `proxy/cert.key` with your private key.
4. Replace every instance of `YOUR_DOMAIN_NAME_HERE` in `proxy/conf` with your domain name. 
5. `docker-compose up --build -d`

That's it!
