package web

import (
	"github.com/kataras/iris/v12/context"
	"github.com/sirupsen/logrus"
	"gzl-tommy/resk-individual/infra"
	"gzl-tommy/resk-individual/infra/base"
	"gzl-tommy/resk-individual/services"
)

func init() {
	infra.RegisterApi(new(AccountApi))
}

// 定义 web api 的时候，对每一个子业务，定义统一的前缀
// 资金账户的根路径定义为：/account
// 版本号：/v1/account
type AccountApi struct {
}

func (a *AccountApi) Init() {
	groupRouter := base.Iris().Party("/v1/account")
	groupRouter.Post("/create", createHandler)
	groupRouter.Post("/transfer", transferHandler)
	groupRouter.Get("/envelope/get", getEnvelopeAccountHandler)
	groupRouter.Get("/get", getAccountHandler)
}

// 账户创建的接口：/v1/account/create
// POST body json
/*
{
	"UserId": "w123456",
	"Username": "测试用户1",
	"AccountName": "测试账户1",
	"AccountType": 0,
	"CurrencyCode": "CNY",
	"Amount": "100.11"
}

{
    "code": 1000,
    "message": "",
    "data": {
        "AccountNo": "1K1hrG0sQw7lDuF6KOQbMBe2o3n",
        "AccountName": "测试账户1",
        "AccountType": 0,
        "CurrencyCode": "CNY",
        "UserId": "w123456",
        "Username": "测试用户1",
        "Balance": "100.11",
        "Status": 1,
        "CreatedAt": "2019-04-18T13:26:51.895+08:00",
        "UpdatedAt": "2019-04-18T13:26:51.895+08:00"
    }
}
*/
func createHandler(ctx context.Context) {
	// 获取请求参数，
	account := services.AccountCreateDTO{}
	err := ctx.ReadJSON(&account)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	// 执行创建账户的代码
	service := services.GetAccountService()
	dto, err := service.CreateAccount(account)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		logrus.Error(err)
	}
	r.Data = dto
	ctx.JSON(&r)
}

//转账的接口 :/v1/account/transfer
/**
{
	"UserId": "w123456-1",
	"Username": "测试用户1",
	"AccountName": "测试账户1",
	"AccountType": 0,
	"CurrencyCode": "CNY",
	"Amount": "100.11"
}
{
	"UserId": "w123456-2",
	"Username": "测试用户2",
	"AccountName": "测试账户2",
	"AccountType": 0,
	"CurrencyCode": "CNY",
	"Amount": "100.11"
}
{
	"TradeNo": "trade123456",
	"TradeBody": {
		"AccountNo": "1K5YdR5Cng5FsBaF95fkcRJE08v",
		"UserId": "w123456-2",
		"Username": "测试用户2"
	},
	"TradeTarget": {
		"AccountNo": "1K5iy4IzhyywntWMeVlxKdxVn4G",
		"UserId": "w123456-1",
		"Username": "测试用户1"
	},
	"AmountStr": "1",

	"ChangeType": -1,
	"ChangeFlag": -1,
	"Decs": "转出"
}
*/
func transferHandler(ctx context.Context) {
	//获取请求参数，
	account := services.AccountTransferDTO{}
	err := ctx.ReadJSON(&account)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	//执行转账逻辑
	service := services.GetAccountService()
	status, err := service.Transfer(account)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		logrus.Error(err)
	}
	if status != services.TransferedStatusSuccess {
		r.Code = base.ResCodeBizError
		r.Message = err.Error()
	}
	r.Data = status
	ctx.JSON(&r)
}

//查询红包账户的web接口: /v1/account/envelope/get
func getEnvelopeAccountHandler(ctx context.Context) {
	userId := ctx.URLParam("userId")
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if userId == "" {
		r.Code = base.ResCodeRequestParamsError
		r.Message = "用户ID不能为空"
		ctx.JSON(&r)
		return
	}
	service := services.GetAccountService()
	account := service.GetEnvelopeAccountByUserId(userId)
	r.Data = account
	ctx.JSON(&r)
}

// 查询账户信息的web接口：/v1/account/get
func getAccountHandler(ctx context.Context) {
	accountNo := ctx.URLParam("accountNo")
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if accountNo == "" {
		r.Code = base.ResCodeRequestParamsError
		r.Message = "账户编号不能为空"
		ctx.JSON(&r)
		return
	}
	service := services.GetAccountService()
	account := service.GetAccount(accountNo)
	r.Data = account
	ctx.JSON(&r)
}
