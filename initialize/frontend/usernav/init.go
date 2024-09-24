package usernav

import (
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/navigate"
)

func init() {
	httpserver.Frontend.Navigate.Add(navigate.Left, LeftNavigate)
	httpserver.Frontend.Navigate.Add(navigate.Top, TopNavigate)
}
