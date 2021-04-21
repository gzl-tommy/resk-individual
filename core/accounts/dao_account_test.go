package accounts

import (
	"database/sql"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"testing"
	"gzl-tommy/resk-individual/infra/base"
	_ "gzl-tommy/resk-individual/testx"
)

func TestAccountDao_GetOne(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{runner: runner}
		convey.Convey("通过编号查询账户数据", t, func() {
			a := &Account{
				Balance:     decimal.NewFromFloat(100),
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "测试资金账户",
				UserId:      ksuid.New().Next().String(),
				Username:    sql.NullString{String: "测试用户", Valid: true},
			}
			id, err := dao.Insert(a)
			convey.So(err, convey.ShouldBeNil)
			convey.So(id, convey.ShouldBeGreaterThan, 0)

			na := dao.GetOne(a.AccountNo)
			convey.So(na, convey.ShouldNotBeNil)
			convey.So(na.Balance.String(), convey.ShouldEqual, a.Balance.String())
			convey.So(na.CreatedAt, convey.ShouldNotBeNil)
			convey.So(na.UpdatedAt, convey.ShouldNotBeNil)
		})

		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}

func TestAccountDao_GetByUserId(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{runner: runner}
		convey.Convey("通过用户ID和账户类型查询账户数据", t, func() {
			a := &Account{
				Balance:     decimal.NewFromFloat(100),
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "测试资金账户",
				UserId:      ksuid.New().Next().String(),
				Username:    sql.NullString{String: "测试用户", Valid: true},
				AccountType: 2,
			}
			id, err := dao.Insert(a)
			convey.So(err, convey.ShouldBeNil)
			convey.So(id, convey.ShouldBeGreaterThan, 0)

			na := dao.GetByUserId(a.UserId, a.AccountType)
			convey.So(na, convey.ShouldNotBeNil)
			convey.So(na.Balance.String(), convey.ShouldEqual, a.Balance.String())
			convey.So(na.CreatedAt, convey.ShouldNotBeNil)
			convey.So(na.UpdatedAt, convey.ShouldNotBeNil)
		})

		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}

func TestAccountDao_UpdateBalance(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{runner: runner}
		balance := decimal.NewFromFloat(100)
		convey.Convey("更新账户余额", t, func() {
			a := &Account{
				Balance:     decimal.NewFromFloat(100),
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "测试资金账户",
				UserId:      ksuid.New().Next().String(),
				Username:    sql.NullString{String: "测试用户", Valid: true},
			}
			id, err := dao.Insert(a)
			convey.So(err, convey.ShouldBeNil)
			convey.So(id, convey.ShouldBeGreaterThan, 0)

			//1.增加余额
			convey.Convey("增加余额", func() {
				amount := decimal.NewFromFloat(10)
				rows, err := dao.UpdateBalance(a.AccountNo, amount)
				convey.So(err, convey.ShouldBeNil)
				convey.So(rows, convey.ShouldEqual, 1)
				na := dao.GetOne(a.AccountNo)
				newBalance := balance.Add(amount)
				convey.So(na, convey.ShouldNotBeNil)
				convey.So(na.Balance.String(), convey.ShouldEqual, newBalance.String())
			})

			//2.扣减余额，余额足够
			convey.Convey("扣减余额，余额足够", func() {
				amount := decimal.NewFromFloat(-10)
				rows, err := dao.UpdateBalance(a.AccountNo, amount)
				convey.So(err, convey.ShouldBeNil)
				convey.So(rows, convey.ShouldEqual, 1)
				na := dao.GetOne(a.AccountNo)
				newBalance := balance.Add(amount)
				convey.So(na, convey.ShouldNotBeNil)
				convey.So(na.Balance.String(), convey.ShouldEqual, newBalance.String())

			})

			//3.扣减余额，余额不够
			convey.Convey("扣减余额，余额不够", func() {
				a1 := dao.GetOne(a.AccountNo)
				convey.So(a1, convey.ShouldNotBeNil)
				amount := decimal.NewFromFloat(-300)
				rows, err := dao.UpdateBalance(a.AccountNo, amount)
				convey.So(err, convey.ShouldBeNil)
				convey.So(rows, convey.ShouldEqual, 0)
				a2 := dao.GetOne(a.AccountNo)
				convey.So(a2, convey.ShouldNotBeNil)
				convey.So(a1.Balance.String(), convey.ShouldEqual, a2.Balance.String())
			})
		})

		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}
