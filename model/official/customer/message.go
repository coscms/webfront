package customer

import (
	"strings"
	"time"

	"github.com/admpub/events"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/param"

	dbschemaNging "github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webcore/model"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/top"
)

/**
 * 1. 支持前台用户向前台用户发送消息
 * 2. 支持后台用户向前台用户发送消息
 * 3. 支持后台用户向后台用户发送消息
 * 4. 支持前台用户向后台用户发送消息
 */

// NewMessage 站内信
func NewMessage(ctx echo.Context) *Message {
	m := &Message{
		OfficialCommonMessage: dbschema.NewOfficialCommonMessage(ctx),
		Viewed:                dbschema.NewOfficialCommonMessageViewed(ctx),
	}
	return m
}

type MsgUser struct {
	Name   string
	Id     uint64
	Avatar string
	Type   string //user-后台用户; customer-前台用户
}

type MessageWithViewed struct {
	*dbschema.OfficialCommonMessage
	MsgFrom  *MsgUser `json:",omitempty"`
	MsgTo    *MsgUser `json:",omitempty"`
	IsViewed bool
}

func (m *MessageWithViewed) MsgUser() *MsgUser {
	if m.MsgFrom != nil {
		return m.MsgFrom
	}
	return m.MsgTo
}

type Message struct {
	*dbschema.OfficialCommonMessage
	Viewed *dbschema.OfficialCommonMessageViewed
}

