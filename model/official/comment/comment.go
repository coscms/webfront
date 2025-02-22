package comment

import (
	"fmt"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"
	"github.com/webx-top/pagination"

	"github.com/admpub/log"
	"github.com/admpub/null"
	dbschemaNging "github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/model"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/top"
	"github.com/coscms/webfront/library/xcommon"
	"github.com/coscms/webfront/library/xrole"
	"github.com/coscms/webfront/library/xrole/xroleutils"
	"github.com/coscms/webfront/middleware/sessdata"
	"github.com/coscms/webfront/model/official"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

func NewComment(ctx echo.Context) *Comment {
	m := &Comment{
		OfficialCommonComment: dbschema.NewOfficialCommonComment(ctx),
	}
	return m
}

type Comment struct {
	*dbschema.OfficialCommonComment
}

func (f *Comment) SetUserID(userID uint) {
	f.OwnerId = uint64(userID)
	f.OwnerType = `user`
}

func (f *Comment) SetCustomerID(customerID uint64) {
	f.OwnerId = customerID
	f.OwnerType = `customer`
}

func (f *Comment) ListCond(typ, subType string, id uint64, flat bool, tableAlias ...string) *db.Compounds {
	var prefix string
	if len(tableAlias) > 0 && len(tableAlias[0]) > 0 {
		prefix = tableAlias[0] + `.`
	}
	cond := db.NewCompounds()
	cond.Add(
		db.Cond{prefix + `target_type`: typ},
		db.Cond{prefix + `target_id`: id},
		db.Cond{prefix + `display`: `Y`},
	)
	if !flat {
		cond.Add(db.Cond{prefix + `reply_comment_id`: 0})
	}
	if len(subType) > 0 {
		cond.AddKV(prefix+`target_subtype`, subType)
	}
	return cond
}

func (f *Comment) ListReplyCond(commentID uint64) *db.Compounds {
	cond := db.NewCompounds()
	cond.Add(
		db.Cond{`display`: `Y`},
		db.Cond{`level`: db.NotEq(0)},
		db.Cond{`root_id`: commentID},
	)
	return cond
}

func (f *Comment) ListReplyCondMutiRoot(commentIDs []uint64) *db.Compounds {
	cond := db.NewCompounds()
	cond.Add(
		db.Cond{`display`: `Y`},
		db.Cond{`level`: db.NotEq(0)},
		db.Cond{`root_id`: db.In(commentIDs)},
	)
	return cond
}

func (f *Comment) CustomerTodayCount(customerID interface{}) (int64, error) {
	startTs, endTs := top.TodayTimestamp()
	return f.Count(nil, db.And(
		db.Cond{`owner_type`: `customer`},
		db.Cond{`owner_id`: customerID},
		db.Cond{`created`: db.Between(startTs, endTs)},
	))
}

func (f *Comment) CustomerPendingCount(customerID interface{}) (int64, error) {
	return f.Count(nil, db.And(
		db.Cond{`owner_type`: `customer`},
		db.Cond{`owner_id`: customerID},
		db.Cond{`display`: `N`},
	))
}

func (f *Comment) CustomerPendingTodayCount(customerID interface{}) (int64, error) {
	startTs, endTs := top.TodayTimestamp()
	return f.Count(nil, db.And(
		db.Cond{`owner_type`: `customer`},
		db.Cond{`owner_id`: customerID},
		db.Cond{`display`: `N`},
		db.Cond{`created`: db.Between(startTs, endTs)},
	))
}

func (f *Comment) check() (func() error, error) {
	f.Content = strings.TrimSpace(f.Content)
	if len(f.Content) < 6 {
		return nil, f.Context().E(`评论内容不可少于6个字`)
	}
	if f.OwnerId < 1 {
		if f.OwnerType == `customer` {
			customer := sessdata.Customer(f.Context())
			if customer == nil {
				return nil, f.Context().NewError(code.Unauthenticated, `请登录后再评论`)
			}
			f.SetCustomerID(customer.Id)
		} else {
			user := backend.User(f.Context())
			if user == nil {
				return nil, f.Context().NewError(code.Unauthenticated, `请登录后再评论`)
			}
			f.SetUserID(user.Id)
		}
	}
	if f.ReplyCommentId > 0 {
		cmtM := dbschema.NewOfficialCommonComment(f.Context())
		err := cmtM.Get(nil, `id`, f.ReplyCommentId)
		if err != nil {
			if err != db.ErrNoMoreRows {
				return nil, err
			}
			return nil, f.Context().E(`您要回复的评论(ID:%d)不存在`, f.ReplyCommentId)
		}
		f.ReplyOwnerId = cmtM.OwnerId
		f.ReplyOwnerType = cmtM.OwnerType
		f.RootId = cmtM.RootId
		f.TargetType = cmtM.TargetType
		f.TargetSubtype = cmtM.TargetSubtype
		f.TargetId = cmtM.TargetId
		if f.RootId < 1 {
			if cmtM.ReplyCommentId == 0 {
				cmtM.RootId = cmtM.Id
				f.RootId = cmtM.RootId
				err := cmtM.UpdateField(nil, `root_id`, f.RootId, db.Cond{`id`: f.Id})
				if err != nil {
					log.Error(err)
				}
			}
			if f.RootId < 1 {
				f.RootId = f.ReplyCommentId
			}
		}
		f.Level = cmtM.Level + 1
		f.Path = cmtM.Path
	}
	if len(f.TargetType) < 1 {
		return nil, f.Context().E(`评论目标类型不能为空`)
	}
	var afterAdd func() error
	typeCfg, ok := CommentAllowTypes[f.TargetType]
	if !ok {
		return afterAdd, f.Context().E(`不支持的评论目标类型:%s`, f.TargetType)
	}
	if err := typeCfg.CheckMaster(f.Context(), f.OfficialCommonComment); err != nil {
		return afterAdd, err
	}
	if typeCfg.AfterAdd != nil {
		afterAdd = func() error {
			return typeCfg.AfterAdd(f.Context(), f.OfficialCommonComment)
		}
	}
	if f.Display != `N` && f.Display != `Y` {
		if CommentReview() {
			f.Display = `N`
		} else {
			f.Display = `Y`
		}
	}

	if len(f.Contype) == 0 || !official.Contype.Has(f.Contype) {
		f.Contype = `text`
	}
	f.Content = common.ContentEncode(f.Content, f.Contype)
	return afterAdd, nil
}

func (f *Comment) checkCustomerAdd(permission *xrole.RolePermission) error {
	err := xcommon.CheckRoleCustomerAdd(f.Context(), permission, BehaviorName, f.OwnerId, f)
	if err == nil {
		return err
	}
	switch err {
	case xcommon.ErrCustomerRoleDisabled:
		return f.Context().E(`当前角色不支持发布评论`)
	case xcommon.ErrCustomerAddClosed:
		return f.Context().E(`评论功能已关闭`)
	case xcommon.ErrCustomerAddMaxPerDay:
		return f.Context().E(`评论失败。您的账号已达到今日最大评论数量`)
	case xcommon.ErrCustomerAddMaxPending:
		return f.Context().E(`评论失败。您的待审核评论数量已达上限，请等待审核通过后再评论`)
	default:
		return err
	}
}

func (f *Comment) Add() (pk interface{}, err error) {
	if !CommentOpen() {
		return nil, f.Context().E(`评论功能已全局关闭`)
	}
	if f.OwnerType == `customer` {
		permission := xroleutils.CustomerPermission(f.Context())
		if err = f.checkCustomerAdd(permission); err != nil {
			return nil, err
		}
	}
	var afterAdd func() error
	afterAdd, err = f.check()
	if err != nil {
		return
	}
	pk, err = f.OfficialCommonComment.Insert()
	if err != nil {
		return
	}
	if afterAdd != nil {
		err = afterAdd()
		if err != nil {
			return
		}
	}
	cmtM := dbschema.NewOfficialCommonComment(f.Context())
	if len(f.Path) > 0 {
		f.Path += `,`
	}
	f.Path += fmt.Sprint(f.Id)
	cond := db.Cond{`id`: f.OfficialCommonComment.Id}
	if f.ReplyCommentId == 0 {
		err = cmtM.UpdateFields(nil, echo.H{
			`path`:    f.Path,
			`root_id`: f.Id,
		}, cond)
		return
	}
	err = cmtM.UpdateField(nil, `path`, f.Path, cond)
	if err != nil {
		return
	}

	err = f.IncrRepliesBy(f.OfficialCommonComment)
	return
}

func (f *Comment) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	_, err := f.check()
	if err != nil {
		return err
	}
	if len(f.Path) > 0 {
		f.Path += `,`
	}
	f.Path += fmt.Sprint(f.Id)
	return f.OfficialCommonComment.Update(mw, args...)
}

