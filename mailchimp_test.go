package mailchimp

import "os"
import "testing"
//import "strings"
//import "time"

var CID = os.Getenv("MAILCHIMPCID")
var RSS = os.Getenv("MAILCHIMPRSS")

var chimp, err = New(os.Getenv("MAILCHIMPKEY"), true)

var schedule = make(chan string, 1)
var unschedule = make(chan string, 1)
var update = make(chan string, 1)
var del = make(chan string, 1)

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

/*
func TestCampaignContent(t *testing.T) {
	parameters := make(map[string]interface{})
	parameters["cid"] = os.Getenv("MAILCHIMPCID")
	parameters["for_archive"] = true
	result, err := chimp.CampaignContent(parameters)
	if err != nil {
		t.Error("mailchimp.CampaignContent:", err)
	}
	if !strings.Contains(result.Html, "<head>") {
		t.Error("mailchimp.CampaignContent: the Html field of the returned struct does not look like html")
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
	cid, err := chimp.CampaignCreate(parameters)
	if err != nil {
		t.Error("mailchimp.CampaignsCreate:", err)
	}
	//send to TestCampaignSchedule
	schedule <- cid
}
*/

/*
func TestCampaignDelete(t *testing.T) {
//wait for CampaignUnschedule
	go func() {
		parameters := make(map[string]interface{})
		cid := <- del
		parameters["cid"] = cid
		_, err := chimp.CampaignDelete(parameters)
		if err != nil {
			t.Error("mailchimp.CampaignsDelete:", err)
		}
	}()
}
*/

func TestCampaignEcommOrderAdd(t *testing.T) {
//untested
}

/*
func TestCampaignPause(t *testing.T) {
	parameters := make(map[string]interface{})
	parameters["cid"] = os.Getenv("MAILCHIMPRSS")
	_, err = chimp.CampaignPause(parameters)
	if err != nil {
		if err.(ChimpError).Err != `Cannot pause this campaign because it is currently "paused"` {
			t.Error("mailchimp.CampaignsPause:", err)
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
		t.Error("mailchimp.CampaignReplicate", err)
	}
}
*/

/*
func TestCampaignResume(t *testing.T) {
	parameters := make(map[string]interface{})
	parameters["cid"] = os.Getenv("MAILCHIMPRSS")
	_, err = chimp.CampaignResume(parameters)
	if err != nil {
		if err.(ChimpError).Err != `Cannot resume this campaign because it is currently "sending"` {
			t.Error("mailchimp.CampaignResume:", err)
		}
	}
}
*/

/*
func TestScheduleCampaign(t *testing.T) {
	//waits for CampaignCreate
	parameters := make(map[string]interface{})
	cid := <- schedule
	parameters["cid"] = cid
	location, err := time.LoadLocation("Local")
	parameters["schedule_time"] = time.Date(2012, 12, 25, 17, 30, 0, 0, location)
	_, err = chimp.CampaignSchedule(parameters)
	if err != nil {
		t.Error("mailchimp.CampaignSchedule:", err)
	}
	unschedule <- cid
}
*/

/*
func TestCampaignSegmentTest(t *testing.T) {
	parameters := make(map[string]interface{})
	parameters["list_id"] = os.Getenv("MAILCHIMPLIST")
	options := make(map[string]interface{})
	options["match"] = "any"
	conditions := make([]map[string]interface{}, 0, 1)
	condition := make(map[string]interface{})
	condition["field"] = "date"
	condition["op"] = "lt"
	condition["value"] = "last_campaign_sent"
	options["conditions"] = append(conditions, condition)
	parameters["options"] = options
	i, err := chimp.CampaignSegmentTest(parameters)
	if err != nil {
		t.Error("mailchimp.CampaignSegmentTest:", err)
	}
	if i <= 0 {
		t.Error("mailchimp.CampaignSegmentTest: expected count to be positive but got", i)
	}
}
*/

/*
func TestCampaignSendNow(t *testing.T) {
	//create a campaign to send
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
	//now send it
	response, err := chimp.CampaignSendNow(map[string]interface{}{"cid": id})
	if err != nil {
		t.Error("mailchimp.CampaignSendNow:", err)
	}
	if !response {
		t.Error("mailchimp.CampaignSendNow failed to send")
	}
	//too soon to delete since it was just sent delete will fail
	//manually delete later from account
}
*/

/*
func TestCampaignSendTest(t *testing.T) {
	//create a campaign to send
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
	//now send it
	p := make(map[string]interface{})
	p["cid"] = id
	p["test_emails"] =[]string{"areed@partitus.com"}
	response, err := chimp.CampaignSendTest(p)
	if err != nil {
		t.Error("mailchimp.CampaignSendTest:", err)
	}
	if !response {
		t.Error("mailchimp.CampaignSendTest failed to send")
	}
	chimp.CampaignDelete(map[string]interface{}{"cid": id})
}
*/

/*
func TestCampaignShareReport(t *testing.T) {
	result, err := chimp.CampaignShareReport(map[string]interface{}{"cid": CID})
	if err != nil {
		t.Error("mailchimp.CampaignShareReport:", err)
	}
	if !strings.HasPrefix(result.Url, "http") {
		t.Error("Expected result.Url to be a valid url but got", result.Url)
	}
}
*/

