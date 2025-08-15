package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/webx-top/com"

	"github.com/coscms/webcore/cmd"
	"github.com/coscms/webfront/library/sitemap"
)

var sitemapCmd = &cobra.Command{
	Use:   "sitemap",
	Short: "sitemap generate",
	Long: `Usage ./webx sitemap <rootURL> [languageCode...]
删除所有: ./webx sitemap --mode=clear
删除指定域名或者指定语言: ./webx sitemap --mode=clear <rootURL> en,zh-CN
删除所有域名下指定语言: ./webx sitemap --mode=clear all en,zh-CN`,
	RunE: sitemapRunE,
}

var sitemapCfg = sitemap.NewConfig()

func sitemapRunE(cmd *cobra.Command, args []string) error {
	var rootURL, langCode string
	com.SliceExtract(args, &rootURL, &langCode)
	return sitemap.CmdGenerate(rootURL, langCode, sitemapCfg)
}

func init() {
	cmd.Add(sitemapCmd)
	sitemapCmd.Flags().StringVar(&sitemapCfg.Mode, `mode`, sitemapCfg.Mode, `模式。支持的值有: full(全量生成) / incr(增量生成) / clear(删除)`)

	slic := sitemap.Registry.Slice()
	groups := make([]string, len(slic))
	for i, v := range slic {
		groups[i] = v.K
	}
	sitemapCmd.Flags().StringVar(&sitemapCfg.Group, `group`, sitemapCfg.Group, `组。指定多个时用半角逗号“,”隔开, 支持的值有: `+strings.Join(groups, ` / `)+` 等。不指定时代表所有组`)
	sitemapCmd.Flags().BoolVar(&sitemapCfg.AllChild, `allChild`, sitemapCfg.AllChild, `是否同时生成所有子页面中的网址`)
}