// Delete 循环删除
func (f *Message) Delete(args ...interface{}) error {
	cnt, err := f.OfficialCommonMessage.ListByOffset(nil, nil, 0, 100, args...)
	if err != nil {
		return err
	}
	for cnt() > 0 {
		for _, row := range f.Objects() {
			err = f.OfficialCommonMessage.Delete(nil, `id`, row.Id)
			if err != nil {
				return err
			}
			err = f.Viewed.Delete(nil, `message_id`, row.Id)
			if err != nil {
				return err
			}
		}
		cnt, err = f.OfficialCommonMessage.ListByOffset(nil, nil, 0, 100, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Message) DecodeContent(row *dbschema.OfficialCommonMessage) *dbschema.OfficialCommonMessage {
	if row.Encrypted == `Y` {
		row.Content = config.FromFile().Decode(row.Content, row.Password)
	}
	return row
}

func (f *Message) check() error {
	if len(f.Content) == 0 {
		return f.Context().E(`内容不能为空`)
	}
	if len(f.Encrypted) == 0 {
		f.Encrypted = common.BoolN
	}
	if f.HasNewReply > 1 {
		f.HasNewReply = 1
	}
	if f.Encrypted == common.BoolY {
		f.Content = config.FromFile().Encode(f.Content, f.Password)
	} else {
		if len(f.Title) == 0 {
			f.Title = com.Substr(com.StripTags(f.Content), `...`, 50)
		}
	}
	f.Title = common.MyCleanText(f.Title)
	if len(f.Contype) == 0 {
		f.Contype = common.ContentTypeText
	}
	if f.CustomerA > 0 {
		f.UserA = 0
	}
	// 允许关联了后台账号的客户的消息在前后台均可以查看
	// if f.CustomerB > 0 {
	// 	f.UserB = 0
	// }
	if f.CustomerGroupId > 0 {
		f.UserRoleId = 0
		f.CustomerB = 0
		f.UserB = 0
	} else if f.UserRoleId > 0 {
		f.CustomerB = 0
		f.UserB = 0
	}
	if f.CustomerA > 0 {
		maxPerDay, interval := f.CustomerMaxPerDay()
		if maxPerDay <= 0 {
			return f.Context().E(`很抱歉！系统已关闭私信功能`)
		}
		todayCount := f.CountTodaySends(f.CustomerA)
		if todayCount >= maxPerDay {
			return f.Context().E(`很抱歉！您今日的发送数量已达上限: %d`, maxPerDay)
		}
		if interval > 0 {
			lastSend, err := f.LastSend(f.CustomerA)
			if err != nil {
				if err != db.ErrNoMoreRows {
					return err
				}
				err = nil
			}
			if lastSend.Created > 0 {
				duration := time.Now().Unix() - int64(lastSend.Created) //间隔时间
				if duration < interval {
					return f.Context().E(`很抱歉！您发送的太快，请等待%d秒之后再发送`, interval-duration)
				}
			}
		}
	}
	return nil
}

func (f *Message) CustomerMaxPerDay() (int64, int64) {
	frequencyCfg := config.Setting().GetStore(`frequency`).GetStore(`message`)
	maxPerDay := frequencyCfg.Int64(`maxPerDay`)
	interval := frequencyCfg.Int64(`interval`)
	return maxPerDay, interval
}

func (f *Message) CountTodaySends(customerID uint64) int64 {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 0, 1)
	n, _ := f.Count(nil, db.And(
		db.Cond{`customer_a`: customerID},
		//db.Cond{`reply_id`: 0},//只统计非回复
		db.Cond{`created`: db.Between(start.Unix(), end.Unix())},
	))
	return n
}

func (f *Message) LastSend(customerID uint64) (*dbschema.OfficialCommonMessage, error) {
	m := dbschema.NewOfficialCommonMessage(f.Context())
	err := m.Get(func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, db.And(
		db.Cond{`customer_a`: customerID},
		//db.Cond{`reply_id`: 0},
	))
	return m, err
}

func (f *Message) SendSystemMessage(customerID uint64, title string, contype string, content string) error {
	m := f.OfficialCommonMessage
	m.Title = title
	m.Content = content
	m.Contype = common.GetContype(contype, common.ContentTypeText)
	m.Encrypted = common.BoolN
	m.CustomerB = customerID
	if m.CustomerB < 1 {
		return f.Context().NewError(code.InvalidParameter, `请选择收信人`).SetZone(`customerId`)
	}
	_, err := f.AddData(nil, nil)
	return err
}

// AddData 添加消息
// * customer 操作客户
// * user 操作后台用户
func (f *Message) AddData(fromCustomer *dbschema.OfficialCustomer, fromUser *dbschemaNging.NgingUser) (pk interface{}, err error) {
	if err = f.check(); err != nil {
		return
	}
	ctx := f.Context()
	ctx.Begin()
	if f.ReplyId > 0 {
		msgM := dbschema.NewOfficialCommonMessage(ctx)
		err = msgM.Get(nil, `id`, f.ReplyId)
		if err != nil {
			ctx.Rollback()
			if err == db.ErrNoMoreRows {
				return nil, ctx.E(`您要回复的消息不存在`)
			}
			return
		}
		if fromCustomer != nil {
			if !((msgM.CustomerB > 0 && fromCustomer.Id == msgM.CustomerB) ||
				(msgM.CustomerA > 0 && fromCustomer.Id == msgM.CustomerA) ||
				(msgM.CustomerGroupId > 0 && fromCustomer.GroupId == msgM.CustomerGroupId)) {
				ctx.Rollback()
				return nil, ctx.E(`您无权回复此消息`)
			}
		} else if fromUser != nil {
			if !((msgM.UserB > 0 && fromUser.Id == msgM.UserB) ||
				(msgM.UserA > 0 && fromUser.Id == msgM.UserA) ||
				(msgM.UserRoleId > 0 && com.InSlice(param.AsString(msgM.UserRoleId), strings.Split(fromUser.RoleIds, `,`)))) {
				ctx.Rollback()
				return nil, ctx.E(`您无权回复此消息`)
			}
		} else {
			ctx.Rollback()
			return nil, nerrors.ErrUserNotLoggedIn
		}
		if msgM.RootId > 0 {
			f.RootId = msgM.RootId
		} else {
			f.RootId = msgM.Id
		}
		if !(f.CustomerA == f.CustomerB || f.UserA == f.UserB) {
			err = msgM.UpdateField(nil, `has_new_reply`, 1, db.Cond{`id`: f.RootId})
			if err != nil {
				ctx.Rollback()
				return
			}
		}
	}
	pk, err = f.OfficialCommonMessage.Insert()
	if err != nil {
		ctx.Rollback()
		return
	}
	echo.FireByNameWithMap(`message.send`, events.Map{`data`: f.OfficialCommonMessage, `fromCustomer`: fromCustomer, `fromUser`: fromUser})
	ctx.Commit()
	return
}

func (f *Message) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.OfficialCommonMessage.Update(mw, args...)
}

