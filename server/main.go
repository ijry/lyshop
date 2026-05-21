package main

import (
	"flag"
	"log"

	"github.com/ijry/lyshop/core/app"

	// Blank-import enabled plugins so their init() registers them.
	_ "github.com/ijry/lyshop/plugins/product"
	_ "github.com/ijry/lyshop/plugins/order"
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