/*
func TestCampaignTemplateContent(t *testing.T) {
	_, err := chimp.CampaignTemplateContent(map[string]interface{}{"cid": CID})
	if err != nil {
		t.Error("mailchimp.CampaignTemplateContent", err)
	}
}
*/

/*
func TestCampaignUnschedule(t *testing.T) {
	//waits for CampaignSchedule
	parameters := make(map[string]interface{})
	cid := <- unschedule
	parameters["cid"] = cid
	_, err := chimp.CampaignUnschedule(parameters)
	if err != nil {
		t.Error(err)
	}
	update <- cid
}
*/

/*
func TestCampaignUpdate(t *testing.T) {
	//waits for CampaignUnschedule
	parameters := make(map[string]interface{})
	cid := <- update
	parameters["cid"] = cid
	parameters["name"] = "from_email"
	parameters["value"] = "areed@partitus.com"
	_, err := chimp.CampaignUpdate(parameters)
	if err != nil {
		t.Error(err)
	}
	del <- cid
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
func TestCampaignAbuseReports(t *testing.T) {
	_, err := chimp.CampaignAbuseReports(map[string]interface{}{"cid": CID})
	if err != nil {
		t.Error("mailchimp.CampaignAbuseReports:", err)
	}
}
*/

/*
func TestCampaignAdvice(t *testing.T) {
	_, err := chimp.CampaignAdvice(map[string]interface{}{"cid": CID})
	if err != nil {
		t.Error("mailchimp.CampaignAdvice", err)
	}
}
*/

/*
func TestCampaignAnalytics(t *testing.T) {
	_, err := chimp.CampaignAnalytics(map[string]interface{}{"cid": CID})
	if err != nil {
		if err.(ChimpError).Err != `Google Analytics Add-on required for this function` {
			t.Error("mailchimp.CampaignAnalytics", err)
		}
	}
}
*/

/*
func TestCampaignBounceMessage(t *testing.T) {
	_, err := chimp.CampaignBounceMessage(map[string]interface{}{"cid": CID, "email": "areed@partitus.com"})
	if err != nil {
		if err.(ChimpError).Code != 319 {
			t.Error("mailchimp.CampaignBounceMessage", err)
		}
	}
}
*/

/*
func TestCampaignBounceMessages(t *testing.T) {
	_, err := chimp.CampaignBounceMessages(map[string]interface{}{"cid": CID})
	if err != nil {
		t.Error("mailchimp.CampaignBounceMessages", err)
	}
}
*/

/*
func TestCampaignClickStats(t *testing.T) {
	_, err := chimp.CampaignClickStats(map[string]interface{}{"cid": CID})
	if err != nil {
		t.Error("mailchimp.CampaignClickStats", err)
	}
}
*/

/*
func TestCampaignEcommOrders(t *testing.T) {
	_, err := chimp.CampaignEcommOrders(map[string]interface{}{"cid": CID})
	if err != nil {
		t.Error("mailchimp.CampaignEcommOrders", err)
	}
}
*/

/*
func TestCampaignEepUrlStats(t *testing.T) {
	_, err := chimp.CampaignEepUrlStats(map[string]interface{}{"cid": CID})
	if err != nil {
		t.Error("mailchimp.CampaignEepUrlStats", err)
	}
}
*/

/*
func TestCampaignEmailDomainPerformance(t *testing.T) {
	_, err := chimp.CampaignEmailDomainPerformance(map[string]interface{}{"cid": CID})
	if err != nil {
		t.Error("mailchimp.CampaignEmailDomainPerformance", err)
	}
}
*/

/*
func TestCampaignGeoOpens(t *testing.T) {
	_, err := chimp.CampaignGeoOpens(map[string]interface{}{"cid": CID})
	if err != nil {
		t.Error("mailchimp.CampaignEmailGeoOpens", err)
	}
}
*/

/*
func TestCampaignGeoOpensForCountry(t *testing.T) {
	_, err := chimp.CampaignGeoOpensForCountry(map[string]interface{}{"cid": CID, "code": "US"})
	if err != nil {
		t.Error("mailchimp.CampaignGeoOpensForCountry", err)
	}
}
*/

/*
func TestCampaignMembers(t *testing.T) {
	result, err := chimp.CampaignMembers(map[string]interface{}{"cid": CID})
	if err != nil {
		t.Error("mailchimp.CampaignMembers", err)
	}
	if result.Total <= 0 {
		t.Error("mailchimp.CampaignMembers: Expected total to be positive but got", result.Total)
	}
	if len(result.Data[0].Email) < 5 {
		t.Error("mailchimp.CampaignMembers: First returned email address appears invalid:", result.Data[0].Email)
	}
}
*/

/*
func TestCampaignStats(t *testing.T) {
	result, err := chimp.CampaignStats(map[string]interface{}{"cid": CID})
	if err != nil {
		t.Error("mailchimp.CampaignStats", err)
	}
	if result.Emails_sent <= 0 {
		t.Error("mailchimp.CampaignStats: expected emails_sent to be positive but got", result.Emails_sent)
	}
}
*/

/*
func TestCampaignUnsubscribes(t *testing.T) {
	_, err := chimp.CampaignUnsubscribes(map[string]interface{}{"cid": CID})
	if err != nil {
		t.Error("mailchimp.CampaignUnsubscribes", err)
	}
}
*/