func (f *Comment) IncrRepliesBy(in *dbschema.OfficialCommonComment) error {
	if in.ReplyCommentId == 0 {
		return nil
	}
	var condVal interface{}
	if in.RootId > 0 && in.ReplyCommentId != in.RootId && in.RootId != in.Id {
		condVal = db.In([]uint64{in.RootId, in.ReplyCommentId})
	} else {
		condVal = in.ReplyCommentId
	}
	return f.UpdateField(nil, `replies`, db.Raw(`replies+1`), `id`, condVal)
}

func (f *Comment) DecrRepliesBy(in *dbschema.OfficialCommonComment) error {
	if in.ReplyCommentId == 0 {
		return nil
	}
	var condVal interface{}
	if in.RootId > 0 && in.ReplyCommentId != in.RootId && in.RootId != in.Id {
		condVal = db.In([]uint64{in.RootId, in.ReplyCommentId})
	} else {
		condVal = in.ReplyCommentId
	}
	return f.UpdateField(nil, `replies`, db.Raw(`replies-1`), db.And(
		db.Cond{`id`: condVal},
		db.Cond{`replies`: db.Gt(0)},
	))
}

func (f *Comment) GetTargetIDs(cond *db.Compounds, limit int, offset int) ([]uint64, error) {
	var targetIDs []uint64
	cond.Add(db.Cond{`display`: `Y`})
	_, err := f.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Group(`target_id`).Select(`target_id`, db.Raw(`MAX(created)`))
	}, offset, limit, cond.And())
	if err != nil {
		return targetIDs, err
	}
	rows := f.Objects()
	targetIDs = make([]uint64, len(rows))
	for k, v := range rows {
		targetIDs[k] = v.TargetId
	}
	return targetIDs, err
}

