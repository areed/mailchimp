package mailchimp

import "bufio"
import "path/filepath"
import "os"
import "testing"
import "encoding/json"

//import "bytes"
//import "strings"
//import "time"
//import "fmt"

var CID = os.Getenv("MAILCHIMPCID")
var RSS = os.Getenv("MAILCHIMPRSS")
var EMAIL = os.Getenv("MAILCHIMPEMAIL")
var STORE = os.Getenv("MAILCHIMPSTORE")
var LIST = os.Getenv("MAILCHIMPLIST")

var schedule = make(chan string, 1)
var unschedule = make(chan string, 1)
var update = make(chan string, 1)
var del = make(chan string, 1)
var orderChannel = make(chan string, 1)
var folderUpdate = make(chan int, 1)
var folderDel = make(chan int, 1)
var interestGroup = make(chan bool, 1)
var interestGrouping = make(chan int, 1)

var chimp, err = New(os.Getenv("MAILCHIMPKEY"), true)

func populate(filename string, response interface{}) error {
	filename = filepath.Join("json", filename+".json")
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	b := make([]byte, 0)
	for {
		line, isPrefix, _ := reader.ReadLine()
		b = append(b, line...)
		if !isPrefix {
			break
		}
	}
	//some response types have an alterJson method that needs
	//to be called before unmarshalling
	switch r := response.(type) {
	case alterJsoner:
		json.Unmarshal(r.alterJson(b), response)
	default:
		json.Unmarshal(b, response)
	}
	return nil
}

func verify(t *testing.T, name string, expected interface{}, actual interface{}) {
	if actual != expected {
		t.Errorf("%s: expected %v but actual value was %d", name, expected, actual)
	}
	return
}

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

/*
func TestCampaignEcommOrderAdd(t *testing.T) {
	//untested
}
*/

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

/*
func TestCampaignClickDetailAIM(t *testing.T) {
	result, err := chimp.CampaignClickDetailAIM(map[string]interface{}{"cid": CID, "url": "http://example.com"})
	if err != nil {
		t.Error("mailchimp.CampaignClickDetailAIM", err)
	}
	t.Error(result)
}
*/

/*
func TestCampaignEmailStatsAIM(t *testing.T) {
	result, err := chimp.CampaignEmailStatsAIM(map[string]interface{}{"cid": CID, "email_address": "areed@partitus.com"})
	if err != nil {
		t.Error("mailchimp.CampaignEmailStatsAIM", err)
	}
	if result.Success != 1 {
		t.Error("mailchimp.CampaignEmailStatsAIM: expected to find 1 matching email but found", result.Success)
	}
}
*/

/*
func TestCampaignEmailStatsAIMAll(t *testing.T) {
	result, err := chimp.CampaignEmailStatsAIMAll(map[string]interface{}{"cid": CID})
	if err != nil {
		t.Error("mailchimp.CampaignEmailStatsAIMAll", err)
	}
	if _, ok := result.Data["areed@partitus.com"]; !ok {
		t.Error("mailchimp.CampaignEmailStatsAIMAll: expected to get data for areed@partitus.com but didn't")
	}
}
*/

/*
func TestCampaignNotOpenedAIM(t *testing.T) {
	result, err := chimp.CampaignNotOpenedAIM(map[string]interface{}{"cid": CID})
	if err != nil {
		t.Error("mailchimp.CampaignNotOpenedAIM", err)
	}
	//this test will fail if a campaign has more unopens than the default page size of 1000
	if len(result.Data) != result.Total {
		t.Error("mailchimp.CampaignNotOpened: the length of the array of email addresses should equal Total count")
	}
}
*/

/*
func TestCampaignOpenedAIM(t *testing.T) {
	result, err := chimp.CampaignOpenedAIM(map[string]interface{}{"cid": CID})
	if err != nil {
		t.Error("mailchimp.CampaignOpenedAIM", err)
	}
	//this test will fail if a campaign has more opens than the default page size of 1000
	if len(result.Data) != result.Total {
		t.Error("mailchimp.CampaignOpenedAIM: the length of the array of email addresses should equal Total count")
	}
}
*/

