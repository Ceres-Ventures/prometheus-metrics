package blockchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/Entrio/subenv"
)

var (
	OwnValidatorHashAddress = subenv.Env("VALIDATOR_KEY_ADDRESS", "C24A7D204E0A07736EAF8A7E76820CD868565B0E")
)

type (
	LatestBlockResponse struct {
		BlockId struct {
			Hash string `json:"hash"`
		} `json:"block_id"`
		Block struct {
			Header struct {
				Height   string    `json:"height"`
				Time     time.Time `json:"time"`
				Proposer string    `json:"proposer_address"`
				Chain    string    `json:"chain_id"`
			} `json:"header"`
			Data struct {
				Transactions []string `json:"txs"`
			} `json:"data"`
		} `json:"block"`
	}
	Validator struct {
		Address string `json:"address"`
	}
)

// GetNumOfTxs returns the number of transaction in the current block
func (lbr LatestBlockResponse) GetNumOfTxs() int {
	return len(lbr.Block.Data.Transactions)
}

func (lbr LatestBlockResponse) GetBlockHeightInt() int {
	v, _ := strconv.Atoi(lbr.Block.Header.Height)
	return v
}

func GetLatestBlockData() (*LatestBlockResponse, error) {
	const op = "GetLatestBlockData"
	baseUrl := subenv.Env("LCD_URL", "http://188.40.140.51:1317")
	validatorUrl := fmt.Sprintf("%s/blocks/latest", baseUrl)
	req, err := http.NewRequest("GET", validatorUrl, nil)
	if err != nil {
		//TODO: Count number of fails, block prometheus response
		return nil, fmt.Errorf("[%s] failed to create request", op)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		//TODO: Same as above. Halting metric collection
		return nil, fmt.Errorf("[%s] failed to execute request", op)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("[%s] error reading response", op)
	}

	var r LatestBlockResponse

	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("[%s] failed to unmarshal response json", op)
	}

	return &r, nil
}
