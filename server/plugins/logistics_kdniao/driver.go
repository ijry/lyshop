package logistics_kdniao

import (
	"context"
	"crypto/md5"
	"encoding/base64"
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

type kdniaoDriver struct {
	APIURL      string
	EBusinessID string
	APIKey      string
}

func (d *kdniaoDriver) Name() string { return "kdniao" }

func (d *kdniaoDriver) Query(ctx context.Context, req logisticsDriver.QueryReq) (*logisticsDriver.TrackResult, error) {
	if strings.TrimSpace(req.TrackingNo) == "" {
		return nil, fmt.Errorf("快递单号不能为空")
	}
	if strings.TrimSpace(req.CompanyCode) == "" {
		return nil, fmt.Errorf("快递公司编码不能为空")
	}
	if strings.TrimSpace(d.EBusinessID) == "" || strings.TrimSpace(d.APIKey) == "" {
		return nil, fmt.Errorf("快递鸟配置不完整")
	}
	requestDataBody := map[string]string{
		"ShipperCode":  normalizeCompanyCodeForKdniao(req.CompanyCode),
		"LogisticCode": strings.TrimSpace(req.TrackingNo),
	}
	if phone := strings.TrimSpace(req.Phone); phone != "" {
		requestDataBody["CustomerName"] = phone
	}
	requestDataRaw, err := json.Marshal(requestDataBody)
	if err != nil {
		return nil, err
	}
	requestData := string(requestDataRaw)
	form := url.Values{}
	form.Set("RequestData", url.QueryEscape(requestData))
	form.Set("EBusinessID", strings.TrimSpace(d.EBusinessID))
	form.Set("RequestType", "1002")
	form.Set("DataSign", url.QueryEscape(d.buildSign(requestData)))
	form.Set("DataType", "2")

	endpoint := strings.TrimSpace(d.APIURL)
	if endpoint == "" {
		endpoint = "https://api.kdniao.com/Ebusiness/EbusinessOrderHandle.aspx"
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
		Success      bool   `json:"Success"`
		Reason       string `json:"Reason"`
		State        string `json:"State"`
		LogisticCode string `json:"LogisticCode"`
		ShipperCode  string `json:"ShipperCode"`
		Traces       []struct {
			AcceptTime    string `json:"AcceptTime"`
			AcceptStation string `json:"AcceptStation"`
			Location      string `json:"Location"`
			Remark        string `json:"Remark"`
		} `json:"Traces"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if !body.Success {
		msg := strings.TrimSpace(body.Reason)
		if msg == "" {
			msg = "快递鸟查询失败"
		}
		return nil, errors.New(msg)
	}
	result := &logisticsDriver.TrackResult{
		Provider:   d.Name(),
		StatusCode: mapKdniaoState(body.State),
		StatusText: strings.TrimSpace(body.Reason),
		Nodes:      make([]logisticsDriver.TrackNode, 0, len(body.Traces)),
	}
	if result.StatusText == "" {
		result.StatusText = result.StatusCode
	}
	for _, item := range body.Traces {
		node := logisticsDriver.TrackNode{
			Time:       parseTrackTime(item.AcceptTime),
			Location:   strings.TrimSpace(item.Location),
			StatusCode: result.StatusCode,
			StatusText: strings.TrimSpace(item.AcceptStation),
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

func (d *kdniaoDriver) buildSign(requestData string) string {
	sum := md5.Sum([]byte(requestData + strings.TrimSpace(d.APIKey)))
	md5String := hex.EncodeToString(sum[:])
	return base64.StdEncoding.EncodeToString([]byte(md5String))
}

func mapKdniaoState(state string) string {
	switch strings.TrimSpace(state) {
	case "2":
		return "in_transit"
	case "3":
		return "signed"
	case "4":
		return "exception"
	default:
		return "shipped"
	}
}

func parseTrackTime(value string) time.Time {
	candidate := strings.TrimSpace(value)
	if candidate == "" {
		return time.Now()
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, candidate, time.Local); err == nil {
			return parsed
		}
	}
	return time.Now()
}

func normalizeCompanyCodeForKdniao(code string) string {
	value := strings.ToUpper(strings.TrimSpace(code))
	switch value {
	case "SF", "SHUNFENG":
		return "SF"
	case "ZTO":
		return "ZTO"
	case "YTO":
		return "YTO"
	case "STO":
		return "STO"
	case "YD", "YUNDA":
		return "YD"
	case "EMS":
		return "EMS"
	case "JD":
		return "JD"
	case "DB", "DBL", "DEPPON":
		return "DBL"
	case "JT", "JTEXPRESS":
		return "JTSD"
	default:
		return strings.ToUpper(strings.TrimSpace(code))
	}
}
