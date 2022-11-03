package external

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Entrio/subenv"
	"github.com/sirupsen/logrus"
)

type (
	CGPriceResponse struct {
		ID         string `json:"id"`
		Symbol     string `json:"symbol"`
		MarketData struct {
			CurrentPrice struct {
				AUD float64 `json:"aud"`
				USD float64 `json:"usd"`
			} `json:"current_price"`
			ATH struct {
				AUD float64 `json:"aud"`
				USD float64 `json:"usd"`
			} `json:"ath"`
			ATL struct {
				AUD float64 `json:"aud"`
				USD float64 `json:"usd"`
			} `json:"atl"`
			High24 struct {
				AUD float64 `json:"aud"`
				USD float64 `json:"usd"`
			} `json:"high_24h"`
			Low24 struct {
				AUD float64 `json:"aud"`
				USD float64 `json:"usd"`
			} `json:"low_24h"`
		} `json:"market_data"`
	}
)

func GetLunaMarketData() (*CGPriceResponse, error) {
	const op = "GetLunaSupply"
	var r CGPriceResponse

	baseUrl := subenv.Env("MARKET_URL", "https://api.coingecko.com/api/v3/coins/terra-luna-2?tickers=false&market_data=true&community_data=false&developer_data=false&sparkline=false")
	//reUrl := fmt.Sprintf("%s/cosmos/bank/v1beta1/supply", baseUrl)
	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("[%s] failed to create request", op)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[%s] failed to execute request", op)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("[%s] error reading response", op)
	}

	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("[%s] failed to unmarshal response json", op)
	}

	logrus.Infof("Requesto to %s was a success", baseUrl)
	if res.StatusCode == 429 || r.MarketData.CurrentPrice.USD == 0 {
		return nil, nil
	}
	return &r, nil
}
