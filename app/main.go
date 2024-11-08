package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"

    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"
)

const (
    healthStatusPass = "pass"
)

// LogEntry represents the structure of the log data received via POST requests.
type LogEntry struct {
    Service   string `json:"service"`
    Endpoint  string `json:"endpoint"`
    Error     string `json:"error"`
    Traceback string `json:"traceback"`
}

var (
    writeAPI api.WriteAPIBlocking
)

func main() {
    // Load InfluxDB configurations from environment variables.
    influxURL := os.Getenv("INFLUXDB_URL")
    token := os.Getenv("INFLUXDB_TOKEN")
    org := os.Getenv("INFLUXDB_ORG")
    bucket := os.Getenv("INFLUXDB_BUCKET")

    // Check for required environment variables and provide detailed error messages.
    missingVars := []string{}
    if influxURL == "" {
        missingVars = append(missingVars, "INFLUXDB_URL")
    }
    if token == "" {
        missingVars = append(missingVars, "INFLUXDB_TOKEN")
    }
    if org == "" {
        missingVars = append(missingVars, "INFLUXDB_ORG")
    }
    if bucket == "" {
        missingVars = append(missingVars, "INFLUXDB_BUCKET")
    }
    if len(missingVars) > 0 {
        log.Fatalf("Missing required environment variables: %v", missingVars)
    }

    // Create InfluxDB client and check its health.
    client := influxdb2.NewClient(influxURL, token)
    defer client.Close()

    // Check InfluxDB client health.
    health, err := client.Health(context.Background())
    if err != nil || health.Status != healthStatusPass {
        log.Fatalf("InfluxDB is not healthy: %v", err)
    }
    log.Println("Connected to InfluxDB successfully.")

    // Initialize the Write API.
    writeAPI = client.WriteAPIBlocking(org, bucket)

    // Register HTTP handlers.
    http.HandleFunc("/health", healthHandler)
    http.HandleFunc("/", logEntryHandler)

    // Start the HTTP server with graceful shutdown.
    srv := &http.Server{
        Addr: ":80",
    }

    // Run the server in a goroutine.
    go func() {
        log.Println("Server is running on port 80...")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed to start: %v", err)
        }
    }()

    // Wait for interrupt signal to gracefully shut down the server.
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    <-quit
    log.Println("Shutting down server...")

    // Context with timeout for server shutdown.
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Println("Server exited gracefully.")
}

// healthHandler responds to health check requests.
func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

// logEntryHandler handles incoming log entries via POST requests.
func logEntryHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
        return
    }

    // Limit the size of the request body to prevent abuse.
    r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1 MB

    var entry LogEntry
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&entry); err != nil {
        http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
        log.Printf("JSON decoding error: %v", err)
        return
    }

    // Validate required fields.
    if entry.Service == "" || entry.Endpoint == "" {
        http.Error(w, "Missing required fields: service and endpoint", http.StatusBadRequest)
        return
    }

    // Create InfluxDB point.
    p := influxdb2.NewPointWithMeasurement("logs").
        AddTag("service", entry.Service).
        AddTag("endpoint", entry.Endpoint).
        AddField("error", entry.Error).
        AddField("traceback", entry.Traceback).
        SetTime(time.Now())

    // Write to InfluxDB with retry logic.
    const maxRetries = 3
    var err error
    for i := 0; i < maxRetries; i++ {
        err = writeAPI.WritePoint(context.Background(), p)
        if err == nil {
            break
        }
        log.Printf("InfluxDB write attempt %d failed: %v", i+1, err)
        time.Sleep(time.Second * time.Duration(i+1))
    }
    if err != nil {
        http.Error(w, "Failed to write to InfluxDB", http.StatusInternalServerError)
        log.Printf("InfluxDB write error after retries: %v", err)
        return
    }

    w.WriteHeader(http.StatusNoContent)
    log.Println("Data inserted into InfluxDB successfully.")
}
