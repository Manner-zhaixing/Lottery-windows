package main

import "net/http"

func main() {
	for i := 0; i < 1000000; i++ {
		http.Get("http://127.0.0.1:5678")
	}
}
