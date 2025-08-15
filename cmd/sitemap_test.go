package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	root := `.`
	dirs, err := os.ReadDir(root)
	if err != nil {
		panic(err.Error())
	}
	for _, info := range dirs {

		t.Log(info.Name())
		if info.IsDir() {
			t.Logf(`dir: %s`, filepath.Join(root, info.Name()))
		}
		//sitemap.RemoveLanguage(_lang, subDir)
	}
}
