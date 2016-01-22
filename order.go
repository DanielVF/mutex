package mutex

import (
	"errors"
	"strings"
	"time"
)

type Order struct {
	Ok          bool         `json:"ok"`
	Symbol      string       `json:"symbol"`
	Venue       string       `json:"venue"`
	Direction   string       `json:"direction"`
	OriginalQty int          `json:"originalQty"`
	Qty         int          `json:"qty"`
	TotalFilled int          `json:"totalFilled"`
	Price       int          `json:"price"`
	OrderType   string       `json:"orderType"`
	Id          int          `json:"id"`
	Account     string       `json:"account"`
	Ts          time.Time    `json:"ts"`
	Fills       []*OrderFill `json:"fills"`
	Open        bool         `json:"open"`
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

func (o *Order) Validate() error {
	if o.Qty <= 0 {
		return errors.New("Qty cannot be zero or less.")
	}
	if o.Qty >= 10000 {
		return errors.New("Qty cannot be more than 10000.")
	}
	if o.Price <= 0 {
		return errors.New("Price cannot be zero or less.")
	}
	if o.Price >= 10000 {
		return errors.New("Price cannot be more than $10,000 (1000000 cents).")
	}
	o.OrderType = strings.ToLower(o.OrderType)
	if o.OrderType != "limit" {
		return errors.New("The only supported order type is 'limit'")
	}
	o.Direction = strings.ToLower(o.Direction)
	if o.Direction != "buy" && o.Direction != "sell" {
		return errors.New("Direction must be buy or sell")
	}
	return nil
}
