package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

var transport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: time.Second,
		DualStack: true,
	}).DialContext,
}

var client = &http.Client{
	Transport: transport,
}

func main() {
	var threads int
	var cookies string
	var wg sync.WaitGroup
	urls := make(chan string)

	flag.IntVar(&threads, "t", 20, "Specify number of threads to run")
	flag.StringVar(&cookies, "c", "", "Specfy Cookie Header")
	flag.Parse()

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
		go worker(urls, client, &wg, cookies)
	}
	wg.Wait()
}

func checkcors(s string, client *http.Client, cookie string) {
	if !checkheaders(s, s, client, cookie) {
		return
	}
	if checkheaders(s, "https://notexisting.com", client, cookie) {
		fmt.Println(s, "is vulnerable to cors")
		return
	}
}

func checkheaders(url, origin string, client *http.Client, cookie string) bool {
	//Initialize Headers
	headers := map[string]string{"Origin": origin,
		"Cache-Control": "no-cache",
		"User-Agent":    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
		"Cookie":        cookie,
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
	defer resp.Body.Close()
	if len(resp.Header.Get("access-control-allow-origin")) == 0 || resp.Header.Get("Access-Control-Allow-Credentials") == "false" {
		return false
	}
	return true
}

func worker(cha chan string, client *http.Client, wg *sync.WaitGroup, cookie string) {
	for i := range cha {
		checkcors(i, client, cookie)
	}
	wg.Done()
}
