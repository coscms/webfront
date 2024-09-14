package minify

import "github.com/coscms/webfront/initialize/frontend"

func init() {
	frontend.TmplCustomParser = TmplCustomParser
}

func TmplCustomParser(_ string, content []byte) []byte {
	return Merge(content)
}
