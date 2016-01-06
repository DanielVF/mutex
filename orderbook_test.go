package mutex

import "testing"

func TestAddFirstLevelToOrderBook(t *testing.T) {
	ob := &OrderBook{
		IsBuy: false,
	}
	o := &Order{
		Price: 1000,
	}
	l := ob.Level(o.Price)
	if l.Price != 1000 {
		t.Error("Level does have the correct price")
	}
	if ob.FirstLevel != l {
		t.Error("Level is not set as first level")
	}

	if ob.LastLevel != l {
		t.Error("Level is not set as last level")
	}
}
func TestLevelAddNewFirstLevel(t *testing.T) {
	ob := &OrderBook{
		IsBuy: true,
	}
	lower_level := ob.Level(500)
	higher_level := ob.Level(800)
	if ob.FirstLevel != higher_level {
		t.Error("First level incorrect")
	}
	if ob.LastLevel != lower_level {
		t.Error("Second level incorrect")
	}
	if ob.FirstLevel.NextLevel != lower_level {
		t.Error("Did not link first level to next level")
	}
}

func TestLevelAddNewLastLevel(t *testing.T) {
	ob := &OrderBook{
		IsBuy: true,
	}
	higher_level := ob.Level(800)
	lower_level := ob.Level(500)
	if ob.FirstLevel != higher_level {
		t.Error("First level incorrect")
	}
	if ob.LastLevel != lower_level {
		t.Error("Second level incorrect")
	}
}

func TestAddThreeLevels(t *testing.T) {
	ob := &OrderBook{
		IsBuy: true,
	}
	ob.Level(500)
	ob.Level(100)
	ob.Level(100)
	// t.Log(ob.String())
	if ob.FirstLevel.PrevLevel != nil && ob.FirstLevel.NextLevel != nil {
		t.Error("Only level somehow has links.")
	}

	if ob.FirstLevel.Price != 500 {
		t.Error("Levels are not ordered correctly")
	}
	if ob.LastLevel.Price != 100 {
		t.Error("Levels are not ordered correctly")
	}
	if ob.LastLevel.PrevLevel.Price != 500 {
		t.Error("Levels are not ordered correctly")
	}
	if ob.FirstLevel.NextLevel == nil {
		t.Error("First level does not link to second level")
		t.FailNow()
	}
	if ob.FirstLevel.NextLevel.Price != 100 {
		t.Error("First level links to wrong second level")
	}
	ob.Level(1000)
	if ob.FirstLevel.Price != 1000 {
		t.Error("Levels are not ordered correctly")
	}
	ob.Level(600)
	if ob.FirstLevel.NextLevel.Price != 600 {
		t.Log(ob.String())
		t.Error("Incorrect sell ordering")
	}
}

func TestRemoveLevels(t *testing.T) {
	ob := &OrderBook{
		IsBuy: true,
	}
	ob.Level(500)
	ob.RemoveLevel(500)
	if ob.FirstLevel != nil {
		t.Error("Did not remove only level")
	}

	// Remove First
	ob.Level(500) // First
	ob.Level(100) // Last
	ob.RemoveLevel(500)
	if ob.FirstLevel.Price != 100 || ob.LastLevel.Price != 100 {
		t.Error("Did not remove first level")
	}
	ob.RemoveLevel(100)

	// Remove Last
	ob.Level(500) // First
	ob.Level(100) // Last
	ob.RemoveLevel(100)
	if ob.FirstLevel.Price != 500 || ob.LastLevel.Price != 500 {
		t.Error("Did not remove last level")
	}
	ob.RemoveLevel(500)

	// Remove Middle Lavel
	ob.Level(500) // First
	ob.Level(200) // Middel
	ob.Level(100) // Last
	ob.RemoveLevel(200)
	if ob.FirstLevel.Price != 500 {
		t.Error("Failed to keep first level")
	}
	if ob.FirstLevel.NextLevel.Price != 100 {
		t.Error("Did not remove middle level")
	}
	if ob.LastLevel.PrevLevel.Price != 500 {
		t.Error("Did not remove middle level")
	}
	ob.RemoveLevel(500)

}

