package usernav

import (
	"github.com/coscms/webcore/library/navigate"
)

var LeftNavigate = &navigate.List{
	{
		Display:   true,
		Name:      `个人资料`,
		Action:    `profile`,
		Icon:      `user`,
		Unlimited: false,
		Children: &navigate.List{
			{
				Display:   true,
				Name:      `个人资料`,
				Action:    `detail`,
				Icon:      `user`,
				Unlimited: true,
				Children:  &navigate.List{},
			},
			{
				Display:   false,
				Name:      `我关注的`,
				Action:    `following`,
				Icon:      ``,
				Unlimited: false,
				Children:  &navigate.List{},
			},
			{
				Display:   false,
				Name:      `关注我的`,
				Action:    `followers`,
				Icon:      ``,
				Unlimited: false,
				Children:  &navigate.List{},
			},
		},
	},
	{
		Display:   false,
		Name:      `文件管理`,
		Action:    `file`,
		Icon:      ``,
		Unlimited: false,
		Children: &navigate.List{
			{
				Display:   false,
				Name:      `文件选择`,
				Action:    `finder`,
				Icon:      ``,
				Unlimited: true,
				Children:  &navigate.List{},
			},
			{
				Display:   false,
				Name:      `文件上传`,
				Action:    `upload/:type`,
				Icon:      ``,
				Unlimited: false,
				Children:  &navigate.List{},
			},
		},
	},
	{
		Display:   true,
		Name:      `我的钱包`,
		Action:    `wallet`,
		Icon:      `wallet iconfont icon-licai`,
		Unlimited: true,
		Children:  &navigate.List{},
	},
	{
		Display:   true,
		Name:      `会员套餐`,
		Action:    `membership`,
		Icon:      `membership iconfont icon-dengji`,
		Unlimited: false,
		Children: &navigate.List{
			{
				Display:   true,
				Name:      `会员套餐`,
				Action:    `index`,
				Icon:      `membership iconfont icon-dengji`,
				Unlimited: false,
				Children:  &navigate.List{},
			},
			{
				Display:   false,
				Name:      `购买会员套餐`,
				Action:    `buy/:packageId`,
				Icon:      `iconfont icon-xianjinliuliangbiao`,
				Unlimited: false,
				Children:  &navigate.List{},
			},
		},
	},
	{
		Display: true,
		Name:    `短链接`,
		Action:  `short_url`,
		Icon:    `link`,
		Children: &navigate.List{
			{
				Display:  true,
				Name:     `我的短链接`,
				Action:   `list`,
				Icon:     `link`,
				Children: &navigate.List{},
			},
		},
	},
	{
		Display: true,
		Name:    `我的收藏`,
		Action:  `favorite`,
		Icon:    `link`,
		Children: &navigate.List{
			{
				Display:  true,
				Name:     `我的收藏`,
				Action:   `index`,
				Icon:     `star`,
				Children: &navigate.List{},
			},
			{
				Display:  false,
				Name:     `删除收藏`,
				Action:   `delete`,
				Icon:     `trash`,
				Children: &navigate.List{},
			},
			{
				Display:   false,
				Name:      `跳转到收藏页`,
				Action:    `go/:id`,
				Icon:      ``,
				Unlimited: true,
				Children:  &navigate.List{},
			},
		},
	},
}
