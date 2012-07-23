package mailchimp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"
)

//format string for time.Format
const ChimpTime = "2006-01-02 15:04:05"

var datacenter = regexp.MustCompile("[a-z]+[0-9]+$")

type API struct {
	Key      string
	endpoint string
}

func New(apikey string, https bool) (*API, error) {
	u := url.URL{}
	if https {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}
	u.Host = fmt.Sprintf("%s.api.mailchimp.com", datacenter.FindString(apikey))
	u.Path = "/1.3/"
	return &API{apikey, u.String() + "?method="}, nil
}

func run(a *API, method string, parameters map[string]interface{}) ([]byte, error) {
	if parameters == nil {
		parameters = make(map[string]interface{})
	}
	parameters["apikey"] = a.Key
	b, err := json.Marshal(parameters)
	if err != nil {
		return nil, err
	}
os.Stdout.Write([]byte(b))
	resp, err := http.Post(a.endpoint+method, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
os.Stdout.Write(body)
	if err = errorCheck(body); err != nil {
		return nil, err
	}
	return body, nil
}

type ChimpError struct {
	Err  string `json:"error"`
	Code int
}

func (e ChimpError) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Err)
}
func errorCheck(body []byte) error {
	var e ChimpError
	json.Unmarshal(body, &e)
	if e.Err != "" || e.Code != 0 {
		return e
	}
	return nil
}

func chimpTime(t interface{}) interface{} {
	switch ti := t.(type) {
	case time.Time:
		return ti.Format(ChimpTime)
	case string:
		return ti
	}
	return t
}

