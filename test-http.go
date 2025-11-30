package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test-http.go <host>")
		os.Exit(1)
	}

	host := os.Args[1]
	url := fmt.Sprintf("http://%s/post", host)

	fmt.Printf("Testing HTTP GET to %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("SUCCESS: Status %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(body))
}
