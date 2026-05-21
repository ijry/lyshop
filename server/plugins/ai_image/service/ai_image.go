package service

import (
	"context"
	"encoding/json"

	aidriver "github.com/ijry/lyshop/core/driver/ai"
	"github.com/ijry/lyshop/core/db"
	aimodel "github.com/ijry/lyshop/plugins/ai_image/model"
)

// Generate creates a task record, then asynchronously calls the AI driver.
func Generate(ctx context.Context, modelID uint64, scene, prompt, negPrompt string, params map[string]any) (*aimodel.AiImageTask, error) {
	paramsJSON, _ := json.Marshal(params)
	task := &aimodel.AiImageTask{
		ModelID:   modelID,
		Scene:     scene,
		Prompt:    prompt,
		NegPrompt: negPrompt,
		Params:    paramsJSON,
		Status:    aimodel.TaskStatusGenerating,
	}
	if err := db.DB.WithContext(ctx).Create(task).Error; err != nil {
		return nil, err
	}

	// Async generation
	go func() {
		var m aimodel.AiModel
		if err := db.DB.First(&m, modelID).Error; err != nil {
			db.DB.Model(task).Updates(map[string]any{
				"status": aimodel.TaskStatusFailed, "error_msg": "model not found",
			})
			return
		}

		width, _ := params["width"].(float64)
		height, _ := params["height"].(float64)
		count, _ := params["count"].(float64)
		style, _ := params["style"].(string)
		if width == 0 { width = 750 }
		if height == 0 { height = 750 }
		if count == 0 { count = 3 }

		d, err := aidriver.Get(m.Driver)
		if err != nil {
			db.DB.Model(task).Updates(map[string]any{
				"status": aimodel.TaskStatusFailed, "error_msg": err.Error(),
			})
			return
		}

		result, err := d.Generate(context.Background(), &aidriver.GenerateParams{
			Prompt: prompt, NegPrompt: negPrompt,
			Width: int(width), Height: int(height),
			Count: int(count), Style: style,
		})
		if err != nil {
			db.DB.Model(task).Updates(map[string]any{
				"status": aimodel.TaskStatusFailed, "error_msg": err.Error(),
			})
			return
		}

		urlsJSON, _ := json.Marshal(result.URLs)
		db.DB.Model(task).Updates(map[string]any{
			"status": aimodel.TaskStatusDone, "result_urls": urlsJSON,
		})
	}()

	return task, nil
}

// GetTask returns a task by ID.
func GetTask(ctx context.Context, id uint64) (*aimodel.AiImageTask, error) {
	var task aimodel.AiImageTask
	err := db.DB.WithContext(ctx).First(&task, id).Error
	return &task, err
}

// ListTasks returns paginated generation history.
func ListTasks(ctx context.Context, page, size int) ([]aimodel.AiImageTask, int64, error) {
	if page <= 0 { page = 1 }
	if size <= 0 || size > 100 { size = 20 }
	var total int64
	db.DB.WithContext(ctx).Model(&aimodel.AiImageTask{}).Count(&total)
	var list []aimodel.AiImageTask
	err := db.DB.WithContext(ctx).Order("id desc").Offset((page-1)*size).Limit(size).Find(&list).Error
	return list, total, err
}

// ListModels returns all configured AI models.
func ListModels(ctx context.Context) ([]aimodel.AiModel, error) {
	var list []aimodel.AiModel
	err := db.DB.WithContext(ctx).Where("status = 1").Find(&list).Error
	return list, err
}

// CreateModel saves a new AI model configuration.
func CreateModel(ctx context.Context, m *aimodel.AiModel) error {
	return db.DB.WithContext(ctx).Create(m).Error
}
