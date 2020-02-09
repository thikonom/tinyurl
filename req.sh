#!/bin/sh

# curl -XPOST -H"Content-Type: application/json" http://localhost:8080/createTiny \
#   -d'{"email": "thikonom@gmail.com", "original_url": "http://myurl/1"}'

# curl -XPOST -H"Content-Type: application/json" http://localhost:8081/generateKey
# curl -XGET -H"Content-Type: application/json" http://localhost:8081/getKey

# curl -XPOST -H"Content-Type: application/json" http://localhost:8082/getCacheKey\
  # -d'{"encoded_url": "NBZxXO"}'

curl -XPOST http://localhost:8080/getTiny -d'{"encoded_url": "EQtUid"}'
