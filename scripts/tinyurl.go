package main

import (
        "os"
        "fmt"
        "log"
        "net/http"
        "encoding/json"
        "time"
        "bytes"

        "github.com/jinzhu/gorm"
        _ "github.com/jinzhu/gorm/dialects/postgres"
        tinyurl "github.com/thikonom/tinyurl"
)

type CreateParams struct {
        Email string `json:"email"`
        OriginalURL string `json:"original_url"`
}

type CacheParams struct {
        OriginalURL string `json:"original_url"`
}

type GetParams struct {
        EncodedURL string `json:"encoded_url"`
}

type URLResource struct {
        EncodedURL  string `json:"encoded_url"`
        CreatedAt *time.Time
}

// type KGSParams struct {
//         OriginalURL string `json:"original_url"`
// }

var db *gorm.DB

func main() {
        // connect to db and load with some users
        db, err := gorm.Open("postgres", "host=postgres port=5432 user=admin password=123 dbname=tinyurl sslmode=disable")
        if err != nil {
                fmt.Println("failed to connect to database", err)
                os.Exit(1)
        }

        defer db.Close()
        db.SetLogger(log.New(os.Stdout, "\r\n", 0))

        db.Debug().DropTable(&tinyurl.TinyURL{})
        db.Debug().DropTable(&tinyurl.User{})
        fmt.Println("Automigrating...")
        db.Debug().AutoMigrate(&tinyurl.TinyURL{})
        db.Debug().AutoMigrate(&tinyurl.User{})
        db.Model(&tinyurl.TinyURL{}).AddForeignKey("user_email", "users(email)", "CASCADE", "CASCADE")

        fmt.Println("Creating users")
        db.Debug().Create(&tinyurl.User{Email: "aaaaaaaa@gmail.com", TinyURLS: []tinyurl.TinyURL{
                {OriginalURL: "http://google.com", ShortenedURL: "http://tinyurl123"},
                {OriginalURL: "http://google.com", ShortenedURL: "http://tinyurl1234"},
        }})
        db.Debug().Create(&tinyurl.User{Email: "mary@gmail.com"})
        db.Debug().Create(&tinyurl.User{Email: "mike@gmail.com"})
        var user tinyurl.User
        db.Debug().First(&user)

        users := []tinyurl.User{}
        db.Debug().Where("email=?", "aaa@gmail.com").Preload("TinyURLS").Find(&users)
        fmt.Println("Users: ", users)
        for _, v := range(users) {
                fmt.Println("TinyURLS: ", v.TinyURLS)
        }


        // register controllers and startup server
        http.HandleFunc("/createTiny", func(w http.ResponseWriter, r *http.Request){
                  decoder := json.NewDecoder(r.Body)
                  var params CreateParams
                  err := decoder.Decode(&params)
                  if err != nil {
                      fmt.Fprintf(w, err.Error())
                  } else {
                        resp, err := http.Get("http://kgs:8081/getKey")
                        if err != nil {
                            panic(err)
                        }
                        defer resp.Body.Close()

                        resource := &URLResource{}
                        err = json.NewDecoder(resp.Body).Decode(resource)
                        if err == nil {
                                if err := db.Debug().Create(&tinyurl.TinyURL{OriginalURL: params.OriginalURL,
                                        ShortenedURL: resource.EncodedURL, UserEmail: params.Email}); err !=nil {
                                          fmt.Fprintf(w, resource.EncodedURL)
                                } else {
                                        fmt.Fprintf(w, "we are fighting some aliens that came to steal our servers")
                                }
                        } else {
                                panic(err)
                        }
                }
        })

        // pass {'encoded_url: ''} and get {original_url: 'http://'}
        http.HandleFunc("/getTiny", func(w http.ResponseWriter, r *http.Request){
                decoder := json.NewDecoder(r.Body)
                var params GetParams
                err := decoder.Decode(&params)
                if err != nil {
                        fmt.Fprintf(w, err.Error())
                } else {
                        postData := map[string]string{"encoded_url": params.EncodedURL}
                        jsonData, err1 := json.Marshal(&postData)
                        if err1 != nil {
                            panic(err1)
                        }
                        resp, err := http.Post("http://cache:8082/getCacheKey", "application/json", bytes.NewBuffer(jsonData))
                        if err != nil {
                            panic(err)
                        }
                        defer resp.Body.Close()

                        decoder := json.NewDecoder(resp.Body)
                        var params CacheParams
                        err = decoder.Decode(&params)
                        if err != nil {
                                fmt.Fprintf(w, err.Error())
                        } else {
                                content := map[string]string{"original_url": params.OriginalURL}
                                jsonData, err1 := json.Marshal(content)
                                if err1 != nil {
                                        panic(err1)
                                }
                                w.Header().Set("Content-Type", "application/json")
                                fmt.Fprintf(w, string(jsonData))
                        }
                }
        })

        http.HandleFunc("/getTinys", func(w http.ResponseWriter, r *http.Request){
                  decoder := json.NewDecoder(r.Body)
                  var params CreateParams
                  err := decoder.Decode(&params)
                  if err != nil {
                      fmt.Fprintf(w, err.Error())
                  } else {
                      var user tinyurl.User
                      db.Debug().Where("email=?", params.Email).Preload("TinyURLS").Find(&user)
                      for _, v := range(user.TinyURLS){
                              fmt.Fprintf(w, v.ShortenedURL + "\n")
                      }
                  }
        })

        // http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        //         var users []tinyurl.User
        //         db.Debug().Find(&users)
        //         for _, v := range(users) {
        //                 fmt.Fprintf(w, v.Email + "\n")
        //         }
        // })

        log.Fatal(http.ListenAndServe(":8080", nil))
}
