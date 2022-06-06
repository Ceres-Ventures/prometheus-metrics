package blockchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Entrio/subenv"
)

type (
	BankSupplyResponse struct {
		Supply []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"supply"`
	}
)

func (bsr *BankSupplyResponse) GetULunaSupply() float64 {
	supply := float64(0)

	for _, v := range bsr.Supply {
		if v.Denom == "uluna" {
			val, e := strconv.ParseFloat(v.Amount, 64)
			if e == nil {
				supply = val
			}
		}
	}

	return supply
}

func (bsr *BankSupplyResponse) GetLunaSupply() float64 {
	return bsr.GetULunaSupply() / 1000000
}

func GetLunaSupply() (*BankSupplyResponse, error) {
	const op = "GetLunaSupply"
	var r BankSupplyResponse

	baseUrl := subenv.Env("LCD_URL", "http://188.40.140.51:1317")
	reUrl := fmt.Sprintf("%s/cosmos/bank/v1beta1/supply", baseUrl)
	req, err := http.NewRequest("GET", reUrl, nil)
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

	return &r, nil
}
