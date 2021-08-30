# db 

<p style="font-size: large; font-style: italic">PostgreSQL database for storing users and posts.</p>

* Runs off the [alpine](https://github.com/docker-library/postgres/blob/517c64f87e6661366b415df3f2273c76cea428b0/13/alpine/Dockerfile) postgres image.
* Sets up the `users` and `posts` table.  
* Creates an `admin` user with a default password of `IAmAn4dm!n` (You'll want to change this right off the bat on your first login)
