package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	intentional build failure
	http.HandleFunc("/", handleRequest)
	fmt.Println("Server is running...")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("message")
	log.Printf("Received message: %s", param)
	if param == "" {
		fmt.Fprintf(w, "No message provided. Use ?message=YourMessage in the URL.")
		return
	}
	fmt.Fprintf(w, "Received message: %s", param)
}
