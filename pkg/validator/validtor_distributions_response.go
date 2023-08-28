package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Entrio/subenv"
)

func QueryValidatorCommissions() (*DistributionCommissionResponse, error) {
	baseUrl := subenv.Env("LCD_URL", "http://188.40.140.51:1317")
	validatorUrl := fmt.Sprintf("%s/cosmos/distribution/v1beta1/validators/%s/commission", baseUrl, subenv.Env("VALIDATOR_ADDRESSES", "terravaloper1q8w4u2wyhx574m70gwe8km5za2ptanny9mnqy3"))
	req, err := http.NewRequest("GET", validatorUrl, nil)
	if err != nil {
		//TODO: Count number of fails, block prometheus response
		return nil, errors.New("failed to create QueryValidatorCommissions request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		//TODO: Same as above. Halting metric collection
		return nil, errors.New("failed to execute QueryValidatorCommissions request")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("error reading QueryValidatorCommissions response")
	}

	var dr DistributionCommissionResponse

	if err := json.Unmarshal(body, &dr); err != nil {
		return nil, errors.New("failed to unmarshal QueryValidatorCommissions response json")
	}

	return &dr, nil
}

func QueryValidatorRewards() (*RewardsResponse, error) {
	const op = "QueryValidatorRewards"
	baseUrl := subenv.Env("LCD_URL", "http://188.40.140.51:1317")
	validatorUrl := fmt.Sprintf("%s/cosmos/distribution/v1beta1/validators/%s/outstanding_rewards", baseUrl, subenv.Env("VALIDATOR_ADDRESSES", "terravaloper1q8w4u2wyhx574m70gwe8km5za2ptanny9mnqy3"))
	req, err := http.NewRequest("GET", validatorUrl, nil)
	if err != nil {
		//TODO: Count number of fails, block prometheus response
		return nil, fmt.Errorf("failed to create %s request", op)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		//TODO: Same as above. Halting metric collection
		return nil, fmt.Errorf("failed to execute %s request", op)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading %s response", op)
	}

	var rr RewardsResponse

	if err := json.Unmarshal(body, &rr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s response json", op)
	}

	return &rr, nil
}
