package mailchimp

import (
	"os"
	"testing"
)

var chimp, err = New(os.Getenv("MAILCHIMPKEY"), true)
var cid = make(chan string, 1)

/*
func TestCampaignContent(t *testing.T) {
	parameters := make(map[string]interface{})
	parameters["cid"] = os.Getenv("MAILCHIMPCID")
	parameters["for_archive"] = true
	_, err = chimp.CampaignContent(parameters)
	if err != nil {
		t.Error("mailchimp.CampaignContent:", err)
	}
}
*/

/*
func TestCampaignCreate(t *testing.T) {
	parameters := make(map[string]interface{})
	parameters["type"] = "regular"
	options := make(map[string]interface{})
	options["list_id"] = os.Getenv("MAILCHIMPLIST")
	options["subject"] = "Go API test"
	options["from_email"] = "support@partitus.com"
	options["from_name"] = "Partitus"
	options["to_name"] = "*|FNAME|*"
	parameters["options"] = options
	content := make(map[string]interface{})
	content["html"] = "<p>Go API Test campaign html content</p>"
	content["text"] = "Go API Test campaign text content"
	parameters["content"] = content
	id, err := chimp.CampaignCreate(parameters)
	if err != nil {
		t.Error("mailchimp.CampaignsCreate:", err)
	}
	//send to TestCampaignDelete
	cid <- id
}
*/

/*
func TestCampaignDelete(t *testing.T) {
	parameters := make(map[string]interface{})
	//block until the cid for the campaign created in
	//TestCampaignCreate is received
	parameters["cid"] = <- cid
	//TODO: figure out why this is happening
	//for some reason double quotes are part of the string
	//and the actual 10 digit id needs to be extracted
	parameters["cid"] = parameters["cid"].(string)[1:11]
	_, err := chimp.CampaignDelete(parameters)
	if err != nil {
		t.Error("mailchimp.CampaignsDelete:", err)
		t.Error(parameters)
	}
}
//*/

func TestCampaignEcommOrderAdd(t *testing.T) {

}

/*
func TestCampaignPause(t *testing.T) {
	parameters := make(map[string]interface{})
	parameters["cid"] = os.Getenv("MAILCHIMPRSS")
	_, err = chimp.CampaignPause(parameters)
	if err != nil {
		if err.(ChimpError).Err != `Cannot pause this campaign because it is currently "paused"` {
			t.Error(err)
		}
	}
}
*/

/*
func TestCampaignReplicate(t *testing.T) {
	parameters := make(map[string]interface{})
	parameters["cid"] = os.Getenv("MAILCHIMPCID")
	_, err = chimp.CampaignReplicate(parameters)
	if err != nil {
		t.Error(err)
	}
}
*/

/*
func TestCampaigns(t *testing.T) {
	parameters := make(map[string]interface{})
	filters := make(map[string]interface{})
	filters["status"] = "sent"
	parameters["filters"] = filters
	result, err := chimp.Campaigns(parameters)
	if err != nil {
		t.Error("mailchimp.Campaigns:", err)
	}
	if result.Data[0].Status != "sent" {
		t.Error("mailchimp.Campaigns: json response did not properly unmarshal in CampaignsResult struct")
	}
}
*/

/*
func TestPing(t *testing.T) {
	result, err := chimp.Ping()
	if err != nil {
		t.Error("mailchimp.Ping", err)
	}
	if result != "Everything's Chimpy!" {
		t.Error(`Expected response "Everything's Chimpy!" but received`, result)
	}
}
*/