func (f *Message) CountUnread(viewerID uint64, groupIDs []uint, isSystemMessage bool, viewerTypes ...string) int64 {
	viewerType := `customer`
	if len(viewerTypes) > 0 {
		viewerType = viewerTypes[0]
	}
	cond := f.makeCond(viewerType, viewerID, groupIDs, isSystemMessage)
	cond = append(cond, db.Or(
		db.Cond{`has_new_reply`: 1},
		db.Raw(`NOT EXISTS (SELECT 1 FROM `+f.ToTable(f.Viewed)+` b WHERE b.message_id=`+f.ToTable(f)+`.id AND b.viewer_id=? AND b.viewer_type=?)`, viewerID, viewerType),
	))
	n, _ := f.Count(nil, db.And(cond...))
	return n
}

func (f *Message) LastUnread(viewerID uint64, groupIDs []uint, isSystemMessage bool, viewerTypes ...string) error {
	viewerType := `customer`
	if len(viewerTypes) > 0 {
		viewerType = viewerTypes[0]
	}
	cond := f.makeCond(viewerType, viewerID, groupIDs, isSystemMessage)
	cond = append(cond, db.Or(
		db.Cond{`has_new_reply`: 1},
		db.Raw(`NOT EXISTS (SELECT 1 FROM `+f.ToTable(f.Viewed)+` b WHERE b.message_id=`+f.ToTable(f)+`.id AND b.viewer_id=? AND b.viewer_type=?)`, viewerID, viewerType),
	))
	err := f.Get(func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, db.And(cond...))
	return err
}

func (f *Message) ToTable(m sqlbuilder.Name_) string {
	return config.FromFile().DB.ToTable(m)
}

func (f *Message) IsSystemMessage() bool {
	return f.UserA == 0 && f.CustomerA == 0
}

func (f *Message) makeCond(viewerType string, viewerID uint64, groupIDs []uint, isSystemMessage bool) []db.Compound {
	cond := []db.Compound{
		//db.Cond{`reply_id`: 0},
	}
	if len(groupIDs) > 0 {
		if viewerType == `customer` {
			cond = append(cond, db.Or(
				db.Cond{viewerType + `_b`: viewerID},
				db.Cond{`customer_group_id`: db.In(groupIDs)},
			))
		} else {
			cond = append(cond, db.Or(
				db.Cond{viewerType + `_b`: viewerID},
				db.Cond{`user_role_id`: db.In(groupIDs)},
			))
		}
	} else {
		cond = append(cond, db.Cond{viewerType + `_b`: viewerID})
	}
	if isSystemMessage {
		cond = append(cond, db.Cond{`customer_a`: 0})
		cond = append(cond, db.Cond{`user_a`: 0})
	} else {
		cond = append(cond, db.Or(
			db.Cond{`customer_a`: db.NotEq(0)},
			db.Cond{`user_a`: db.NotEq(0)},
		))
	}
	return cond
}

