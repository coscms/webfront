package top

import (
	"github.com/coscms/webcore/library/httpserver"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/subdomains"
)

func URLFor(purl string) string {
	return subdomains.Default.URL(purl, httpserver.KindFrontend)
}

func RelativeURL(purl string, args ...string) string {
	return subdomains.Default.RelativeURL(purl, httpserver.KindFrontend)
}

func URLByName(name string, params ...interface{}) string {
	return subdomains.Default.URLByName(name, params...)
}

func RelativeURLByName(name string, params ...interface{}) string {
	return subdomains.Default.RelativeURLByName(name, params...)
}

func BackendURLFor(purl string) string {
	return subdomains.Default.URL(purl, httpserver.KindBackend)
}

func BackendRelativeURL(purl string, args ...string) string {
	return subdomains.Default.RelativeURL(purl, httpserver.KindBackend)
}

func FrontendURLByName(name string, params ...interface{}) string {
	return subdomains.Default.URLByNamex(httpserver.KindFrontend, name, params...)
}

func FrontendRelativeURLByName(name string, params ...interface{}) string {
	return subdomains.Default.RelativeURLByNamex(httpserver.KindFrontend, name, params...)
}

func BackendURLByName(name string, params ...interface{}) string {
	return subdomains.Default.URLByNamex(httpserver.KindBackend, name, params...)
}

func BackendRelativeURLByName(name string, params ...interface{}) string {
	return subdomains.Default.RelativeURLByNamex(httpserver.KindBackend, name, params...)
}

func AbsoluteURL(purl string) string {
	if !com.IsFullURL(purl) {
		return URLFor(purl)
	}
	return purl
}
