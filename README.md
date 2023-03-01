# bulk-http-check

Very fast concurrent check of many HTTP/s URLs.

## Usage examples
`urls.txt` is simple one url per line, e.g.:
~~~
https://httpbin.org/delay/10
https://httpbin.org/status/200
...
~~~

Simplest case, just get status (if no errors):
```shell
$ ./bulk-http-check < urls.txt 
https://httpbin.org/status/404 OK 404
https://httpbin.org/status/201 OK 201
https://httpbin.org/status/500 OK 500
https://asdf2sdcjsdsd.zs/sdf ERR Head "https://asdf2sdcjsdsd.zs/sdf": dial tcp: lookup asdf2sdcjsdsd.zs on 10.0.0.254:53: no such host
https://httpbin.org/status/202 OK 202
https://httpbin.org/json OK 200
https://httpbin.org/status/200 OK 200
https://www.google.com/1 OK 404
https://www.google.com/ OK 200
https://ifconfig.me/ OK 200
https://google.com/1 OK 404
https://httpbin.org/delay/10 ERR Head "https://httpbin.org/delay/10": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```

Show specific HTTP header (and use 20 concurrent connections):
```
$ ./bulk-http-check -n 20 -s content-type < urls.txt 
https://asdf2sdcjsdsd.zs/sdf ERR Head "https://asdf2sdcjsdsd.zs/sdf": dial tcp: lookup asdf2sdcjsdsd.zs on 10.0.0.254:53: no such host
https://ifconfig.me/ OK 200 text/plain; charset=utf-8
https://google.com/1 OK 404 text/html; charset=UTF-8
https://www.google.com/ OK 200 text/html; charset=ISO-8859-1
https://httpbin.org/json OK 200 application/json
https://httpbin.org/status/404 OK 502 text/html
... 
```

If you want to know content-size, it could be little more tricky. Most reliable is to combine `-l` flag with `-g` (use HTTP GET method instead of HEAD). Sometimes Content-Length (reported from server) does not reflect real payload size because of encoding/gzipping (in this case, it's unavailable in headers).

```
$ ./bulk-http-check -l -g < urls.txt 
https://httpbin.org/status/200 OK 200 0
https://www.google.com/1 OK 404 1562
https://httpbin.org/json OK 200 429
```

## Command-line options

```
$ ./bulk-http-check -h
Usage of ./bulk-http-check:
  -1	Disable HTTP/2 support
  -g	use GET method instead of HEAD
  -l	Show Content-length
  -n int
    	Number of connections (default 5)
  -s string
    	Show this header
  -t int
    	Timeout (5) (default 5)
```


