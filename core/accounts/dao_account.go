package accounts

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type AccountDao struct {
	runner *dbx.TxRunner
}

// 查询数据库持久化对象的单实例，获取一行数据
func (dao *AccountDao) GetOne(accountNo string) *Account {
	a := &Account{AccountNo: accountNo} // 这里的 AccountNo 字段必须是唯一索引
	ok, err := dao.runner.GetOne(a)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		logrus.Error("get Account not ok! ", accountNo)
		return nil
	}
	return a
}

// 通过用户ID和账户类型来查询账户信息
func (dao *AccountDao) GetByUserId(userId string, accountType int) *Account {
	a := &Account{}
	ok, err := dao.runner.Get(a, "select * from account where user_id=? and account_type=?", userId, accountType)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return a
}

// 账号数据的插入
func (dao *AccountDao) Insert(a *Account) (id int64, err error) {
	rs, err := dao.runner.Insert(a)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}

// 账户余额的更新
// amount 如果是负数，就是扣减；如果是正数，就是增加
func (dao *AccountDao) UpdateBalance(accountNo string, amount decimal.Decimal) (rows int64, err error) {
	rs, err := dao.runner.Exec(
		"update account set balance=balance+CAST(? AS DECIMAL(30,6)) where account_no=? and balance>=-1*CAST(? AS DECIMAL(30,6))",
		amount.String(), accountNo, amount.String())
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

// 账号状态更新
func (dao *AccountDao) UpdateStatus(accountNo string, status int) (rows int64, err error) {
	rs, err := dao.runner.Exec("update account set status=? where account_no=?", status, accountNo)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}
