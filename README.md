# u1m

A loader and API for the Umbrella Top 100 Domains

https://twitter.com/opendns/status/809143307821518848

http://s3-us-west-1.amazonaws.com/umbrella-static/index.html

## Usage

This API is available as a service at [http://u1m.jn.gl](http://u1m.jn.gl)!

```
$ curl -s http://u1m.jn.gl/rank/21675 | jq
{
  "domain": "news.ycombinator.com",
  "rank": 21675
}
```

Instructions for building your own instance below.

## Build

This requires a MySQL instance with the schema applied. See `api/db/migrations` for the schema.
You can apply this by hand, or modify `db/dbconf.yml` and use [goose](https://bitbucket.org/liamstask/goose/).

Load the data (warning: this is a 30MB compressed zipfile, with 1M rows of data)

```
$ make
$ docker run -it \
    -e APP_DB="root:secret@tcp(db.yourdomain.com:3306)/u1m" \
    u1m/loader
```

Run an instance of the API

```
$ make
$ docker run \
    -e APP_DB="root:secret@tcp(db.yourdomain.com:3306)/u1m"
    -e APP_BIND=":8080"
    -p 8080:8080 \
    u1m/u1mapi
```

Test it out

```
$ curl http://localhost:8080/domain/github.com
```
