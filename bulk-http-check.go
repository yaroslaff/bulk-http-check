package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex
var processed = 0

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}

func checkurl(id int, client *http.Client, ch chan string, useget bool, header string, clen bool) {

	var res *http.Response
	var err error
	var content_len int64

	for url := range ch {
		if useget {
			res, err = client.Get(url)
		} else {
			res, err = client.Head(url)
		}

		mu.Lock()
		processed++
		mu.Unlock()

		if err != nil {
			fmt.Printf("%s ERR %s\n", url, err)
		} else {

			if clen && useget {
				content_len, err = io.Copy(ioutil.Discard, res.Body)
			} else {
				content_len = res.ContentLength
			}

			if len(header) > 0 {
				// print header
				fmt.Printf("%s OK %d %s\n", url, res.StatusCode, res.Header.Get(header))
			} else if clen {
				fmt.Printf("%s OK %d %d\n", url, res.StatusCode, content_len)
			} else {
				fmt.Printf("%s OK %d\n", url, res.StatusCode)
			}

		}
	}
}

func main() {
	var wg sync.WaitGroup

	ch := make(chan string)

	var nconn = flag.Int("n", 5, "Number of connections")
	var clen = flag.Bool("l", false, "Show Content-length")
	var no2 = flag.Bool("1", false, "Disable HTTP/2 support")
	var header = flag.String("s", "", "Show this header")
	var useget = flag.Bool("g", false, "use GET method instead of HEAD")
	timeout := flag.Int("t", 5, "Timeout")
	var benchmark_period = flag.Int("b", 0, "print benchmark every N seconds")
	var exit = flag.Int("x", 0, "exit after N seconds")

	var started = time.Now().Unix()
	var last_printed = started

	flag.Parse()

	println("started")

	// make default client
	client := &http.Client{
		Timeout: (time.Second * time.Duration(*timeout))}

	if *no2 {
		client = &http.Client{
			Transport: &http.Transport{
				TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			},
		}
	}

	for n := 1; n <= *nconn; n++ {
		wg.Add(1)
		n := n
		go func() {
			defer wg.Done()
			checkurl(n, client, ch, *useget, *header, *clen)
		}()

	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url := scanner.Text()
		if strings.HasPrefix(url, "http") {
			ch <- url
		}

		now := time.Now().Unix()
		uptime := now - started
		if *benchmark_period > 0 && now >= last_printed+int64(*benchmark_period) {
			rate := float64(processed) / float64(uptime)
			last_printed = now
			fmt.Fprintf(os.Stderr, "# runs %d seconds, processed %d, rate: %.2f/sec\n", uptime, processed, rate)
		}

		if *exit > 0 && int(uptime) >= *exit {
			os.Exit(0)
		}

	}
	close(ch)

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	wg.Wait()
}
