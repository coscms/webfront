package comment

import (
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/middleware/sessdata"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

// QueryCommentList 查询评论列表
// @param c echo.Context
// @param articleID uint64 文章ID
// @param articleSN string 文章SN
// @param targetType string 评论目标类型
// @param subType string 评论子类型
// @param flat bool 是否平铺评论
// @param urlLayout string 分页的 URL 格式
// @param pagingVarSuffix ...string 分页的参数后缀
// @return []*CommentAndExtra 评论列表
// @return error 错误
func QueryCommentList(c echo.Context, articleID uint64, articleSN string, targetType string, subType string, flat bool, urlLayout string, pagingVarSuffix ...string) ([]*CommentAndExtra, error) {
	tp, ok := CommentAllowTypes[targetType]
	if !ok {
		return nil, c.NewError(code.Unsupported, `不支持评论目标: %v`, targetType).SetZone(`type`)
	}
	var err error
	if len(articleSN) > 0 && tp.SN2ID != nil {
		articleID, err = tp.SN2ID(c, articleSN)
		if err != nil {
			return nil, err
		}
	}
	cmtM := NewComment(c)
	cond := cmtM.ListCond(targetType, subType, articleID, flat)
	list := []*CommentAndReplyTarget{}
	p, err := common.NewLister(cmtM, &list, func(r db.Result) db.Result {
		return r.OrderBy(`id`).Relation(`ReplyTarget`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
			if NeedWithQuoteComment(c) {
				return sel
			}
			return nil
		})
	}, cond.And()).Paging(c, pagingVarSuffix...)
	if err != nil {
		return nil, err
	}
	if len(urlLayout) > 0 {
		p.SetURL(urlLayout)
	}
	var rowNums map[uint64]int
	var replyCommentIndexes map[uint64][]int
	if flat {
		replyCommentIndexes = map[uint64][]int{}
		replyCommentIds := []uint64{}
		for index, row := range list {
			if row.ReplyCommentId > 0 {
				if _, ok := replyCommentIndexes[row.ReplyCommentId]; !ok {
					replyCommentIndexes[row.ReplyCommentId] = []int{}
					replyCommentIds = append(replyCommentIds, row.ReplyCommentId)
				}
				replyCommentIndexes[row.ReplyCommentId] = append(replyCommentIndexes[row.ReplyCommentId], index)
			}
		}
		if len(replyCommentIds) > 0 {
			rowNums, err = cmtM.RowNums(targetType, subType, articleID, replyCommentIds)
			if err != nil {
				return nil, err
			}
		}
	}

	rows, err := cmtM.WithExtra(list, sessdata.Customer(c), backend.User(c), p)
	if err != nil {
		return nil, err
	}
	for id, rowNum := range rowNums {
		for _, index := range replyCommentIndexes[id] {
			rows[index].ReplyFloorNumber = rowNum
		}
	}
	return rows, err
}

// QueryCommentReplyList 根据评论ID查询评论下的回复列表
//
//	commentID: 评论ID
//	urlLayout:  URL前缀
//	pagingVarSuffix:  分页变量后缀
//
//	返回评论下的回复列表,如果评论不存在,则返回DataNotFound错误
func QueryCommentReplyList(c echo.Context, commentID uint64, urlLayout string, pagingVarSuffix ...string) ([]*CommentAndExtra, error) {
	if commentID == 0 {
		return nil, c.NewError(code.InvalidParameter, `commentId无效`)
	}
	cmtM := NewComment(c)
	err := cmtM.Get(nil, `id`, commentID)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, c.NewError(code.DataNotFound, `评论不存在`)
		}
		return nil, err
	}
	cond := cmtM.ListReplyCond(commentID)
	list := []*CommentAndReplyTarget{}
	p, err := common.NewLister(cmtM, &list, func(r db.Result) db.Result {
		return r.OrderBy(`id`).Relation(`ReplyTarget`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
			if NeedWithQuoteComment(c) {
				return sel
			}
			return nil
		})
	}, cond.And()).Paging(c, pagingVarSuffix...)
	if err != nil {
		return nil, err
	}
	if len(urlLayout) > 0 {
		p.SetURL(urlLayout)
	}
	return cmtM.WithExtra(list, sessdata.Customer(c), backend.User(c), p)
}