func (f *Comment) RowNums(targetType, subType string, targetID uint64, ids []uint64) (map[uint64]int, error) {
	if len(ids) == 0 {
		return map[uint64]int{}, nil
	}
	idss := make([]string, len(ids))
	for i, v := range ids {
		idss[i] = param.AsString(v)
	}
	sqls := []string{
		`t.target_type = ?`,
		`t.target_id = ?`,
	}
	args := []interface{}{
		targetType,
		targetID,
	}
	if len(subType) > 0 {
		sqls = append(sqls, `t.target_subtype = ?`)
		args = append(args, subType)
	}
	r, err := f.NewParam().SetCollection(`SELECT b.rownum,b.id FROM
	(
		SELECT t.*, @rownum := @rownum + 1 AS rownum
		FROM (SELECT @rownum := 0) r, ` + dbschema.WithPrefix(`official_common_comment`) + ` AS t
		WHERE ` + strings.Join(sqls, ` AND `) + ` ORDER BY t.id ASC
	) AS b WHERE b.id IN (` + strings.Join(idss, `,`) + `)`).SetArgs(args...).Query()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	result := map[uint64]int{}
	for r.Next() {
		var rownum null.Int
		var id null.Uint64
		err = r.Scan(&rownum, &id)
		if err != nil {
			return nil, err
		}
		result[id.Uint64] = rownum.Int
	}
	return result, nil
}

