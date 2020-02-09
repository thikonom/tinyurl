package main

import (
    "net/http"
    "encoding/json"
    "fmt"
    "log"
    "time"
    "math/rand"

    "github.com/go-redis/redis"
)

type Key string

type CreateParams struct {
        EncodedURL  string `json:"encoded_url"`
        OriginalURL string `json:"original_url"`
}

type URLResource struct {
      EncodedURL string     `json:"encoded_url"`
      CreatedAt  time.Time  `json:"created_at"`
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
  "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
  rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}

func randString(length int) string {
  return StringWithCharset(length, charset)
}

func main() {
        client := redis.NewClient(&redis.Options{
            Addr:     "redis:6379",
            Password: "",
            DB:       0,
        })

        // create a new random string(6) and save it to redis
        http.HandleFunc("/generateKey", func(w http.ResponseWriter, r *http.Request){
                  for i := 0; i < 10; i++ {
                          encoded := randString(6)
                          err := client.LPush("availableUrls", encoded).Err()
                          if err != nil {
                              panic(err)
                          }
                  }

                  content := map[string]int{"op": 1}
                  jsonData, err1 := json.Marshal(content)
                  if err1 != nil {
                          panic(err1)
                  }
                  w.Header().Set("Content-Type", "application/json")
                  fmt.Fprintf(w, string(jsonData))
        })

        // ask redis for the next available random string
        http.HandleFunc("/getKey", func(w http.ResponseWriter, r *http.Request){
                encoded, err1 := client.RPop("availableUrls").Result()
                if err1 != nil {
                      panic(err1)
                }
                url := &URLResource{EncodedURL: encoded, CreatedAt: time.Now()}
                var jsonData []byte
                jsonData, err := json.Marshal(url)
                if err != nil {
                      panic(err)
                }
                w.Header().Set("Content-Type", "application/json")
                fmt.Fprintf(w, string(jsonData))

        })

        log.Fatal(http.ListenAndServe(":8081", nil))
}
