package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/admpub/imageproxy"
	"github.com/coscms/webcore/cmd"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
)

var minifyCmd = &cobra.Command{
	Use:   "minify",
	Short: "minify file",
	Long:  `Usage ` + executable + ` minify src.jpg dest.jpg`,
	RunE:  minifyRunE,
}

var minifyIMGOptions = imageproxy.Options{
	Quality: 70,
}

func minifyRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return cmd.Usage()
	}
	file := args[0]
	src, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	var res []byte
	ext := filepath.Ext(file)
	ext = strings.ToLower(ext)
	switch ext {
	case `.css`:
		m := minify.New()
		m.AddFunc("text/css", css.Minify)
		res, err = m.Bytes("text/css", src)
	case `.js`:
		m := minify.New()
		m.AddFunc("application/javascript", js.Minify)
		res, err = m.Bytes("application/javascript", src)
	case `.jpeg`, `.webp`, `.bmp`, `.gif`, `.png`, `.tiff`:
		res, err = imageproxy.Transform(src, minifyIMGOptions)
	default:
		err = fmt.Errorf(`unsupported file: %v`, filepath.Base(file))
	}
	if err != nil {
		return err
	}
	err = os.WriteFile(args[1], res, 0644)
	return err
}

func init() {
	cmd.Add(minifyCmd)
	minifyCmd.Flags().IntVar(&minifyIMGOptions.Quality, `quality`, minifyIMGOptions.Quality, `图片保留质量`)
}
