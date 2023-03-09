package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ParseBody(req *http.Request, x interface{}) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(body), x)
	return err

}
