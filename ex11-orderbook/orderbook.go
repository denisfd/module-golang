package orderbook

import (
//	"fmt"
)

type Orderbook struct {
	LastID int
	Bids   []*Order
	Asks   []*Order
}

func New() *Orderbook {
	ob := &Orderbook{}

	ob.Bids = []*Order{}
	ob.Asks = []*Order{}

	return ob
}

func (orderbook *Orderbook) Match(order *Order) ([]*Trade, *Order) {
	//fmt.Printf("Order: %+v\n", *order)
	//fmt.Printf("State:\n")
	//for num, bid := range orderbook.Bids {
	//fmt.Printf(" #%d: %+v\n", num, *bid)
	//}
	//for num, ask := range orderbook.Asks {
	//fmt.Printf(" #%d: %+v\n", num, ask)
	//}

	var (
		trades []*Trade
		ord    *Order
	)
	switch {

	case order.Side == SideAsk && order.Kind == KindMarket:
		trades, ord = orderbook.MarketAsk(order)

	case order.Side == SideBid && order.Kind == KindMarket:
		trades, ord = orderbook.MarketBid(order)

	case order.Side == SideAsk && order.Kind == KindLimit:
		trades, ord = orderbook.LimitAsk(order)

	case order.Side == SideBid && order.Kind == KindLimit:
		trades, ord = orderbook.LimitBid(order)
	}

	//fmt.Printf("Trades generated: %+v\n", trades)
	//for num, trade := range trades {
	//fmt.Printf(" #%d: %+v\n", num, *trade)
	//}
	//fmt.Printf("\n")

	return trades, ord
}

func (ob *Orderbook) MarketBid(order *Order) ([]*Trade, *Order) {
	return ob.LimitBid(order)
}

func (ob *Orderbook) MarketAsk(order *Order) ([]*Trade, *Order) {
	return ob.LimitAsk(order)
}

func (ob *Orderbook) LimitAsk(order *Order) ([]*Trade, *Order) {
	//	println("Limit Ask called")
	trades := []*Trade{}

	for i := 0; i < len(ob.Asks); i++ {
		ask := ob.Asks[i]
		if order.Price == 0 || ask.Price >= order.Price {
			trade := &Trade{}
			trade.Price = ask.Price
			trade.Bid = order
			trade.Ask = ask

			if ask.Volume > order.Volume {
				trade.Volume = order.Volume
				ask.Volume -= order.Volume
				order.Volume = 0
			} else { //Ask shoud be removed from asks
				trade.Volume = ask.Volume
				order.Volume -= ask.Volume
				ask.Volume = 0
				ob.Asks = append(ob.Asks[:i], ob.Asks[i+1:]...)
				i -= 1
			}

			trades = append(trades, trade)
			if order.Volume == 0 {
				break
			}
		}
	}

	if order.Volume > 1 { //adding resting Bid
		if order.Price == 0 {
			return trades, order
		}
		ob.Bids = append(ob.Bids, order)
		ob.SortBids()
		//fmt.Printf("Bid added: %+v\n", *order)
	}

	return trades, nil
}

func (ob *Orderbook) LimitBid(order *Order) ([]*Trade, *Order) {
	//	println("Limit Bid called")
	trades := []*Trade{}

	for i := 0; i < len(ob.Bids); i++ {
		bid := ob.Bids[i]
		if order.Price == 0 || bid.Price <= order.Price {
			trade := &Trade{}
			trade.Price = bid.Price
			trade.Bid = bid
			trade.Ask = order

			if bid.Volume > order.Volume {
				trade.Volume = order.Volume
				bid.Volume -= order.Volume
				order.Volume = 0
			} else { //Ask shoud be removed from asks
				trade.Volume = bid.Volume
				order.Volume -= bid.Volume
				bid.Volume = 0
				ob.Bids = append(ob.Bids[:i], ob.Bids[i+1:]...)
				i -= 1
			}

			trades = append(trades, trade)
			if order.Volume == 0 {
				break
			}
		}
	}

	if order.Volume > 1 { //adding resting Bid
		if order.Price == 0 {
			return trades, order
		}
		//	fmt.Printf("Ask added: %+v\n", *order)
		ob.Asks = append(ob.Asks, order)
		ob.SortAsks()
	}

	return trades, nil
}

func (ob *Orderbook) SortBids() { //loswest first
	for i := 0; i < len(ob.Bids); i++ {
		for j := len(ob.Bids) - 1; j > i; j-- {
			if ob.Bids[j-1].Price > ob.Bids[j].Price {
				ob.Bids[j-1], ob.Bids[j] = ob.Bids[j], ob.Bids[j-1]
			}
		}
	}
}

func (ob *Orderbook) SortAsks() { //loswest first
	for i := 0; i < len(ob.Asks); i++ {
		for j := len(ob.Asks) - 1; j > i; j-- {
			if ob.Asks[j-1].Price < ob.Asks[j].Price {
				ob.Asks[j-1], ob.Asks[j] = ob.Asks[j], ob.Asks[j-1]
			}
		}
	}
}
