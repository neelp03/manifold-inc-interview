package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type LogEntry struct {
    Service   string `json:"service"`
    Endpoint  string `json:"endpoint"`
    Error     string `json:"error"`
    Traceback string `json:"traceback"`
}

func main() {
    // InfluxDB configurations
    influxURL := os.Getenv("INFLUXDB_URL")
    token := os.Getenv("INFLUXDB_TOKEN")
    org := os.Getenv("INFLUXDB_ORG")
    bucket := os.Getenv("INFLUXDB_BUCKET")

    // Check for required environment variables
    if influxURL == "" || token == "" || org == "" || bucket == "" {
        log.Fatal("InfluxDB configuration not set in environment variables")
    }

    // Create InfluxDB client
    client := influxdb2.NewClient(influxURL, token)
    defer client.Close()
    writeAPI := client.WriteAPIBlocking(org, bucket)
    
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    
    // HTTP handler
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
            return
        }

        var entry LogEntry
        decoder := json.NewDecoder(r.Body)
        if err := decoder.Decode(&entry); err != nil {
            http.Error(w, "Bad request", http.StatusBadRequest)
            log.Printf("JSON decoding error: %v", err)
            return
        }

        // Create InfluxDB point
        p := influxdb2.NewPointWithMeasurement("logs").
            AddTag("service", entry.Service).
            AddTag("endpoint", entry.Endpoint).
            AddField("error", entry.Error).
            AddField("traceback", entry.Traceback).
            SetTime(time.Now())

        // Write to InfluxDB
        if err := writeAPI.WritePoint(context.Background(), p); err != nil {
            http.Error(w, fmt.Sprintf("Failed to write to InfluxDB: %v", err), http.StatusInternalServerError)
            log.Printf("InfluxDB write error: %v", err)
            return
        }

        w.WriteHeader(http.StatusNoContent)
        fmt.Println("Data inserted into InfluxDB")
    })

    // Start server on port 80
    fmt.Println("Server is running on port 80...")
    if err := http.ListenAndServe(":80", nil); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