// ListWithViewedByRecipient 收件箱列表
func (f *Message) ListWithViewedByRecipient(viewerID uint64, groupIDs []uint, isSystemMessage bool, onlyUnread bool, otherCond db.Compound, viewerTypes ...string) ([]*MessageWithViewed, error) {
	viewerType := `customer`
	if len(viewerTypes) > 0 {
		viewerType = viewerTypes[0]
	}
	cond := f.makeCond(viewerType, viewerID, groupIDs, isSystemMessage)
	if onlyUnread {
		cond = append(cond,
			db.Or(
				db.Cond{`has_new_reply`: 1},
				db.Raw(`NOT EXISTS (SELECT 1 FROM `+f.ToTable(f.Viewed)+` b WHERE b.message_id=`+f.ToTable(f)+`.id AND b.viewer_id=? AND b.viewer_type=?)`, viewerID, viewerType),
			),
		)
	}
	if otherCond != nil {
		cond = append(cond, otherCond)
	}
	sorts := common.Sorts(f.Context(), `official_message`)
	if len(sorts) == 0 {
		sorts = []interface{}{
			`-has_new_reply`,
			db.Raw(`IF(view_progress=100,1,0)`),
		}
	}
	sorts = append(sorts, `-id`)
	_, err := common.NewLister(f.OfficialCommonMessage, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, db.And(cond...)).Paging(f.Context())
	if err != nil {
		return nil, err
	}

	rows := f.Objects()
	return f.WithViewedByRecipient(rows, viewerID, viewerType)
}

func (f *Message) WithViewedByRecipient(rows []*dbschema.OfficialCommonMessage, viewerID uint64, viewerType string) ([]*MessageWithViewed, error) {
	var err error
	list := make([]*MessageWithViewed, len(rows))
	msgIDs := make([]uint64, len(rows))
	cstIDs := []uint64{}
	usrIDs := []uint64{}
	for idx, row := range rows {
		row.Password = ``
		if row.Encrypted == `Y` {
			row.Content = `[本消息内容已加密]`
			if len(row.Title) == 0 {
				row.Title = `[这是一条加密消息]`
			}
		}
		item := &MessageWithViewed{
			OfficialCommonMessage: row,
			MsgFrom:               &MsgUser{},
			IsViewed:              false,
		}
		//有新回复的情况，始终标为未读
		if row.HasNewReply == 0 { //没有新回复的情况才查询我是否已读
			msgIDs[idx] = row.Id
		}
		if row.CustomerA > 0 {
			if !com.InUint64Slice(row.CustomerA, cstIDs) {
				cstIDs = append(cstIDs, row.CustomerA)
			}
			item.MsgFrom.Type = `customer`
			item.MsgFrom.Id = row.CustomerA
		} else if row.UserA > 0 {
			uid := uint64(row.UserA)
			if !com.InUint64Slice(uid, usrIDs) {
				usrIDs = append(usrIDs, uid)
			}
			item.MsgFrom.Type = `user`
			item.MsgFrom.Id = uid
		}
		list[idx] = item
	}
	if len(cstIDs) > 0 {
		customerM := NewCustomer(f.Context())
		_, err = customerM.ListByOffset(nil, nil, 0, -1, db.Cond{`id`: db.In(cstIDs)})
		if err != nil {
			return list, err
		}
		for _, customer := range customerM.Objects() {
			for idx, row := range list {
				if row.CustomerA == customer.Id {
					list[idx].MsgFrom.Name = customer.Name
					list[idx].MsgFrom.Avatar = customer.Avatar
				}
			}
		}
	}
	if len(usrIDs) > 0 {
		userM := model.NewUser(f.Context())
		_, err = userM.ListByOffset(nil, nil, 0, -1, db.Cond{`id`: db.In(usrIDs)})
		if err != nil {
			return list, err
		}
		for _, user := range userM.Objects() {
			for idx, row := range list {
				if row.UserA == user.Id {
					list[idx].MsgFrom.Name = user.Username
					list[idx].MsgFrom.Avatar = user.Avatar
				}
			}
		}
	}
	if len(msgIDs) > 0 {
		_, err = f.Viewed.ListByOffset(nil, nil, 0, -1, db.And(
			db.Cond{`message_id`: db.In(msgIDs)},
			db.Cond{`viewer_id`: viewerID},
			db.Cond{`viewer_type`: viewerType},
		))
		if err != nil {
			return list, err
		}
		for _, val := range f.Viewed.Objects() {
			for idx, row := range list {
				if row.Id == val.MessageId {
					list[idx].IsViewed = true
					break
				}
			}
		}
	}
	return list, err
}

