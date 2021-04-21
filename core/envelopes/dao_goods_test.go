package envelopes

import (
	"database/sql"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"testing"
	"time"
	"gzl-tommy/resk-individual/infra/base"
	"gzl-tommy/resk-individual/services"
	_ "gzl-tommy/resk-individual/testx"
	"github.com/segmentio/ksuid"
)

func TestRedEnvelopeGoodsDao_Insert(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &RedEnvelopeGoodsDao{
			runner: runner,
		}
		Convey("插入红包数据", t, func() {
			re := RedEnvelopeGoods{
				EnvelopeNo:   ksuid.New().Next().String(),
				EnvelopeType: services.LuckyEnvelopeType,
				Username:     sql.NullString{String: "测试红包账号", Valid: true},
				UserId:       ksuid.New().Next().String(),
				Blessing:     sql.NullString{String: "恭喜发财，大吉大利！", Valid: true},
				Amount:       decimal.NewFromFloat(100),
				//AmountOne:,
				Quantity:       10,
				RemainAmount:   decimal.NewFromFloat(100),
				RemainQuantity: 10,
				ExpiredAt:      time.Now().Add(24 * time.Hour),
				Status:         services.OrderCreate,
				OrderType:      services.OrderTypeSending,
				PayStatus:      services.Payed,
			}
			id, err := dao.Insert(&re)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)

			na := dao.GetOne(re.EnvelopeNo)
			So(na, ShouldNotBeNil)
			So(na.Amount.String(), ShouldEqual, re.Amount.String())
			So(na.CreatedAt, ShouldNotBeNil)
			So(na.UpdatedAt, ShouldNotBeNil)
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
