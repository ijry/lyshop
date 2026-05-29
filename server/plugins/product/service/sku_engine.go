package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	productmodel "github.com/ijry/lyshop/plugins/product/model"
)

const (
	defaultSkuKey       = "__default__"
	maxSkuCartesianSize = 300
)

type SpecSchemaGroup struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type SkuOverride struct {
	SkuKey  string   `json:"sku_key"`
	SkuCode string   `json:"sku_code"`
	Price   *float64 `json:"price"`
	Stock   *int     `json:"stock"`
}

type SkuDiffSummary struct {
	Added       int `json:"added"`
	Kept        int `json:"kept"`
	Inactivated int `json:"inactivated"`
}

func DecodeSkuAttrs(raw json.RawMessage) ([]productmodel.SkuAttr, error) {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return []productmodel.SkuAttr{}, nil
	}
	payload := []byte(trimmed)
	if strings.HasPrefix(trimmed, "\"") {
		var encoded string
		if err := json.Unmarshal(payload, &encoded); err != nil {
			return nil, fmt.Errorf("解析SKU属性字符串失败: %w", err)
		}
		payload = []byte(encoded)
	}
	var attrs []productmodel.SkuAttr
	if err := json.Unmarshal(payload, &attrs); err != nil {
		return nil, fmt.Errorf("解析SKU属性失败: %w", err)
	}
	return NormalizeSkuAttrs(attrs), nil
}

func EncodeSkuAttrs(attrs []productmodel.SkuAttr) (json.RawMessage, error) {
	normalized := NormalizeSkuAttrs(attrs)
	raw, err := json.Marshal(normalized)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func NormalizeSkuAttrs(attrs []productmodel.SkuAttr) []productmodel.SkuAttr {
	out := make([]productmodel.SkuAttr, 0, len(attrs))
	for _, attr := range attrs {
		name := strings.TrimSpace(attr.Name)
		value := strings.TrimSpace(attr.Value)
		if name == "" || value == "" {
			continue
		}
		out = append(out, productmodel.SkuAttr{Name: name, Value: value})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name == out[j].Name {
			return out[i].Value < out[j].Value
		}
		return out[i].Name < out[j].Name
	})
	return out
}

func CanonicalSkuKey(attrs []productmodel.SkuAttr) string {
	normalized := NormalizeSkuAttrs(attrs)
	if len(normalized) == 0 {
		return defaultSkuKey
	}
	parts := make([]string, 0, len(normalized))
	for _, attr := range normalized {
		parts = append(parts, attr.Name+":"+attr.Value)
	}
	return strings.Join(parts, "|")
}

func BuildSkusFromSpecSchema(schema []SpecSchemaGroup, basePrice float64, overrides []SkuOverride) ([]productmodel.ProductSku, error) {
	normalized := normalizeSpecSchema(schema)
	if len(normalized) == 0 {
		return []productmodel.ProductSku{
			{
				Price:  basePrice,
				Stock:  0,
				SkuKey: defaultSkuKey,
				Status: productmodel.ProductSkuStatusActive,
				Attrs:  json.RawMessage("[]"),
			},
		}, nil
	}
	total := 1
	for _, group := range normalized {
		total *= len(group.Values)
	}
	if total > maxSkuCartesianSize {
		return nil, fmt.Errorf("SKU组合数量过大: %d，超过上限%d", total, maxSkuCartesianSize)
	}

	overridesByKey := make(map[string]SkuOverride, len(overrides))
	for _, row := range overrides {
		key := strings.TrimSpace(row.SkuKey)
		if key == "" {
			continue
		}
		overridesByKey[key] = row
	}

	combos := make([][]productmodel.SkuAttr, 0, total)
	var walk func(idx int, prefix []productmodel.SkuAttr)
	walk = func(idx int, prefix []productmodel.SkuAttr) {
		if idx >= len(normalized) {
			attrs := make([]productmodel.SkuAttr, len(prefix))
			copy(attrs, prefix)
			combos = append(combos, attrs)
			return
		}
		group := normalized[idx]
		for _, value := range group.Values {
			next := append(prefix, productmodel.SkuAttr{Name: group.Name, Value: value})
			walk(idx+1, next)
		}
	}
	walk(0, make([]productmodel.SkuAttr, 0, len(normalized)))

	rows := make([]productmodel.ProductSku, 0, len(combos))
	for _, attrs := range combos {
		key := CanonicalSkuKey(attrs)
		raw, err := EncodeSkuAttrs(attrs)
		if err != nil {
			return nil, err
		}
		row := productmodel.ProductSku{
			Attrs:   raw,
			Price:   basePrice,
			Stock:   0,
			SkuKey:  key,
			Status:  productmodel.ProductSkuStatusActive,
			SkuCode: "",
		}
		if override, ok := overridesByKey[key]; ok {
			row.SkuCode = strings.TrimSpace(override.SkuCode)
			if override.Price != nil {
				row.Price = *override.Price
			}
			if override.Stock != nil {
				row.Stock = *override.Stock
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func normalizeSpecSchema(schema []SpecSchemaGroup) []SpecSchemaGroup {
	out := make([]SpecSchemaGroup, 0, len(schema))
	for _, group := range schema {
		name := strings.TrimSpace(group.Name)
		if name == "" {
			continue
		}
		dedup := make(map[string]struct{}, len(group.Values))
		values := make([]string, 0, len(group.Values))
		for _, raw := range group.Values {
			value := strings.TrimSpace(raw)
			if value == "" {
				continue
			}
			if _, ok := dedup[value]; ok {
				continue
			}
			dedup[value] = struct{}{}
			values = append(values, value)
		}
		if len(values) == 0 {
			continue
		}
		sort.Strings(values)
		out = append(out, SpecSchemaGroup{Name: name, Values: values})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func normalizeIncomingSkus(skus []productmodel.ProductSku) ([]productmodel.ProductSku, error) {
	rows := make([]productmodel.ProductSku, 0, len(skus))
	seen := make(map[string]struct{}, len(skus))
	for _, row := range skus {
		attrs, err := DecodeSkuAttrs(row.Attrs)
		if err != nil {
			return nil, err
		}
		key := CanonicalSkuKey(attrs)
		if _, ok := seen[key]; ok {
			return nil, errors.New("SKU属性组合重复")
		}
		seen[key] = struct{}{}
		raw, err := EncodeSkuAttrs(attrs)
		if err != nil {
			return nil, err
		}
		normalized := row
		normalized.Attrs = raw
		normalized.SkuKey = key
		normalized.SkuCode = strings.TrimSpace(row.SkuCode)
		normalized.Status = productmodel.ProductSkuStatusActive
		rows = append(rows, normalized)
	}
	return rows, nil
}
