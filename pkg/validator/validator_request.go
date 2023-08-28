package validator

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Entrio/subenv"
)

func GetValidatorData() *[]ValidatorResponse {
	const op = "GetValidatorData"
	log.Debug(op)
	baseUrl := subenv.Env("LCD_URL", "http://188.40.140.51:1317")

	log.Debugf("Got validators: %s", subenv.Env("VALIDATOR_ADDRESSES", "terravaloper1q8w4u2wyhx574m70gwe8km5za2ptanny9mnqy3"))
	validators := strings.Split(subenv.Env("VALIDATOR_ADDRESSES", "terravaloper1q8w4u2wyhx574m70gwe8km5za2ptanny9mnqy3"), ",")

	responses := make([]ValidatorResponse, 0)

	for _, validatorAddress := range validators {
		var response ValidatorResponse
		validatorUrl := fmt.Sprintf("%s/cosmos/staking/v1beta1/validators/%s", baseUrl, validatorAddress)
		req, err := http.NewRequest("GET", validatorUrl, nil)
		if err != nil {
			log.Warnf("Failed to create request. Error: %s", err.Error())
			continue
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Warnf("Failed to execute request. Error: %s", err.Error())
			continue
		}
		log.Debugf("Status code: %d", res.StatusCode)

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Warnf("Failed to read the response. Error: %s", err.Error())
			continue
		}

		if err := json.Unmarshal(body, &response); err != nil {
			log.Warnf("Failed to unmarshal response body. Error: %s", err.Error())
			continue
		}

		_ = res.Body.Close()
		log.Debugf("Added %s to validator list", response.Validator.OperatorAddress)
		responses = append(responses, response)
	}

	log.Debugf("%s end", op)

	return &responses
}
