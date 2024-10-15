package main

import (
	"io"
	"log"
	"net/http"
	"slices"

	"github.com/pb33f/libopenapi"

	"github.com/ghodss/yaml"
	"gorm.io/gorm"
)

type ListToolResponse struct {
	Tools []Tool `json:"tools"`
}

type ListToolOutput struct {
	Body ListToolResponse
}

func ListTools(db gorm.DB) ListToolResponse {
	var tools []Tool
	db.Find(&tools)

	return ListToolResponse{Tools: tools}
}

type PathFilter struct {
	Type  string   `json:"type,omitempty" enum:"only,except"`
	Value []string `json:"value,omitempty"`
}

type IngestToolRequest struct {
	URL     string            `json:"url"`
	Filter  PathFilter        `json:"filter,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

type IngestToolResponse struct {
	ToolCount   int `json:"ingested_tools,omitempty"`
	ServerCount int `json:"ingested_servers,omitempty"`
}

type IngestToolInput struct {
	Body IngestToolRequest
}

type IngestToolOutput struct {
	Body IngestToolResponse
}

func IngestTool(db gorm.DB, req IngestToolRequest) (IngestToolResponse, error) {

	request, err := http.Get(req.URL)
	if err != nil {
		log.Println(err.Error())
		return IngestToolResponse{}, err
	}
	defer request.Body.Close()
	log.Printf("Grabbed URL")

	spec, err := io.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		return IngestToolResponse{}, err
	}
	log.Printf("Read Data")

	document, err := libopenapi.NewDocument(spec)
	if err != nil {
		log.Println(err.Error())
		return IngestToolResponse{}, err
	}
	log.Printf("Parsed Document")

	model, errs := document.BuildV3Model()
	if errs != nil {
		for _, e := range errs {
			log.Println(e.Error())
		}
		return IngestToolResponse{}, err
	}
	log.Printf("Built Model")

	items := model.Model.Paths.PathItems
	paths := items.KeysFromNewest()

	toolCount := 0
	serverCount := 0
	for _, server := range model.Model.Servers {
		server_ := Server{
			Description: server.Description,
			URL:         server.URL,
		}
		db.Create(&server_)
		serverCount += 1

		for path := range paths {
			pathItem, _ := items.Get(path)
			strPath, _ := pathItem.Render()
			json_path, _ := yaml.YAMLToJSON(strPath)
			tool := Tool{Endpoint: path, Schema: json_path, Server: server_}

			if req.Filter.Type == "only" && slices.Contains(req.Filter.Value, path) {
				db.Create(&tool)
				toolCount += 1
			}

			if req.Filter.Type == "except" && !slices.Contains(req.Filter.Value, path) {
				db.Create(&tool)
				toolCount += 1
			}

			if req.Filter.Type == "" {
				db.Create(&tool)
				toolCount += 1
			}
		}

	}

	return IngestToolResponse{ToolCount: toolCount, ServerCount: serverCount}, err
}

type SearchToolsInput struct {
	Term string `path:"term"`
}

func SearchTools(db gorm.DB, req SearchToolsInput) ListToolResponse {
	tools := []Tool{}
	db.Table("tools").Where("endpoint LIKE ?", "%"+req.Term+"%").Scan(&tools)

	return ListToolResponse{Tools: tools}
}
