# hits

yet another whoami container with rate counter and remote ip and path

```console
$ docker run -p 8080:8080 ghcr.io/matti/hits:latest
2022/12/29 11:09:09 hits listens in :8080
2022/12/29 11:09:13 85af1865188b hit at path / from 172.17.0.1:65116 - rate is 1 requests in 10s
2022/12/29 11:09:16 85af1865188b hit at path /hello from 172.17.0.1:65118 - rate is 2 requests in 10s
```
