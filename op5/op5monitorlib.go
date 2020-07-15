package op5monitor

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

type ConfigTree struct {
	URL    string
	User   string
	Pass   string
	Input  string
	Output string
}

// Check structure to replace using a map.
type CheckResults struct {
	Hostname           string `json:"host_name"`
	StatusCode         int    `json:"status_code"`
	PluginOutput       string `json:"plugin_output"`
	ServiceDescription string `json:"service_description,omitempty"`
}

func buildString(delimiter string, strList ...string) string {
	return strings.Join(strList, delimiter)
}

func Setup(server, user, pass, input, output string) (config ConfigTree) {
	log.Printf("Setting up config objects.\n")

	if server == "" {
		config.URL = server
	} else {
		config.URL = buildString("/", "https:/", server, "api")
	}
	config.User = user
	config.Pass = pass
	config.Input = input
	config.Output = output

	return config
}

func execRequest(
	config ConfigTree,
	requestType string,
	endpoint string,
	payload []byte,
) (res *http.Response, body string, err error) {
	// Making sure the requestType is always uppercase.
	requestType = strings.ToUpper(requestType)

	headerAccept := buildString("/", "application", config.Input)
	headerContentType := buildString("/", "application", config.Output)

	log.Printf("Request Type: %s\n", requestType)

	req, err := http.NewRequest(requestType, endpoint, bytes.NewBuffer(payload))

	if err != nil {
		log.Fatalf("Fatal error, %v\n", err)
		os.Exit(1)
	}

	query := req.URL.Query()
	query.Add("format", config.Output)
	req.URL.RawQuery = query.Encode()

	// This doesn't do anything, but it should.
	req.Header.Add("accept", headerAccept)
	if requestType == "POST" || requestType == "PUT" || requestType == "PATCH" {
		req.Header.Add("content-type", headerContentType)
	}

	req.SetBasicAuth(config.User, config.Pass)

	log.Printf("Complete URL: %s\n", req.URL.String())
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("HTTP Request:\n%s", reqDump)
	// Execute the request and set the response.
	res, err = http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		log.Fatalf("Something went wrong with the %s request, %s\n", requestType, err)
	}
	resDump, err := httputil.DumpResponse(res, true)

	if err != nil {
		log.Fatalf("Problem dumping the response, %s\n", err)
	}

	bodyExtract, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("There is a problem with the body, %s\n", err)
	}

	body = string(bodyExtract)

	log.Printf("Result Dump:\n%s\n", resDump)

	return res, body, err
}

func HostFetch(config ConfigTree, host string) {
	endpoint := buildString("/", config.URL, "config", "host", host)

	execRequest(config, "GET", endpoint, nil)
}

//func HostReplace(config ConfigTree, data map[string]interface{}) {
//	execRequest(config, "PUT", endpoint, jsonPayload)
//}

//func HostUpdate(config ConfigTree, data map[string]interface{}) {
//	execRequest(config, "PATCH", endpoint, jsonPayload)
//}

func HostCreate(config ConfigTree, data map[string]interface{}, replaceHost bool) (ok bool) {
	endpoint := buildString("/", config.URL, "config", "host")

	jsonPayload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("JSON encoding failed, %v\n", err)
		return false
	}

	execRequest(config, "POST", endpoint, jsonPayload)

	if replaceHost {
		log.Printf("(Not Implemented) Replacing host.\n")
	}

	return true
}

func HostDelete(config ConfigTree, hostname string) bool {
	endpoint := buildString("/", config.URL, "config", "host", hostname)

	_, _, err := execRequest(config, "DELETE", endpoint, nil)

	if err != nil {
		return false
	}

	return true
}

func ServiceCreate(config ConfigTree, data map[string]interface{}) (ok bool) {
	endpoint := buildString("/", config.URL, "config", "service")

	jsonPayload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("JSON encoding failed, %v\n", err)
		return false
	}

	execRequest(config, "POST", endpoint, jsonPayload)

	return true
}

func QueueSave(config ConfigTree) bool {
	endpoint := buildString("/", config.URL, "config", "change")
	// Empty JSON object to post.
	data := make(map[string]interface{})

	jsonPayload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("JSON encoding failed, %v\n", err)
		return false
	}

	execRequest(config, "POST", endpoint, jsonPayload)

	return true
}

func SendCheck(config ConfigTree, checkType string, data CheckResults) {
	command := ""
	checkType = strings.ToLower(checkType)

	if checkType == "host" {
		command = "PROCESS_HOST_CHECK_RESULT"
	} else if checkType == "service" {
		command = "PROCESS_SERVICE_CHECK_RESULT"
	} else {
		// Throw an error and return is the better option.
		command = ""
		log.Fatalf("Check type not specified.\n")
		os.Exit(1)
	}

	endpoint := buildString("/", config.URL, "command", command)

	jsonPayload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("JSON encoding failed, %v\n", err)
		os.Exit(1)
	}

	execRequest(config, "POST", endpoint, jsonPayload)
}

func SendCheckHost(
	config ConfigTree,
	hostname string,
	statusCode int,
	output string,
) {

	data := CheckResults{
		Hostname:     hostname,
		StatusCode:   statusCode,
		PluginOutput: output,
	}

	SendCheck(config, "host", data)
}

func SendCheckService(
	config ConfigTree,
	hostname string,
	service string,
	statusCode int,
	output string,
) {
	data := CheckResults{
		Hostname:           hostname,
		StatusCode:         statusCode,
		PluginOutput:       output,
		ServiceDescription: service,
	}

	SendCheck(config, "service", data)
}
