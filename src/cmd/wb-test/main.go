package main

import (
	"bufio"
	"log"
	"os"
	"sync"
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
)

func init() {
	log.SetFlags(0)
}

// make http.Get and count appearance of a substring "Go" in resp.Body
func count(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return strings.Count(string(b), "Go")
}

func worker(urls <-chan string, nums chan<- int, wg *sync.WaitGroup) {
	for u := range urls {
		n := count(u)
		nums <- n
		fmt.Println("Count for " + u + ":", n)
		wg.Done()
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	urls := make(chan string, 10)
	nums := make(chan int, 10)

	wg := sync.WaitGroup{}

	// creating five workers
	for k := 0; k < 5; k++ {
		go worker(urls, nums, &wg)
	}

	for scanner.Scan() {
		wg.Add(1)
		urls <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	// close channel urls to stop range in workers
	close(urls)

	total := 0
	wg1 := sync.WaitGroup{}
	go func(wg *sync.WaitGroup) {
		wg.Add(1)
		for r := range nums {
			total += r
		}
		wg.Done()
	}(&wg1)

	// sleep in main goroutine to wait executing all tasks
	wg.Wait()

	// close channel nums to stop range
	close(nums)
	wg1.Wait()

	log.Printf("Total: %v", total)
}