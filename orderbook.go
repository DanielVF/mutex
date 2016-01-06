package mutex

import "strconv"
import "bytes"

type OrderBook struct {
	IsBuy      bool
	Last       int
	Price      int
	Depth      int
	TotalDepth int
	FirstLevel *Level
	LastLevel  *Level
}

func (ob *OrderBook) String() string {
	var buffer bytes.Buffer
	if ob.IsBuy {
		buffer.WriteString("BUY ")
	} else {
		buffer.WriteString("SELL ")
	}
	for current_level := ob.FirstLevel; current_level != nil; current_level = current_level.NextLevel {
		buffer.WriteString(strconv.Itoa(current_level.Price))
		buffer.WriteString(" ")
	}
	return buffer.String()
}

func (ob *OrderBook) Level(price int) *Level {
	// Check for no levels
	if ob.FirstLevel == nil {
		l := &Level{
			Price: price,
		}
		ob.FirstLevel = l
		ob.LastLevel = l
		return l
	} else {
		// Find the next level
		for current_level := ob.FirstLevel; current_level != nil; current_level = current_level.NextLevel {
			if current_level.Price == price {
				return current_level
			}
			if (ob.IsBuy && current_level.Price < price) || (!ob.IsBuy && current_level.Price > price) {
				l := &Level{
					Price: price,
				}
				if current_level.PrevLevel == nil {
					// New first level
					ob.FirstLevel = l
					l.NextLevel = current_level
					current_level.PrevLevel = l
				} else {
					// New middle level
					l.NextLevel = current_level
					l.PrevLevel = current_level.PrevLevel
					current_level.PrevLevel.NextLevel = l
					current_level.PrevLevel = l
				}
				return l
			}
		}

		// If we get here, we need to add to the end of levels
		l := &Level{
			Price: price,
		}
		l.PrevLevel = ob.LastLevel
		ob.LastLevel.NextLevel = l
		ob.LastLevel = l
		return l

	}
}

func (ob *OrderBook) RemoveLevel(price int) {
	l := ob.Level(price)
	if l.FirstOrder != nil {
		return
	}
	// Only order
	if ob.FirstLevel == l && ob.LastLevel == l {
		ob.FirstLevel = nil
		ob.LastLevel = nil

	} else if ob.FirstLevel == l { // First Order
		ob.FirstLevel = ob.FirstLevel.NextLevel
		ob.FirstLevel.PrevLevel = nil

	} else if ob.LastLevel == l { // Last Order
		ob.LastLevel = ob.LastLevel.PrevLevel
		ob.LastLevel.NextLevel = nil

	} else { // Middle Level
		l.NextLevel.PrevLevel = l.PrevLevel
		l.PrevLevel.NextLevel = l.NextLevel
	}
}

func (ob *OrderBook) AddOrder(o *Order) {
	if o.Qty > 0 {
		o.Open = true
	}

	l := ob.Level(o.Price)
	l.AddOrder(o)
	ob.TotalDepth += o.Qty
	ob.UpdateFirstStats()
}

func (ob *OrderBook) RemoveOrder(o *Order) {
	l := ob.Level(o.Price)
	did_remove := l.RemoveOrder(o)
	if did_remove {
		ob.TotalDepth -= o.Qty
	}
	if l.Qty == 0 && l.FirstOrder == nil {
		ob.RemoveLevel(l.Price)
	}
	ob.UpdateFirstStats()
}

func (ob *OrderBook) UpdateFirstStats() {
	if ob.FirstLevel != nil {
		ob.Depth = ob.FirstLevel.Qty
		ob.Price = ob.FirstLevel.Price
	} else {
		ob.Depth = 0
		ob.Price = 0
	}
}

func (ob *OrderBook) Crosses(resting int, incoming int) bool {
	return (ob.IsBuy && incoming <= resting) || (!ob.IsBuy && incoming >= resting)
}

func (ob *OrderBook) ExecuteOrder(o *Order) {
	if ob.FirstLevel == nil {
		return
	}
	if !ob.Crosses(ob.Price, o.Price) {
		return
	}
	for l := ob.FirstLevel; l != nil && ob.Crosses(l.Price, o.Price); l = l.NextLevel {
		for lo := l.FirstOrder; lo != nil; lo = lo.NextOrder {
			fill := lo.Order.ApplyIncoming(o)
			if fill != nil {
				l.Qty -= fill.Qty
				ob.Last = fill.Price
				if lo.Order.Open == false {
					l.RemoveFilledFirstOrder(lo.Order, fill.Qty)
					ob.TotalDepth -= fill.Qty
				}
				if o.Open == false {
					break
				}
			}
		}
		if l.Qty == 0 && l.FirstOrder == nil {
			ob.RemoveLevel(l.Price)
		}
	}
	ob.UpdateFirstStats()
}