func TestSellBookOrdering(t *testing.T) {
	ob := &OrderBook{
		IsBuy: false,
	}
	ob.Level(500)
	ob.Level(100)
	ob.Level(1000)
	ob.Level(600)
	// t.Log(ob.String())
	if ob.FirstLevel.Price != 100 {
		t.Error("Incorrect sell ordering")
	}
	if ob.FirstLevel.NextLevel.Price != 500 {
		t.Error("Incorrect sell ordering")
	}
	if ob.FirstLevel.NextLevel.NextLevel.Price != 600 {
		t.Error("Incorrect sell ordering")
	}
	if ob.FirstLevel.NextLevel.NextLevel.NextLevel.Price != 1000 {
		t.Log(ob.String())
		t.Error("Incorrect sell ordering")
	}

}

func TestOrderbookDebugging(t *testing.T) {
	ob := &OrderBook{
		IsBuy: false,
	}
	ob.AddOrder(&Order{
		Price:     100,
		Qty:       200,
		Direction: "buy",
	})
	ob.String()     // Just happy for no crashes
	ob.IsBuy = true // Changing the orderbook type would be nuts
	ob.String()
}

func TestAddSingleOrder(t *testing.T) {
	ob := &OrderBook{
		IsBuy: true,
	}
	o := &Order{
		Price:     100,
		Qty:       200,
		Direction: "buy",
	}
	ob.AddOrder(o)
	if ob.FirstLevel == nil {
		t.Error("Did not create level")
	}
	if ob.FirstLevel.Price != 100 {
		t.Error("Did not create level with correct price")
	}
	if ob.FirstLevel.FirstOrder.Order != o {
		t.Error("Did not store order")
	}
}

func TestAddMultipleOrdersWithStats(t *testing.T) {
	ob := &OrderBook{
		IsBuy: true,
	}
	ob.AddOrder(&Order{
		Price:     100,
		Qty:       200,
		Direction: "buy",
	})
	ob.AddOrder(&Order{
		Price:     102,
		Qty:       10,
		Direction: "buy",
	})
	ob.AddOrder(&Order{
		Price:     98,
		Qty:       15,
		Direction: "buy",
	})
	if ob.FirstLevel.Price != 102 {
		t.Log(ob.String())
		t.Error("Did not create level with correct price")
	}
	if ob.TotalDepth != 225 {
		t.Error("Total depth not set correctly")
	}
	if ob.Depth != 10 {
		t.Error("depth not set correctly")
	}
	if ob.Price != 102 {
		t.Error("Orderbook price not set correctly")
	}
}

func TestRemoveOrdersWithStats(t *testing.T) {
	ob := &OrderBook{
		IsBuy: true,
	}
	o1 := &Order{
		Price:     1000,
		Qty:       10,
		Direction: "buy",
	}
	ob.AddOrder(o1)
	o2 := &Order{
		Price:     1000,
		Qty:       20,
		Direction: "buy",
	}
	ob.AddOrder(o2)
	if ob.TotalDepth != 30 {
		t.Error("Total depth not set correctly")
	}
	if ob.Depth != 30 {
		t.Error("Depth not set correctly")
	}
	ob.RemoveOrder(o2)
	if ob.TotalDepth != 10 {
		t.Error("Total depth not set correctly")
	}
	if ob.Depth != 10 {
		t.Error("Depth not set correctly")
	}
	ob.RemoveOrder(o1)
	if ob.TotalDepth != 0 {
		t.Error("Total depth not set correctly")
	}
	if ob.Depth != 0 {
		t.Error("Depth not set correctly")
	}
	if ob.FirstLevel != nil {
		t.Error("Empty level not removed")
	}
	if ob.LastLevel != nil {
		t.Error("Empty level not removed")
	}
}

func TestCrossing(t *testing.T) {
	ob := &OrderBook{
		IsBuy: true,
	}
	if !ob.Crosses(1000, 900) {
		t.Error("A resting buy order for 1000 failed to cross a sell order for 900")
	}
	if ob.Crosses(1000, 1100) {
		t.Error("A resting buy order for 1000 incorrectly crossed a sell order for 1100")
	}
}