/*
func TestEcommOrderAdd(t *testing.T) {
	parameters := make(map[string]interface{})
	order := make(map[string]interface{})
	//need a unique order id each time test is run
	order_id := fmt.Sprint(time.Now())
	order["id"] = order_id
	order["email"] = EMAIL
	order["total"] = 100.10
	order["store_id"] = STORE
	items := make([]map[string]interface{}, 1)
	item := make(map[string]interface{})
	item["product_id"] = "1000"
	item["product_name"] = "widget_1"
	item["category_id"] = 10
	item["category_name"] = "widgets"
	item["qty"] = 10.0
	item["cost"] = 10.01
	items = append(items, item)
	order["items"] = items
	parameters["order"] = order
	go func() {
		result, err := chimp.EcommOrderAdd(parameters)
		if err != nil {
			t.Error("mailchimp.EcommOrderAdd",err)
		}
		if !result {
			t.Error("mailchimp.EcommOrderAdd: expected return value to be true but got", result)
		} else {
			orderChannel <- order_id
		}
	}()
}
*/

/*
func TestEcommOrderDel(t *testing.T) {
	go func() {
		parameters := make(map[string]interface{})
		parameters["store_id"] = STORE
		parameters["order_id"] = <- orderChannel
		result, err := chimp.EcommOrderDel(parameters)
		if err != nil {
			t.Error("mailchimp.EcommOrderDel",err)
		}
		if !result {
			t.Error("mailchimp.EcommOrderDel: expected return value to be true but got", result)
		}
	}()
}
*/

/*
func TestEcommOrders(t *testing.T) {
	result, err := chimp.EcommOrders(nil)
	if err != nil {
		t.Error("mailchimp.EcommOrders", err)
	}
	if result.Data[0].Lines[0].Line_num != 1 {
		t.Error("mailchimp.EcommOrders: expected first line_num of first order returned to be 1 but got", result.Data[0].Lines[0].Line_num)
	}
}
*/

/*
func TestFolderAdd(t *testing.T) {
	result, err := chimp.FolderAdd(map[string]interface{}{"name": "TesterFolder"})
	if err != nil {
		t.Error("mailchimp.FolderAdd", err)
	}
	if result <= 0 {
		t.Error("Expected the folder_id to be a positive integer but got", result)
	}
	folderUpdate <- result
}

func TestFolderUpdate(t *testing.T) {
	fid := <- folderUpdate
	result, err := chimp.FolderUpdate(map[string]interface{}{"fid": fid, "name": "UpdatedTesterFolder"})
	if err != nil {
		t.Error("mailchimp.FolderUpdate", err)
	}
	if !result {
		t.Error("mailchimp.FolderUpdate: expected result to be true but got", result)
	}
	folderDel <- fid
}

func TestFolderDel(t *testing.T) {
	result, err := chimp.FolderDel(map[string]interface{}{"fid": <- folderDel})
	if err != nil {
		t.Error("mailchimp.FolderDel", err)
	}
	if !result {
		t.Error("mailchimp.FolderDel expected result to be true but got", result)
	}
}
*/

/*
func TestFolders(t *testing.T) {
	result, err := chimp.Folders(nil)
	if err != nil {
		t.Error("mailchimp.Folders", err)
	}
	if result[0].Folder_id <= 0 {
		t.Error("mailchimp.Folders: expected first Folder_id to be a positive integer but got", result[0].Folder_id)
	}
}
*/

/*
func TestGmonkeyActivity(t *testing.T) {
	_, err := chimp.GmonkeyActivity(nil)
	if err != nil {
		t.Error("mailchimp.GmonkeyActivity", err)
	}
}
*/

