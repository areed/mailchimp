package mailchimp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
//	"log"
	"net/http"
	"net/url"
//	"os"
	"regexp"
	"strconv"
	"time"
)

const ChimpTime = "2006-01-02 15:04:05"

var datacenter = regexp.MustCompile("[a-z]+[0-9]+$")

type api struct {
	Key      string
	endpoint string
}

type ChimpError struct {
	Err string `json:"error"`
	Code int
}
func (e ChimpError) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Err)
}

func New(apikey string, https bool) (*api, error) {
	u := url.URL{}
	if https {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}
	u.Host = fmt.Sprintf("%s.api.mailchimp.com", datacenter.FindString(apikey))
	u.Path = "/1.3/"
	return &api{apikey, u.String() + "?method="}, nil
}

func run(a *api, method string, parameters map[string]interface{}) ([]byte, error) {
	if parameters == nil {
		parameters = make(map[string]interface{})
	}
	parameters["apikey"] = a.Key
	b, err := json.Marshal(parameters)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(a.endpoint + method, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func verify(body []byte) error {
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
	return t;
}

func parseInt(body []byte) (int, error) {
	if err := verify(body); err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(string(body), 10, 0)
	return int(i), err
}

func parseString(body []byte) (string, error) {
	if err := verify(body); err != nil {
		return "", err
	}
	return strconv.Unquote(string(body))
}

func parseStruct(body []byte, ret interface{}) error {
	if err := verify(body); err != nil {
		return err
	}
	json.Unmarshal(body, ret)
	return nil
}

func parseBoolean(body []byte) (bool, error) {
	if err := verify(body); err != nil {
		return false, err
	}
	return strconv.ParseBool(string(body))
}

type CampaignContentResult struct {
	Html string
	Text string
}
func (a *api) CampaignContent(parameters map[string]interface{}) (*CampaignContentResult, error) {
	body, err := run(a, "campaignContent", parameters)
	if err != nil {
		return nil, err
	}
	var ccr CampaignContentResult
	err = parseStruct(body, &ccr)
	return &ccr, nil
}

func (a *api) CampaignCreate(parameters map[string]interface{}) (string, error) {
	body, err := run(a, "campaignCreate", parameters)
	if err != nil {
		return "", err
	}
	return parseString(body)
}

func (a *api) CampaignDelete(parameters map[string]interface{}) (bool, error) {
	body, err := run(a, "campaignDelete", parameters)
	if err != nil {
		return false, err
	}
	return parseBoolean(body)
}

//not tested
func (a *api) CampaignEcommOrderAdd(parameters map[string]interface{}) (bool, error) {
	body, err := run(a, "campaignEcommOrderAdd", parameters)
	if err != nil {
		return false, err
	}
	return parseBoolean(body)
}

func (a *api) CampaignPause(parameters map[string]interface{}) (bool, error) {
	body, err := run(a, "campaignPause", parameters)
	if err != nil {
		return false, err
	}
	return parseBoolean(body)
}


func (a *api) CampaignReplicate(parameters map[string]interface{}) (string, error) {
	resp, err := run(a, "campaignReplicate", parameters)
	if err != nil {
		return "", err
	}
	return parseString(resp)
}

func (a *api) CampaignResume(parameters map[string]interface{}) (bool, error) {
	body, err := run(a, "campaignResume", parameters)
	if err != nil {
		return false, err
	}
	return parseBoolean(body)
}

func (a *api) CampaignSchedule(parameters map[string]interface{}) (bool, error) {
	if parameters == nil {
		parameters = make(map[string]interface{})
	}
	parameters["schedule_time"] = chimpTime(parameters["schedule_time"])
	if t, exist := parameters["schedule_time_b"]; exist {
		parameters["schedule_time_b"] = chimpTime(t)
	}
	body, err := run(a, "campaignSchedule", parameters)
	if err != nil {
		return false, err
	}
	return parseBoolean(body)
}

func (a *api) CampaignSegmentTest(parameters map[string]interface{}) (int, error) {
	body, err := run(a, "campaignSegmentTest", parameters)
	if err != nil {
		return 0, err
	}
	return parseInt(body)
}

func (a *api) CampaignSendNow(parameters map[string]interface{}) (bool, error) {
	body, err := run(a, "campaignSendNow", parameters)
	if err != nil {
		return false, err
	}
	return parseBoolean(body)
}

func (a *api) CampaignSendTest(parameters map[string]interface{}) (bool, error) {
	body, err := run(a, "campaignSendTest", parameters)
	if err != nil {
		return false, err
	}
	return parseBoolean(body)
}

type CampaignShareReportResult struct {
	Title string
	Url string
	Secure_url string
	Password string
}
func (a *api) CampaignShareReport(parameters map[string]interface{}) (*CampaignShareReportResult, error) {
	body, err := run(a, "campaignShareReport", parameters)
	if err != nil {
		return nil, err
	}
	var csrr CampaignShareReportResult
	if err = parseStruct(body, &csrr); err != nil {
		return nil, err
	}
	return &csrr, nil
}

func (a *api) CampaignTemplateContent(parameters map[string]interface{}) (map[string]interface{}, error) {
	body, err := run(a, "campaignTemplateContent", parameters)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	if err = json.Unmarshal(body, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (a *api) CampaignUnschedule(parameters map[string]interface{}) (bool, error) {
	body, err := run(a, "campaignUnschedule", parameters)
	if err != nil {
		return false, err
	}
	return parseBoolean(body)
}

func (a *api) CampaignUpdate(parameters map[string]interface{}) (bool, error) {
	body, err := run(a, "campaignUpdate", parameters)
	if err != nil {
		return false, err
	}
	return parseBoolean(body)
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
func (a *api) Campaigns(parameters map[string]interface{}) (*CampaignsResult, error) {
	body, err := run(a, "campaigns", parameters)
	var cr CampaignsResult
	if err = parseStruct(body, &cr); err != nil {
		return nil, err
	}
	return &cr, nil
/*
	//mock response for development
	file, err := os.Open("/home/ec2-user/go/src/github.com/areed/mailchimp/go.json")
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(file)
*/
}

func (a *api) Ping() (result string, err error) {
	body, err := run(a, "ping", make(map[string]interface{}))
	return parseString(body)
}
