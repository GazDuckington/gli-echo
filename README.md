# stack

- docker-compose (to setup redis server)
- Go (Echo framework & Gorm)
- Sqlite
- Redis

# setup

- run redis server

```sh
docker-compose up -d
```

- run bin/main

# answers/features

1. DATABASE:
check for existing data in database on first fetch. if nothing is found fetch data directly from provided API, then store the response into database.

2.CACHE:
check for cached data on response, return cached data immediately if found. else go back to step 1

# documentations

go to: localhost:8000/swagger-ui/

1. get all breeds
- localhost:8000/
- localhost:8000/breeds/
- if breed is "sheepdog" append as "sheepdog-{subbreed}"
- if breed is "terrier" append as "terrier-{subbreed}" which is key that contains an array of that subbreed's images

2. get breed images
- localhost:8000/breeds/{breed}/images
- if breed is "shiba" get only images that are labeled with odd numbers
