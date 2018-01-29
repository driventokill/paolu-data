package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adshao/go-binance"
	// "time"
)

type BookItem struct {
	Bid      bool   `json:"bid"`
	Ask      bool   `json:"ask"`
	Symbol   string `json:"symbol"`
	UpdateID int64  `json:"update_id"`
	UpdateAt int64  `json:"update_at"`
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

type PriceQuantity struct {
	Price    string
	Quantity string
}

func toList(symbol string, updateAt int64, updateID int64, priceList []PriceQuantity, isAsk bool) []BookItem {
	list := make([]BookItem, len(priceList))
	for i, p := range priceList {
		item := new(BookItem)
		item.Ask = isAsk
		item.Bid = !isAsk
		item.Symbol = symbol
		item.UpdateAt = updateAt
		item.UpdateID = updateID
		item.Price = p.Price
		item.Quantity = p.Quantity
		list[i] = *item
	}
	return list
}

func main() {

	symbol := "LTCBTC"

	if len(os.Args) > 1 {
		symbol = os.Args[1]
	}

	wsDepthHandler := func(event *binance.WsDepthEvent) {
		s := event.Symbol
		uID := event.UpdateID
		updateAt := event.Time

		pqs := make([]PriceQuantity, len(event.Asks))
		for i, a := range event.Asks {
			pqs[i] = PriceQuantity(a)
		}
		items := toList(s, updateAt, uID, pqs, true)

		for _, i := range items {
			j, _ := json.Marshal(i)
			fmt.Println(string(j))
		}

		pqs1 := make([]PriceQuantity, len(event.Bids))
		for i, a := range event.Bids {
			pqs1[i] = PriceQuantity(a)
		}
		items1 := toList(s, updateAt, uID, pqs1, false)

		for _, i := range items1 {
			j, _ := json.Marshal(i)
			fmt.Println(string(j))
		}
	}
	done, err := binance.WsDepthServe(symbol, wsDepthHandler)
	if err != nil {
		fmt.Println(err)
		return
	}

	<-done

}
