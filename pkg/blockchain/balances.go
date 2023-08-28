package blockchain

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Entrio/subenv"
)

type BalanceResponse struct {
	WalletAddress string `json:"wallet_address"`
	Balances      []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"balances"`
}

/*
Parse the environmental variable of provided wallets
*/
func getWallets() []string {
	wallets := strings.Split(subenv.Env("WALLET_ADDRESSES", "terra1gefdegujr5urhxtn4e9m9sh3sw8dy9pc24atfu"), ",")
	return wallets
}

func GetBalances() (*[]BalanceResponse, error) {
	const op = "GetBalances"
	baseUrl := subenv.Env("LCD_URL", "http://188.40.140.51:1317")

	wallets := getWallets()
	returnWallets := make([]BalanceResponse, 0)

	for _, wallet := range wallets {
		validatorUrl := fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s", baseUrl, wallet)
		log.Debugf("[%s] Fetching balance for %s", op, wallet)
		req, err := http.NewRequest("GET", validatorUrl, nil)
		if err != nil {
			//TODO: Count number of fails, block prometheus response
			return nil, fmt.Errorf("[%s][%s] failed to create request", op, wallet)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			//TODO: Same as above. Halting metric collection
			return nil, fmt.Errorf("[%s][%s] failed to execute request (%s)", op, wallet, err.Error())
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("[%s][%s] error reading response (%s)", op, wallet, err.Error())
		}

		var r BalanceResponse
		if err := json.Unmarshal(body, &r); err != nil {
			return nil, fmt.Errorf("[%s][%s] failed to unmarshal response json (%s)", op, wallet, err.Error())
		}
		r.WalletAddress = wallet
		returnWallets = append(returnWallets, r)
		_ = res.Body.Close()
	}

	return &returnWallets, nil
}
