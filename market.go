package mutex

import "time"
import "errors"

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

func (m *Market) CancelOrder(id int) error {
	o, err := m.Order(id)
	if err != nil {
		return err
	}
	if o.Direction == "buy" {
		m.buyBook.RemoveOrder(o)
	} else if o.Direction == "sell" {
		m.sellBook.RemoveOrder(o)
	}
	o.Open = false
	return nil
}

func (m *Market) Quote() {
	// Grab from totals, (should be fast)
}

func (m *Market) Orderbook() {
	// Summaries
}

func (m *Market) Order(id int) (*Order, error) {
	if id <= 0 {
		return nil, errors.New("Order id must be greater than 0")
	}
	if id > len(m.orders) {
		return nil, errors.New("Order does not exist")
	}
	order := m.orders[id-1]
	return order, nil
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
