package mailchimp

import (
	"bytes"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

var datacenter = regexp.MustCompile("[a-z]+[0-9]+$")

type api struct {
	Key      string
	endpoint string
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
	return &api{apikey, u.String()}, nil
}

type CampaignsResult struct {
	Total int
	Data []CampaignsResultData
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
	Tracking CampaignsResultDataTracking
	Segment_text string
	Segment_opts []CampaignsResultDataSegment_opts
	Type_opts []map[string]interface{}
}

type CampaignsResultDataTracking struct {
	Html_clicks bool
	Text_clicks bool
	Opens bool
}

type CampaignsResultDataSegment_opts struct {
	Match string
	Conditions []map[string]interface{}
}

func (a *api) Campaigns(parameters map[string]interface{}) (*CampaignsResult, error) {
/*
	//all passed in parameters are optional, but apikey is required
	if parameters == nil {
		parameters = make(map[string]interface{})
	}
	parameters["apikey"] = a.Key
	b, err := json.Marshal(parameters)
	if err != nil {
		return nil, err
	}
	//don't modify a directly because it may be used by other methods
	u := a.endpoint + "?method=campaigns"
	resp, err := http.Post(u, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
*/
//mock response for development
	file, err := os.Open("/home/ec2-user/go/src/github.com/areed/mailchimp/campaigns.json")
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	//create a temporary source map where the json blob can be unmarshalled and accessed;
	//data will be taken from source map, type asserted, and placed in destination 
	//struct that will be returned by Campaigns method
	var sourceMap map[string]interface{}
	if err = json.Unmarshal([]byte(line), &sourceMap); err != nil {
		return nil, err
	}
	var result CampaignsResult
	result.Total = int(sourceMap["total"].(float64))
	//assert data in source map is a slice so we can loop over each item
	data := sourceMap["data"].([]interface{})
	for _, item := range data {
		//assert each item in data is a map so it's accessible
		datum := item.(map[string]interface{})
		//the destination struct for this datum
		var crd CampaignsResultData
		//the form of type assertion that returns two values causes
		//the first value to be the zero value for the type if the assertion
		//fails, which it will if the value is nil (json null);
		//type assertion of nil to string with a single return value
		//would panic upon failure
		crd.Id, _ = datum["id"].(string)
		f64, _ := datum["web_id"].(float64)
		crd.Web_id = int(f64)
		crd.List_id, _ = datum["list_id"].(string)
		f64, _ = datum["folder_id"].(float64)
		crd.Folder_id = int(f64)
		f64, _ = datum["template_id"].(float64)
		crd.Template_id = int(f64)
		crd.Content_type, _ = datum["content_type"].(string)
		crd.Title, _ = datum["title"].(string)
		crd.Type, _ = datum["type"].(string)
		crd.Create_time, _ = datum["create_time"].(string)
		crd.Send_time, _ = datum["send_time"].(string)
		f64, _ = datum["emails_sent"].(float64)
		crd.Emails_sent = int(f64)
		crd.Status, _  = datum["status"].(string)
		crd.From_name, _ = datum["from_name"].(string)
		crd.From_email, _ = datum["from_email"].(string)
		crd.Subject, _ = datum["subject"].(string)
		crd.To_name, _ = datum["to_name"].(string)
		crd.Archive_url, _ = datum["archive_url"].(string)
		crd.Inline_css, _ = datum["inline_css"].(bool)
		crd.Analytics, _ = datum["analytics"].(string)
		crd.Analytics_tag, _ = datum["analytics_tag"].(string)
		crd.Authenticate, _ = datum["authenticate"].(bool)
		crd.Ecomm360, _ = datum["ecomm360"].(bool)
		crd.Auto_tweet, _ = datum["auto_tweet"].(bool)
		crd.Auto_fb_post, _ = datum["auto_fb_post"].(string)
		crd.Auto_footer, _ = datum["auto_footer"].(bool)
		crd.Timewarp, _ = datum["timewarp"].(bool)
		crd.Timewarp_schedule, _  = datum["timewarp_schedule"].(string)
		//populate tracking struct
		//TODO: can probably directly do crd.Tracking.Html_clicks = tracking["html_clicks"].(bool)
		tracking := datum["tracking"].(map[string]interface{})
		var crdt CampaignsResultDataTracking
		crdt.Html_clicks, _ = tracking["html_clicks"].(bool)
		crdt.Text_clicks, _ = tracking["text_clicks"].(bool)
		crdt.Opens, _ = tracking["opens"].(bool)
		crd.Tracking = crdt
		crd.Segment_text, _ = datum["segment_text"].(string)
		//populate segment options array
		segment_opts := datum["segment_opts"].([]interface{})
		for _, item := range segment_opts {
			segment_opt := item.(map[string]interface{})
			var crds CampaignsResultDataSegment_opts
			crds.Match = segment_opt["match"].(string)
			conditions := segment_opt["conditions"].([]interface{})
			for _, item := range conditions {
				condition := item.(map[string]interface{})
				crds.Conditions = append(crds.Conditions, condition)
			}
		}
		//populate type options array
		type_opts := datum["type_opts"].([]interface{})
		for _, item := range type_opts {
			type_opt := item.(map[string]interface{})
			crd.Type_opts = append(crd.Type_opts, type_opt)
		}
		result.Data = append(result.Data, crd)
	}
	return &result, nil
}

func (a *api) Ping() (result string, err error) {
	u := a.endpoint + "?method=ping"
	b, err := json.Marshal(map[string]string{"apikey": a.Key})
	if err != nil {
		return "", err
	}
	resp, err := http.Post(u, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
