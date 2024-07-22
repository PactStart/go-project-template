package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"orderin-server/pkg/common/config"
	"orderin-server/pkg/common/constant"
	"orderin-server/pkg/common/email"
	"orderin-server/pkg/common/file_store"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/sms"
)

type RootCmdPt interface {
	GetPortFromConfig(portType string) int
}
type RootCmd struct {
	Command        cobra.Command
	Name           string
	port           int
	prometheusPort int
	cmdItf         RootCmdPt
}

type CmdOpts struct {
	logFilePrefixName string
}

func NewRootCmd(name string, opts ...func(*CmdOpts)) *RootCmd {
	rootCmd := &RootCmd{Name: name}
	cmd := cobra.Command{
		Use:     name,
		Short:   fmt.Sprintf("Start %s server", name),
		Example: fmt.Sprintf("orderin %s -c config/config.yaml", name),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.persistentPreRun(cmd, opts...)
		},
	}
	rootCmd.Command = cmd
	rootCmd.addConfFlag()
	return rootCmd
}

func (rc *RootCmd) persistentPreRun(cmd *cobra.Command, opts ...func(*CmdOpts)) error {
	//初始化配置
	if err := rc.initializeConfiguration(cmd); err != nil {
		return fmt.Errorf("failed to get configuration from command: %w", err)
	}
	cmdOpts := rc.applyOptions(opts...)
	//初始化logger
	if err := rc.initializeLogger(cmdOpts); err != nil {
		return fmt.Errorf("failed to initialize from config: %w", err)
	}

	//初始化file store
	file_store.InitFileStore()

	//初始化sms vendor
	sms.InitSmsVendor()

	//初始化email vendor
	email.InitEmailVendor()
	return nil
}

func (rc *RootCmd) initializeConfiguration(cmd *cobra.Command) error {
	return rc.getConfFromCmdAndInit(cmd)
}

func (rc *RootCmd) applyOptions(opts ...func(*CmdOpts)) *CmdOpts {
	cmdOpts := defaultCmdOpts()
	for _, opt := range opts {
		opt(cmdOpts)
	}

	return cmdOpts
}

func (rc *RootCmd) initializeLogger(cmdOpts *CmdOpts) error {
	logConfig := config.Config.Log

	return log.InitFromConfig(
		cmdOpts.logFilePrefixName,
		rc.Name,
		logConfig.RemainLogLevel,
		logConfig.IsStdout,
		logConfig.IsJson,
		logConfig.StorageLocation,
		logConfig.RemainRotationCount,
		logConfig.RotationTime,
	)
}

func defaultCmdOpts() *CmdOpts {
	return &CmdOpts{
		logFilePrefixName: "orderin.log.all",
	}
}

func (r *RootCmd) SetRootCmdPt(cmdItf RootCmdPt) {
	r.cmdItf = cmdItf
}

func (r *RootCmd) addConfFlag() {
	r.Command.Flags().StringP(constant.FlagConf, "c", "", "path to config file folder")
}

func (r *RootCmd) AddPortFlag() {
	r.Command.Flags().IntP(constant.FlagPort, "p", 0, "server listen port")
}

func (r *RootCmd) getPortFlag(cmd *cobra.Command) int {
	port, err := cmd.Flags().GetInt(constant.FlagPort)
	if err != nil {
		fmt.Println("Error getting ws port flag:", err)
	}
	if port == 0 {
		port = r.PortFromConfig(constant.FlagPort)
	}
	return port
}

func (r *RootCmd) GetPortFlag() int {
	return r.port
}

func (r *RootCmd) AddPrometheusPortFlag() {
	r.Command.Flags().IntP(constant.FlagPrometheusPort, "", 0, "server prometheus listen port")
}

func (r *RootCmd) getPrometheusPortFlag(cmd *cobra.Command) int {
	port, _ := cmd.Flags().GetInt(constant.FlagPrometheusPort)
	if port == 0 {
		port = r.PortFromConfig(constant.FlagPrometheusPort)
	}
	return port
}

func (r *RootCmd) GetPrometheusPortFlag() int {
	return r.prometheusPort
}

func (r *RootCmd) getConfFromCmdAndInit(cmdLines *cobra.Command) error {
	configFolderPath, _ := cmdLines.Flags().GetString(constant.FlagConf)
	fmt.Println("configFolderPath:", configFolderPath)
	return config.InitConfig(configFolderPath)
}

func (r *RootCmd) Execute() error {
	return r.Command.Execute()
}

func (r *RootCmd) AddCommand(cmds ...*cobra.Command) {
	r.Command.AddCommand(cmds...)
}

func (r *RootCmd) PortFromConfig(portType string) int {
	return r.cmdItf.GetPortFromConfig(portType)
}
