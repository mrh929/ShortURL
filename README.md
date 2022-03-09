# ShortURL

This project is an open source golang server for generating short URLs.

Aims to write a self-hosted URL Shortener like [URL Shortener (shorturl.at)](https://www.shorturl.at/).

### example

https://m1ku.in/s/PGMDr

## Run

### deploy http server

ShortURL Webapp can be easily deployed by docker-compose if you do not need customization.

```bash
docker-compose up -d
```

Then you can access server using **http port** (https need to set up nginx which will be discussed later).

### env

| Name              | Type   | Description                                                                          |
| ----------------- | ------ | ------------------------------------------------------------------------------------ |
| **SRV_PASSWD**    | string | User need to type password to access this web app.                                   |
| SRV_HOST          | string | IP that **web app** listens to. **Default: 0.0.0.0**                                 |
| SRV_PORT          | string | Port that **web app** listens to. **Default: 8000**                                  |
| SRV_PROTO         | string | Generated link protocal. **Default: http**                                           |
| SRV_BASE_PATH     | string | Generated link base path. If empty, this will be set to the value of header['HOST']. |
| SQL_ROOT_PASSWD   | string | Password to access mysql. **Default: test**                                          |
| SQL_HOST          | string | Mysql host. **Default: 127.0.0.1**                                                   |
| SQL_PORT          | string | Mysql port. **Default: 3306**                                                        |
| SQL_DATABASE_NAME | string | Generated database name. **Default: db**                                             |

### set up https using nginx

We can set up a reverse proxy for users to access. Usually I will set up nginx reverse proxy in host machine to access guest machines.

Notice that the header **X-Forwarded-Proto** is necessary if you want to set up an https server because it is used to distinguish protocals.

```nginx
server {
        listen 443 ssl;
        server_name your_domain;

        ssl_certificate /etc/letsencrypt/live/your_domain/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/your_domain/privkey.pem;
        ssl_session_timeout 5m;
        ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
        ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:HIGH:!aNULL:!MD5:!RC4:!DHE;
        ssl_prefer_server_ciphers on;

        # proxy to ShortURL
        location / {
            proxy_pass https://127.0.0.1:8000;

         # necessary. we need header X-Forwarded-Proto to distinguish 
         # bewteen HTTPS and HTTP server
            proxy_set_header   X-Forwarded-Proto https;

            proxy_set_header   X-Forwarded-Port  443;
            proxy_set_header   Host             $host;
            proxy_set_header   X-Real-IP        $remote_addr;
            proxy_set_header   X-Forwarded-For  $proxy_add_x_for
warded_for;
        }
}
```

## dev

```bash
.
├── data
├── docker
│   └── web
│       ├── Dockerfile
│       └── entrypoint.sh
├── docker-compose.yml
├── README.md
└── src
    ├── db.go           # access sql
    ├── go.mod
    ├── go.sum
    ├── handler.go      # handle different requests
    ├── main.go         # entrypoint
    ├── sql
    │   └── table.sql   # sql structure
    ├── tool.go         # other tools
    └── web
        └── index.html  # html page
```
