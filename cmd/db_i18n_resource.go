package cmd

import (
	"fmt"
	"strconv"

	"github.com/coscms/webcore/cmd"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webfront/model/i18nm"
	"github.com/spf13/cobra"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/defaults"
)

var dbI18nResourceCmd = &cobra.Command{
	Use:   "dbI18nResource",
	Short: "official_i18n_resource table initialize and translate",
	Long: `Usage ` + executable + ` dbI18nResource init
` + executable + ` dbI18nResource translate your_table_name 100`,
	Example: executable + ` dbI18nResource init
` + executable + ` dbI18nResource translate your_table_name 100
` + executable + ` dbI18nResource translate your_table_name 100 --gtID=2 --queryAll=true --translateAll=true
` + executable + ` dbI18nResource translate your_table_name 100 --eqID=1 --queryAll=true --translateAll=true`,
	RunE: dbI18nResourceRunE,
}

type dbI18nResourceTranslateOptions struct {
	table                  string
	chunks                 int
	queryAll, translateAll bool
	eqID, gtID             uint64
}

var dbI18nResourceTranslateCfg = dbI18nResourceTranslateOptions{}

func dbI18nResourceRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmd.Usage()
	}
	err := config.ParseConfig()
	if err != nil {
		return err
	}
	operation := args[0]
	switch operation {
	case `init`:
		err = i18nm.Initialize(defaults.NewMockContext())
	case `translate`:
		cfg := dbI18nResourceTranslateCfg
		if len(args) > 0 {
			var chunksNum string
			com.SliceExtract(args, &cfg.table, &chunksNum)
			if len(chunksNum) > 0 {
				cfg.chunks, err = strconv.Atoi(chunksNum)
				if err != nil {
					return err
				}
			}
		}
		cfg.gtID, err = i18nm.AutoTranslate(defaults.NewMockContext(), cfg.table, cfg.queryAll, cfg.translateAll, cfg.eqID, cfg.gtID, cfg.chunks)
	default:
		err = fmt.Errorf(`unsupported operation: %v`, operation)
	}
	return err
}

func init() {
	cmd.Add(dbI18nResourceCmd)
	dbI18nResourceCmd.Flags().StringVar(&dbI18nResourceTranslateCfg.table, `table`, dbI18nResourceTranslateCfg.table, `指定要翻译的表名`)
	dbI18nResourceCmd.Flags().BoolVar(&dbI18nResourceTranslateCfg.queryAll, `queryAll`, dbI18nResourceTranslateCfg.queryAll, `是否查询所有记录（忽略已翻译记录）`)
	dbI18nResourceCmd.Flags().BoolVar(&dbI18nResourceTranslateCfg.translateAll, `translateAll`, dbI18nResourceTranslateCfg.translateAll, `是否强制翻译所有内容`)
	dbI18nResourceCmd.Flags().Uint64Var(&dbI18nResourceTranslateCfg.eqID, `eqID`, dbI18nResourceTranslateCfg.eqID, `指定ID等于该值的记录`)
	dbI18nResourceCmd.Flags().Uint64Var(&dbI18nResourceTranslateCfg.gtID, `gtID`, dbI18nResourceTranslateCfg.gtID, `指定ID大于该值的记录`)
	dbI18nResourceCmd.Flags().IntVar(&dbI18nResourceTranslateCfg.chunks, `chunks`, dbI18nResourceTranslateCfg.chunks, `指定分块大小`)
}
