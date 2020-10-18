package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	var threads int
	var wg sync.WaitGroup
	urls := make(chan string)

	flag.IntVar(&threads, "t", 32, "Specify number of threads to run")
	flag.Parse()

	//Make The Client
	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    time.Second,
		DisableCompression: true,
		DialContext: (&net.Dialer{
			Timeout:   time.Second * 10,
			KeepAlive: time.Second,
		}).DialContext,
	}
	client := &http.Client{Transport: tr}

	//Recieve input from stdin
	input := bufio.NewScanner(os.Stdin)
	go func() {
		for input.Scan() {
			urls <- input.Text()
		}
		close(urls)
	}()

	//workers
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker(urls, client, &wg)
	}
	wg.Wait()
}

func checkcors(s string, client *http.Client) {
	if !checkheaders(s, s, client) {
		return
	}
	if checkheaders(s, "https://notexisting.com", client) {
		fmt.Println(s, "is vulnerable to cors")
		return
	}
}

func checkheaders(url, origin string, client *http.Client) bool {
	//Initialize Headers
	headers := map[string]string{"Origin": origin,
		"Cache-Control": "no-cache",
		"User-Agent":    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	}
	//Make Request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}

	//Add headers to request
	for a, b := range headers {
		request.Header.Add(a, b)
	}

	//Do the request
	resp, err := client.Do(request)
	if err != nil {
		return false
	}
	if len(resp.Header.Get("access-control-allow-origin")) == 0 || resp.Header.Get("Access-Control-Allow-Credentials") == "false" {
		return false
	}
	return true
}

func worker(cha chan string, client *http.Client, wg *sync.WaitGroup) {
	for i := range cha {
		checkcors(i, client)
	}
	wg.Done()
}
