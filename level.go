package mutex

type Level struct {
	Price      int
	Qty        int
	NextLevel  *Level
	PrevLevel  *Level
	FirstOrder *LevelOrder
	LastOrder  *LevelOrder
}

type LevelOrder struct {
	Order     *Order
	Level     *Level
	NextOrder *LevelOrder
	PrevOrder *LevelOrder
}

func (l *Level) AddOrder(o *Order) {
	lo := &LevelOrder{
		Order: o,
		Level: l,
	}
	if l.LastOrder == nil {
		// We have no orders on this level yet
		l.FirstOrder = lo
		l.LastOrder = lo
	} else {
		l.LastOrder.NextOrder = lo
		lo.PrevOrder = l.LastOrder
		l.LastOrder = lo
	}
	l.Qty += o.Qty
}

func (l *Level) RemoveFilledFirstOrder(o *Order, filled int) {
	if l.FirstOrder.Order != o {
		panic("Attempted to remove wrong first level")
	}
	l.FirstOrder = l.FirstOrder.NextOrder
	if l.FirstOrder != nil {
		l.FirstOrder.PrevOrder = nil
	}
}

func (l *Level) RemoveOrder(o *Order) bool {
	if l.LastOrder == nil {
		return false
	}
	if l.LastOrder.Order == o {
		// Special case handling the last order
		if l.LastOrder.PrevOrder == nil {
			// If this is the only order on the level
			l.LastOrder = nil
			l.FirstOrder = nil
		} else {
			l.LastOrder.PrevOrder.NextOrder = nil
			l.LastOrder = l.LastOrder.PrevOrder
		}
	} else {
		for lo := l.LastOrder.PrevOrder; lo != nil; lo = lo.PrevOrder {
			if lo.Order == o {
				if l.FirstOrder == lo {
					l.FirstOrder = lo.NextOrder
					// Special case for first order
				} else {
					// Handle a middle order by removing me from links on either side
					lo.NextOrder.PrevOrder = lo.PrevOrder
					lo.PrevOrder.NextOrder = lo.NextOrder
				}
				break
			}
		}
	}
	l.Qty -= o.Qty
	return true
}
