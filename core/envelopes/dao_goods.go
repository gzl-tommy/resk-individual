package envelopes

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"time"
	"gzl-tommy/resk-individual/services"
)

type RedEnvelopeGoodsDao struct {
	runner *dbx.TxRunner
}

// 插入
func (dao *RedEnvelopeGoodsDao) Insert(po *RedEnvelopeGoods) (int64, error) {
	rs, err := dao.runner.Insert(po)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	return rs.LastInsertId()
}

// 查询，根据红包编号
func (dao *RedEnvelopeGoodsDao) GetOne(envelopeNo string) *RedEnvelopeGoods {
	po := &RedEnvelopeGoods{EnvelopeNo: envelopeNo}
	ok, err := dao.runner.GetOne(po)
	if err != nil || !ok {
		logrus.Error(err)
		return nil
	}
	return po
}

// 更新红包余额和数量
// 不在使用事务行锁来更新红包余额和数量
// 使用乐观锁来保证更新红包余额和数量的安全，避免负库存问题
// 通过在where子句中判断红包剩余金额和数量来解决2个问题：
// 1：负库存问题，避免红包剩余金额和数量不够时进行扣减更新
// 2：减少实际数据更新，也就是过滤掉部分无效的更新，提高总体性能
func (dao *RedEnvelopeGoodsDao) UpdateBalance(envelopeNo string, amount decimal.Decimal) (int64, error) {
	sql := " update red_envelope_goods " +
		" set remain_amount=remain_amount-CAST(? AS DECIMAL(30,6)), " +
		" remain_quantity=remain_quantity-1 " +
		" where envelope_no=? " +
		// 最重要的，乐观锁的关键
		" and remain_quantity>0 " +
		" and remain_amount >= CAST(? AS DECIMAL(30,6))"
	rs, err := dao.runner.Exec(sql, amount.String(), envelopeNo, amount.String())
	if err != nil {
		logrus.Error(err)
		return 0, nil
	}
	return rs.RowsAffected()
}

// 跟新订单状态
func (dao *RedEnvelopeGoodsDao) UpdateOrderStatus(envelopeNo string, status services.OrderStatus) (int64, error) {
	sql := " update red_envelope_goods" +
		" set order_status=? " +
		" where envelope_no=?"
	rs, err := dao.runner.Exec(sql, status, envelopeNo)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

// 过期，把过期的所有红包都查询出来，分页，limit,offset siz
func (dao *RedEnvelopeGoodsDao) FindExpired(offset, size int) []RedEnvelopeGoods {
	var goods []RedEnvelopeGoods
	sql := " select * from red_envelope_goods " +
		" where expired_at=? " +
		" limit ?,?"
	err := dao.runner.Find(&goods, sql, time.Now(), offset, size)
	if err != nil {
		logrus.Error(err)
	}
	return goods
}

//
func (dao *RedEnvelopeGoodsDao) Find(po *RedEnvelopeGoods, offset, limit int) []RedEnvelopeGoods {
	var redEnvelopeGoodss []RedEnvelopeGoods
	err := dao.runner.FindExample(po, &redEnvelopeGoodss)
	if err != nil {
		logrus.Error(err)
	}
	return redEnvelopeGoodss
}

func (dao *RedEnvelopeGoodsDao) FindByUser(userId string, offset, limit int) []RedEnvelopeGoods {
	var goods []RedEnvelopeGoods

	sql := " select * from red_envelope_goods " +
		" where  user_id=?  order by created_at desc limit ?,?"
	err := dao.runner.Find(&goods, sql, userId, offset, limit)
	if err != nil {
		logrus.Error(err)
	}
	return goods
}

func (dao *RedEnvelopeGoodsDao) ListReceivable(offset, size int) []RedEnvelopeGoods {
	var goods []RedEnvelopeGoods
	now := time.Now()
	sql := " select * from red_envelope_goods " +
		" where  remain_quantity>0  and expired_at>? order by created_at desc limit ?,?"
	err := dao.runner.Find(&goods, sql, now, offset, size)
	if err != nil {
		logrus.Error(err)
	}
	return goods
}