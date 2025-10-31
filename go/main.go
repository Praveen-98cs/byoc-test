package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/trigger", handleTrigger)
	http.HandleFunc("/crash", handleCrash)
	log.Printf("Server starting ...")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OOMKill Simulator\n")
	io.WriteString(w, "Send a POST request to /trigger to simulate OOMKill\n")
}

func handleTrigger(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	go simulateOOMKill()
	io.WriteString(w, "OOMKill simulation triggered\n")
}

func handleCrash(w http.ResponseWriter, r *http.Request) {
	log.Println("crashing server...")
	os.Exit(1)
}

func simulateOOMKill() {
	log.Println("Starting OOMKill simulation")
	log.Println("Initializing memory allocation process...")
	var memory [][]byte
	chunkSize := 500 * 1024 * 1024 // 100MB chunks
	for {
		log.Printf("Allocating memory chunk %d (%d MB)...", len(memory)+1, chunkSize/(1024*1024))
		memory = append(memory, make([]byte, chunkSize))
		allocatedMB := len(memory) * 500
		log.Printf("Allocated memory: %d MB", allocatedMB)
		// Fill the allocated memory with non-zero values
		for i := range memory[len(memory)-1] {
			memory[len(memory)-1][i] = byte(i % 256)
		}
		log.Printf("Memory chunk %d filled with data", len(memory))
	}
}
