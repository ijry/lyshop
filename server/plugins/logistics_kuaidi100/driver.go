package logistics_kuaidi100

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	logisticsDriver "github.com/ijry/lyshop/core/driver/logistics"
)

type kuaidi100Driver struct {
	APIURL   string
	Customer string
	Key      string
}

func (d *kuaidi100Driver) Name() string { return "kuaidi100" }

func (d *kuaidi100Driver) Query(ctx context.Context, req logisticsDriver.QueryReq) (*logisticsDriver.TrackResult, error) {
	if strings.TrimSpace(req.TrackingNo) == "" {
		return nil, fmt.Errorf("快递单号不能为空")
	}
	if strings.TrimSpace(req.CompanyCode) == "" {
		return nil, fmt.Errorf("快递公司编码不能为空")
	}
	if strings.TrimSpace(d.Customer) == "" || strings.TrimSpace(d.Key) == "" {
		return nil, fmt.Errorf("快递100配置不完整")
	}
	param := map[string]string{
		"com":      normalizeCompanyCodeForKuaidi100(req.CompanyCode),
		"num":      strings.TrimSpace(req.TrackingNo),
		"phone":    strings.TrimSpace(req.Phone),
		"resultv2": "1",
	}
	paramJSON, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	sign := d.sign(string(paramJSON))
	form := url.Values{}
	form.Set("customer", strings.TrimSpace(d.Customer))
	form.Set("sign", sign)
	form.Set("param", string(paramJSON))

	endpoint := strings.TrimSpace(d.APIURL)
	if endpoint == "" {
		endpoint = "https://poll.kuaidi100.com/poll/query.do"
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		State   string `json:"state"`
		IsCheck string `json:"ischeck"`
		Com     string `json:"com"`
		Nu      string `json:"nu"`
		Data    []struct {
			Time     string `json:"time"`
			FTime    string `json:"ftime"`
			Context  string `json:"context"`
			Location string `json:"location"`
			Status   string `json:"status"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if strings.TrimSpace(body.Status) != "200" && strings.TrimSpace(body.Status) != "ok" {
		msg := strings.TrimSpace(body.Message)
		if msg == "" {
			msg = "快递100查询失败"
		}
		return nil, errors.New(msg)
	}

	result := &logisticsDriver.TrackResult{
		Provider:   d.Name(),
		StatusCode: mapKuaidi100State(body.State),
		StatusText: strings.TrimSpace(body.Message),
		Nodes:      make([]logisticsDriver.TrackNode, 0, len(body.Data)),
	}
	if result.StatusText == "" {
		result.StatusText = result.StatusCode
	}
	for _, item := range body.Data {
		node := logisticsDriver.TrackNode{
			Time:       parseTrackTime(item.FTime, item.Time),
			Location:   strings.TrimSpace(item.Location),
			StatusCode: result.StatusCode,
			StatusText: strings.TrimSpace(item.Context),
		}
		raw, _ := json.Marshal(item)
		node.RawPayload = raw
		result.Nodes = append(result.Nodes, node)
	}
	if result.StatusCode == "signed" && len(result.Nodes) > 0 {
		signed := result.Nodes[0].Time
		result.SignedAt = &signed
	}
	return result, nil
}

func (d *kuaidi100Driver) sign(param string) string {
	sum := md5.Sum([]byte(param + strings.TrimSpace(d.Key) + strings.TrimSpace(d.Customer)))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

func mapKuaidi100State(state string) string {
	switch strings.TrimSpace(state) {
	case "0", "1", "5":
		return "in_transit"
	case "2", "4":
		return "exception"
	case "3":
		return "signed"
	default:
		return "shipped"
	}
}

func parseTrackTime(values ...string) time.Time {
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	}
	for _, value := range values {
		candidate := strings.TrimSpace(value)
		if candidate == "" {
			continue
		}
		for _, layout := range layouts {
			if parsed, err := time.ParseInLocation(layout, candidate, time.Local); err == nil {
				return parsed
			}
		}
	}
	return time.Now()
}

func normalizeCompanyCodeForKuaidi100(code string) string {
	value := strings.ToUpper(strings.TrimSpace(code))
	switch value {
	case "SF", "SHUNFENG":
		return "shunfeng"
	case "ZTO":
		return "zhongtong"
	case "YTO":
		return "yuantong"
	case "STO":
		return "shentong"
	case "YD", "YUNDA":
		return "yunda"
	case "EMS":
		return "ems"
	case "JD":
		return "jd"
	case "DB", "DBL", "DEPPON":
		return "debangwuliu"
	case "JT", "JTEXPRESS":
		return "jtexpress"
	default:
		return strings.ToLower(strings.TrimSpace(code))
	}
}
