package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type TodoItem struct {
	Id          int        `json:"id"`
	Title       string     `json:"title`
	Description string     `json:"description"`
	Status      string     `json:"status`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func main() {
	now := time.Now().UTC()
	item := TodoItem{
		Id:          1,
		Title:       "Task 1",
		Description: "Content 1",
		Status:      "Doing",
		CreatedAt:   &now,
		UpdatedAt:   nil,
	}

	if jsData, err := json.Marshal(item); err != nil {
		log.Fatalln(err)
	} else {
		log.Println(string(jsData))
	}

	jsStr := "{\"id\":1,\"Title\":\"Task 1\"}"

	var item2 TodoItem
	if err := json.Unmarshal([]byte(jsStr), &item2); err != nil {
		log.Fatalln(err)
	}
	log.Println(item2)
	//////////////////////////////////////////////////
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		items := v1.Group("/items")
	}
	if err := r.Run(":3000"); err != nil {
		log.Fatalln(err)
	}
}
