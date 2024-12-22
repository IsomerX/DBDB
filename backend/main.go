package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Go backend!")
}

func main() {
	http.HandleFunc("/data", handler)
	http.ListenAndServe(":5000", nil)
}
