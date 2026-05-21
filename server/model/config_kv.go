package model

// ConfigKV stores plugin-namespaced key-value configuration.
type ConfigKV struct {
	ID     uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Plugin string `gorm:"size:64;not null;uniqueIndex:uk_plugin_key"  json:"plugin"`
	Key    string `gorm:"size:128;not null;uniqueIndex:uk_plugin_key" json:"key"`
	Value  string `gorm:"type:text"                                  json:"value"`
}

func (ConfigKV) TableName() string { return "configs" }
