package entity

type Investor struct {
	ID string
	Name string
	AssetPosition []*InvestorAssetPosition
}

func (i *Investor) GetAssetPosition(assetId string) *InvestorAssetPosition {
	for _, assetPosition := range i.AssetPosition {
		if assetPosition.ID == assetId {
			return assetPosition
		}
	}
	return nil
}

func (i *Investor) AddAssetPosition(assetPosition *InvestorAssetPosition) {
	i.AssetPosition = append(i.AssetPosition, assetPosition)
}

func (i *Investor) UpdateAssetPosition(assetId string, qtdShares int) {
	assetPosition := i.GetAssetPosition(assetId)
	if assetPosition == nil {
		i.AssetPosition = append(i.AssetPosition, NewInvestorAssetPosition(assetId, int8(qtdShares)))
	} else {
		assetPosition.Shares += qtdShares
	}
}

func NewInvestor(id string) *Investor {
	return &Investor{
		ID: id,
		AssetPosition: []*InvestorAssetPosition{},
	}
}

type InvestorAssetPosition struct {
	ID string
	Shares int
}

func NewInvestorAssetPosition(assetId string, shares int8) *InvestorAssetPosition {
	return &InvestorAssetPosition{
		ID: assetId,
		Shares: int(shares),
	}
}