package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	admin "orderin-server/cmd/admin-api"
	app "orderin-server/cmd/app-api"
	"orderin-server/pkg/common/utils"
	"os"
)

var rootCmd = &cobra.Command{
	Use:          "orderin",
	Short:        "orderin",
	SilenceUsage: true,
	Long:         `orderin`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New(utils.Red("requires at least one arg"))
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

func tip() {
	usageStr := `欢迎使用 ` + utils.Green(`orderin 1.0.0`) + ` 可以使用 ` + utils.Red(`-h`) + ` 查看命令`
	fmt.Printf("%s\n", usageStr)
}

func init() {
	rootCmd.AddCommand(&admin.AdminApiCmd.Command)
	rootCmd.AddCommand(&app.AppApiCmd.Command)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
