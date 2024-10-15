package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const version = "0.0.1"

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Tool{})
	// Create a new router & API
	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig("Laridae", version))

	huma.Get(api, "/list/tools", func(ctx context.Context, input *struct{}) (*ListToolOutput, error) {
		resp := ListTools(*db)
		out := ListToolOutput{Body: resp}

		return &out, nil
	})

	huma.Post(api, "/ingest/tools", func(ctx context.Context, input *IngestToolInput) (*IngestToolOutput, error) {
		req, err := json.Marshal(input)
		if err != nil {
			log.Println(err)
		}
		log.Println("request body: ", string(req))
		log.Printf("Starting Ingest for URL" + input.Body.URL)
		resp, err := IngestTool(*db, input.Body)

		out := IngestToolOutput{Body: resp}

		return &out, err
	})

	huma.Get(api, "/search/tools/{term}", func(ctx context.Context, input *SearchToolsInput) (*ListToolOutput, error) {
		resp := SearchTools(*db, *input)
		out := ListToolOutput{Body: resp}

		return &out, nil
	})

	// Start the server!
	http.ListenAndServe("127.0.0.1:8888", router)
}
