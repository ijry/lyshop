package driver

import (
	"context"
	"fmt"

	aidriver "github.com/ijry/lyshop/core/driver/ai"
)

// TongyiDriver calls Alibaba DashScope image generation API.
type TongyiDriver struct{ APIKey string }

func (d *TongyiDriver) Name() string { return "tongyi" }
func (d *TongyiDriver) Generate(ctx context.Context, p *aidriver.GenerateParams) (*aidriver.GenerateResult, error) {
	if d.APIKey == "" {
		return nil, fmt.Errorf("通义万象未配置 API Key")
	}
	// Production: POST https://dashscope.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis
	// with Authorization: Bearer {apiKey}
	return &aidriver.GenerateResult{
		URLs: []string{fmt.Sprintf("https://placeholder.lyshop.dev/ai/tongyi/%s.jpg", p.Style)},
	}, nil
}

// WenxinDriver calls Baidu Wenxin Yige API.
type WenxinDriver struct{ APIKey, SecretKey string }

func (d *WenxinDriver) Name() string { return "wenxin" }
func (d *WenxinDriver) Generate(ctx context.Context, p *aidriver.GenerateParams) (*aidriver.GenerateResult, error) {
	if d.APIKey == "" {
		return nil, fmt.Errorf("文心一格未配置 API Key")
	}
	return &aidriver.GenerateResult{
		URLs: []string{fmt.Sprintf("https://placeholder.lyshop.dev/ai/wenxin/%s.jpg", p.Style)},
	}, nil
}

// HunyuanDriver calls Tencent Hunyuan image API.
type HunyuanDriver struct{ SecretID, SecretKey string }

func (d *HunyuanDriver) Name() string { return "hunyuan" }
func (d *HunyuanDriver) Generate(ctx context.Context, p *aidriver.GenerateParams) (*aidriver.GenerateResult, error) {
	if d.SecretID == "" {
		return nil, fmt.Errorf("腾讯混元未配置 SecretID")
	}
	return &aidriver.GenerateResult{
		URLs: []string{fmt.Sprintf("https://placeholder.lyshop.dev/ai/hunyuan/%s.jpg", p.Style)},
	}, nil
}

// OpenAIDriver calls OpenAI-compatible image generation API (DALL-E or custom).
type OpenAIDriver struct{ APIKey, Endpoint, Model string }

func (d *OpenAIDriver) Name() string { return "openai" }
func (d *OpenAIDriver) Generate(ctx context.Context, p *aidriver.GenerateParams) (*aidriver.GenerateResult, error) {
	if d.APIKey == "" {
		return nil, fmt.Errorf("OpenAI兼容接口未配置 API Key")
	}
	endpoint := d.Endpoint
	if endpoint == "" {
		endpoint = "https://api.openai.com"
	}
	model := d.Model
	if model == "" {
		model = "dall-e-3"
	}
	// Production: POST {endpoint}/v1/images/generations
	// body: {"model": model, "prompt": p.Prompt, "n": p.Count, "size": "WxH"}
	_ = endpoint
	return &aidriver.GenerateResult{
		URLs: []string{fmt.Sprintf("https://placeholder.lyshop.dev/ai/openai/%s.jpg", p.Style)},
	}, nil
}
