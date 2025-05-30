package comment

import (
	"fmt"

	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/model/official/article"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func commentArticleGetTarget(ctx echo.Context, targetID uint64) (res echo.H, breadcrumb []echo.KV, detailURL string, err error) {
	articleM := article.NewArticle(ctx)
	err = articleM.Get(nil, `id`, targetID)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, `文章不存在`).SetZone(`id`)
		}
		return
	}
	res = articleM.AsMap()
	id := fmt.Sprint(targetID)
	breadcrumb = []echo.KV{
		{K: `official/article/index`, V: ctx.T(`文章列表`)},
		{K: `official/article/edit?id=` + id, V: articleM.Title},
	}
	detailURL = `article/` + id + ctx.DefaultExtension()
	return
}

func commentArticleCheck(ctx echo.Context, f *dbschema.OfficialCommonComment) error {
	if f.TargetId < 1 {
		return ctx.NewError(code.InvalidParameter, `文章ID无效`).SetZone(`TargetId`)
	}
	newsM := dbschema.NewOfficialCommonArticle(ctx)
	err := newsM.Get(nil, `id`, f.TargetId)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		return ctx.NewError(code.DataNotFound, `您要评论的文章(ID:%d)不存在`, f.TargetId).SetZone(`id`)
	}
	if newsM.CommentAutoDisplay == `Y` {
		f.Display = `Y`
	} else {
		f.Display = `N`
	}
	f.TargetOwnerId = newsM.OwnerId
	f.TargetOwnerType = newsM.OwnerType
	return nil
}

func commentArticleAfterAdd(ctx echo.Context, f *dbschema.OfficialCommonComment) error {
	if f.ReplyCommentId == 0 { //不包含对评论的回复统计，只包含根评论
		newsM := dbschema.NewOfficialCommonArticle(ctx)
		return newsM.UpdateField(nil, `comments`, db.Raw(`comments+1`), db.Cond{`id`: f.TargetId})
	}
	return nil
}

func commentArticleAfterDelete(ctx echo.Context, f *dbschema.OfficialCommonComment) error {
	if f.ReplyCommentId == 0 { //不包含对评论的回复统计，只包含根评论
		newsM := dbschema.NewOfficialCommonArticle(ctx)
		return newsM.UpdateField(nil, `comments`, db.Raw(`comments-1`), db.And(
			db.Cond{`id`: f.TargetId},
			db.Cond{`comments`: db.Gt(0)},
		))
	}
	return nil
}