/*
func TestGmonkeyAdd(t *testing.T) {
	parameters := make(map[string]interface{})
	parameters["id"] = LIST
	parameters["email_address"] = []string{EMAIL}
	result, err := chimp.GmonkeyAdd(parameters)
	if err != nil {
		t.Error("mailchimp.GmonkeyAdd", err)
	}
	if result.Success + result.Errors != 1 {
		t.Errorf("mailchimp.GmonkeyAdd: expected either 1 success or 1 error but got %d successes and %d errors", result.Success, result.Errors)
	}
	if result.Errors != len(result.Data) {
		t.Error("mailchimp:GmonkeyAdd: There should be one item in the data array for each error")
	}
}

func TestGmonkeyDel(t *testing.T) {
	parameters := make(map[string]interface{})
	parameters["id"] = LIST
	parameters["email_address"] = []string{EMAIL}
	result, err := chimp.GmonkeyDel(parameters)
	if err != nil {
		t.Error("mailchimp.GmonkeyDel", err)
	}
	if result.Success + result.Errors != 1 {
		t.Errorf("mailchimp.GmonkeyDel: expected either 1 success or 1 error but got %d successes and %d errors", result.Success, result.Errors)
	}
	if result.Errors != len(result.Data) {
		t.Error("mailchimp:GmonkeyDel: There should be one item in the data array for each error")
	}
}

func TestGmonkeyMembers(t *testing.T) {
	_, err := chimp.GmonkeyMembers(nil)
	if err != nil {
		t.Error("mailchimp.GmonkeyMembers", err)
	}
}
*/

/*
func TestCampaignsForEmail(t *testing.T) {
	result, err := chimp.CampaignsForEmail(map[string]interface{}{"email_address": EMAIL})
	if err != nil {
		t.Error("mailchimp.CampaignsForEmail", err)
	}
	if len(result[0]) < 5 {
		t.Error("mailchimp.CampaignsForEmail: expected a 10-character campaign ID but got", result[0])
	}
}
*/

/*
func TestChimpChatter(t *testing.T) {
	result, err := chimp.ChimpChatter(nil)
	if err != nil {
		t.Error("mailchimp.ChimpChatter", err)
	}
	if len(result[0].Message) <= 5 {
		t.Error("mailchimp.ChimpChatter: first message looks too short to be a real message:", result[0].Message)
	}
}
*/

/*
func TestGenerateText(t *testing.T) {
	result, err := chimp.GenerateText(map[string]interface{}{"type": "html", "content": "<div><p>Paragraph in a div</p></div>"})
	if err != nil {
		t.Error("mailchimp.GenerateText", err)
	}
	if !strings.HasPrefix(result, "Paragraph in a div") {
		t.Error("mailchimp.GenerateText: Expected result of `Paragraph in a div` but got", result)
	}
}
*/

/*
func TestGetAccountDetails(t *testing.T) {
	result, err := chimp.GetAccountDetails(nil)
	if err != nil {
		t.Error("mailchimp.GetAccountDetails", err)
	}
	if result.Contact.Email != EMAIL {
		t.Error("mailchimp.GetAccountDetails: was expecting a differenct contact email address")
	}
}
*/

/*
func TestGetVerifiedDomains(t *testing.T) {
	result, err := chimp.GetVerifiedDomains(nil)
	if err != nil {
		t.Error("mailchimp.GetVerifiedDomains")
	}
	if result[0].Status != "verified" && result[0].Status != "pending" {
		t.Error("mailchimp.GetVerifiedDomains: Expected first status to be \"verified\" or \"pending\" but got", result[0].Status)
	}
}
*/

/*
Not Sure what the html parameter is supposed to be, so this method is untested
func TestInlineCss(t *testing.T) {
}
*/

/*
func TestListForEmail(t *testing.T) {
	result, err := chimp.ListsForEmail(map[string]interface{}{"email_address": EMAIL})
	if err != nil {
		t.Error("mailchimp.ListsForEmail", err)
	}
	if len(result[0]) < 5 {
		t.Error("mailchimp.ListsForEmail: Expected the first list id returned to be 10 characters but got", result[0])
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

/*
func TestListAbuseReportsResponse(t *testing.T) {
	response := new(ListAbuseReportsResponse)
	populate("listAbuseReports", response)

	a := response.Total
	if a != 2 {
		t.Error("ListAbuseReportsResponse a: expected 2 but got", a)
	}
	b := response.Data[0].Type
	if b != "AOL" {
		t.Error("ListAbuseReportsResponse b: expected \"AOL\" but got", b)
	}
	c := response.Data[1].Date
	if c.Day() != 3 {
		t.Error("ListAbuseReportsResponse c: expected 3 got", c)
	}
}
func TestListAbuseReports(t *testing.T) {
	result, err := chimp.ListAbuseReports(map[string]interface{}{"id": LIST})
	if err != nil {
		t.Error("mailchimp.ListAbuseReport", err)
	}
	if result.Total != len(result.Data) {
		t.Error("mailchimp.ListAbuseReport: Expected total to equal the number of abuse reports in Data")
	}
}
*/

