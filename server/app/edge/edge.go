package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/logx"

	"server/app/edge/internal/config"
	"server/app/edge/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/edge.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	_ = svc.NewServiceContext(c)

	logx.DisableStat()

}
