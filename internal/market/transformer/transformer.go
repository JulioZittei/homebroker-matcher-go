package transformer

import (
	"github.com/JulioZittei/homebroker-matcher-go/internal/market/dto"
	"github.com/JulioZittei/homebroker-matcher-go/internal/market/entity"
)

func TransformInput(input dto.TradeInput) *entity.Order {
	asset := entity.NewAsset(input.AssetId, input.AssetId, 1000)
	investor := entity.NewInvestor(input.InvestorId)
	order := entity.NewOrder(input.OrderId, investor, asset, input.Shares, input.Price, input.OrderType)

	if input.CurrentShares > 0 {
		assetPosition := entity.NewInvestorAssetPosition(input.AssetId, int8(input.CurrentShares))
		investor.AddAssetPosition(assetPosition)
	}

	return order
}

func TransformOutput(order *entity.Order) *dto.OrderOutput {
	output := &dto.OrderOutput{
		OrderId: order.ID,
		InvestorId: order.Investor.ID,
		AssetId: order.Asset.ID,
		OrderType: order.OrderType,
		Status: order.Status,
		Partial: order.PendingShares,
		Shares: order.Shares,
	}

	var transactions []*dto.TransactionOutput

	for _, t := range order.Transactions {
		transactionOutput := &dto.TransactionOutput{
			TransactionId: t.ID,
			BuyerId: t.BuyingOrder.ID,
			SellerId: t.SellingOrder.ID,
			AssetId: t.SellingOrder.Asset.ID,
			Price: t.Price,
			Shares: t.SellingOrder.Shares - t.SellingOrder.PendingShares,
		}
		transactions = append(transactions, transactionOutput)
	}

	output.Transactions = transactions

	return output
}