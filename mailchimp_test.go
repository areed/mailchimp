package mailchimp

import (
	"os"
	"testing"
)

func TestCampaigns(t *testing.T) {
	apikey := os.Getenv("MAILCHIMPKEY")
	parameters := make(map[string]interface{})
	filters := make(map[string]interface{})
	filters["status"] = "sent"
	parameters["filters"] = filters
	api := API{"1.3", "campaigns", apikey, parameters}
	_, err := api.Run()
	if err != nil {
		t.Fail()
	}
}
