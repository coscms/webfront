package apiutils

import (
	"fmt"
	"net/url"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/oauth2"

	"github.com/coscms/sdk/sdk_options"
	"github.com/coscms/webfront/library/oauthutils"
)

type OAuthProvider struct {
	Name      string `json:"name" xml:"name"`
	English   string `json:"english" xml:"english"`
	IconClass string `json:"iconClass" xml:"iconClass"`
	IconImage string `json:"iconImage" xml:"iconImage"`
	WrapClass string `json:"wrapClass" xml:"wrapClass"`
	LoginURL  string `json:"loginURL" xml:"loginURL"`
}

func OAuthProvidersFrom(accounts []oauth2.Account) []*OAuthProvider {
	var providers []*OAuthProvider
	for _, item := range accounts {
		if !item.On {
			continue
		}
		provider := &OAuthProvider{
			Name:      item.Name,
			English:   item.Name,
			IconClass: ``,
			IconImage: ``,
			LoginURL:  item.LoginURL,
		}
		if item.Extra != nil {
			provider.IconImage = item.Extra.String(`iconImage`)
			provider.IconClass = item.Extra.String(`iconClass`)
			provider.WrapClass = item.Extra.String(`wrapClass`)
			title := item.Extra.String(`title`)
			if len(title) > 0 {
				provider.Name = title
			}
		}
		providers = append(providers, provider)
	}
	return providers
}

type OAuthProvidersResponse struct {
	List []*OAuthProvider `json:"list"`
}

type OAuthOption interface {
	GetAccountID() uint64
	ApplySetting() (err error)
	GetAppID() string
	OAuthProviderListURL() (string, error)
}

var OAuthOptionsCreater = func(ctx echo.Context, typ Type, generators ...sdk_options.URLValuesGenerator) OAuthOption {
	return NewOptions(ctx, TypeOAuth, generators...)
}

func OAuthProviders(ctx echo.Context) ([]*OAuthProvider, error) {
	apiOpt := OAuthOptionsCreater(ctx, TypeOAuth)
	accountID := apiOpt.GetAccountID()
	if accountID <= 0 {
		return OAuthProvidersFrom(oauthutils.Accounts()), nil
	}
	if err := apiOpt.ApplySetting(); err != nil {
		return nil, err
	}
	appID := apiOpt.GetAppID()
	if len(appID) == 0 {
		return OAuthProvidersFrom(oauthutils.Accounts()), nil
	}
	apiURL, err := apiOpt.OAuthProviderListURL()
	if err != nil {
		return nil, err
	}
	platformList := &OAuthProvidersResponse{}
	apiResp := echo.NewData(ctx)
	apiResp.Data = platformList
	_, err = sdk_options.SubmitWithRecv(ctx, apiResp, apiURL, url.Values{})
	if err == nil && apiResp.Code.Int() != 1 {
		err = fmt.Errorf(`OauthProviders: %v`, apiResp.Info)
	}
	return platformList.List, err
}

func GetOAuthProviderTitle(list []*OAuthProvider, name string) string {
	for _, v := range list {
		if v.English == name {
			return v.Name
		}
	}
	return ``
}

func GetOAuthProvider(list []*OAuthProvider, name string) *OAuthProvider {
	for _, v := range list {
		if v.English == name {
			return v
		}
	}
	return nil
}
