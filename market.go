package mutex

import "time"

type Market struct {
	// Orders
	orders   []*Order
	buyBook  OrderBook
	sellBook OrderBook
	// last
	lastId int
}

func (m *Market) PlaceOrder(o *Order) {
	m.prepOrder(o)
	// Check for crosses
	m.applyOrder(o)
	// Add to books if needed
	m.addToBooks(o)
	// One day will fire off execution and ticker events
}

func (m *Market) CancelOrder(id int) {
	// Remove it, adjust stats
}

func (m *Market) Quote() {
	// Grab from totals, (should be fast)
}

func (m *Market) Orderbook() {
	// Summaries
}

func (m *Market) OrderStatus() {

}

func (m *Market) prepOrder(o *Order) {
	m.lastId += 1
	o.Id = m.lastId
	o.Ts = time.Now()
	o.OriginalQty = o.Qty
	o.TotalFilled = 0
	o.Ok = true
	o.Open = true
	m.orders = append(m.orders, o)
}

func (m *Market) applyOrder(o *Order) {
	if o.Direction == "buy" {
		m.sellBook.ExecuteOrder(o)
	} else if o.Direction == "sell" {
		m.buyBook.ExecuteOrder(o)
	}
}

func (m *Market) addToBooks(o *Order) {

	if o.OrderType == "immediate-or-cancel" {
		return
	}
	if o.Qty <= 0 {
		return
	}
	if o.Direction == "buy" {
		m.buyBook.AddOrder(o)
	} else if o.Direction == "sell" {
		m.sellBook.AddOrder(o)
	}
}