// ListWithViewedBySender 发件箱列表
func (f *Message) ListWithViewedBySender(senderID uint64, onlyUnread bool, otherCond db.Compound, viewerTypes ...string) ([]*MessageWithViewed, error) {
	viewerType := `customer`
	if len(viewerTypes) > 0 {
		viewerType = viewerTypes[0]
	}
	cond := []db.Compound{
		//db.Cond{`reply_id`: 0},
		db.Cond{viewerType + `_a`: senderID},
	}
	if onlyUnread {
		cond = append(cond, db.Cond{`view_progress`: db.NotEq(100)})
	}
	if otherCond != nil {
		cond = append(cond, otherCond)
	}
	_, err := common.NewLister(f.OfficialCommonMessage, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, db.And(cond...)).Paging(f.Context())
	if err != nil {
		return nil, err
	}

	rows := f.Objects()
	return f.WithViewedBySender(rows)
}

func (f *Message) WithViewedBySender(rows []*dbschema.OfficialCommonMessage) ([]*MessageWithViewed, error) {
	var err error
	list := make([]*MessageWithViewed, len(rows))
	msgIDs := make([]uint64, len(rows))
	cstIDs := []uint64{}
	usrIDs := []uint64{}
	for idx, row := range rows {
		row.Password = ``
		if row.Encrypted == `Y` {
			row.Content = `[本消息内容已加密]`
			if len(row.Title) == 0 {
				row.Title = `[这是一条加密消息]`
			}
		}
		item := &MessageWithViewed{
			OfficialCommonMessage: row,
			MsgTo:                 &MsgUser{},
			IsViewed:              row.ViewProgress == 100,
		}
		msgIDs[idx] = row.Id
		if row.CustomerB > 0 {
			if !com.InUint64Slice(row.CustomerB, cstIDs) {
				cstIDs = append(cstIDs, row.CustomerB)
			}
			item.MsgTo.Type = `customer`
			item.MsgTo.Id = row.CustomerB
		} else if row.UserB > 0 {
			uid := uint64(row.UserB)
			if !com.InUint64Slice(uid, usrIDs) {
				usrIDs = append(usrIDs, uid)
			}
			item.MsgTo.Type = `user`
			item.MsgTo.Id = uid
		}
		list[idx] = item
	}
	if len(cstIDs) > 0 {
		customerM := NewCustomer(f.Context())
		_, err = customerM.ListByOffset(nil, nil, 0, -1, db.Cond{`id`: db.In(cstIDs)})
		if err != nil {
			return list, err
		}
		for _, customer := range customerM.Objects() {
			for idx, row := range list {
				if row.CustomerB == customer.Id {
					list[idx].MsgTo.Name = customer.Name
					list[idx].MsgTo.Avatar = customer.Avatar
				}
			}
		}
	}
	if len(usrIDs) > 0 {
		userM := model.NewUser(f.Context())
		_, err = userM.ListByOffset(nil, nil, 0, -1, db.Cond{`id`: db.In(usrIDs)})
		if err != nil {
			return list, err
		}
		for _, user := range userM.Objects() {
			for idx, row := range list {
				if row.UserB == user.Id {
					list[idx].MsgTo.Name = user.Username
					list[idx].MsgTo.Avatar = user.Avatar
				}
			}
		}
	}
	return list, err
}

