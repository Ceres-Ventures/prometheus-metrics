package blockchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Entrio/subenv"
)

type BalanceResponse struct {
	WalletAddress string `json:"wallet_address"`
	Balances      []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"balances"`
}

func GetBalances() (*BalanceResponse, error) {
	const op = "GetBalances"
	baseUrl := subenv.Env("LCD_URL", "http://188.40.140.51:1317")
	validatorUrl := fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s", baseUrl, subenv.Env("WALLET_ADDRESSES", "terra1gefdegujr5urhxtn4e9m9sh3sw8dy9pc24atfu"))
	req, err := http.NewRequest("GET", validatorUrl, nil)
	if err != nil {
		//TODO: Count number of fails, block prometheus response
		return nil, fmt.Errorf("[%s] failed to create request", op)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		//TODO: Same as above. Halting metric collection
		return nil, fmt.Errorf("[%s] failed to execute request (%s)", op, err.Error())
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("[%s] error reading response (%s)", op, err.Error())
	}

	var r BalanceResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("[%s] failed to unmarshal response json (%s)", op, err.Error())
	}
	r.WalletAddress = subenv.Env("WALLET_ADDRESSES", "terra1gefdegujr5urhxtn4e9m9sh3sw8dy9pc24atfu")

	return &r, nil
}
