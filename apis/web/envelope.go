package web

import (
	"github.com/kataras/iris/v12/context"
	"gzl-tommy/resk-individual/infra"
	"gzl-tommy/resk-individual/infra/base"
	"gzl-tommy/resk-individual/services"
)

func init() {
	infra.RegisterApi(&EnvelopeApi{})
}

// 红包API
type EnvelopeApi struct {
	service services.RedEnvelopeService
}

func (e *EnvelopeApi) Init() {
	e.service = services.GetRedEnvelopeService()
	groupRouter := base.Iris().Party("/v1/envelope")
	groupRouter.Post("/sendout", e.sendOutHandler)
	groupRouter.Post("/receive", e.receiveHandler)
}

//{
//	"envelopeType": 0,
//	"username": "",
//	"userId": "",
//	"blessing": "",
//	"amount": "0",
//	"quantity": 0
//}0
func (e *EnvelopeApi) sendOutHandler(ctx context.Context) {
	dto := services.RedEnvelopeSendingDTO{}
	err := ctx.ReadJSON(&dto)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	activity, err := e.service.SendOut(dto)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = activity
	ctx.JSON(r)
}
/*
{
"envelopeNo" :"1qpAGFKrxpX2bPVGdDOLwsfs2jL",
"recvUsername":"测试用户10",
"recvUserId" :"1qulNd2v7jAilY7gFVHysja1Ddl",
"accountNo":"1qulNdcoFWOW6V3gjbGWGpnsmUM"
}
*/
func (e *EnvelopeApi) receiveHandler(ctx context.Context) {
	dto := services.RedEnvelopeReceiveDTO{}
	err := ctx.ReadJSON(&dto)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	item, err := e.service.Receive(dto)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = item
	ctx.JSON(r)
}
