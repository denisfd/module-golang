package orderbook

import "container/list"

type Orderbook struct {
	LastID int
	Bids   *list.List
	Asks   *list.List
}

func New() *Orderbook {
	ob := &Orderbook{}

	ob.Bids = list.New()
	ob.Asks = list.New()

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

/*func (ob *Orderbook) Cancel(ID int) bool {
	for i, order := range ob.Bids {
		if order.ID == ID {
			//ob.Bids = append(ob.Bids[:i], ob.Bids[i+1:]...)
			return true
		}
	}

	for i, order := range ob.Asks {
		if order.ID == ID {
			//ob.Asks = append(ob.Asks[:i], ob.Asks[i+1:]...)
			return true
		}
	}

	return false
}*/

func (ob *Orderbook) LimitAsk(order *Order) ([]*Trade, *Order) {
	trades := []*Trade{}
	var next *list.Element

	for el := ob.Asks.Front(); el != nil; el = next {
		next = el.Next()
		ask := el.Value.(*Order)
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
				ob.Asks.Remove(el)
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
	var next *list.Element

	for el := ob.Bids.Front(); el != nil; el = next {
		next = el.Next()
		bid := el.Value.(*Order)
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
				ob.Bids.Remove(el)
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

func (ob *Orderbook) AddBid(new *Order) {
	var bid *Order
	for el := ob.Bids.Front(); el != nil; el = el.Next() {
		bid = el.Value.(*Order)
		if new.Price < bid.Price {
			ob.Bids.InsertBefore(new, el)
			return
		}
	}
	ob.Bids.PushBack(new)
}

func (ob *Orderbook) AddAsk(new *Order) {
	var ask *Order
	for el := ob.Asks.Front(); el != nil; el = el.Next() {
		ask = el.Value.(*Order)
		if new.Price > ask.Price {
			ob.Asks.InsertBefore(new, el)
			return
		}
	}
	ob.Asks.PushBack(new)
}
