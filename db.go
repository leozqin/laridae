package main

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type GormModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
} //@name models.gormModel

type Header struct {
	Key   string
	Value string
}

type Tool struct {
	GormModel `json:"gorm_model,omitempty"`
	Endpoint  string         `json:"endpoint,omitempty"`
	Schema    datatypes.JSON `json:"schema,omitempty"`
	ServerId  uint           `json:"server_id"`
	Server    Server         `json:"server,omitempty"`
}

type Server struct {
	GormModel
	Description string         `json:"description,omitempty"`
	URL         string         `json:"url,omitempty"`
	Variables   datatypes.JSON `json:"variables,omitempty"`
}