func (f *Comment) WithExtra(list []*CommentAndReplyTarget, customer *dbschema.OfficialCustomer, user *dbschemaNging.NgingUser, p *pagination.Pagination) ([]*CommentAndExtra, error) {
	c := f.Context()
	listx := make([]*CommentAndExtra, len(list))
	var (
		customerIds       []uint64
		userIds           []uint64
		productIdOwnerIds = map[string]map[string]map[string][]uint64{}
		targets           = map[string][]uint64{}
		targetObjects     = map[string]map[uint64][]int{}
	)
	commentIds := make([]uint64, len(list))
	var err error
	blankClickFlow := dbschema.NewOfficialCommonClickFlow(c)
	for k, row := range list {
		extra := echo.H{
			`isBought`:                      false,
			`isTargetAuthor`:                false,
			`isAdmin`:                       false,
			`clickFlow`:                     blankClickFlow,
			`repliedCustomerIsTargetAuthor`: false,
			`repliedCustomerIsBought`:       false,
			`repliedCustomerIsAdmin`:        false,
		}
		if _, _y := targets[row.TargetType]; !_y {
			targets[row.TargetType] = []uint64{}
		}
		if !com.InUint64Slice(row.TargetId, targets[row.TargetType]) {
			targets[row.TargetType] = append(targets[row.TargetType], row.TargetId)
		}
		if _, _y := targetObjects[row.TargetType]; !_y {
			targetObjects[row.TargetType] = map[uint64][]int{}
		}
		if _, _y := targetObjects[row.TargetType][row.TargetId]; !_y {
			targetObjects[row.TargetType][row.TargetId] = []int{}
		}
		targetObjects[row.TargetType][row.TargetId] = append(targetObjects[row.TargetType][row.TargetId], k)
		extra[`targetObject`] = echo.H{
			`ownerId`:   0,
			`ownerType`: ``,
			`productId`: 0,
			`detailURL`: ``,
			`title`:     ``,
		}
		commentIds[k] = row.Id
		listx[k] = &CommentAndExtra{
			CommentAndReplyTarget: row,
			FloorNumber:           common.FloorNumber(p.Page(), p.Size(), k),
			Extra:                 extra,
		}
		if row.OwnerId > 0 {
			if row.OwnerType == `user` {
				if !com.InUint64Slice(row.OwnerId, userIds) {
					userIds = append(userIds, row.OwnerId)
				}
			} else {
				if !com.InUint64Slice(row.OwnerId, customerIds) {
					customerIds = append(customerIds, row.OwnerId)
				}
			}
		}
		if row.ReplyOwnerId > 0 {
			if row.ReplyOwnerType == `user` {
				if !com.InUint64Slice(row.ReplyOwnerId, userIds) {
					userIds = append(userIds, row.ReplyOwnerId)
				}
			} else {
				if !com.InUint64Slice(row.ReplyOwnerId, customerIds) {
					customerIds = append(customerIds, row.ReplyOwnerId)
				}
			}
		}
	}
	for targetType, row := range targetObjects {
		tp, ok := CommentAllowTypes[targetType]
		if !ok || tp.WithTarget == nil {
			continue
		}
		listx, err = tp.WithTarget(f.Context(), listx, productIdOwnerIds, targets, row)
		if err != nil {
			return listx, err
		}
	}
	if len(customerIds) > 0 {
		custM := modelCustomer.NewCustomer(c)
		_, err = custM.ListByOffset(nil, nil, 0, -1, db.Cond{`id`: db.In(customerIds)})
		if err == nil {
			for _, v := range custM.Objects() {
				for kk, vv := range listx {
					targetAuthorUID := listx[kk].Extra.GetStore(`targetObject`).Uint(`ownerId`)
					if vv.OwnerType == `customer` && vv.OwnerId == v.Id {
						listx[kk].Extra[`owner_id`] = v.Id
						listx[kk].Extra[`owner_type`] = `customer`
						listx[kk].Extra[`name`] = v.Name
						listx[kk].Extra[`gender`] = v.Gender
						listx[kk].Extra[`avatar`] = v.Avatar
						listx[kk].Extra[`agent_level`] = v.AgentLevel
						if targetAuthorUID > 0 {
							listx[kk].Extra[`isTargetAuthor`] = v.Uid > 0 && v.Uid == targetAuthorUID
						}
						listx[kk].Extra[`isAdmin`] = v.Uid > 0
					}
					if vv.ReplyOwnerId > 0 && vv.ReplyOwnerType == `customer` && vv.ReplyOwnerId == v.Id {
						listx[kk].Extra[`repliedCustomerName`] = v.Name
						if targetAuthorUID > 0 {
							listx[kk].Extra[`repliedCustomerIsTargetAuthor`] = v.Uid > 0 && v.Uid == targetAuthorUID
						}
						listx[kk].Extra[`repliedCustomerIsAdmin`] = v.Uid > 0
					}
				}
			}
		}
	}
	if len(userIds) > 0 {
		userM := model.NewUser(c)
		_, err = userM.ListByOffset(nil, nil, 0, -1, db.Cond{`id`: db.In(userIds)})
		if err == nil {
			for _, v := range userM.Objects() {
				for kk, vv := range listx {
					targetAuthorUID := listx[kk].Extra.GetStore(`targetObject`).Uint(`ownerId`)
					uid := uint64(v.Id)
					if vv.OwnerType == `user` && vv.OwnerId == uid {
						listx[kk].Extra[`owner_id`] = uid
						listx[kk].Extra[`owner_type`] = `user`
						listx[kk].Extra[`name`] = v.Username
						listx[kk].Extra[`gender`] = v.Gender
						listx[kk].Extra[`avatar`] = v.Avatar
						listx[kk].Extra[`agent_level`] = 0
						if targetAuthorUID > 0 {
							listx[kk].Extra[`isTargetAuthor`] = v.Id == targetAuthorUID
						}
						listx[kk].Extra[`isAdmin`] = true
					}
					if vv.ReplyOwnerId > 0 && vv.ReplyOwnerType == `user` && vv.ReplyOwnerId == uid {
						listx[kk].Extra[`repliedCustomerName`] = v.Username
						if targetAuthorUID > 0 {
							listx[kk].Extra[`repliedCustomerIsTargetAuthor`] = v.Id == targetAuthorUID
						}
						listx[kk].Extra[`repliedCustomerIsAdmin`] = true
					}
				}
			}
		}
		if customer != nil || user != nil {
			flowM := official.NewClickFlow(c)
			var clickedFlow map[uint64]*dbschema.OfficialCommonClickFlow
			if customer != nil {
				clickedFlow, err = flowM.ListByCustomerTargets(`comment`, commentIds, customer.Id)
			} else {
				clickedFlow, err = flowM.ListByAdminTargets(`comment`, commentIds, user.Id)
			}
			if err == nil {
				for kk, vv := range listx {
					cf, ok := clickedFlow[vv.Id]
					if ok {
						listx[kk].Extra[`clickFlow`] = cf
					}
				}
			}
		}
		for sourceTable, sourceList := range productIdOwnerIds {
			st, ok := official.SourceTables[sourceTable]
			if !ok || st == nil || st.QueryBought == nil {
				continue
			}
			for sourceID, owners := range sourceList {
				boughtCustomerIDs, err := st.QueryBought(c, sourceID, owners[`customer`])
				if err != nil {
					return listx, err
				}
				if len(boughtCustomerIDs) == 0 {
					continue
				}
				for kk, vv := range listx {
					if vv.OwnerType == `customer` {
						isBought := boughtCustomerIDs[vv.OwnerId]
						listx[kk].Extra[`isBought`] = isBought
					}
					if vv.ReplyOwnerId > 0 && vv.ReplyOwnerType == `customer` {
						isBought := boughtCustomerIDs[vv.ReplyOwnerId]
						listx[kk].Extra[`repliedCustomerIsBought`] = isBought
					}
				}
			}
		}
	}
	return listx, err
}

func (f *Comment) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.OfficialCommonComment.Delete(mw, args...)
}
