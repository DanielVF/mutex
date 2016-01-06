package mutex

import "time"

type Order struct {
	Symbol      string
	Venue       string
	Direction   string
	OriginalQty int
	Qty         int
	TotalFilled int
	Price       int
	OrderType   string
	Id          int
	Account     string
	Ts          time.Time
	Fills       []*OrderFill
	Open        bool
}

type OrderFill struct {
	Price int
	Qty   int
	Ts    time.Time
}

func (o *Order) ApplyIncoming(incoming *Order) *OrderFill {
	fillQty := 0
	if o.Qty < incoming.Qty {
		fillQty = o.Qty
	} else {
		fillQty = incoming.Qty
	}
	fill := &OrderFill{
		Qty:   fillQty,
		Price: o.Price,
		Ts:    time.Now(),
	}
	o.Fill(fill)
	incoming.Fill(fill)
	return fill
}

func (o *Order) Fill(fill *OrderFill) {
	o.Qty -= fill.Qty
	o.TotalFilled += fill.Qty
	if o.Qty == 0 {
		o.Open = false
	}
	o.Fills = append(o.Fills, fill)
}
