package article

import (
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webfront/dbschema"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

var AllowComment = map[string]func(echo.Context, *dbschema.OfficialCustomer, *dbschema.OfficialCommonArticle) error{
	`buyer`:  allowBuyerComment,
	`author`: allowAuthorComment,
	`admin`:  allowAdminComment,
	`none`:   disallowComment,
}

func allowBuyerComment(c echo.Context, customer *dbschema.OfficialCustomer, articleM *dbschema.OfficialCommonArticle) error {
	if customer == nil {
		return c.E(`很抱歉，您不是当前商品买家，本文只有当前商品买家才能评论`)
	}
	boughtDetector := Source.GetBoughtDetector(articleM.SourceTable)
	if boughtDetector != nil {
		if boughtDetector(customer, articleM.SourceId) {
			return nil
		}
		return c.E(`很抱歉，您不是当前商品买家，本文只有当前商品买家才能评论`)
	}
	return OrderQuerier(c, customer, articleM.SourceId, articleM.SourceTable)
}

func allowAuthorComment(c echo.Context, customer *dbschema.OfficialCustomer, articleM *dbschema.OfficialCommonArticle) error {
	user := backend.User(c)
	if articleM.OwnerType == `user` {
		if user == nil || articleM.OwnerId != uint64(user.Id) {
			customerM := modelCustomer.NewCustomer(c)
			err := customerM.Get(nil, db.And(
				db.Cond{`id`: customer.Id},
				db.Cond{`uid`: articleM.OwnerId},
			))
			if err != nil {
				if err == db.ErrNoMoreRows {
					return c.E(`很抱歉，您不是本文作者，本文只有作者才能评论`)
				}
				return err
			}
		}
	} else { //客户发布的文章，后台用户的作者可以评论
		if user == nil && (customer == nil || customer.Id != articleM.OwnerId) {
			return c.E(`很抱歉，您不是本文作者，本文只有作者才能评论`)
		}
	}
	return nil
}

func allowAdminComment(c echo.Context, customer *dbschema.OfficialCustomer, articleM *dbschema.OfficialCommonArticle) error {
	user := backend.User(c)
	if user == nil || user.Id != 1 {
		return c.E(`很抱歉，本文只有管理员才能评论`)
	}
	return nil
}

func disallowComment(c echo.Context, customer *dbschema.OfficialCustomer, articleM *dbschema.OfficialCommonArticle) error {
	return c.E(`很抱歉，本文暂未开放评论`)
}
