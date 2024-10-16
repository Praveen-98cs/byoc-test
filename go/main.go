package main

import (
	"io"
	"log"
	"net/http"
	"runtime"
	"time"
)

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/trigger", handleTrigger)

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

func simulateOOMKill() {
	log.Println("Starting OOMKill simulation")

	var mem [][]int
	for {
		mem = append(mem, make([]int, 1024*1024))
		log.Printf("Allocated memory: %d MB", len(mem))
		runtime.GC()
		time.Sleep(100 * time.Millisecond)
	}
}
