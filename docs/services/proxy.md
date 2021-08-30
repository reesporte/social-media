# proxy

<p style="font-size: large; font-style: italic">NGINX reverse proxy that handles all incoming requests and enforces https.</p>

* Runs off the [alpine](https://github.com/nginxinc/docker-nginx/blob/f958fbacada447737319e979db45a1da49123142/mainline/alpine/Dockerfile) nginx image. 
* Redirects all non-http traffic to https.
* Rate limits requests.
* Proxies requests to the `social_media` service
