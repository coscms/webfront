package cmd

import (
	"fmt"

	"github.com/coscms/webcore/cmd"
	"github.com/coscms/webfront/model/i18nm"
	"github.com/spf13/cobra"
	"github.com/webx-top/echo/defaults"
)

var i18nResourceCmd = &cobra.Command{
	Use:   "i18nResource",
	Short: "official_i18n_resource table initialize",
	Long:  `Usage ./webx i18nResource init`,
	RunE:  i18nResourceRunE,
}

func i18nResourceRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmd.Usage()
	}
	var err error
	operation := args[0]
	switch operation {
	case `init`:
		err = i18nm.Initialize(defaults.NewMockContext())
	default:
		err = fmt.Errorf(`unsupported operation: %v`, operation)
	}
	return err
}

func init() {
	cmd.Add(i18nResourceCmd)
}
