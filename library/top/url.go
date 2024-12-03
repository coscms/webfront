package top

import (
	"github.com/webx-top/com"
	"github.com/webx-top/echo/subdomains"
)

func URLFor(purl string) string {
	return subdomains.Default.URL(purl, `frontend`)
}

func URLByName(name string, params ...interface{}) string {
	return subdomains.Default.URLByName(name, params...)
}

func AbsoluteURL(purl string) string {
	if !com.IsFullURL(purl) {
		return URLFor(purl)
	}
	return purl
}
