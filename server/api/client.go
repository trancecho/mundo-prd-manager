package api

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
	"github.com/spf13/viper"
	"github.com/trancecho/ragnarok/fastgpt"
	"net/http"
	"time"
)

var Cli *Client

type Client struct {
	Client    *lark.Client    // 飞书客户端
	WsCli     *larkws.Client  // 飞书WebSocket客户端
	Fcli      *fastgpt.Client // FastGPT客户端
	appId     string
	appSecret string
}

func ClientInit() {
	AppID, AppSecret := viper.GetString("feishu.appID"), viper.GetString("feishu.appSecret")
	sclient := lark.NewClient(AppID, AppSecret, // 默认配置为自建应用
		// lark.WithMarketplaceApp(), // 可设置为商店应用
		lark.WithLogLevel(larkcore.LogLevelDebug),
		lark.WithReqTimeout(3*time.Second),
		lark.WithEnableTokenCache(true),
		lark.WithHelpdeskCredential("id", "token"),
		lark.WithHttpClient(http.DefaultClient),
	)
	api, url := viper.GetString("fastgpt.api_key"), viper.GetString("fastgpt.base_url")
	fclient := fastgpt.NewFastClient(api, url)
	Cli = &Client{
		Client:    sclient,
		Fcli:      fclient,
		appId:     AppID,
		appSecret: AppSecret,
	}
}
