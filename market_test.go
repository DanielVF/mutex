package mutex

import "testing"

func TestPlaceOrder(t *testing.T) {
	m := &Market{}
	o := &Order{
		Direction: "buy",
		Qty:       10,
	}
	m.PlaceOrder(o)
	if o.Id != 1 {
		t.FailNow()
	}
}

func TestCancelOrder(t *testing.T) {
	m := new(Market)
	o := new(Order)
	m.PlaceOrder(o)
	m.CancelOrder(o.Id)
}
