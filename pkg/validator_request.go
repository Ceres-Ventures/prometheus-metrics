package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Entrio/subenv"
	"github.com/ceres-ventures/prometheus-metrics/pkg/models"
)

func MakeRequest() (*models.ValidatorResponse, error) {
	baseUrl := subenv.Env("LCD_URL", "")
	validatorUrl := fmt.Sprintf("%s/%s", baseUrl, subenv.Env("VALIDATOR_ADDRESS", ""))
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

	var vr models.ValidatorResponse

	if err := json.Unmarshal(body, &vr); err != nil {
		return nil, errors.New("failed to unmarshal response json")
	}

	return nil, nil
}
