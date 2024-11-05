package main

import (
    "bytes"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/brianvoe/gofakeit/v7"
)

type LogEntry struct {
    Service   string `json:"service"`
    Endpoint  string `json:"endpoint"`
    Error     string `json:"error"`
    Traceback string `json:"traceback"`
}

func main() {
    numEntries := flag.Int("n", 100, "Number of log entries to generate")
    serverURL := flag.String("url", "http://app:80", "Server URL")
    flag.Parse()

    gofakeit.Seed(time.Now().UnixNano())
    var wg sync.WaitGroup

    for i := 0; i < *numEntries; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            entry := generateRandomLogEntry()
            err := sendLogEntry(*serverURL, entry)
            if err != nil {
                log.Printf("Error sending log entry: %v", err)
            } else {
                log.Printf("Log entry sent successfully")
            }
        }()
    }
    wg.Wait()
}

func generateRandomLogEntry() LogEntry {
    services := []string{"user-service", "auth-service", "payment-service"}
    endpoints := []string{"/api/users", "/api/login", "/api/payments"}

    errorMessage := gofakeit.HackerPhrase()
    filePath := fmt.Sprintf("/%s/%s.%s", gofakeit.Word(), gofakeit.Word(), gofakeit.FileExtension())
    traceback := fmt.Sprintf(`File "%s", line %d, in %s\n    %s`,
        filePath, gofakeit.Number(10, 100), gofakeit.BuzzWord(), errorMessage)

    return LogEntry{
        Service:   gofakeit.RandomString(services),
        Endpoint:  gofakeit.RandomString(endpoints),
        Error:     errorMessage,
        Traceback: traceback,
    }
}

func sendLogEntry(url string, entry LogEntry) error {
    data, err := json.Marshal(entry)
    if err != nil {
        return fmt.Errorf("failed to marshal log entry: %w", err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
    if err != nil {
        return fmt.Errorf("failed to create HTTP request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send HTTP request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent {
        return fmt.Errorf("received unexpected status code: %d", resp.StatusCode)
    }

    return nil
}
