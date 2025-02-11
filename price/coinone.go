package price

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TradeHistory struct {
	Trades []Trade `json:"trades"`
}

type Trade struct {
	Timestamp uint64 `json:"timestamp"`
	Price string `json:"price"`
	Volume string `json:"volume"`
	IsSellerMaker bool `json:"is_seller_maker"`
}

func (ps *PriceService) coinoneToLuna(logger log.Logger) {
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("Unknown error", r)
				}

				time.Sleep(5 * time.Second)
			}()

			resp, err := http.Get("https://tb.coinone.co.kr/api/v1/tradehistory/recent/?market=krw&target=luna")
			if (err != nil) {
				logger.Error("Fail to fetch from coinone", err.Error())
				return
			}
			defer func() {
				resp.Body.Close()
			}()
			body, err := ioutil.ReadAll(resp.Body)
			if (err != nil) {
				logger.Error("Fail to read body", err.Error())
				return
			}

			th := TradeHistory{}
			err = json.Unmarshal(body, &th)
			if (err != nil) {
				logger.Error("Fail to unmarshal json", err.Error())
				return
			}

			trades := th.Trades
			recent := trades[len(trades) - 1]
			logger.Info(fmt.Sprintf("Recent luna/krw: %s, timestamp: %d", recent.Price, recent.Timestamp))
			
			amount, ok := sdk.NewIntFromString(recent.Price)
			if !ok {
				logger.Error("Fail to parse price to int")
			}
			ps.setPrice("luna/krw", sdk.NewDecCoinFromCoin(sdk.NewCoin("krw", amount)))
		}()
	}
}
