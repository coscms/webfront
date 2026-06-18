package oauthutils

import (
	"github.com/admpub/goth"
	"github.com/webx-top/echo/handler/oauth2"

	// - oauth2 provider
	"github.com/admpub/goth/providers/microsoftonline"
	"github.com/coscms/webfront/library/oauth2client/providers/webx"
)

// RegisterProvider 注册Provider
func RegisterProvider(c *oauth2.Config) {

	oauth2.Register(`microsoft`, func(account *oauth2.Account) goth.Provider {
		if len(account.CallbackURL) == 0 {
			account.CallbackURL = c.CallbackURL(account.Name)
		}
		m := microsoftonline.New(account.Key, account.Secret, account.CallbackURL)
		m.SetName(`microsoft`)
		return m
	})

	oauth2.Register(`coscms`, func(account *oauth2.Account) goth.Provider {
		hostURL := account.Extra.String(`hostURL`, `https://www.coscms.com`)
		if len(account.CallbackURL) == 0 {
			account.CallbackURL = oauth2.DefaultPath + "/callback/" + account.Name
		}
		m := webx.New(account.Key, account.Secret, account.CallbackURL, hostURL, `profile`)
		m.SetName(`coscms`)
		return m
	})

	/*
		oauth2.Register(`apple`, func(account *oauth2.Account) goth.Provider {
			if len(account.CallbackURL) == 0 {
				account.CallbackURL = c.CallbackURL(account.Name)
			}
			return apple.New(account.Key, account.Secret, account.CallbackURL)
		})
	*/
}
