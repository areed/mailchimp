Mailchimp.go
=========

A golang wrapper for the Mailchimp API

## Usage
Create an API struct with the API version, method, apikey, and a map with all the parameters for the method except apikey.
Refer to the Mailchimp API Docs for the parameter options for each method.
If successful, API.Run() returns the response as a map.

### Example
``` go
filters := make(map[string]interface{})
filters["status"] = "sent"
paramaters := make(map[string]interface{})
parameters["filters"] = filters
api := mailchimp.API{"1.3", "campaigns", "abc123-us1", parameters}
resp, err := api.Run()
