package models

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/hublabs/common/auth"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/ctxbase"
	"github.com/pangpanglabs/goutils/jwtutil"
)

type FieldType string
type ConditionType string

const (
	FieldTypeChannel       FieldType = "channel"
	FieldTypePrepareCoupon FieldType = "prepareCoupon"
)

const (
	ConditionTypeMemberId      = "member_id"
	ConditionTypeGradeId       = "grade_id"
	ConditionTypeBirthdayMonth = "birthday_month"
	ConditionTypeStoreId       = "store_id"
	ConditionTypeBrandId       = "brand_id"
)

// Readonly
var customerConditionList = []struct {
	ConditionType string
}{
	{ConditionTypeMemberId},
	{ConditionTypeGradeId},
	{ConditionTypeBirthdayMonth},
	{ConditionTypeBrandId},
}

// Readonly
var comparerList = []struct {
	Comparer   ComparerType
	ExistsWord string
	Operator   string
}{
	{Comparer: ComparerTypeIn,
		ExistsWord: "EXISTS",
		Operator:   "="},
	{Comparer: ComparerTypeNin,
		ExistsWord: "NOT EXISTS",
		Operator:   "="},
}

type FieldTypeList []FieldType

func (l FieldTypeList) Contains(f ...FieldType) bool {
	for _, fi := range f {
		for _, li := range l {
			if fi == li {
				return true
			}
		}
	}
	return false
}

func getTenantCode(ctx context.Context) string {
	user := auth.UserClaim{}.FromCtx(ctx)
	return user.Issuer
}

func GetColleagueId(ctx context.Context) int64 {
	user := auth.UserClaim{}.FromCtx(ctx)
	return user.ColleagueId
}

func setSortOrder(q xorm.Interface, sortby, order []string, table ...string) error {
	connect := func(col string) string {
		if len(table) > 0 {
			return table[0] + "." + col
		}
		return col
	}

	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				v = connect(v)
				if order[i] == "desc" {
					q.Desc(v)
				} else if order[i] == "asc" {
					q.Asc(v)
				} else {
					return errors.New("Invalid order. Must be either [asc|desc]")
				}
			}
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				v = connect(v)
				if order[0] == "desc" {
					q.Desc(v)
				} else if order[0] == "asc" {
					q.Asc(v)
				} else {
					return errors.New("Invalid order. Must be either [asc|desc]")
				}
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return errors.New("'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return errors.New("unused 'order' fields")
		}
	}
	return nil
}

//获取传入的时间所在月份的第一天，即某月第一天的0点
func GetFirstDateOfMonth(d time.Time) time.Time {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return d
	}
	d = d.In(location)
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

//获取传入的时间所在月份的最后一天，即某月最后一天的0点
func GetLastDateOfMonth(d time.Time) time.Time {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return d
	}
	d = GetFirstDateOfMonth(d.In(location)).AddDate(0, 1, -1)
	return GetLastTime(d)
}

//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return d
	}
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, location)
}

//获取某一天的23点时间
func GetLastTime(d time.Time) time.Time {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return d
	}
	return time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, location)
}

//生成一个随机数
func random(num float64) float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rd := r.Float64()
	return rd * num
}

//查找随机数落在哪个分批内
func findIndex(num float64, list []float64) int {
	flag := float64(0)
	for i := range list {
		if num >= flag && num < flag+list[i] {
			return i
		}
		flag += list[i]
	}
	return -1
}

func AddToCtx(ctx context.Context, tenantCode string) context.Context {
	token, _ := jwtutil.NewToken(map[string]interface{}{"iss": "colleague", "tenantCode": tenantCode})
	behaviorLogContext := behaviorlog.LogContext{
		AuthToken: token,
		Service:   "coupon-api",
		ActionID:  ctxbase.NewID(),
		RequestID: ctxbase.NewID()}
	return behaviorLogContext.ToCtx(ctx)
}
