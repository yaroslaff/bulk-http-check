# bulk http check

Very fast concurrent check of many HTTP/s URLs. (Few thousands requests per seconds, depending on hardware and network bandwidth)

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
...
```

Show specific HTTP header (and use 20 concurrent connections):
```
$ ./bulk-http-check -n 20 -s content-type < urls.txt 
https://httpbin.org/json OK 200 application/json
... 
```

If you want to know content-size, it could be little more tricky. Most reliable is to combine `-l` flag with `-g` (use HTTP GET method instead of HEAD). Sometimes Content-Length (reported from server) does not reflect real payload size because of encoding/gzipping (in this case, it's unavailable in headers).

```
$ ./bulk-http-check -l -g < urls.txt 
https://httpbin.org/status/200 OK 200 0
https://www.google.com/1 OK 404 1562
https://httpbin.org/json OK 200 429
```

## Benchmarks
If specify `-b N`, bulk-http-check will print benchmark results on stderr, like:
~~~
# runs 10 seconds, processed 14302, rate: 1430.20/sec
# runs 20 seconds, processed 49920, rate: 2496.00/sec
# runs 30 seconds, processed 125943, rate: 4198.10/sec
~~~

Option `-x N` to eXit automatically after N seconds.

core 2 duo is my home desktop with 100Mbps Internet. CX11 is cheapest hetzner VPS with 2Gb RAM, AX51-NVMe is dedicated Hetzner server with 8 cores, 16 threads and 64Gb.


| Connections  | core2duo    | CX11 |  AX51-NVMe |
|---           |---          |---   |---         |
| 1            | 5           | 23   |         24 |
| 10           | 50          |113   |        208 |
| 100          | 255         |1170  |       1829 |
| 1000         | 1188        |1743 *|       2540 |
| 10000        | 3098 *      |--    |       3458 |

`*` - Errors happened during tests, mostly timeouts because hit bandwidth limit.

## Install

### Install via go install
If you have Golang installed:
~~~
go install github.com/yaroslaff/bulk-http-check
~~~

### Install from repo
~~~
git clone https://github.com/yaroslaff/bulk-http-check
cd bulk-http-check
go buld
cp bulk-http-check /usr/local/bin
~~~

### Download static binary 
Visit [Latest release](https://github.com/yaroslaff/bulk-http-check/releases/latest) and download `bulk-http-check` from there.
Then do:
~~~
chmod +x bulk-http-check
cp bulk-http-check /usr/local/bin
~~~



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


