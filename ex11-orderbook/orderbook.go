package orderbook

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
	switch order.Side {
	case SideAsk:
		return orderbook.LimitAsk(order)
	case SideBid:
		return orderbook.LimitBid(order)
	}

	return nil, nil
}

func (ob *Orderbook) Cancel(ID int) bool {
	for i, order := range ob.Bids {
		if order.ID == ID {
			ob.Bids = append(ob.Bids[:i], ob.Bids[i+1:]...)
			return true
		}
	}

	for i, order := range ob.Asks {
		if order.ID == ID {
			ob.Asks = append(ob.Asks[:i], ob.Asks[i+1:]...)
			return true
		}
	}

	return false
}

func (ob *Orderbook) LimitAsk(order *Order) ([]*Trade, *Order) {
	trades := []*Trade{}

	for i := 0; i < len(ob.Asks); i++ {
		ask := ob.Asks[i]
		if order.Price == 0 || order.Price <= ask.Price {
			trade := &Trade{
				Price: ask.Price,
				Bid:   order,
				Ask:   ask,
			}

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
		} else {
			break
		}
	}

	if order.Volume > 1 { //adding resting Bid
		if order.Price == 0 {
			return trades, order
		}
		ob.AddBid(order)
	}

	return trades, nil
}

func (ob *Orderbook) LimitBid(order *Order) ([]*Trade, *Order) {
	trades := []*Trade{}

	for i := 0; i < len(ob.Bids); i++ {
		bid := ob.Bids[i]
		if order.Price == 0 || bid.Price <= order.Price {
			trade := &Trade{
				Price: bid.Price,
				Bid:   bid,
				Ask:   order,
			}

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
		} else {
			break
		}
	}

	if order.Volume > 1 { //adding resting Bid
		if order.Price == 0 {
			return trades, order
		}
		ob.AddAsk(order)
	}

	return trades, nil
}

func (ob *Orderbook) AddBid(bid *Order) {
	ob.Bids = append(ob.Bids, bid)
	for i := len(ob.Bids) - 1; i > 0; i-- {
		if ob.Bids[i].Price < ob.Bids[i-1].Price {
			ob.Bids[i], ob.Bids[i-1] = ob.Bids[i-1], ob.Bids[i]
		} else {
			break
		}
	}
}

func (ob *Orderbook) AddAsk(ask *Order) {
	ob.Asks = append(ob.Asks, ask)
	for i := len(ob.Asks) - 1; i > 0; i-- {
		if ob.Asks[i].Price > ob.Asks[i-1].Price {
			ob.Asks[i], ob.Asks[i-1] = ob.Asks[i-1], ob.Asks[i]
		} else {
			break
		}
	}
}
