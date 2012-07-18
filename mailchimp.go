package mailchimp

import (
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
	"regexp"
	"io/ioutil"
)

func main() {

}

type API struct {
	Version string
	Method  string
	Key string
	Parameters map[string]interface{}
}

func (api *API) Run() (interface{}, error) {
	url := endpoint(api)
	api.Parameters["apikey"] = api.Key
	b, err := json.Marshal(api.Parameters)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
	var f interface{}
	err = json.Unmarshal(body, &f)
	m := f.(map[string]interface{})
	return m, nil
}

func endpoint(api *API) string {
	dc := regexp.MustCompile("[a-z]+[0-9]+$")
	return fmt.Sprintf("http://%s.api.mailchimp.com/%s/?method=%s", dc.FindString(api.Key), api.Version, api.Method)
}
