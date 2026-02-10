package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// Register HTTP handlers for different endpoints
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/trigger", handleTrigger)
	http.HandleFunc("/crash", handleCrash)
	log.Println("[STARTUP] OOMKill Simulator server starting on :9090")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatalf("[ERROR] Server failed to start: %v", err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Read environment value for crash control
	envValue := os.Getenv("SHOULD_CRASH")
	log.Printf("[INFO] Health check received, SHOULD_CRASH=%s", envValue)

	io.WriteString(w, "OOMKill Simulator\n")

	//	if envValue == "YES" {
	//		go os.Exit(1)
	//	}
}

// handleTrigger initiates the OOM simulation in a separate goroutine
func handleTrigger(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[WARN] Invalid method received: %s, expected POST", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Println("[ACTION] Starting OOM simulation in background")
	go simulateOOMKill()
	io.WriteString(w, "OOMKill simulation triggered\n")
}

// handleCrash forces an immediate server exit
func handleCrash(w http.ResponseWriter, r *http.Request) {
	log.Println("[CRITICAL] Crash triggered, server exiting with code 1")
	go os.Exit(1)
	io.WriteString(w, "Crash triggered\n")
}

// simulateOOMKill continuously allocates memory until the system runs out
func simulateOOMKill() {
	log.Println("[ACTION] OOMKill simulation started - allocating memory in 500MB chunks")
	var memory [][]byte
	chunkSize := 500 * 1024 * 1024 // 500MB chunks
	for {
		memory = append(memory, make([]byte, chunkSize))
		allocatedMB := len(memory) * 500
		log.Printf("[MEMORY] Allocated: %d MB, Total chunks: %d", allocatedMB, len(memory))
		// Fill the allocated memory with non-zero values to prevent compiler optimization
		for i := range memory[len(memory)-1] {
			memory[len(memory)-1][i] = byte(i % 256)
		}
	}
}
