package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Orders []*Order
	Transactions []*Transaction
	OrderChanIn chan *Order
	OrderChanOut chan *Order
	Wg *sync.WaitGroup
}

func NewBook(orderChanIn chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Orders: []*Order{},
		Transactions: []*Transaction{},
		OrderChanIn: orderChanIn,
		OrderChanOut: orderChanOut,
		Wg: wg,
	}
}

func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares

	minShares := sellingShares
	if buyingShares < minShares {
		minShares = buyingShares
	}
	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.AddSellingOrderPendingShares(-minShares)
	
	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.ID, minShares)
	transaction.AddBuyingOrderPendingShares(-minShares)

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)

	transaction.CloseBuyingOrder()
	transaction.CloseSellingOrder()

	b.Transactions = append(b.Transactions, transaction)
}

func (b *Book) Trade() {
	buyingOrders := NewOrderQueue()
	sellingOrders := NewOrderQueue()

	heap.Init(buyingOrders)
	heap.Init(sellingOrders)

	for order := range b.OrderChanIn {
		if order.OrderType == "BUY" {
			buyingOrders.Push(order)
			if sellingOrders.Len() > 0 && sellingOrders.Orders[0].Price <= order.Price {
				sellOrder := sellingOrders.Pop().(*Order)
				if sellOrder.PendingShares > 0 {
					transaction := NewTransaction(sellOrder, order, order.Shares, sellOrder.Price)
					b.AddTransaction(transaction, b.Wg)
					sellOrder.Transactions = append(sellOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)
					b.OrderChanOut <- sellOrder
					b.OrderChanOut <- order
					if sellOrder.PendingShares > 0 {
						sellingOrders.Push(sellOrder)
					}
				}
			}
		} else if order.OrderType == "SELL" {
			sellingOrders.Push(order)
			if buyingOrders.Len() > 0 && buyingOrders.Orders[0].Price <= order.Price {
				buyOrder := buyingOrders.Pop().(*Order)
				if buyOrder.PendingShares > 0 {
					transaction := NewTransaction(order, buyOrder, order.Shares, buyOrder.Price)
					b.AddTransaction(transaction, b.Wg)
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)
					b.OrderChanOut <- buyOrder
					b.OrderChanOut <- order
					if buyOrder.PendingShares > 0 {
						buyingOrders.Push(buyOrder)
					} 
				}
			}
		}
	}
}