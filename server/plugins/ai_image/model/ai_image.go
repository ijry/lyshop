package model

import (
	"encoding/json"

	"github.com/ijry/lyshop/model"
)

// AiModel stores configuration for each AI image generation model.
type AiModel struct {
	model.Base
	Name      string          `gorm:"size:64;not null"  json:"name"`
	Driver    string          `gorm:"size:32;not null"  json:"driver"` // tongyi|wenxin|hunyuan|openai
	Endpoint  string          `gorm:"size:255"          json:"endpoint"`
	ApiKey    string          `gorm:"size:500"          json:"api_key"`
	Params    json.RawMessage `gorm:"type:json"         json:"params"`
	IsDefault int8            `gorm:"not null;default:0" json:"is_default"`
	Status    int8            `gorm:"not null;default:1" json:"status"`
	SupportsRefImage int8     `gorm:"not null;default:0" json:"supports_ref_image"`
}

const (
	TaskStatusGenerating int8 = 1
	TaskStatusDone       int8 = 2
	TaskStatusFailed     int8 = 3
)

// AiImageTask records a generation request and its result.
type AiImageTask struct {
	model.Base
	ModelID    uint64          `gorm:"not null"          json:"model_id"`
	Scene      string          `gorm:"size:32;not null"  json:"scene"` // carousel|detail
	BizType    string          `gorm:"size:32;not null;default:'detail'" json:"biz_type"` // cover|gallery|detail|intro
	TargetProductID uint64     `gorm:"not null;default:0;index" json:"target_product_id"`
	RefImageURL string         `gorm:"size:500"          json:"ref_image_url"`
	Prompt     string          `gorm:"type:text;not null" json:"prompt"`
	NegPrompt  string          `gorm:"type:text"         json:"neg_prompt"`
	Params     json.RawMessage `gorm:"type:json"         json:"params"`
	Status     int8            `gorm:"not null"          json:"status"`
	ResultURLs json.RawMessage `gorm:"type:json"         json:"result_urls"`
	ErrorMsg   string          `gorm:"size:255"          json:"error_msg"`
}