// ListAll 列出所有人的消息
func (f *Message) ListAll(onlyUnread bool, otherCond db.Compound) ([]*MessageWithViewed, error) {
	cond := []db.Compound{}
	if onlyUnread {
		cond = append(cond, db.Cond{`view_progress`: db.NotEq(100)})
	}
	if otherCond != nil {
		cond = append(cond, otherCond)
	}
	_, err := common.NewLister(f.OfficialCommonMessage, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, db.And(cond...)).Paging(f.Context())
	if err != nil {
		return nil, err
	}

	rows := f.Objects()
	return f.WithViewedByAll(rows)
}

func (f *Message) WithViewedByAll(rows []*dbschema.OfficialCommonMessage) ([]*MessageWithViewed, error) {
	var err error
	list := make([]*MessageWithViewed, len(rows))
	msgIDs := make([]uint64, len(rows))
	cstIDs := []uint64{}
	usrIDs := []uint64{}
	for idx, row := range rows {
		row.Password = ``
		if row.Encrypted == `Y` {
			row.Content = `[本消息内容已加密]`
			if len(row.Title) == 0 {
				row.Title = `[这是一条加密消息]`
			}
		}
		item := &MessageWithViewed{
			OfficialCommonMessage: row,
			MsgFrom:               &MsgUser{},
			MsgTo:                 &MsgUser{},
			IsViewed:              row.ViewProgress == 100,
		}
		msgIDs[idx] = row.Id
		if row.CustomerA > 0 {
			if !com.InUint64Slice(row.CustomerA, cstIDs) {
				cstIDs = append(cstIDs, row.CustomerA)
			}
			item.MsgFrom.Type = `customer`
			item.MsgFrom.Id = row.CustomerA
		} else if row.UserA > 0 {
			uid := uint64(row.UserA)
			if !com.InUint64Slice(uid, usrIDs) {
				usrIDs = append(usrIDs, uid)
			}
			item.MsgFrom.Type = `user`
			item.MsgFrom.Id = uid
		}
		if row.CustomerB > 0 {
			if !com.InUint64Slice(row.CustomerB, cstIDs) {
				cstIDs = append(cstIDs, row.CustomerB)
			}
			item.MsgTo.Type = `customer`
			item.MsgTo.Id = row.CustomerB
		} else if row.UserB > 0 {
			uid := uint64(row.UserB)
			if !com.InUint64Slice(uid, usrIDs) {
				usrIDs = append(usrIDs, uid)
			}
			item.MsgTo.Type = `user`
			item.MsgTo.Id = uid
		}
		list[idx] = item
	}
	if len(cstIDs) > 0 {
		customerM := NewCustomer(f.Context())
		_, err = customerM.ListByOffset(nil, nil, 0, -1, db.Cond{`id`: db.In(cstIDs)})
		if err != nil {
			return list, err
		}
		for _, customer := range customerM.Objects() {
			for idx, row := range list {
				if row.CustomerB == customer.Id {
					list[idx].MsgTo.Name = customer.Name
					list[idx].MsgTo.Avatar = customer.Avatar
				}
				if row.CustomerA == customer.Id {
					list[idx].MsgFrom.Name = customer.Name
					list[idx].MsgFrom.Avatar = customer.Avatar
				}
			}
		}
	}
	if len(usrIDs) > 0 {
		userM := model.NewUser(f.Context())
		_, err = userM.ListByOffset(nil, nil, 0, -1, db.Cond{`id`: db.In(usrIDs)})
		if err != nil {
			return list, err
		}
		for _, user := range userM.Objects() {
			for idx, row := range list {
				if row.UserB == user.Id {
					list[idx].MsgTo.Name = user.Username
					list[idx].MsgTo.Avatar = user.Avatar
				}
				if row.UserA == user.Id {
					list[idx].MsgFrom.Name = user.Username
					list[idx].MsgFrom.Avatar = user.Avatar
				}
			}
		}
	}
	return list, err
}

