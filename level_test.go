package mutex

import "testing"

func TestAddOrderToLevel(t *testing.T) {
	l := &Level{}
	o := &Order{}
	l.AddOrder(o)
	if l.FirstOrder.Order != o {
		t.FailNow()
	}
	if l.LastOrder.Order != o {
		t.FailNow()
	}
	o2 := &Order{}
	l.AddOrder(o2)
	if l.FirstOrder.Order != o {
		t.FailNow()
	}
	if l.LastOrder.Order != o2 {
		t.Error("Second order was not inserted at end of level")
	}
	if l.FirstOrder.NextOrder.Order != o2 {
		t.Error("Link to second order not established")
	}
	if l.FirstOrder.NextOrder.PrevOrder.Order != o {
		t.Error("Link back to first order not established")
	}
}

func TestRemoveOnlyFromLevel(t *testing.T) {
	l := &Level{}
	o := &Order{}
	l.AddOrder(o)
	l.RemoveOrder(o)
	if l.FirstOrder != nil {
		t.Error("Did not remove first link")
	}
	if l.LastOrder != nil {
		t.Error("Did not remove last link")
	}
}

func TestRemoveLastFromLevel(t *testing.T) {
	l := &Level{}
	first_order := &Order{}
	o := &Order{}
	l.AddOrder(first_order)
	l.AddOrder(o)
	l.RemoveOrder(o)
	if l.FirstOrder.Order != first_order {
		t.Error("Did not keep first link to first order")
	}
	if l.LastOrder.Order != first_order {
		t.Error("Did not change last link to first order")
	}
	if l.FirstOrder.NextOrder != nil {
		t.Error("Did not clear link to second order")
	}
}

func TestRemoveMiddleFromLevel(t *testing.T) {
	l := &Level{}
	first_order := &Order{}
	middle_order := &Order{}
	last_order := &Order{}
	l.AddOrder(first_order)
	l.AddOrder(middle_order)
	l.AddOrder(last_order)
	l.RemoveOrder(middle_order)
	if l.FirstOrder.Order != first_order {
		t.Error("First link should be untouched")
	}
	if l.LastOrder.Order != last_order {
		t.Error("Last link should be untouched")
	}
	if l.FirstOrder.NextOrder.Order != last_order {
		t.Error("Middle order should be removed going forward")
	}
	if l.LastOrder.PrevOrder.Order != first_order {
		t.Error("Middle order should be removed going backwards")
	}
}

func TestRemoveFirstFromLevel(t *testing.T) {
	l := &Level{}
	first_order := &Order{}
	middle_order := &Order{}
	last_order := &Order{}
	l.AddOrder(first_order)
	l.AddOrder(middle_order)
	l.AddOrder(last_order)
	l.RemoveOrder(first_order)
	if l.FirstOrder.Order != middle_order {
		t.Fail()
	}
	if l.FirstOrder.NextOrder.Order != last_order {
		t.Fail()
	}
	if l.FirstOrder.NextOrder.PrevOrder.Order != middle_order {
		t.Fail()
	}
}
