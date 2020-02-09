package models

import (
    "database/sql"
)

func initDB(endpoint string) {
    db, err := sql.Open("postgres", endpoint)
    if err != nil {
          fmt.Println("failed to talk to db %v", err)
    }
}