func (f *Message) GetWithViewed(row *dbschema.OfficialCommonMessage) (*MessageWithViewed, error) {
	var err error
	cstIDs := []uint64{}
	usrIDs := []uint64{}
	row.Password = ``
	if row.Encrypted == `Y` {
		row.Content = `[本消息内容已加密]`
		if len(row.Title) == 0 {
			row.Title = `[这是一条加密消息]`
		}
	}
	item := &MessageWithViewed{
		OfficialCommonMessage: row,
		MsgFrom:               &MsgUser{},
		MsgTo:                 &MsgUser{},
		IsViewed:              row.ViewProgress == 100,
	}
	if row.CustomerA > 0 {
		if !com.InUint64Slice(row.CustomerA, cstIDs) {
			cstIDs = append(cstIDs, row.CustomerA)
		}
		item.MsgFrom.Type = `customer`
		item.MsgFrom.Id = row.CustomerA
	} else if row.UserA > 0 {
		uid := uint64(row.UserA)
		if !com.InUint64Slice(uid, usrIDs) {
			usrIDs = append(usrIDs, uid)
		}
		item.MsgFrom.Type = `user`
		item.MsgFrom.Id = uid
	}
	if row.CustomerB > 0 {
		if !com.InUint64Slice(row.CustomerB, cstIDs) {
			cstIDs = append(cstIDs, row.CustomerB)
		}
		item.MsgTo.Type = `customer`
		item.MsgTo.Id = row.CustomerB
	} else if row.UserB > 0 {
		uid := uint64(row.UserB)
		if !com.InUint64Slice(uid, usrIDs) {
			usrIDs = append(usrIDs, uid)
		}
		item.MsgTo.Type = `user`
		item.MsgTo.Id = uid
	}
	if len(cstIDs) > 0 {
		customerM := NewCustomer(f.Context())
		_, err = customerM.ListByOffset(nil, nil, 0, -1, db.Cond{`id`: db.In(cstIDs)})
		if err != nil {
			return item, err
		}
		for _, customer := range customerM.Objects() {
			if item.CustomerB == customer.Id {
				item.MsgTo.Name = customer.Name
				item.MsgTo.Avatar = customer.Avatar
			}
			if item.CustomerA == customer.Id {
				item.MsgFrom.Name = customer.Name
				item.MsgFrom.Avatar = customer.Avatar
			}
		}
	}
	if len(usrIDs) > 0 {
		userM := model.NewUser(f.Context())
		_, err = userM.ListByOffset(nil, nil, 0, -1, db.Cond{`id`: db.In(usrIDs)})
		if err != nil {
			return item, err
		}
		for _, user := range userM.Objects() {
			if item.UserB == user.Id {
				item.MsgTo.Name = user.Username
				item.MsgTo.Avatar = user.Avatar
			}
			if item.UserA == user.Id {
				item.MsgFrom.Name = user.Username
				item.MsgFrom.Avatar = user.Avatar
			}
		}
	}
	return item, err
}

func (f *Message) View(row *dbschema.OfficialCommonMessage, viewerID uint64, groupIDs []uint, viewerTypes ...string) error {
	viewerType := `customer`
	if len(viewerTypes) > 0 {
		viewerType = viewerTypes[0]
	}
	switch viewerType {
	case `customer`:
		if row.CustomerB > 0 {
			if row.CustomerB != viewerID {
				return nil
			}
		} else if row.CustomerGroupId > 0 {
			if !com.InUintSlice(row.CustomerGroupId, groupIDs) {
				return nil
			}
		}

	case `user`:
		if row.UserB > 0 {
			if uint64(row.UserB) != viewerID {
				return nil
			}
		} else if row.UserRoleId > 0 {
			if com.InUintSlice(row.UserRoleId, groupIDs) {
				return nil
			}
		}

	default:
		return nil
	}
	err := f.Viewed.Get(nil, db.And(
		db.Cond{`message_id`: row.Id},
		db.Cond{`viewer_id`: viewerID},
		db.Cond{`viewer_type`: viewerType},
	))
	if err == nil {
		return f.setViewProgress(row)
	}
	if err != db.ErrNoMoreRows {
		return err
	}
	f.Viewed.Reset()
	f.Viewed.MessageId = row.Id
	f.Viewed.ViewerId = viewerID
	f.Viewed.ViewerType = viewerType
	_, err = f.Viewed.Insert()
	if err != nil {
		return err
	}
	return f.setViewProgress(row)
}

