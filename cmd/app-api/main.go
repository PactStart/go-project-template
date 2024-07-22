package app_api

import (
	"context"
	"fmt"
	"net"
	"orderin-server/internal/components"
	"orderin-server/internal/routers"
	"orderin-server/pkg/common/cache"
	"orderin-server/pkg/common/cmd"
	"orderin-server/pkg/common/config"
	"orderin-server/pkg/common/db/relation"
	"orderin-server/pkg/common/log"
	"strconv"
)

var (
	AppApiCmd = cmd.NewApiCmd(cmd.AppApi)
)

func init() {
	AppApiCmd.AddPortFlag()
	AppApiCmd.AddPrometheusPortFlag()
	AppApiCmd.AddApi(run)
}

func run(port int, proPort int) error {
	log.ZInfo(context.Background(), "app-api port:", "port", port, "proPort", proPort)

	if port == 0 || proPort == 0 {
		err := "port or proPort is empty:" + strconv.Itoa(port) + "," + strconv.Itoa(proPort)
		log.ZError(context.Background(), err, nil)

		return fmt.Errorf(err)
	}
	rdb, err := cache.NewRedis()
	if err != nil {
		log.ZError(context.Background(), "Failed to initialize Redis", err)
		return err
	}
	db, err := relation.NewGormDB()
	if err != nil {
		log.ZError(context.Background(), "Failed to initialize mysql db", err)
		return err
	}
	components.RegisterComponents(db, rdb)

	router := routers.NewH5GinRouter(db, rdb)
	log.ZInfo(context.Background(), "app-api init routers success")
	var address string
	if config.Config.AdminApi.ListenIP != "" {
		address = net.JoinHostPort(config.Config.AdminApi.ListenIP, strconv.Itoa(port))
	} else {
		address = net.JoinHostPort("0.0.0.0", strconv.Itoa(port))
	}
	log.ZInfo(context.Background(), "start app-api server", "address", address, "version", config.Version)

	err = router.Run(address)
	if err != nil {
		log.ZError(context.Background(), "app-api run failed", err, "address", address)

		return err
	}

	return nil
}
