# zero-http

This is a replacement for `python -m SimpleHTTPServer` or `python -m http.server`.

## Why
It explicitly makes sure not to cache any results and probably some other things eventually

## Config
It reads a `.zero-http` folder either from the working directory or `~/.zero-http` see [this](.zero-http)
folder for an example

## Docker

`docker run -it -p 9000:9000 -v $(pwd)/config:/zero/.zero-http -v $(pwd):/zero/srv dylanowen/zero-http:latest zero-http`

### docker-compose.yml
```yaml
version: "2"
services:
  serve-site:
    image: dylanowen/zero-http:latest
    volumes:
    - $PWD:/zero/srv
    - $PWD/config:/zero/.zero-http
    ports:
    - 9000:9000
    command: zero-http
```

### Publish
`make docker`  
`docker push dylanowen/zero-http:latest`