func (f *Message) setViewProgress(row *dbschema.OfficialCommonMessage) (err error) {
	if row.ViewProgress == 100 && row.HasNewReply == 0 {
		return
	}
	var viewProgress int
	if row.CustomerB > 0 || row.UserB > 0 { //发送给单个人
		viewProgress = 100
	} else { //发送给某个群体
		// 查询已浏览人数
		viewers, err := f.Viewed.Count(nil, db.Cond{`message_id`: row.Id})
		if err != nil {
			return err
		}
		if row.CustomerGroupId > 0 {
			customerM := NewCustomer(f.Context())
			total, err := customerM.Count(nil, db.Cond{`group_id`: row.CustomerGroupId})
			if err != nil {
				return err
			}
			viewProgress = param.AsInt(tplfunc.NumberTrim(float64(viewers)/float64(total), 2)) * 100
		} else if row.UserRoleId > 0 {
			userM := model.NewUser(f.Context())
			total, err := userM.Count(nil, top.CondFindInSet(`role_ids`, row.UserRoleId))
			if err != nil {
				return err
			}
			viewProgress = param.AsInt(tplfunc.NumberTrim(float64(viewers)/float64(total), 2)) * 100
		}
		if viewProgress > 100 {
			viewProgress = 100
		}
	}

	err = row.UpdateFields(nil, echo.H{
		`view_progress`: viewProgress,
		`has_new_reply`: 0,
	}, db.Cond{`id`: row.Id})
	return
}

func (f *Message) MsgUser() *MsgUser {
	msgUser := &MsgUser{}
	if f.CustomerA > 0 {
		customerM := NewCustomer(f.Context())
		customerM.Get(nil, `id`, f.CustomerA)
		msgUser.Id = f.CustomerA
		msgUser.Name = customerM.Name
		msgUser.Type = `customer`
		msgUser.Avatar = customerM.Avatar
	} else if f.UserA > 0 {
		userM := model.NewUser(f.Context())
		userM.Get(nil, `id`, f.UserA)
		msgUser.Id = uint64(f.UserA)
		msgUser.Name = userM.Username
		msgUser.Type = `user`
		msgUser.Avatar = userM.Avatar
	}
	return msgUser
}

// CheckRecvPerm 检查是否为收信人权限
func (f *Message) CheckRecvPerm(customer *dbschema.OfficialCustomer) bool {
	return f.CheckRecvCustomerPerm(customer) || f.CheckRecvUserPerm(customer)
}

func (f *Message) CheckRecvCustomerPerm(customer *dbschema.OfficialCustomer) bool {
	if f.CustomerGroupId == 0 {
		return f.CustomerB == customer.Id
	}
	return f.CustomerGroupId == customer.GroupId || f.CustomerB == customer.Id
}

func (f *Message) CheckRecvUserPerm(customer *dbschema.OfficialCustomer) bool {
	if customer.Uid == 0 {
		return false
	}
	if f.UserRoleId == 0 {
		return f.UserB == customer.Uid
	}
	userM := &dbschemaNging.NgingUser{}
	err := userM.Get(nil, `id`, customer.Uid)
	if err != nil {
		return f.UserB == customer.Uid
	}
	roleIDs := param.Split(userM.RoleIds, `,`).Uint(func(idx int, val uint) bool {
		return val > 0
	})
	return f.UserA == customer.Uid || com.InUintSlice(f.UserRoleId, roleIDs)
}
