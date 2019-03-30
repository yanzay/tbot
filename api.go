package tbot

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type responseParameters struct {
	MigrateToChatID int `json:"migrate_to_chat_id"`
	ReplyAfter      int `json:"retry_after"`
}

type apiResponse struct {
	OK          bool                `json:"ok"`
	Result      json.RawMessage     `json:"result"`
	Description string              `json:"description"`
	ErrorCode   int                 `json:"error_code"`
	Parameters  *responseParameters `json:"parameters"`
}

func (c *Client) doRequest(method string, request url.Values, response interface{}) error {
	endpoint := fmt.Sprintf(c.url, method)
	var resp *http.Response
	var err error
	if request == nil {
		resp, err = c.httpClient.Post(endpoint, "application/x-www-form-urlencoded", nil)
	} else {
		resp, err = c.httpClient.PostForm(endpoint, request)
	}
	if err != nil {
		return fmt.Errorf("unable to send message: %v", err)
	}

	apiResp := &apiResponse{}
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return fmt.Errorf("unable to decode sendMessage response: %v", err)
	}
	err = resp.Body.Close()
	if err != nil {
		c.logger.Errorf("unable to close response body: %v", err)
	}
	if !apiResp.OK {
		return fmt.Errorf(apiResp.Description)
	}
	return json.Unmarshal(apiResp.Result, response)
}

func (c *Client) doRequestWithFiles(method string, request url.Values, response interface{}, files ...inputFile) error {
	endpoint := fmt.Sprintf(c.url, method)
	r, w := io.Pipe()

	done := make(chan struct{})
	var resp *http.Response
	var err error

	mw := multipart.NewWriter(w)

	go func() {
		defer close(done)
		req, err := http.NewRequest(http.MethodPost, endpoint, r)
		if err != nil {
			c.logger.Error(err)
			return
		}
		req.Header.Set("Content-Type", mw.FormDataContentType())
		resp, err = c.httpClient.Do(req)
	}()

	for k := range request {
		mw.WriteField(k, request.Get(k))
	}
	for _, file := range files {
		f, err := os.Open(file.name)
		if err != nil {
			return err
		}
		fileWriter, err := mw.CreateFormFile(file.field, file.name)
		if err != nil {
			return err
		}

		io.Copy(fileWriter, f)
		f.Close()
	}

	mw.Close()
	w.Close()

	<-done // post request is done
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}
	apiResp := &apiResponse{}
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return fmt.Errorf("unable to decode sendMessage response: %v", err)
	}
	err = resp.Body.Close()
	if err != nil {
		c.logger.Errorf("unable to close response body: %v", err)
	}
	if !apiResp.OK {
		return fmt.Errorf(apiResp.Description)
	}
	return json.Unmarshal(apiResp.Result, response)
}
