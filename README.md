# tinyurl
A simple TinyURL app.
Go practice

* /createTiny with {"original_url": <url_to_shorten>} will return {"encoded_url": <shortened_url}
* /getTiny with {"encoded_url": <shortened_url} will return {"original_url": <url_to_shorten>}

system design

A go app that serves above requests on 8080

A go app that acts as a cache (using Memcached) to the DB (postgresql) listening on 8082

A go app that acts as the key generation service listening on 8081 and storing random 6letter strings in redis

docker-compose up with bring tinyurl, kgs, cache, memcached, redis, postgresql containers up

req.sh has some example requests to test the app
