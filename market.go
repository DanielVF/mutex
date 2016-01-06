package mutex

import "time"

type Market struct {
	// Orders
	orders   []*Order
	buyBook  []*Order
	sellBook []*Order
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
}

func (m *Market) applyOrder(o *Order) {

}

func (m *Market) addToBooks(o *Order) {
	m.orders = append(m.orders, o)
	if o.OrderType == "immediate-or-cancel" {
		return
	}
	if o.Qty <= 0 {
		return
	}
	if o.Direction == "buy" {
		m.buyBook = append(m.buyBook, o)
	} else if o.Direction == "sell" {
		m.sellBook = append(m.sellBook, o)
	}
}
