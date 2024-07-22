package cmd

import (
	"github.com/spf13/cobra"
	"orderin-server/pkg/common/config"
	"orderin-server/pkg/common/constant"
)

const (
	AdminApi = "admin-api"
	AppApi   = "app-api"
)

type ApiCmd struct {
	*RootCmd
}

func NewApiCmd(name string) *ApiCmd {
	ret := &ApiCmd{NewRootCmd(name)}
	ret.SetRootCmdPt(ret)
	return ret
}

func (a *ApiCmd) GetPortFromConfig(portType string) int {
	switch a.Name {
	case AdminApi:
		if portType == constant.FlagPort {
			return config.Config.AdminApi.Port[0]
		} else if portType == constant.FlagPrometheusPort {
			return config.Config.Prometheus.AdminPrometheusPort[0]
		}
	case AppApi:
		if portType == constant.FlagPort {
			return config.Config.AppApi.Port[0]
		} else if portType == constant.FlagPrometheusPort {
			return config.Config.Prometheus.AppPrometheusPort[0]
		}
	}
	return 0
}

func (a *ApiCmd) AddApi(f func(port int, promPort int) error) {
	a.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return f(a.getPortFlag(cmd), a.getPrometheusPortFlag(cmd))
	}
}
