package settings

import (
	"github.com/coscms/webcore/library/config"
)

type SettingsMultilingual struct {
	SiteName            string `json:"siteName"`
	SiteSlogan          string `json:"siteSlogan"`
	SiteMetaKeywords    string `json:"siteMetaKeywords"`
	SiteMetaDescription string `json:"siteMetaDescription"`
	SiteAnnouncement    string `json:"siteAnnouncement"`
}

func (s SettingsMultilingual) Get(name string, defaults ...string) string {
	switch name {
	case `siteName`:
		if len(s.SiteName) > 0 {
			return s.SiteName
		}
	case `siteSlogan`:
		if len(s.SiteSlogan) > 0 {
			return s.SiteSlogan
		}
	case `siteMetaKeywords`:
		if len(s.SiteMetaKeywords) > 0 {
			return s.SiteMetaKeywords
		}
	case `siteMetaDescription`:
		if len(s.SiteMetaDescription) > 0 {
			return s.SiteMetaDescription
		}
	case `siteAnnouncement`:
		if len(s.SiteAnnouncement) > 0 {
			return s.SiteAnnouncement
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ``
}

type SettingsMultilinguals map[string]SettingsMultilingual // key is language code, value is SettingsMultilingual

func GetBaseMultilinguals() *SettingsMultilinguals {
	v, _ := config.FromFile().Settings().Base.Get(`multilingual`).(*SettingsMultilinguals)
	return v
}
