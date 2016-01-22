package mutex

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestOrders(t *testing.T) {
	Convey("Validation", t, func() {
		order := Order{
			Symbol:    "BREQ",
			Venue:     "MUTEX",
			Qty:       10,
			Price:     1000,
			Direction: "buy",
			OrderType: "Limit",
		}

		Convey(`Quantities should be sane`, func() {
			Convey(`Not zero or less`, func() {
				order.Qty = 0
				So(order.Validate(), ShouldNotEqual, nil)
			})
			Convey("Not more than 10K", func() {
				order.Qty = 10001
				So(order.Validate(), ShouldNotEqual, nil)
			})

			Convey("Otherwise good", func() {
				order.Qty = 100
				So(order.Validate(), ShouldEqual, nil)
			})
		})

		Convey("Price should be sane", func() {

			Convey(`Not zero or less`, func() {
				order.Price = 0
				So(order.Validate(), ShouldNotEqual, nil)
			})
			Convey("Not more than 10K", func() {
				order.Price = 1000001
				So(order.Validate(), ShouldNotEqual, nil)
			})

			Convey("Otherwise good", func() {
				order.Price = 1200
				So(order.Validate(), ShouldEqual, nil)
			})

		})

		Convey("OrderType", func() {

			Convey("Handle strange case", func() {
				order.OrderType = "liMIT"
				So(order.Validate(), ShouldEqual, nil)
			})

			Convey("Forbid unsupported", func() {
				order.OrderType = "doublecross"
				So(order.Validate(), ShouldNotEqual, nil)
			})

			Convey("Allow Limit orders", func() {
				order.OrderType = "limit"
				So(order.Validate(), ShouldEqual, nil)
			})

		})

		Convey("Direction", func() {

			Convey("Handle any case", func() {
				order.Direction = "buY"
				So(order.Validate(), ShouldEqual, nil)
			})

			Convey("Forbid unsupported", func() {
				order.Direction = "maybe"
				So(order.Validate(), ShouldNotEqual, nil)
			})

			Convey("Allow sell", func() {
				order.Direction = "sell"
				So(order.Validate(), ShouldEqual, nil)
			})

			Convey("Allow buy", func() {
				order.Direction = "buy"
				So(order.Validate(), ShouldEqual, nil)
			})

		})

	})
}
