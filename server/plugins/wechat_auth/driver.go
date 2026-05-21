package wechat_auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ijry/lyshop/core/driver/oauth"
)

// WechatAuthDriver implements oauth.Driver for WeChat.
type WechatAuthDriver struct {
	MiniAppID     string
	MiniAppSecret string
	H5AppID       string
	H5AppSecret   string
}

func (d *WechatAuthDriver) Name() string { return "wechat" }

// GetAuthURL returns WeChat H5 OAuth2 redirect URL.
func (d *WechatAuthDriver) GetAuthURL(state string) string {
	return fmt.Sprintf(
		"https://open.weixin.qq.com/connect/oauth2/authorize"+
			"?appid=%s&redirect_uri=%%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect",
		d.H5AppID, state)
}

// HandleCallback handles miniapp wx.login code → openid/session_key via code2session.
func (d *WechatAuthDriver) HandleCallback(ctx context.Context, code string) (*oauth.UserInfo, error) {
	if d.MiniAppID == "" {
		return nil, fmt.Errorf("微信小程序未配置 AppID")
	}
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		d.MiniAppID, d.MiniAppSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		OpenID     string `json:"openid"`
		UnionID    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("wx code2session error %d: %s", result.ErrCode, result.ErrMsg)
	}
	return &oauth.UserInfo{
		OpenID:  result.OpenID,
		UnionID: result.UnionID,
	}, nil
}

// GetUserInfo fetches WeChat user profile via access token (H5 OAuth2 flow).
func (d *WechatAuthDriver) GetUserInfo(ctx context.Context, accessToken string) (*oauth.UserInfo, error) {
	// Production: GET https://api.weixin.qq.com/sns/userinfo?access_token=...&openid=...&lang=zh_CN
	return &oauth.UserInfo{Nickname: "微信用户", Avatar: ""}, nil
}
