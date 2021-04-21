package envelopes

import (
	"context"
	"database/sql"
	"errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"gzl-tommy/resk-individual/core/accounts"
	"gzl-tommy/resk-individual/infra/algo"
	"gzl-tommy/resk-individual/infra/base"
	"gzl-tommy/resk-individual/services"
)

var multiple = decimal.NewFromFloat(100.0)

//收红包业务逻辑代码
func (d *goodsDomain) Receive(ctx context.Context, dto services.RedEnvelopeReceiveDTO) (item *services.RedEnvelopeItemDTO, err error) {
	//1.创建收红包的订单明细 preCreatItem
	d.preCreateItem(dto)
	//2.查询出当前红包的剩余数量和剩余金额信息
	goods := d.Get(dto.EnvelopeNo)
	//3. 效验剩余红包和剩余金额：
	//- 如果没有剩余，直接返回无可用红包金额
	if goods.RemainQuantity <= 0 || goods.RemainAmount.Cmp(decimal.NewFromFloat(0)) <= 0 {
		logrus.Errorf("%+v", goods)
		return nil, errors.New("没有足够的红包和金额了")
	}
	//4. 使用红包算法计算红包金额
	nextAmount := d.nextAmount(goods)
	err = base.Tx(func(runner *dbx.TxRunner) error {
		//5. 使用乐观锁更新语句，尝试更新剩余数量和剩余金额：
		dao := RedEnvelopeGoodsDao{runner: runner}
		rows, err := dao.UpdateBalance(goods.EnvelopeNo, nextAmount)
		// - 如果更新失败，也就是返回0，表示无可用红包数量和金额，抢红包失败
		if rows <= 0 || err != nil {
			logrus.Errorf("rows=%d, %s", rows, err.Error())
			return errors.New("没有足够的红包和金额了")
		}

		// - 如果更新成功，也就是返回1，表示抢到红包
		//6. 保存订单明细数据
		d.item.Quantity = 1
		d.item.PayStatus = services.Paying
		d.item.AccountNo = dto.AccountNo
		d.item.RemainAmount = goods.RemainAmount.Sub(nextAmount)
		d.item.Amount = nextAmount
		txCtx := base.WithValueContext(ctx, runner)
		_, err = d.item.Save(txCtx)
		//insert into `red_envelope_item`(`item_no`,`envelope_no`,`amount`,`quantity`,`remain_amount`,`account_no`,`recv_username`,`recv_user_id`,`pay_status`)
		// values(?,?,?,?,?,?,?,?,?)
		// 1qsdPkkobtbkeZDojucTyYPxeXz
		// 1LlU0eDA6WNt16ZbLKp3RwpFWcp
		if err != nil {
			logrus.Error(err)
			return err
		}
		//7. 将抢到的红包金额从系统红包中间账户转入当前用户的资金账户
		// transfer
		status, err := d.transfer(txCtx, dto)
		if status == services.TransferedStatusSuccess {
			return nil
		}
		return err
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return d.item.ToDTO(), err
}

//预创建收红包订单明细
func (d *goodsDomain) preCreateItem(dto services.RedEnvelopeReceiveDTO) {
	d.item.AccountNo = dto.AccountNo
	d.item.EnvelopeNo = dto.EnvelopeNo
	d.item.RecvUsername = sql.NullString{String: dto.RecvUsername, Valid: true}
	d.item.RecvUserId = dto.RecvUserId
	d.item.createItemNo()
}

func (d *goodsDomain) transfer(ctx context.Context, dto services.RedEnvelopeReceiveDTO) (status services.TransferedStatus, err error) {
	systemAccount := base.GetSystemAccount()
	body := services.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		UserId:    systemAccount.UserId,
		Username:  systemAccount.Username,
	}
	target := services.TradeParticipator{
		AccountNo: dto.AccountNo,
		UserId:    dto.RecvUserId,
		Username:  dto.RecvUsername,
	}

	adomain := accounts.NewAccountDomain()
	//从系统红包资金账户扣减
	transfer := services.AccountTransferDTO{
		TradeBody:   body,
		TradeTarget: target,
		TradeNo:     dto.EnvelopeNo,
		Amount:      d.item.Amount,
		ChangeType:  services.EnvelopeOutgoing,
		ChangeFlag:  services.FlagTransferOut,
		Decs:        "红包扣减：" + dto.EnvelopeNo,
	}
	status, err = adomain.TransferWithContextTx(ctx, transfer)
	if err != nil || status != services.TransferedStatusSuccess {
		return status, err
	}

	//从系统红包资金账户转入当前用户
	transfer = services.AccountTransferDTO{
		TradeBody:   target,
		TradeTarget: body,
		TradeNo:     dto.EnvelopeNo,
		Amount:      d.item.Amount,
		ChangeType:  services.EnvelopeIncoming,
		ChangeFlag:  services.FlagTransferIn,
		Decs:        "红包收入：" + dto.EnvelopeNo,
	}
	return adomain.TransferWithContextTx(ctx, transfer)
}

//计算红包金额
func (d *goodsDomain) nextAmount(goods *RedEnvelopeGoods) (amount decimal.Decimal) {
	if goods.RemainQuantity == 1 {
		return goods.RemainAmount
	}
	if goods.EnvelopeType == services.GeneralEnvelopeType {
		return goods.AmountOne
	} else if goods.EnvelopeType == services.LuckyEnvelopeType {
		cent := goods.RemainAmount.Mul(multiple).IntPart()
		next := algo.DoubleAverage(int64(goods.RemainQuantity), cent)
		amount = decimal.NewFromFloat(float64(next)).Div(multiple)
	} else {
		logrus.Error("不支持的红包类型")
	}
	return amount
}
