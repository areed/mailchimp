/*
	Package mailchimp implements version 1.3 of Mailchimp's API

	Routines are implemented as methods on type API, which should
	be created with the New function

	chimp := mailchimp.New("apikey123apikey123-us1")

	The comment for each method contains a link to the corresonding
	Mailchimp routine documentation in lieu of a description of parameters
	used with each method

	TODO: implement timeouts and corresponding error
*/
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

type API struct {
	Key      string
	endpoint string
}

var datacenter = regexp.MustCompile("[a-z]+[0-9]+$")

func New(apikey string, https ...bool) (*API, error) {
	u := url.URL{}
	if https[0] {
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
	//os.Stdout.Write([]byte(b))
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

//ChimpTime is a named struct with a single anonymous field of type time.Time
//and inherits all its methods except UnmarshalJSON, which is overridden since
//Mailchimp does not adhere to the RFC3339 format
type ChimpTime struct {
	time.Time
}

func (t *ChimpTime) UnmarshalJSON(data []byte) (err error) {
	s := string(data)
	l := len(s)
	switch {
	case l == 12:
		t.Time, err = time.Parse(`"2006-01-02"`, s)
	case l == 21:
		t.Time, err = time.Parse(`"2006-01-02 15:04:05"`, s)
	case l == 9:
		t.Time, err = time.Parse(`"2006-01"`, s)
	}
	return
}

//format string for time.Format
const ChimpTimeFormat = "2006-01-02 15:04:05"

func chimpTime(t interface{}) interface{} {
	switch ti := t.(type) {
	case time.Time:
		return ti.Format(ChimpTimeFormat)
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

type alterJsoner interface {
	alterJson(b []byte) []byte
}

func parseJson(a *API, method string, parameters map[string]interface{}, retVal interface{}) error {
	body, err := run(a, method, parameters)
	if err != nil {
		return err
	}
	switch r := retVal.(type) {
	case alterJsoner:
		blob := r.alterJson(body)
		os.Stdout.Write(blob)
		json.Unmarshal(r.alterJson(body), retVal)
	default:
		json.Unmarshal(body, retVal)
	}
	return nil
}

type CampaignContentResult struct {
	Html string
	Text string
}

func (a *API) CampaignContent(parameters map[string]interface{}) (retVal *CampaignContentResult, err error) {
	retVal = new(CampaignContentResult)
	err = parseJson(a, "campaignContent", parameters, retVal)
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
	err = parseJson(a, "campaignShareReport", parameters, retVal)
	return
}

//CampaignTemplateContent method returns a map[string]interface{} of all content sections for the campaign
//Section names are dependent upon the template used and thus can't be documented
//TODO: If all values in the resulting map are string, change return type to map[string]string to obviate type assertions
func (a *API) CampaignTemplateContent(parameters map[string]interface{}) (retVal map[string]interface{}, err error) {
	err = parseJson(a, "campaignTemplateContent", parameters, &retVal)
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
	err = parseJson(a, "campaigns", parameters, retVal)
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
	err = parseJson(a, "campaignAbuseReports", parameters, retVal)
	return
}

type CampaignAdviceResultItem struct {
	Msg  string
	Type string
}

func (a *API) CampaignAdvice(parameters map[string]interface{}) (retVal []CampaignAdviceResultItem, err error) {
	err = parseJson(a, "campaignAdvice", parameters, &retVal)
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
	err = parseJson(a, "campaignAnalytics", parameters, retVal)
	return
}

type CampaignBounceMessageResult struct {
	Date    string
	Email   string
	Message string
}

func (a *API) CampaignBounceMessage(parameters map[string]interface{}) (retVal *CampaignBounceMessageResult, err error) {
	retVal = new(CampaignBounceMessageResult)
	err = parseJson(a, "campaignBounceMessage", parameters, retVal)
	return
}

type CampaignBounceMessagesResult struct {
	Total int
	Data  []CampaignBounceMessageResult
}

func (a *API) CampaignBounceMessages(parameters map[string]interface{}) (retVal *CampaignBounceMessagesResult, err error) {
	retVal = new(CampaignBounceMessagesResult)
	err = parseJson(a, "campaignBounceMessages", parameters, retVal)
	return
}

//CampaignClickStats method returns a map where the keys are urls extracted from the campaign
type CampaignClickStatsResultItem struct {
	Clicks int
	Unique int
}

func (a *API) CampaignClickStats(parameters map[string]interface{}) (retVal map[string]CampaignClickStatsResultItem, err error) {
	err = parseJson(a, "campaignClickStats", parameters, &retVal)
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
	err = parseJson(a, "campaignEcommOrders", parameters, retVal)
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
	err = parseJson(a, "campaignEepUrlStats", parameters, &retVal)
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
	err = parseJson(a, "campaignEmailDomainPerformance", parameters, &retVal)
	return
}

type CampaignGeoOpensResultItem struct {
	Code          string
	Name          string
	Opens         int
	Region_detail bool
}

func (a *API) CampaignGeoOpens(parameters map[string]interface{}) (retVal []CampaignGeoOpensResultItem, err error) {
	err = parseJson(a, "campaignGeoOpens", parameters, &retVal)
	return
}

type CampaignGeoOpensForCountryReturnItem struct {
	Code  string
	Name  string
	Opens int
}

func (a *API) CampaignGeoOpensForCountry(parameters map[string]interface{}) (retVal []CampaignGeoOpensForCountryReturnItem, err error) {
	err = parseJson(a, "campaignGeoOpensForCountry", parameters, &retVal)
	return
}

type CampaignMembersResult struct {
	Total int
	Data  []struct {
		Email         string
		Status        string
		Absplit_group string
		Tz_group      string
	}
}

func (a *API) CampaignMembers(parameters map[string]interface{}) (retVal *CampaignMembersResult, err error) {
	retVal = new(CampaignMembersResult)
	err = parseJson(a, "campaignMembers", parameters, retVal)
	return
}

//CampaignStatsResult method has only been tested with limited return data
//The nested structs in the return struct in particular may be incorrect
type CampaignStatsResult struct {
	Syntax_errors     int
	Hard_bounces      int
	Soft_bounces      int
	Unsubscribes      int
	Abuse_reports     int
	Forwards          int
	Forwards_opens    int
	Opens             int
	Last_open         string
	Unique_opens      int
	Clicks            int
	Unique_clicks     int
	Last_click        string
	Users_who_clicked int
	Emails_sent       int
	Unique_likes      int
	Recipient_likes   int
	Facebook_likes    int
	Absplit           struct {
		Bounces_a          int
		Bounces_b          int
		Forwards_a         int
		Forwards_b         int
		Abuse_reports_a    int
		Abuse_reports_b    int
		Unsubs_a           int
		Unsubs_b           int
		Recipients_click_a int
		Recipients_click_b int
		Forwards_opens_a   int
		Forwards_opens_b   int
	}
	Timewarp map[string]struct {
		Opens        int
		Last_open    string
		Unique_opens int
		Clicks       int
		Last_click   string
		Bounces      int
		Total        int
		Sent         int
	}
	Timeseries []struct {
		Timestamp        string
		Emails_sent      int
		Unique_opens     int
		Recipients_click int
	}
}

func (a *API) CampaignStats(parameters map[string]interface{}) (retVal *CampaignStatsResult, err error) {
	retVal = new(CampaignStatsResult)
	err = parseJson(a, "campaignStats", parameters, retVal)
	return
}

type CampaignUnsubscribesResult struct {
	Total int
	Data  []struct {
		Email       string
		Reason      string
		Reason_text string
	}
}

func (a *API) CampaignUnsubscribes(parameters map[string]interface{}) (retVal *CampaignUnsubscribesResult, err error) {
	retVal = new(CampaignUnsubscribesResult)
	err = parseJson(a, "campaignUnsubscribes", parameters, retVal)
	return
}

type CampaignClickDetailAIMResult struct {
	Total int
	Data  []struct {
		Email  string
		Clicks int
	}
}

func (a *API) CampaignClickDetailAIM(parameters map[string]interface{}) (retVal *CampaignClickDetailAIMResult, err error) {
	retVal = new(CampaignClickDetailAIMResult)
	err = parseJson(a, "campaignClickDetailAIM", parameters, retVal)
	return
}

type CampaignEmailStatsAIMResult struct {
	Success int
	Error   int
	Data    []struct {
		Action    string
		Timestamp string
		Url       string
	}
}

func (a *API) CampaignEmailStatsAIM(parameters map[string]interface{}) (retVal *CampaignEmailStatsAIMResult, err error) {
	retVal = new(CampaignEmailStatsAIMResult)
	err = parseJson(a, "campaignEmailStatsAIM", parameters, retVal)
	return
}

type CampaignEmailStatsAIMAllResult struct {
	Total int
	Data  map[string][]struct {
		Action    string
		Timestamp string
		Url       string
	}
}

func (a *API) CampaignEmailStatsAIMAll(parameters map[string]interface{}) (retVal *CampaignEmailStatsAIMAllResult, err error) {
	retVal = new(CampaignEmailStatsAIMAllResult)
	err = parseJson(a, "campaignEmailStatsAIMAll", parameters, retVal)
	return
}

type CampaignNotOpenedAIMResult struct {
	Total int
	Data  []string
}

func (a *API) CampaignNotOpenedAIM(parameters map[string]interface{}) (retVal *CampaignNotOpenedAIMResult, err error) {
	retVal = new(CampaignNotOpenedAIMResult)
	err = parseJson(a, "campaignNotOpenedAIM", parameters, retVal)
	return
}

type CampaignOpenedAIMResult struct {
	Total int
	Data  []struct {
		Email      string
		Open_count int
	}
}

func (a *API) CampaignOpenedAIM(parameters map[string]interface{}) (retVal *CampaignOpenedAIMResult, err error) {
	retVal = new(CampaignOpenedAIMResult)
	err = parseJson(a, "campaignOpenedAIM", parameters, retVal)
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
	Data  []struct {
		Store_id    string
		Store_name  string
		Order_id    string
		Email       string
		Order_total float64
		Tax_total   float64
		Ship_total  float64
		Order_date  string
		Lines       []struct {
			Line_num              int
			Product_id            int
			Product_name          string
			Product_sku           string
			Product_category_id   int
			Product_category_name string
			Qty                   int
			Cost                  float64
		}
	}
}

func (a *API) EcommOrders(parameters map[string]interface{}) (retVal *EcommOrdersResult, err error) {
	retVal = new(EcommOrdersResult)
	err = parseJson(a, "ecommOrders", parameters, retVal)
	return
}

func (a *API) FolderAdd(parameters map[string]interface{}) (int, error) {
	return parseInt(run(a, "folderAdd", parameters))
}

func (a *API) FolderDel(parameters map[string]interface{}) (bool, error) {
	return parseBoolean(run(a, "folderDel", parameters))
}

func (a *API) FolderUpdate(parameters map[string]interface{}) (bool, error) {
	return parseBoolean(run(a, "folderUpdate", parameters))
}

type FoldersResultItem struct {
	Folder_id    int
	Name         string
	Date_created string
	Type         string
}

func (a *API) Folders(parameters map[string]interface{}) (retVal []FoldersResultItem, err error) {
	err = parseJson(a, "folders", parameters, &retVal)
	return
}

type GmonkeyActivityResultItem struct {
	Action        string
	Timestamp     string
	Url           string
	Unique_id     string
	Title         string
	List_name     string
	Email         string
	Fname         string
	Lname         string
	Member_rating int
	Member_since  string
	Geo           struct {
		Latitude  string
		Longitude string
		Gmtoff    string
		Dstoff    string
		Timezone  string
		Cc        string
		Region    string
	}
}

func (a *API) GmonkeyActivity(parameters map[string]interface{}) (retVal []GmonkeyActivityResultItem, err error) {
	err = parseJson(a, "gmonkeyActivity", parameters, &retVal)
	return
}

type GmonkeyAddResult struct {
	Success int
	Errors  int
	Data    []struct {
		Email_address string
		Error         string
	}
}

func (a *API) GmonkeyAdd(parameters map[string]interface{}) (retVal *GmonkeyAddResult, err error) {
	retVal = new(GmonkeyAddResult)
	err = parseJson(a, "gmonkeyAdd", parameters, retVal)
	return
}

type GmonkeyDelResult struct {
	Success int
	Errors  int
	Data    []struct {
		Email_address string
		Error         string
	}
}

func (a *API) GmonkeyDel(parameters map[string]interface{}) (retVal *GmonkeyDelResult, err error) {
	retVal = new(GmonkeyDelResult)
	err = parseJson(a, "gmonkeyDel", parameters, retVal)
	return
}

type GmonkeyMembersItem struct {
	List_id       string
	List_name     string
	Email         string
	Fname         string
	Lname         string
	Member_rating int
	Member_since  int
}

func (a *API) GmonkeyMembers(parameters map[string]interface{}) (retVal []GmonkeyMembersItem, err error) {
	err = parseJson(a, "gmonkeyMembers", parameters, &retVal)
	return
}

func (a *API) CampaignsForEmail(parameters map[string]interface{}) (retVal []string, err error) {
	err = parseJson(a, "campaignsForEmail", parameters, &retVal)
	return
}

type ChimpChatterResultItem struct {
	Message     string
	Type        string
	Url         string
	List_id     string
	Campaign_id string
	Update_time string
}

func (a *API) ChimpChatter(parameters map[string]interface{}) (retVal []ChimpChatterResultItem, err error) {
	err = parseJson(a, "chimpChatter", parameters, &retVal)
	return
}

func (a *API) GenerateText(parameters map[string]interface{}) (string, error) {
	return parseString(run(a, "generateText", parameters))
}

type GetAccountDetailsResult struct {
	Username        string
	User_id         string
	Is_trial        bool
	Is_approved     bool
	Has_activated   bool
	Timezone        string
	Plan_type       string
	Plan_low        int
	Plan_high       int
	Plan_start_date string
	Emails_left     int
	Pending_monthly bool
	First_payment   string
	Last_payment    string
	Times_logged_in int
	Last_login      string
	Affiliate_link  string
	Contact         struct {
		Fname    string
		Lname    string
		Email    string
		Company  string
		Address1 string
		Address2 string
		City     string
		State    string
		Zip      string
		Country  string
		Url      string
		Phone    string
		Fax      string
	}
	Modules []struct {
		Name  string
		Added string
	}
	Orders []struct {
		Order_id     int
		Type         string
		Amount       float64
		Date         string
		Credits_used float64
	}
	Rewards struct {
		Referrals_this_month int
		Notify_on            string
		Notify_email         string
		Credits              struct {
			This_month   int
			Total_earned int
			Remaining    int
		}
		Inspections struct {
			This_month   int
			Total_earned int
			Remaining    int
		}
		Referrals []struct {
			Name        string
			Email       string
			Signup_date string
			Type        string
		}
		Applied []struct {
			Value      int
			Date       string
			Order_id   int
			Order_desc string
		}
	}
}

func (a *API) GetAccountDetails(parameters map[string]interface{}) (retVal *GetAccountDetailsResult, err error) {
	retVal = new(GetAccountDetailsResult)
	err = parseJson(a, "getAccountDetails", parameters, retVal)
	return
}

type GetVerifiedDomainsResultItem struct {
	Domain string
	Status string
	Emails string
}

func (a *API) GetVerifiedDomains(parameters map[string]interface{}) (retVal []GetVerifiedDomainsResultItem, err error) {
	err = parseJson(a, "getVerifiedDomains", parameters, &retVal)
	return
}

func (a *API) InlineCss(parameters map[string]interface{}) (string, error) {
	return parseString(run(a, "inlineCss", parameters))
}

func (a *API) ListsForEmail(parameters map[string]interface{}) (retVal []string, err error) {
	err = parseJson(a, "listsForEmail", parameters, &retVal)
	return
}

func (a *API) Ping() (string, error) {
	return parseString(run(a, "ping", nil))
}

//ListAbuseReportsResponse is the type for values returned from the ListAbuseReports method
type ListAbuseReportsResponse struct {
	Total int
	Data  []struct {
		Date        ChimpTime
		Email       string
		Campaign_id string
		Type        string
	}
}

//ListAbuseReports gets all email addresses that complained about a given campaign
func (a *API) ListAbuseReports(parameters map[string]interface{}) (retVal *ListAbuseReportsResponse, err error) {
	retVal = new(ListAbuseReportsResponse)
	err = parseJson(a, "listAbuseReports", parameters, retVal)
	return
}

//ListActivityElement is the type of elements in the slice returned from the ListActivity method
type ListActivityElement struct {
	User_id          int //not documented in the api docs; not sure what it means
	Day              ChimpTime
	Emails_sent      int
	Unique_opens     int
	Recipient_clicks int
	Hard_bounce      int
	Soft_bounce      int
	Abuse_reports    int
	Subs             int
	Unsubs           int
	Other_adds       int
	Other_removes    int
}

//ListActivity accesses up to the previous 180 days of daily detailed aggregated
//activity stats for a given list
func (a *API) ListActivity(parameters map[string]interface{}) (retVal []ListActivityElement, err error) {
	err = parseJson(a, "listActivity", parameters, &retVal)
	return
}

//ListBatchSubscribeResponse is the type for values returned from the ListBatchSubscribe method
type ListBatchSubscribeResponse struct {
	Add_count    int
	Update_count int
	Error_count  int
	Errors       []struct {
		Email   string
		Code    int
		Message string
	}
}

//ListBatchSubscribe subscribes a batch of email address to a list at once.
//You should limit batches to 5k - 10k records.
//http://apidocs.mailchimp.com/api/1.3/listbatchsubscribe.func.php
func (a *API) ListBatchSubscribe(parameters map[string]interface{}) (retVal *ListBatchSubscribeResponse, err error) {
	retVal = new(ListBatchSubscribeResponse)
	err = parseJson(a, "listBatchSubscribe", parameters, retVal)
	return
}

//ListBatchUnsubscribeResponse is the type for values returned from the ListBatchUnsubscribe method
type ListBatchUnsubsribeResponse struct {
	Success_count int
	Error_count   int
	Errors        []struct {
		Email   string
		Code    int
		Message string
	}
}
//ListBatchUnsubscribe unsubscribes a batch of email addresses from a list
//http://apidocs.mailchimp.com/api/1.3/listbatchunsubscribe.func.php
func (a *API) ListBatchUnsubscribe(parameters map[string]interface{}) (retVal *ListBatchUnsubsribeResponse, err error) {
	retVal = new(ListBatchUnsubsribeResponse)
	err = parseJson(a, "listBatchUnsubscribe", parameters, retVal)
	return
}

//ListClientsResponse is the type for values returned from method ListClients
//It implements the alterJsoner interface, changing the members property
//from string to int before unmarshaling
type ListClientsResponse struct {
	Desktop struct {
		Penetration float64
		Clients     []struct {
			Client  string
			Icon    string
			Percent float64
			Members int
		}
	}
	Mobile struct {
		Penetration float64
		Clients     []struct {
			Client  string
			Icon    string
			Percent float64
			Members int
		}
	}
}

var listClientsRX = regexp.MustCompile(`"members":"([0-9]*)"`)

func (r *ListClientsResponse) alterJson(b []byte) []byte {
	return listClientsRX.ReplaceAll(b, []byte(`"members":$1`))
}

//ListClients retrieves the clients that the list's subscribers have been
//tagged as being used based on user agents seen e.g. hotmail or iPhone.
func (a *API) ListClients(parameters map[string]interface{}) (retVal *ListClientsResponse, err error) {
	retVal = new(ListClientsResponse)
	err = parseJson(a, "listClients", parameters, retVal)
	return
}

//ListGrowthHistoryResponse is the type for values returned by
//the ListGrowthHistory method. It implements the alterJsoner interface,
//changing the existing, imports, and optins properties from strings to ints
//before unmarshaling, which is why the method has to return a named type
//instaed of a slice of named types as other slice-returning methods do
type ListGrowthHistoryResponse []ListGrowthHistoryElement
type ListGrowthHistoryElement struct {
	Month    ChimpTime
	Existing int
	Imports  int
	Optins   int
}
var listGrowthHistoryRX = regexp.MustCompile(`"(existing|imports|optins)":"([0-9]*)"`)
func (r *ListGrowthHistoryResponse) alterJson(b []byte) []byte {
	return listGrowthHistoryRX.ReplaceAll(b, []byte(`"$1":$2`))
}

//ListGrowthHistory accesses the growth history by month for a given list
//http://apidocs.mailchimp.com/api/1.3/listgrowthhistory.func.php
func (a *API) ListGrowthHistory(parameters map[string]interface{}) (retVal *ListGrowthHistoryResponse, err error) {
	retVal = new(ListGrowthHistoryResponse)
	err = parseJson(a, "listGrowthHistory", parameters, retVal)
	return
}

//ListInterestGroupAdd adds a single interest group, enabling interest groups
//for the list if necessary. http://apidocs.mailchimp.com/api/1.3/listinterestgroupadd.func.php
func (a *API) ListInterestGroupAdd(parameters map[string]interface{}) (bool, error) {
	return parseBoolean(run(a, "listInterestGroupAdd", parameters))
}

//ListInterestGroupDel deletes a single interest group and turns off groups
//for the list if it was the last group.  
//http://apidocs.mailchimp.com/api/1.3/listinterestgroupdel.func.php
func (a *API) ListInterestGroupDel(parameters map[string]interface{}) (bool, error) {
	return parseBoolean(run(a, "listInterestGroupDel", parameters))
}

//ListInterestGroupUpdate changes the name of an interest group
//http://apidocs.mailchimp.com/api/1.3/listinterestgroupupdate.func.php
func (a *API) ListInterestGroupUpdate(parameters map[string]interface{}) (bool, error) {
	return parseBoolean(run(a, "listInterestGroupUpdate", parameters))
}

//ListInterestGroupingAdd adds a new interest grouping, automatically enabling
//interest groups for the list if necessary
func (a *API) ListInterestGroupingAdd(parameters map[string]interface{}) (int, error) {
	return parseInt(run(a, "listInterestGroupingAdd", parameters))
}

//ListInterestGroupingUpdate updates an existing interest grouping
func (a *API) ListInterestGroupingUpdate(parameters map[string]interface{}) (bool, error) {
  return parseBoolean(run(a, "listInterestGroupingUpdate", parameters))
}

//ListInterestGroupingDel deletes an existing interest grouping, including all
//contained interest groups
func (a *API) ListInterestGroupingDel(parameters map[string]interface{}) (bool, error) {
  return parseBoolean(run(a, "listInterestGroupingDel", parameters))
}

//ListInterestGroupingsElement is the type of elements in the slice returned from the ListActivity method
type ListInterestGroupingsElement struct {
  Id int
  Name string
  Form_fields string
  Groups []struct {
    Bit string
    Name string
    Display_order string
    Subscribers int
  }
}
//ListInterestGroupings gets the list of interest groupings for a given list,
//including the lable, form information, and included groups for each
func (a *API) ListInterestGroupings(parameters map[string]interface{}) (retVal []ListInterestGroupingsElement, err error) {
  err = parseJson(a, "listInterestGroupings", parameters, &retVal)
  return
}