func parseInt(body []byte, err error) (int, error) {
	i, err := strconv.ParseInt(string(body), 10, 0)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func parseString(body []byte, err error) (string, error) {
	if err != nil {
		return "", err
	}
	return strconv.Unquote(string(body))
}

func parseBoolean(body []byte, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(string(body))
}

func parseStruct(a *API, method string, parameters map[string]interface{}, retVal interface{}) error {
	body, err := run(a, method, parameters)
	if err != nil {
		return err
	}
	json.Unmarshal(body, retVal)
	return nil
}

func (a *API) Ping() (string, error) {
	return parseString(run(a, "ping", nil))
}

type CampaignContentResult struct {
	Html string
	Text string
}

func (a *API) CampaignContent(parameters map[string]interface{}) (retVal *CampaignContentResult, err error) {
	retVal = new(CampaignContentResult)
	err = parseStruct(a, "campaignContent", parameters, retVal)
	return
}

func (a *API) CampaignCreate(parameters map[string]interface{}) (string, error) {
	return parseString(run(a, "campaignCreate", parameters))
}

func (a *API) CampaignDelete(parameters map[string]interface{}) (bool, error) {
	return parseBoolean(run(a, "campaignDelete", parameters))
}

//CampaignEcommOrderAdd method has not been tested with real return data
func (a *API) CampaignEcommOrderAdd(parameters map[string]interface{}) (bool, error) {
	return parseBoolean(run(a, "campaignEcommOrderAdd", parameters))
}

func (a *API) CampaignPause(parameters map[string]interface{}) (bool, error) {
	return parseBoolean(run(a, "campaignPause", parameters))
}

func (a *API) CampaignReplicate(parameters map[string]interface{}) (string, error) {
	return parseString(run(a, "campaignReplicate", parameters))
}

func (a *API) CampaignResume(parameters map[string]interface{}) (bool, error) {
	return parseBoolean(run(a, "campaignResume", parameters))
}

func (a *API) CampaignSchedule(parameters map[string]interface{}) (bool, error) {
	//convert times to Mailchimp's format
	if parameters == nil {
		return false, errors.New("missing required parameters")
	}
	parameters["schedule_time"] = chimpTime(parameters["schedule_time"])
	parameters["schedule_time_b"] = chimpTime(parameters["schedule_time_b"])
	return parseBoolean(run(a, "campaignSchedule", parameters))
}

func (a *API) CampaignSegmentTest(parameters map[string]interface{}) (int, error) {
	return parseInt(run(a, "campaignSegmentTest", parameters))
}

func (a *API) CampaignSendNow(parameters map[string]interface{}) (bool, error) {
	return parseBoolean(run(a, "campaignSendNow", parameters))
}

func (a *API) CampaignSendTest(parameters map[string]interface{}) (bool, error) {
	return parseBoolean(run(a, "campaignSendTest", parameters))
}

type CampaignShareReportResult struct {
	Title      string
	Url        string
	Secure_url string
	Password   string
}

func (a *API) CampaignShareReport(parameters map[string]interface{}) (retVal *CampaignShareReportResult, err error) {
	retVal = new(CampaignShareReportResult)
	err = parseStruct(a, "campaignShareReport", parameters, retVal)
	return
}

//CampaignTemplateContent method returns a map[string]interface{} of all content sections for the campaign
//Section names are dependent upon the template used and thus can't be documented
//TODO: If all values in the resulting map are string, change return type to map[string]string to obviate type assertions
func (a *API) CampaignTemplateContent(parameters map[string]interface{}) (retVal map[string]interface{}, err error) {
	err = parseStruct(a, "campaignTemplateContent", parameters, &retVal)
	return
}

func (a *API) CampaignUnschedule(parameters map[string]interface{}) (bool, error) {
	return parseBoolean(run(a, "campaignUnschedule", parameters))
}

func (a *API) CampaignUpdate(parameters map[string]interface{}) (bool, error) {
	return parseBoolean(run(a, "campaignUpdate", parameters))
}

type CampaignsResult struct {
	Total int
	Data  []CampaignsResultData
}
type CampaignsResultData struct {
	Id                string
	Web_id            int
	List_id           string
	Folder_id         int
	Template_id       int
	Content_type      string
	Title             string
	Type              string
	Create_time       string
	Send_time         string
	Emails_sent       int
	Status            string
	From_name         string
	From_email        string
	Subject           string
	To_name           string
	Archive_url       string
	Inline_css        bool
	Analytics         string
	Analytics_tag     string
	Authenticate      bool
	Ecomm360          bool
	Auto_tweet        bool
	Auto_fb_post      string
	Auto_footer       bool
	Timewarp          bool
	Timewarp_schedule string
	Tracking          CampaignsResultDataTracking
	Segment_text      string
	Segment_opts      CampaignsResultDataSegment_opts
	Type_opts         map[string]interface{}
}
type CampaignsResultDataTracking struct {
	Html_clicks bool
	Text_clicks bool
	Opens       bool
}
type CampaignsResultDataSegment_opts struct {
	Match      string
	Conditions []map[string]interface{}
}

func (a *API) Campaigns(parameters map[string]interface{}) (retVal *CampaignsResult, err error) {
	retVal = new(CampaignsResult)
	err = parseStruct(a, "campaigns", parameters, retVal)
	return
}

type CampaignAbuseReportsResultDataItem struct {
	Date  string
	Email string
	Type  string
}
type CampaignAbuseReportsResult struct {
	Total int
	Data  []CampaignAbuseReportsResultDataItem
}

func (a *API) CampaignAbuseReports(parameters map[string]interface{}) (retVal *CampaignAbuseReportsResult, err error) {
	retVal = new(CampaignAbuseReportsResult)
	err = parseStruct(a, "campaignAbuseReports", parameters, retVal)
	return
}

type CampaignAdviceResultItem struct {
	Msg  string
	Type string
}

func (a *API) CampaignAdvice(parameters map[string]interface{}) (retVal []CampaignAdviceResultItem, err error) {
	err = parseStruct(a, "campaignAdvice", parameters, &retVal)
	return
}

type CampaignAnalyticsResultGoals struct {
	Name        string
	Conversions int
}
type CampaignAnalyticsResult struct {
	Visits            int
	Pages             int
	New_visits        int
	Bounces           int
	Time_on_site      float64
	Goal_conversions  int
	Goal_value        float64
	Revenue           float64
	Transactions      int
	Ecomm_conversions int
	Goals             CampaignAnalyticsResultGoals
}

func (a *API) CampaignAnalytics(parameters map[string]interface{}) (retVal *CampaignAnalyticsResult, err error) {
	retVal = new(CampaignAnalyticsResult)
	err = parseStruct(a, "campaignAnalytics", parameters, retVal)
	return
}

type CampaignBounceMessageResult struct {
	Date    string
	Email   string
	Message string
}

func (a *API) CampaignBounceMessage(parameters map[string]interface{}) (retVal *CampaignBounceMessageResult, err error) {
	retVal = new(CampaignBounceMessageResult)
	err = parseStruct(a, "campaignBounceMessage", parameters, retVal)
	return
}

type CampaignBounceMessagesResult struct {
	Total int
	Data  []CampaignBounceMessageResult
}

func (a *API) CampaignBounceMessages(parameters map[string]interface{}) (retVal *CampaignBounceMessagesResult, err error) {
	retVal = new(CampaignBounceMessagesResult)
	err = parseStruct(a, "campaignBounceMessages", parameters, retVal)
	return
}

//CampaignClickStats method returns a map where the keys are urls extracted from the campaign
type CampaignClickStatsResultItem struct {
	Clicks int
	Unique int
}

func (a *API) CampaignClickStats(parameters map[string]interface{}) (retVal map[string]CampaignClickStatsResultItem, err error) {
	err = parseStruct(a, "campaignClickStats", parameters, &retVal)
	return
}

//CampaignEcommOrders method has not been tested with real response data
//The json returned by this routine might not unmarshal correctly into the return struct for this method
type CampaignEcommOrdersResultDataItemLinesItem struct {
	Line_num              int
	Product_id            int
	Product_name          string
	Product_sku           string
	Product_category_id   int
	Product_category_name int
	Qty                   int
	Cost                  float64
}
type CampaignEcommOrdersResultDataItem struct {
	Store_id    string
	Store_name  string
	Order_id    string
	Email       string
	Order_total float64
	Tax_total   float64
	Ship_total  float64
	Order_date  string
	Lines       []CampaignEcommOrdersResultDataItemLinesItem
}
type CampaignEcommOrdersResult struct {
	Total int
	Data  []CampaignEcommOrdersResultDataItem
}

func (a *API) CampaignEcommOrders(parameters map[string]interface{}) (retVal *CampaignEcommOrdersResult, err error) {
	retVal = new(CampaignEcommOrdersResult)
	err = parseStruct(a, "campaignEcommOrders", parameters, retVal)
	return
}

//Mailchimp's documentation for CampaignEepUrlStats is incorrect and I don't have any examples of non-empty JSON
//return values from this routine, so there is no way to design accurate types to unmarshal responses into.
//Therefore, this method returns an interface{} and users will have to use type assertion or type
//switching to unpack the response. The custom types for this method are not used yet.
type CampaignEepUrlStatsResultTwitterClicksLocations struct {
	Country string
	Region  string
	Total   int
}
type CampaignEepUrlStatsResultTwitterClicksReferrers struct {
	Referrer    string
	Clicks      int
	First_click string
	Last_click  string
}
type CampaignEepUrlStatsResultTwitterClicks struct {
	Clicks      int
	First_click string
	Last_click  string
	Locations   CampaignEepUrlStatsResultTwitterClicksLocations
	Referrers   []CampaignEepUrlStatsResultTwitterClicksReferrers
}
type CampaignEepUrlStatsResultTwitterStatuses struct {
	Status      string
	Screen_name string
	Status_id   string
	Datetime    string
	Is_retweet  bool
}
type CampaignEepUrlStatsResultTwitter struct {
	Tweets        int
	First_tweet   string
	Last_tweet    string
	Retweets      int
	First_retweet string
	Last_retweet  string
	Statuses      CampaignEepUrlStatsResultTwitterStatuses
	Clicks        CampaignEepUrlStatsResultTwitterClicks
}
type CampaignEepUrlStatsResult struct {
	Twitter CampaignEepUrlStatsResultTwitter
}

func (a *API) CampaignEepUrlStats(parameters map[string]interface{}) (retVal interface{}, err error) {
	err = parseStruct(a, "campaignEepUrlStats", parameters, &retVal)
	return
}

type CampaignEmailDomainPerformanceResultItem struct {
	Domain     string
	Total_sent int
	Email      int
	Bounces    int
	Opens      int
	Clicks     int
	Unsubs     int
	Delivered  int
	Emails_pct int
	Opens_pct  int
	Clicks_pct int
	Unsubs_pct int
}

func (a *API) CampaignEmailDomainPerformance(parameters map[string]interface{}) (retVal []CampaignEmailDomainPerformanceResultItem, err error) {
	err = parseStruct(a, "campaignEmailDomainPerformance", parameters, &retVal)
	return
}

type CampaignGeoOpensResultItem struct {
	Code          string
	Name          string
	Opens         int
	Region_detail bool
}

func (a *API) CampaignGeoOpens(parameters map[string]interface{}) (retVal []CampaignGeoOpensResultItem, err error) {
	err = parseStruct(a, "campaignGeoOpens", parameters, &retVal)
	return
}

type CampaignGeoOpensForCountryReturnItem struct {
	Code  string
	Name  string
	Opens int
}

func (a *API) CampaignGeoOpensForCountry(parameters map[string]interface{}) (retVal []CampaignGeoOpensForCountryReturnItem, err error) {
	err = parseStruct(a, "campaignGeoOpensForCountry", parameters, &retVal)
	return
}

type CampaignMembersResult struct {
	Total int
	Data []struct{
		Email string
		Status string
		Absplit_group string
		Tz_group string
	}
}
func (a *API) CampaignMembers(parameters map[string]interface{}) (retVal *CampaignMembersResult, err error) {
	retVal = new(CampaignMembersResult)
	err = parseStruct(a, "campaignMembers", parameters, retVal)
	return
}

//CampaignStatsResult method has only been tested with limited return data
//The nested structs in the return struct in particular may be incorrect
type CampaignStatsResult struct {
	Syntax_errors int
	Hard_bounces int
	Soft_bounces int
	Unsubscribes int
	Abuse_reports int
	Forwards int
	Forwards_opens int
	Opens int
	Last_open string
	Unique_opens int
	Clicks int
	Unique_clicks int
	Last_click string
	Users_who_clicked int
	Emails_sent int
	Unique_likes int
	Recipient_likes int
	Facebook_likes int
	Absplit struct {
		Bounces_a int
		Bounces_b int
		Forwards_a int
		Forwards_b int
		Abuse_reports_a int
		Abuse_reports_b int
		Unsubs_a int
		Unsubs_b int
		Recipients_click_a int
		Recipients_click_b int
		Forwards_opens_a int
		Forwards_opens_b int
	}
	Timewarp map[string]struct{
		Opens int
		Last_open string
		Unique_opens int
		Clicks int
		Last_click string
		Bounces int
		Total int
		Sent int
	}
	Timeseries[]struct{
		Timestamp string
		Emails_sent int
		Unique_opens int
		Recipients_click int
	}
}
func (a *API) CampaignStats(parameters map[string]interface{}) (retVal *CampaignStatsResult, err error) {
	retVal = new(CampaignStatsResult)
	err = parseStruct(a, "campaignStats", parameters, retVal)
	return
}

type CampaignUnsubscribesResult struct {
	Total int
	Data []struct{
		Email string
		Reason string
		Reason_text string
	}
}
func (a *API) CampaignUnsubscribes(parameters map[string]interface{}) (retVal *CampaignUnsubscribesResult, err error) {
	retVal = new(CampaignUnsubscribesResult)
	err = parseStruct(a, "campaignUnsubscribes", parameters, retVal)
	return
}

type CampaignClickDetailAIMResult struct {
	Total int
	Data []struct{
		Email string
		Clicks int
	}
}
func (a *API) CampaignClickDetailAIM(parameters map[string]interface{}) (retVal *CampaignClickDetailAIMResult, err error) {
	retVal = new(CampaignClickDetailAIMResult)
	err = parseStruct(a, "campaignClickDetailAIM", parameters, retVal)
	return
}

type CampaignEmailStatsAIMResult struct {
	Success int
	Error int
	Data []struct{
		Action string
		Timestamp string
		Url string
	}
}
func (a *API) CampaignEmailStatsAIM(parameters map[string]interface{}) (retVal *CampaignEmailStatsAIMResult, err error) {
	retVal = new(CampaignEmailStatsAIMResult)
	err = parseStruct(a, "campaignEmailStatsAIM", parameters, retVal)
	return
}

type CampaignEmailStatsAIMAllResult struct {
	Total int
	Data map[string][]struct{
		Action string
		Timestamp string
		Url string
	}
}
func (a *API) CampaignEmailStatsAIMAll(parameters map[string]interface{}) (retVal *CampaignEmailStatsAIMAllResult, err error) {
	retVal = new(CampaignEmailStatsAIMAllResult)
	err = parseStruct(a, "campaignEmailStatsAIMAll", parameters, retVal)
	return
}

type CampaignNotOpenedAIMResult struct {
	Total int
	Data []string
}
func (a *API) CampaignNotOpenedAIM(parameters map[string]interface{}) (retVal *CampaignNotOpenedAIMResult, err error) {
	retVal = new(CampaignNotOpenedAIMResult)
	err = parseStruct(a, "campaignNotOpenedAIM", parameters, retVal)
	return
}

type CampaignOpenedAIMResult struct {
	Total int
	Data []struct{
		Email string
		Open_count int
	}
}
func (a *API) CampaignOpenedAIM(parameters map[string]interface{}) (retVal *CampaignOpenedAIMResult, err error) {
	retVal = new(CampaignOpenedAIMResult)
	err = parseStruct(a, "campaignOpenedAIM", parameters, retVal)
	return
}

func (a *API) EcommOrderAdd(parameters map[string]interface{}) (bool, error) {
	return parseBoolean(run(a, "ecommOrderAdd", parameters))
}

func (a *API) EcommOrderDel(parameters map[string]interface{}) (bool, error) {
	return parseBoolean(run(a, "ecommOrderDelete", parameters))
}

//EcommOrdersResult tested with data; result unmarshals correctly into this struct
type EcommOrdersResult struct {
	Total int
	Data []struct{
		Store_id string
		Store_name string
		Order_id string
		Email string
		Order_total float64
		Tax_total float64
		Ship_total float64
		Order_date string
		Lines []struct{
			Line_num int
			Product_id int
			Product_name string
			Product_sku string
			Product_category_id int
			Product_category_name string
			Qty int
			Cost float64
		}
	}
}
func (a *API) EcommOrders(parameters map[string]interface{}) (retVal *EcommOrdersResult, err error) {
	retVal = new(EcommOrdersResult)
	err = parseStruct(a, "ecommOrders", parameters, retVal)
	return
}


