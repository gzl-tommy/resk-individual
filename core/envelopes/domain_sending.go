package envelopes

import (
	"context"
	"github.com/tietang/dbx"
	"path"
	"gzl-tommy/resk-individual/core/accounts"
	"gzl-tommy/resk-individual/infra/base"
	"gzl-tommy/resk-individual/services"
)

// 发红包业务领域代码
func (d *goodsDomain) SendOut(goods services.RedEnvelopeGoodsDTO) (activity *services.RedEnvelopeActivity, err error) {
	// 创建红包商品
	d.Create(goods)

	// 创建活动
	activity = new(services.RedEnvelopeActivity)

	// 红包连接，格式：http://域名/v1/envelope/{id}/link/
	link := base.GetEnvelopeActivityLink()
	domain := base.GetEnvelopeDomain()
	activity.Link = path.Join(domain, link, d.EnvelopeNo)
	accountDomain := accounts.NewAccountDomain()

	err = base.Tx(func(runner *dbx.TxRunner) error {
		// 事务逻辑问题：
		// 保存红包商品和红包金额的支付必须保证全部成功或者全部失败

		// 保存红包商品
		ctx := base.WithValueContext(context.Background(), runner)
		id, err := d.Save(ctx)
		if id <= 0 || err != nil {
			return err
		}

		// 红包金额支付
		// 1.需要红包中间商的红包资金账户，定义在配置文件中，事先初始化到资金账户表中
		// 2.从红包发送人的资金账户中扣减红包金额，把红包金额从红包发送人的资金账户里扣除
		body := services.TradeParticipator{
			AccountNo: goods.AccountNo,
			UserId:    goods.UserId,
			Username:  goods.Username,
		}
		systemAccount := base.GetSystemAccount()
		target := services.TradeParticipator{
			AccountNo: systemAccount.AccountNo,
			Username:  systemAccount.Username,
			UserId:    systemAccount.UserId,
		}
		transfer := services.AccountTransferDTO{
			TradeBody:   body,
			TradeTarget: target,
			TradeNo:     d.EnvelopeNo,
			Amount:      d.Amount,
			ChangeType:  services.EnvelopeOutgoing,
			ChangeFlag:  services.FlagTransferOut,
			Decs:        "红包金额支付",
		}
		_, err = accountDomain.TransferWithContextTx(ctx, transfer)
		if err != nil {
			return err
		}

		//3. 将扣减的红包总金额转入红包中间商的红包资金账户
		//入账
		transfer = services.AccountTransferDTO{
			TradeBody:   target,
			TradeTarget: body,
			TradeNo:     d.EnvelopeNo,
			Amount:      d.Amount,
			ChangeType:  services.EnvelopeIncoming,
			ChangeFlag:  services.FlagTransferIn,
			Decs:        "红包金额转入",
		}
		_, err = accountDomain.TransferWithContextTx(ctx, transfer)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	//扣减金额没有问题，返回活动
	activity.RedEnvelopeGoodsDTO = *d.RedEnvelopeGoods.ToDTO()

	return activity, err
}
