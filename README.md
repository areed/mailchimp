Mailchimp.go
=========

A golang wrapper for the Mailchimp API.  Refer to the [Mailchimp API Docs](http://apidocs.mailchimp.com/api/1.3) for the correct parameters to pass to each routine and the expected return value.

## Usage
``` go
//first call the constructor
//func New (apikey string, useHttps bool) *api
chimp := mailchimp.New("abcdefg-us1", true)

//then call methods on the constructor
chimp.Ping()
//"Everything's Chimpy"

//refer to the Mailchimp API Docs for the required and optional parameters for each
//routine and assemble them into a map[string]interface{}, excluding apikey
filters := make(map[string]interface{})
filters["status"] = "sent"
parameters := make(map[string]interface{})
parameters["filters"] = filters
parameters["start"] = 0
parameters["limit"] = 10

//pass the parameter map to the method
result, err := chimp.Campaigns(parameters)

//Campaigns returns a struct with all constant return values correctly typed
if result.Data[0].Status != "Sent" {
	panic("should not panic unless there were no matching campaigns")
}
//variable return objects will be returned as map[string]interface{}
field := result.Data[0].Segment_opts.Conditions[0]["field"].(string)
if field == "rating" {
	value := result.Data[0].Segment_opts.Conditions[0]["value"].(int)
}	

//Fortunately, most routines return simple string, boolean, or int values
result, err := chimp.CampaignUnschedule(map[string]interface{}{"cid": "abcdefghij"})
//true, nil
