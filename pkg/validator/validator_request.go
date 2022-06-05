package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Entrio/subenv"
)

func MakeRequest() (*ValidatorResponse, error) {
	baseUrl := subenv.Env("LCD_URL", "http://188.40.140.51:1317")
	validatorUrl := fmt.Sprintf("%s/cosmos/staking/v1beta1/validators/%s", baseUrl, subenv.Env("VALIDATOR_ADDRESS", "terravaloper1q8w4u2wyhx574m70gwe8km5za2ptanny9mnqy3"))
	req, err := http.NewRequest("GET", validatorUrl, nil)
	if err != nil {
		//TODO: Count number of fails, block prometheus response
		return nil, errors.New("failed to create request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		//TODO: Same as above. Halting metric collection
		return nil, errors.New("failed to execute request")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("error reading response")
	}

	var vr ValidatorResponse

	if err := json.Unmarshal(body, &vr); err != nil {
		return nil, errors.New("failed to unmarshal response json")
	}

	return &vr, nil
}
