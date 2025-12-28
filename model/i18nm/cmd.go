package i18nm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
)

// AutoTranslate 自动翻译指定表中的内容
//
// 参数:
//
//	ctx: echo上下文
//	table: 要翻译的表名
//	queryAll: 是否查询所有记录（忽略已翻译记录）
//	translateAll: 是否强制翻译所有内容
//	eqID: 指定ID等于该值的记录
//	gtID: 指定ID大于该值的记录
//	chunks: 每次查询的记录数限制
//
// 返回值:
//
//	uint64: 最后处理的记录ID
//	error: 错误信息
//
// 功能:
//  1. 检查表是否存在
//  2. 获取表中需要翻译的字段
//  3. 构建查询条件并执行查询
//  4. 对查询结果进行批量翻译
//  5. 支持内容类型(contype)的特殊处理
//  6. 支持分块查询和翻译
func AutoTranslate(ctx echo.Context, table string, queryAll bool, translateAll bool, eqID uint64, gtID uint64, chunks int) (uint64, error) {
	fields, exists := dbschema.DBI.Fields[table]
	if !exists {
		return gtID, fmt.Errorf(`Table %s does not exist in this system.`, table)
	}
	rIDs, rKeys, err := getResourceIDs(ctx, table)
	if err != nil {
		return gtID, err
	}
	if len(rIDs) == 0 {
		return gtID, fmt.Errorf(`Translation is not supported for the fields in Table %s.`, table)
	}
	_, hasContype := fields[`contype`]
	columns := make([]string, 0, len(rKeys))
	resourceIDsByField := map[string]uint{}
	contentTypes := map[string]string{}
	for resourceID, field := range rKeys {
		columns = append(columns, `AA.`+field)
		resourceIDsByField[field] = resourceID
		contentTypes[field] = `text`
	}
	var contentField string
	if hasContype {
		columns = append(columns, `AA.contype`)
		for _, field := range []string{`content`, `description`} {
			if _, ok := resourceIDsByField[field]; ok {
				contentField = field
				break
			}
		}
	}

	mdl := dbschema.DBI.NewModel(com.PascalCase(table))
	options := []func(*SaveModelTranslationsOptions){
		OptionAutoTranslate(true),
		OptionResourceIDsByField(resourceIDsByField),
		OptionContentTypes(contentTypes),
	}
	if translateAll {
		options = append(options,
			OptionForceTranslate(true),
			OptionAllowForceTranslate(func(ctx echo.Context) bool {
				return true
			}),
		)
	}
	queryDefault := `SELECT AA.id,` + strings.Join(columns, `,`) + ` FROM ` + dbschema.WithPrefix(table) + ` AS AA`
	var whereDefault string
	if !queryAll {
		whereDefault = `NOT EXISTS ( SELECT 1 FROM ` + dbschema.WithPrefix(`official_i18n_translation`) + ` AS TT WHERE TT.row_id=AA.id AND TT.resource_id IN (` + com.JoinNumbers(rIDs, `,`) + `) )`
	}
	sb := strings.Builder{}
	f := func() error {
		var err error
		where := make([]string, 0, 3)
		if eqID > 0 {
			where = append(where, `AA.id = `+strconv.FormatUint(eqID, 10))
		}
		if gtID > 0 {
			where = append(where, `AA.id > `+strconv.FormatUint(gtID, 10))
		}
		if len(whereDefault) > 0 {
			where = append(where, whereDefault)
		}
		sb.WriteString(queryDefault)
		if len(where) > 0 {
			sb.WriteString(` WHERE ` + strings.Join(where, ` AND `))
		}
		sb.WriteString(` ORDER BY AA.id ASC`)
		sb.WriteString(` LIMIT ` + strconv.Itoa(chunks))
		p := factory.ParamPoolGet()
		defer p.Release()
		rows, err := p.SetCollection(sb.String()).Query()
		sb.Reset()
		if err != nil {
			return err
		}
		defer rows.Close()
		var qColumns []string
		qColumns, err = rows.Columns()
		if err != nil {
			return err
		}
		for rows.Next() {
			row := make([]interface{}, len(qColumns))
			for i := range row {
				row[i] = new(interface{})
			}
			err = rows.Scan(row...)
			if err != nil {
				return err
			}
			data := map[string]interface{}{}
			for i, v := range row {
				field := qColumns[i]
				data[field] = v
			}
			mdl.FromRow(data)
			id := GetRowID(mdl)
			if id > 0 {
				gtID = id
				if hasContype {
					contype, ok := mdl.GetField(`Contype`).(string)
					if ok && contentField != `` && contype != `` {
						err = SaveModelTranslations(ctx, mdl, id, append([]func(*SaveModelTranslationsOptions){OptionContentType(contentField, contype)}, options...)...)
						if err != nil {
							return err
						}
					}
				}
				err = SaveModelTranslations(ctx, mdl, id, options...)
				if err != nil {
					return err
				}
			}
			// Reset the data to avoid re-saving the same data
			for field := range data {
				data[field] = nil
			}
			mdl.FromRow(data)
		}
		return err
	}
	for {
		err = f()
		if err != nil {
			break
		}
	}
	return gtID, err
}
