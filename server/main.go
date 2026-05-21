package main

import (
	"flag"
	"log"

	"github.com/ijry/lyshop/core/app"

	// Blank-import enabled plugins so their init() registers them.
	_ "github.com/ijry/lyshop/plugins/product"
	_ "github.com/ijry/lyshop/plugins/order"
	_ "github.com/ijry/lyshop/plugins/wms"
	_ "github.com/ijry/lyshop/plugins/marketing"
	_ "github.com/ijry/lyshop/plugins/im"
	_ "github.com/ijry/lyshop/plugins/ai_image"
	_ "github.com/ijry/lyshop/plugins/decor"
	_ "github.com/ijry/lyshop/plugins/storage_local"
	_ "github.com/ijry/lyshop/plugins/storage_oss"
	_ "github.com/ijry/lyshop/plugins/storage_cos"
	_ "github.com/ijry/lyshop/plugins/storage_qiniu"
	_ "github.com/ijry/lyshop/plugins/wechat_pay"
	_ "github.com/ijry/lyshop/plugins/alipay"
	_ "github.com/ijry/lyshop/plugins/sms"
	_ "github.com/ijry/lyshop/plugins/wechat_auth"
)

func main() {
	cfg := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	if err := app.Init(*cfg); err != nil {
		log.Fatalf("init: %v", err)
	}
	if err := app.Run(); err != nil {
		log.Fatalf("run: %v", err)
	}
}
