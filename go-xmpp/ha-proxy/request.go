package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func request(url string, count int) {
	for i := 0; i < 10; i++ {
		go func() {
			for {
				resp, err := http.Get(url)

				if err != nil {
					log.Fatal(err)
				}

				body, err := ioutil.ReadAll(resp.Body)

				fmt.Printf("Response: %s, Count: %d\n", string(body), count)
				count++
				resp.Body.Close()
			}
		}()
	}
}

func main() {
	done := make(chan bool, 1)
	count := 0

	request("http://lucas.com.br:8080", count)
	//request("http://charabanaia.com.br:8080", count)

	<-done
}