/*
func TestListActivityElement(t *testing.T) {
	response := make([]ListActivityElement, 0)
	populate("listActivity", &response)

	a := response[0].User_id
	if a != 1234567 {
		t.Error("ListActivity: expected 1234567 but got", a)
	}
	b := response[0].Recipient_clicks
	if b != 40 {
		t.Error("ListActivity: expected 40 but got", b)
	}
	c := response[1].Other_adds
	if c != 0 {
		t.Error("ListActivity: expected 0 but got", c)
	}
}
func TestListActivity(t *testing.T) {
	response, err := chimp.ListActivity(map[string]interface{}{"id": LIST})
	if err != nil {
		t.Error("mailchimp.ListActivity", err)
	}
	if response[0].Day.Year() < 2000 {
		t.Error("mailchimp.ListActivity: the year of the first day's values returned appears incorrect:", response[0].Day.Year())
	}
}
*/

/*
func TestListBatchSubscribeResponse(t *testing.T) {
	_ := new(ListBatchSubscribeResponse)
	//need real json response for mocking
}
func TestListBatchSubscribe(t *testing.T) {
	parameters := make(map[string]interface{})
	parameters["id"] = LIST
	type Email struct {
		EMAIL string
		EMAIL_TYPE string
	}
	parameters["batch"] = []Email{{EMAIL, "html"}}
	parameters["double_optin"] = false
	parameters["update_existing"] = true
	result, err := chimp.ListBatchSubscribe(parameters)
	if err != nil {
		t.Error("mailchimp.ListBatchSubscribe", err)
	}
	if result.Add_count + result.Update_count != 1 {
		t.Error("mailchimp.ListBatchSubscribe: Expected email to be subscribed or updated but it was not")
	}
	if len(result.Errors) != result.Error_count {
		t.Error("mailchimp.ListBatchSubscribe: The error_count should equal the length of errors array")
	}
}
*/

/* NO!
func TestListBatchUnsubscribeResponse(t *testing.T) {
	_ := new(ListBatchUnsubscribeResponse)
	//need real json response for mocking
}
//Don't run this test - can't be resubscribed via API
func TestListBatchUnsubscribe(t *testing.T) {
	parameters := make(map[string]interface{})
	parameters["id"] = LIST
	parameters["emails"] = []string{EMAIL}
	parameters["send_goodbye"] = false
	result, err := chimp.ListBatchUnsubscribe(parameters)
	if err != nil {
		t.Error("mailchimp.ListBatchUnsubscribe", err)
	}
	if result.Success_count + result.Error_count != 1 {
		t.Error("mailchimp.ListBatchUnsubscribe: Expected Success_count or Error_count to be 1")
	}
	if len(result.Errors) != result.Error_count {
		t.Error("mailchimp.ListBatchUnsubscribe: Number of items in array doesn't equal Error_count")
	}
}
*/

/*
func TestListClientsResponse(t *testing.T) {
	response := new(ListClientsResponse)
	populate("listClients", response)

	a := response.Mobile.Clients[0].Members
	if a != 9 {
		t.Error("ListClientsResponse a: expected 9 but got", a)
	}
	b := response.Desktop.Penetration
	if b != 0.83050847457627 {
		t.Error("ListClientsResponse b: expected 0.83050847457627 but got", b)
	}
	c := response.Desktop.Clients[8].Percent
	if c != 0.016949152542373 {
		t.Error("ListClientsResponse c: expected 0.016949152542373 but got", c)
	}
	d := response.Mobile.Clients[1].Client
	if d != "Android" {
		t.Error("ListClientsResponse d: expected \"Android\" but got", d)
	}
}
func TestListClients(t *testing.T) {
	_, err := chimp.ListClients(map[string]interface{}{"id": LIST})
	if err != nil {
		t.Error("ListClients", err)
	}
}
*/

