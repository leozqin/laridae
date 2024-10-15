package main

import (
	"context"
	"encoding/json"
	"fmt"
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

// GreetingOutput represents the greeting operation response.
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

// GreetingInput represents the greeting operation request.
type GreetingInput struct {
	Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
}

func Greeting(input GreetingInput) GreetingOutput {
	output := &GreetingOutput{}
	output.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)

	return *output
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Tool{})
	// Create a new router & API
	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig("Laridae", version))

	huma.Get(api, "/greeting/{name}", func(ctx context.Context, input *GreetingInput) (*GreetingOutput, error) {
		resp := Greeting(*input)

		return &resp, nil
	})

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
