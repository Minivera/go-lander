package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Printf("Server serving content at http://localhost:8080")
	err := http.ListenAndServe(":8080", http.FileServer(http.Dir("./example")))
	if err != nil {
		panic(err)
	}
}
