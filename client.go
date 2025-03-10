package appwrite

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

// Client is the client struct to access Appwrite services
type Client struct {
	client     *http.Client
	endpoint   string
	headers    map[string]string
	selfSigned bool
}

// SetEndpoint sets the default endpoint to which the Client connects to
func (clt *Client) SetEndpoint(endpoint string) {
	clt.endpoint = endpoint
}

// SetSelfSigned sets the condition that specify if the Client should allow connections to a server using a self-signed certificate
func (clt *Client) SetSelfSigned(status bool) {
	clt.selfSigned = status
}

// AddHeader add a new custom header that the Client should send on each request
func (clt *Client) AddHeader(key string, value string) {
	clt.headers[key] = value
}

// Your project ID
func (clt *Client) SetProject(value string) {
	clt.headers["X-Appwrite-Project"] = value
}

// Your secret API key
func (clt *Client) SetKey(value string) {
	clt.headers["X-Appwrite-Key"] = value
}

func (clt *Client) SetLocale(value string) {
	clt.headers["X-Appwrite-Locale"] = value
}

func (clt *Client) SetMode(value string) {
	clt.headers["X-Appwrite-Mode"] = value
}

// Call an API using Client
func (clt *Client) Call(method string, path string, headers map[string]interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	clt.ensureClientInitialized()

	urlPath := clt.endpoint + path
	isGet := strings.ToUpper(method) == "GET"

	var reqBody io.Reader
	if !isGet {
		reqBody = prepareRequestBody(params)
	}

	req, err := http.NewRequest(method, urlPath, reqBody)
	if err != nil {
		return nil, err
	}

	setHeaders(req, clt.headers, headers)

	if isGet {
		updateQueryParameters(req, params)
	} else {
		// Set the Content-Type header for non-GET requests
		req.Header.Set("Content-Type", "application/json")
	}

	response, err := clt.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	jsonResponse, err := parseJSONResponse(response)
	if err != nil {
		return nil, err
	}

	return jsonResponse, nil
}

func (clt *Client) ensureClientInitialized() {
	if clt.client == nil {
		// Create HTTP client if it's not initialized
		clt.client = &http.Client{}
	}
}

func prepareRequestBody(params map[string]interface{}) io.Reader {
	jsonData, err := json.Marshal(params)
	if err != nil {
		// Handle the error
		return nil
	}
	log.Println(string(jsonData))
	return bytes.NewReader(jsonData)
}

func setHeaders(req *http.Request, clientHeaders map[string]string, customHeaders map[string]interface{}) {
	// Set Client headers
	for key, val := range clientHeaders {
		req.Header.Set(key, val)
	}

	// Set Custom headers
	for key, val := range customHeaders {
		req.Header.Set(key, ToString(val))
	}
}

func updateQueryParameters(req *http.Request, params map[string]interface{}) {
	q := req.URL.Query()
	for key, val := range params {
		q.Add(key, ToString(val))
	}
	req.URL.RawQuery = q.Encode()
}

func parseJSONResponse(response *http.Response) (map[string]interface{}, error) {
	var jsonResponse map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&jsonResponse)
	if err != nil {
		return nil, err
	}
	return jsonResponse, nil
}
