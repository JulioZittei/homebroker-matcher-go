package entity

type OrderQueue struct {
	Orders []*Order
}

func (q *OrderQueue) Less(i, j int) bool {
	return q.Orders[i].Price < q.Orders[j].Price
}

func (q *OrderQueue) Swap(i, j int) {
	q.Orders[i], q.Orders[j] = q.Orders[j], q.Orders[i]
}

func (q *OrderQueue) Len() int {
	return len(q.Orders)
}

func (q *OrderQueue) Push(x interface{}) {
	q.Orders = append(q.Orders, x.(*Order))
}

func (q *OrderQueue) Pop() interface{} {
	old := q.Orders
	length := len(old)
	item := old[length - 1]
	q.Orders = old[0 : length - 1]
	return item
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{}
}