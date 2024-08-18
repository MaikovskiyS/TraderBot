package bybit_provider

import (
	"context"
	"fmt"

	"github.com/MaikovskiyS/TraderBot/internal/domain"
	trading "github.com/MaikovskiyS/TraderBot/internal/trader/services/traiding_service"
	"github.com/hirokisan/bybit/v2"
)

func (p *BybitProvider) GetPositionInfo(ctx context.Context) (*trading.GetPositionInfoResponse, error) {
	settel := bybit.CoinUSDT
	infos, err := p.BybitClient.V5().Position().GetPositionInfo(bybit.V5GetPositionInfoParam{
		Category:   bybit.CategoryV5Linear,
		Symbol:     nil,
		BaseCoin:   nil,
		SettleCoin: &settel,
		Limit:      nil,
		Cursor:     nil,
	})
	if err != nil {
		return nil, fmt.Errorf("get account info: %w", err)
	}

	if infos.Result.List == nil {
		return nil, ErrEmptyResponse
	}

	positions := make([]*domain.Position, len(infos.Result.List))
	for i, info := range infos.Result.List {
		positions[i] = convertPositionInfoToDomain(info)
	}

	return &trading.GetPositionInfoResponse{
		Positions: positions,
	}, nil

}

func (p *BybitProvider) GetBalance(ctx context.Context) (float64, error) {
	tickers := []bybit.Coin{bybit.CoinUSDT}

	resp, err := p.BybitClient.V5().Account().GetWalletBalance(bybit.AccountTypeV5CONTRACT, tickers)
	if err != nil {
		return 0, fmt.Errorf("get balance: %w", err)
	}

	if resp.Result.List == nil {
		return 0, ErrInvalidResponseType
	}

	return convertStrToFloat(resp.Result.List[0].Coin[0].WalletBalance), nil
}