func TestSimpleExecution(t *testing.T) {
	ob := &OrderBook{
		IsBuy: true,
	}
	resting_order := &Order{
		Direction: "buy",
		Price:     1000,
		Qty:       20,
		Open:      true,
	}
	ob.AddOrder(resting_order)
	incoming_order := &Order{
		Direction: "sell",
		Price:     990,
		Qty:       15,
		Open:      true,
	}
	if ob.FirstLevel.Qty != 20 {
		t.Error("Level Qty did not adjust correctly")
		t.Log(ob.FirstLevel.Qty)
	}

	ob.ExecuteOrder(incoming_order)
	if ob.FirstLevel.Qty != 5 {
		t.Error("Level Qty did not adjust correctly")
		t.Log(ob.FirstLevel.Qty)
	}
	if incoming_order.Qty != 0 {
		t.Log(resting_order.Qty)
		t.Log(incoming_order.Qty)
		t.Error("Incoming order qty incorrect")
	}
	if incoming_order.TotalFilled != 15 {
		t.Error("Incoming TotalFilled qty incorrect")
	}
	if incoming_order.Qty != 0 {
		t.Error("Filled incoming order did not go to zero qty")
	}
	if incoming_order.Open {
		t.Error("Filled incoming order did not close")
	}
	if resting_order.Qty != 5 {
		t.Error("Resting qty incorrect")
	}
	if resting_order.TotalFilled != 15 {
		t.Error("Resting total filled incorrect")
	}

	// Check fills
	if resting_order.Fills[0].Qty != 15 {
		t.Error("Resting fill qty incorrect")
	}
	if resting_order.Fills[0].Price != 1000 {
		t.Error("Resting fill price incorrect")
	}
	if incoming_order.Fills[0].Qty != 15 {
		t.Error("Incoming fill qty incorrect")
	}
	if incoming_order.Fills[0].Price != 1000 {
		t.Error("Incoming fill qty incorrect")
	}

	// Next order will clear out initial order
	next_order := &Order{
		Direction: "sell",
		Price:     990,
		Qty:       16,
		Open:      true,
	}
	ob.ExecuteOrder(next_order)
	if ob.FirstLevel != nil {
		t.Error("Level did not remove from Orderbook")
		if ob.FirstLevel.Qty != 0 {
			t.Log(ob.FirstLevel.Qty)
			t.Error("Level Qty did not adjust correctly")
		}
		if ob.FirstLevel.FirstOrder != nil {
			t.Error("Order did not remove from level")
		}
	}

	if next_order.Qty != 11 {
		t.Error("Next order did not partially fill")
	}
	if next_order.TotalFilled != 5 {
		t.Error("Next order did not partially fill")
	}
	if next_order.Open != true {
		t.Error("Unfilled incoming should not be marked closed")
	}
	if resting_order.Qty != 0 {
		t.Error("Resting qty incorrect")
	}
	if resting_order.TotalFilled != 20 {
		t.Error("Resting qty incorrect")
	}
	if resting_order.Open != false {
		t.Error("Resting order did not close out when filled")
	}
	if ob.Price != 0 {
		t.Error("Failed to reset orderbook price")
	}

}

func TestMultipleOrdersOnLevel(t *testing.T) {
	ob := &OrderBook{
		IsBuy: false,
	}

	restingOrdersQty := []int{100, 150, 200}
	for _, qty := range restingOrdersQty {
		ob.AddOrder(&Order{
			Direction: "sell",
			Price:     1000,
			Qty:       qty,
			Open:      true,
		})
	}
	if ob.FirstLevel.Qty != 450 {
		t.Error("Level Qty did not adjust correctly")
		t.Log(ob.FirstLevel.Qty)
	}
	incomingOrder := &Order{
		Direction: "sell",
		Price:     1000,
		Qty:       500,
		Open:      true,
	}
	ob.ExecuteOrder(incomingOrder)
	if ob.FirstLevel != nil {
		t.Error("First level should have been wiped out")
		t.Log(ob.FirstLevel.Qty)
	}
	if incomingOrder.Qty != 50 {
		t.Error("incoming order should still have some remaining")
	}

}