/*
func TestListGrowthHistoryElement(t *testing.T) {
	response := make(ListGrowthHistoryResponse, 0)
	populate("listGrowthHistory", &response)

	verify(t, "ListGrowthHistory", 5, int(response[0].Month.Month()))
	verify(t, "ListGrowthHistory", 2, response[1].Existing)
	verify(t, "ListGrowthHistory", 1, response[2].Imports)
	verify(t, "ListGrowthHistory", 1, response[3].Optins)
}
func TestListGrowthHistory(t *testing.T) {
	_, err := chimp.ListGrowthHistory(map[string]interface{}{"id": LIST})
	if err != nil {
		t.Error("ListGrowthHistory", err)
	}
}
*/

/*
func TestListInterestGroupAdd(t *testing.T) {
	parameters := make(map[string]interface{})
	parameters["id"] = LIST
	parameters["group_name"] = "Test Interest Group"
	result, err := chimp.ListInterestGroupAdd(parameters)
	if err != nil {
		t.Error("ListInterestGroupAdd", err)
	}
	if !result {
		t.Error("ListInterestGroupAdd: expected true but actual value was false")
	}
	interestGroup <- result
}

func TestListInterestGroupUpdate(t *testing.T) {
	<- interestGroup
	parameters := make(map[string]interface{})
	parameters["id"] = LIST
	parameters["old_name"] = "Test Interest Group"
	parameters["new_name"] =  "Updated Interest Group"
	result, err := chimp.ListInterestGroupUpdate(parameters)
	if err != nil {
		t.Error("ListInterestGroupUpdate", err)
	}
	if !result {
		t.Error("ListInterestGroupUpdate: expected true but actual value was false")
	}
	interestGroup <- result
}

func TestListInterestGroupDel(t *testing.T) {
	<- interestGroup
	parameters := make(map[string]interface{})
	parameters["id"] = LIST
	parameters["group_name"] = "Updated Interest Group"
	result, err := chimp.ListInterestGroupDel(parameters)
	if err != nil {
		t.Error("ListInterstGroupDel", err)
	}
	if !result {
		t.Error("ListInterestGroupDel: expected true but actual value was false")
	}
}
*/

/*
func TestListInterestGrouingAdd (t *testing.T) {
  parameters := make(map[string]interface{})
  parameters["id"] = LIST
  parameters["name"] = "Diet"
  parameters["type"] = "radio"
  parameters["groups"] = []string{"vegetarian", "carnivore"}
  result, err := chimp.ListInterestGroupingAdd(parameters)
  if err != nil {
    t.Error("ListInterestGroupDel", err)
  }
	if result <= 0 {
		t.Error("Expected the folder_id to be a positive integer but got", result)
	}
  interestGrouping <- result
}

func TestListInterestGroupingUpdate (t *testing.T) {
  id := <- interestGrouping
  parameters := make(map[string]interface{})
  parameters["grouping_id"] = id
  parameters["name"] = "name"
  parameters["value"] = "Dietary Preferences"
  result, err := chimp.ListInterestGroupingUpdate(parameters);
  if err != nil {
    t.Error("ListInterestGroupingUpdate", err)
  }
  if !result {
    t.Error("ListInterestGroupingUpdate: expected true but got", result)
  }
  interestGrouping <- id
}

func TestListInterestGroupings (t *testing.T) {
  id := <- interestGrouping
  result, err := chimp.ListInterestGroupings(map[string]interface{}{"id": LIST})
  if err != nil {
    t.Error("ListInterestGroupings", err)
  }
  if result[0].Id != id {
    t.Error("ListInterestGroupings: Expected", id, "but got", result[0].Id)
  }
  if result[0].Groups[0].Name != "vegetarian" {
    t.Error("ListInterestGroupings: Expected vegetarian but got", result[0].Groups[0].Name)
  }
  if result[0].Groups[1].Subscribers != 0 {
    t.Error("ListInterestGroupings: Expected 0 but got", result[0].Groups[1].Subscribers)
  }
  interestGrouping <- id
}

func TestListInterestGroupingDel (t *testing.T) {
  result, err := chimp.ListInterestGroupingDel(map[string]interface{}{"grouping_id": <- interestGrouping})
  if err != nil {
    t.Error("ListInterestGroupingDel", err)
  }
  if !result {
    t.Error("ListInterestGroupingDel", err)
  }
}
*/
