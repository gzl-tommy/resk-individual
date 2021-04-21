package accounts

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"sync"
	"gzl-tommy/resk-individual/infra/base"
	"gzl-tommy/resk-individual/services"
)

var _ services.AccountService = new(accountService)
var once sync.Once

func init() {
	once.Do(func() {
		services.IAccountService = new(accountService)
	})
}

type accountService struct {
}

func (a *accountService) CreateAccount(dto services.AccountCreateDTO) (*services.AccountDTO, error) {
	domain := accountDomain{}
	// 验证输入参数
	err := base.ValidateStruct(&dto)
	if err != nil {
		return nil, err
	}

	// 验证账号是否存在和幂等性
	acc := domain.GetAccountByUserIdAndType(dto.UserId, services.AccountType(dto.AccountType))
	if acc != nil {
		return acc, errors.New(fmt.Sprintf("用户的该类型账户已经存在：username=%s[%s],账户类型=%d",
			dto.Username, dto.UserId, dto.AccountType))
	}
	// 执行账号创建的业务逻辑
	amount, err := decimal.NewFromString(dto.Amount)
	if err != nil {
		return nil, err
	}
	account := services.AccountDTO{
		//AccountNo:    "",
		AccountName:  dto.AccountName,
		AccountType:  dto.AccountType,
		CurrencyCode: dto.CurrencyCode,
		UserId:       dto.UserId,
		Username:     dto.Username,
		Balance:      amount,
		Status:       1,
		//CreatedAt:    time.Time{},
		//UpdatedAt:    time.Time{},
	}
	rdto, err := domain.Create(account)
	return rdto, err
}

func (a *accountService) Transfer(dto services.AccountTransferDTO) (services.TransferedStatus, error) {
	domain := accountDomain{}
	// 验证输入参数
	err := base.ValidateStruct(&dto)
	if err != nil {		
		return services.TransferedStatusFailure, err
	}

	// 执行转账逻辑
	amount, err := decimal.NewFromString(dto.AmountStr)
	if err != nil {
		return services.TransferedStatusFailure, err
	}
	dto.Amount = amount
	if dto.ChangeFlag == services.FlagTransferOut {
		if dto.ChangeType > 0 {
			return services.TransferedStatusFailure, errors.New("如果changeFlag为支出，那么changeType必须小于0")
		}
	} else {
		if dto.ChangeType < 0 {
			return services.TransferedStatusFailure, errors.New("如果changeFlag为收入,那么changeType必须大于0")
		}
	}

	status, err := domain.Transfer(dto)
	if status == services.TransferedStatusSuccess {
		backwardDto := dto
		backwardDto.TradeBody = dto.TradeTarget
		backwardDto.TradeTarget = dto.TradeBody
		backwardDto.ChangeType = -dto.ChangeType
		backwardDto.ChangeFlag = -dto.ChangeFlag
		status, err := domain.Transfer(backwardDto)
		return status, err
	}
	return status, err
}

func (a *accountService) StoreValue(dto services.AccountTransferDTO) (services.TransferedStatus, error) {
	dto.TradeTarget = dto.TradeBody
	dto.ChangeFlag = services.FlagTransferIn
	dto.ChangeType = services.AccountStoreValue
	return a.Transfer(dto)
}

func (a *accountService) GetEnvelopeAccountByUserId(userId string) *services.AccountDTO {
	domain := accountDomain{}
	account := domain.GetEnvelopeAccountByUserId(userId)
	return account
}

func (a *accountService) GetAccount(accountNo string) *services.AccountDTO {
	domain := accountDomain{}
	return domain.GetAccount(accountNo)
}
