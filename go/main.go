package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Config holds the application configuration
type Config struct {
	Server struct {
		Port     int    `json:"port"`
		LogLevel string `json:"logLevel"`
	} `json:"server"`
	Proxy struct {
		DefaultHost           string `json:"defaultHost"`
		DefaultPath           string `json:"defaultPath"`
		RequestTimeoutSeconds int    `json:"requestTimeoutSeconds"`
	} `json:"proxy"`
	Features struct {
		EnableStatusEndpoint bool `json:"enableStatusEndpoint"`
		EnableConfigEndpoint bool `json:"enableConfigEndpoint"`
	} `json:"features"`
}

var startTime time.Time
var appConfig Config

func init() {
	startTime = time.Now()
}

// loadConfig loads configuration from file and environment variables
func loadConfig() Config {
	config := Config{}
	
	// Set default values
	config.Server.Port = 9090
	config.Server.LogLevel = "info"
	config.Proxy.DefaultHost = "http://postman-echo.com"
	config.Proxy.DefaultPath = "get?foo1=bar1&foo2=bar2"
	config.Proxy.RequestTimeoutSeconds = 30
	config.Features.EnableStatusEndpoint = true
	config.Features.EnableConfigEndpoint = true
	
	// Load from config file if CONFIG_FILE_PATH is set
	if configPath := os.Getenv("CONFIG_FILE_PATH"); configPath != "" {
		log.Printf("Loading configuration from file: %s", configPath)
		if data, err := ioutil.ReadFile(configPath); err == nil {
			if err := json.Unmarshal(data, &config); err != nil {
				log.Printf("Warning: Failed to parse config file: %v", err)
			} else {
				log.Println("Configuration loaded from file successfully")
			}
		} else {
			log.Printf("Warning: Failed to read config file: %v", err)
		}
	}
	
	// Override with environment variables (highest priority)
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Server.Port = p
			log.Printf("Server port overridden by env var: %d", p)
		}
	}
	
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.Server.LogLevel = logLevel
		log.Printf("Log level overridden by env var: %s", logLevel)
	}
	
	if host := os.Getenv("PROXY_DEFAULT_HOST"); host != "" {
		config.Proxy.DefaultHost = host
		log.Printf("Proxy default host overridden by env var: %s", host)
	}
	
	if path := os.Getenv("PROXY_DEFAULT_PATH"); path != "" {
		config.Proxy.DefaultPath = path
		log.Printf("Proxy default path overridden by env var: %s", path)
	}
	
	if timeout := os.Getenv("REQUEST_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			config.Proxy.RequestTimeoutSeconds = t
			log.Printf("Request timeout overridden by env var: %d", t)
		}
	}
	
	if enable := os.Getenv("ENABLE_STATUS_ENDPOINT"); enable != "" {
		config.Features.EnableStatusEndpoint = enable == "true"
		log.Printf("Status endpoint enabled: %v", config.Features.EnableStatusEndpoint)
	}
	
	if enable := os.Getenv("ENABLE_CONFIG_ENDPOINT"); enable != "" {
		config.Features.EnableConfigEndpoint = enable == "true"
		log.Printf("Config endpoint enabled: %v", config.Features.EnableConfigEndpoint)
	}
	
	return config
}


func main() {
	// Load configuration
	appConfig = loadConfig()
	
	log.Printf("Starting server with configuration:")
	log.Printf("  Port: %d", appConfig.Server.Port)
	log.Printf("  Log Level: %s", appConfig.Server.LogLevel)
	log.Printf("  Proxy Default Host: %s", appConfig.Proxy.DefaultHost)
	
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprintf(w, "{\"active\": true}")
	})
	http.HandleFunc("/healthz/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprintf(w, "{\"healthy\": true}")
	})
	http.HandleFunc("/hello/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello %s", req.URL.Query().Get("name"))
	})
	
	// Conditionally register status endpoint
	if appConfig.Features.EnableStatusEndpoint {
		http.HandleFunc("/status/", handleStatus)
		log.Println("Status endpoint enabled at /status/")
	}
	
	// Conditionally register config endpoint
	if appConfig.Features.EnableConfigEndpoint {
		http.HandleFunc("/config/", handleConfig)
		log.Println("Config endpoint enabled at /config/")
	}
	
	http.HandleFunc("/proxy/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			decoder := json.NewDecoder(req.Body)
			var data map[string]string
			err := decoder.Decode(&data)
			if err != nil {
				w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
			}

			host := data["host"]
			path := data["path"]

			if len(host) == 0 {
				// Use config default
				host = appConfig.Proxy.DefaultHost
			}
			if len(path) == 0 {
				// Use config default
				path = appConfig.Proxy.DefaultPath
			}

			resp, err := http.Get(fmt.Sprintf("%s/%s", strings.TrimRight(host, "/"), strings.TrimLeft(path, "/")))
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(body)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	fmt.Printf("listening on %v\n", appConfig.Server.Port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", appConfig.Server.Port), logRequest(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	uptime := time.Since(startTime)

	status := map[string]interface{}{
		"server": map[string]interface{}{
			"status":    "running",
			"uptime":    uptime.String(),
			"startTime": startTime.Format(time.RFC3339),
		},
		"memory": map[string]interface{}{
			"allocatedMB":  memStats.Alloc / 1024 / 1024,
			"totalAllocMB": memStats.TotalAlloc / 1024 / 1024,
			"sysMB":        memStats.Sys / 1024 / 1024,
			"numGC":        memStats.NumGC,
		},
		"runtime": map[string]interface{}{
			"numGoroutines": runtime.NumGoroutine(),
			"goVersion":     runtime.Version(),
			"numCPU":        runtime.NumCPU(),
		},
		"endpoints": []string{
			"/",
			"/healthz/",
			"/hello/",
			"/status/",
			"/proxy/",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Return current configuration
	configResponse := map[string]interface{}{
		"server": map[string]interface{}{
			"port":     appConfig.Server.Port,
			"logLevel": appConfig.Server.LogLevel,
		},
		"proxy": map[string]interface{}{
			"defaultHost":           appConfig.Proxy.DefaultHost,
			"defaultPath":           appConfig.Proxy.DefaultPath,
			"requestTimeoutSeconds": appConfig.Proxy.RequestTimeoutSeconds,
		},
		"features": map[string]interface{}{
			"enableStatusEndpoint": appConfig.Features.EnableStatusEndpoint,
			"enableConfigEndpoint": appConfig.Features.EnableConfigEndpoint,
		},
		"configSource": map[string]interface{}{
			"filePathEnvVar": os.Getenv("CONFIG_FILE_PATH"),
			"loadedFromFile": os.Getenv("CONFIG_FILE_PATH") != "",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(configResponse)
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
