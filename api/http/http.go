package http

import (
	"encoding/json"
	"errors"
	"github.com/changsongl/delay-queue-client/api"
	"github.com/changsongl/delay-queue-client/common"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

// respBody http response body
type respBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *Data  `json:"data"`
}

// Data http response data
type Data struct {
	Topic string `json:"topic"`
	ID    string `json:"id"`
	Body  string `json:"body"`
	Delay uint64 `json:"delay"`
	TTR   uint64 `json:"ttr"`
}

// requester implemented api.Request
type requester struct {
	req  *resty.Client
	host string
}

// NewRequester create a new api.Request by requester
func NewRequester(host string) api.Request {
	httpReq := resty.New()
	r := &requester{req: httpReq, host: host}

	return r
}

// get resty request
func (r *requester) getRequest() *resty.Request {
	return r.req.R()
}

// AddDelayJob add job to delay queue, if job is already exist, it will return error
func (r *requester) AddDelayJob(topic, id, body string, delay, ttr uint) error {
	return r.addDelayJob(topic, id, body, delay, ttr, false)
}

// ReplaceDelayJob replace job in delay queue, it will replace or add job to delay queue
func (r *requester) ReplaceDelayJob(topic, id, body string, delay, ttr uint) error {
	return r.addDelayJob(topic, id, body, delay, ttr, true)
}

// add job to delay queue
func (r *requester) addDelayJob(topic, id, body string, delay, ttr uint, override bool) error {
	reqBody := r.getAddRequestMap(id, body, delay, ttr, override)
	url := r.getHost() + addJobPath(topic)
	resp, err := r.getRequest().SetBody(reqBody).Post(url)
	if err != nil {
		return err
	}

	_, err = r.getRespBody(resp)
	if err != nil {
		return err
	}

	return nil
}

// FinishDelayJob finish job, when have processed job
func (r *requester) FinishDelayJob(topic string, id string) error {
	url := r.getHost() + finishJobPath(topic, id)
	resp, err := r.getRequest().Put(url)
	if err != nil {
		return err
	}

	_, err = r.getRespBody(resp)
	if err != nil {
		return err
	}

	return nil
}

// DeleteDelayJob delete job from delay queue
func (r *requester) DeleteDelayJob(topic string, id string) error {
	url := r.getHost() + finishJobPath(topic, id)
	resp, err := r.getRequest().Delete(url)
	if err != nil {
		return err
	}

	_, err = r.getRespBody(resp)
	if err != nil {
		return err
	}

	return nil
}

// PopDelayJob pop job from delay queue
func (r *requester) PopDelayJob(topic string, timeout time.Duration) (topicName, id, body string, delay, ttr uint64, err error) {
	url := r.getURL(popJobPath(topic, timeout))
	resp, err := r.getRequest().Get(url)
	if err != nil {
		return "", "", "", 0, 0, err
	}

	respBody, err := r.getRespBody(resp)
	if err != nil {
		return "", "", "", 0, 0, err
	}

	if respBody.Data == nil {
		return "", "", "", 0, 0, common.ErrorNoAvailableJob
	}

	return respBody.Data.Topic, respBody.Data.ID, respBody.Data.Body,
		respBody.Data.Delay, respBody.Data.TTR, nil
}

// get host of delay queue
func (r *requester) getHost() string {
	return r.host
}

// get response body
func (r *requester) getRespBody(resp *resty.Response) (*respBody, error) {
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("request failed")
	}

	body := &respBody{}
	err := json.Unmarshal(resp.Body(), body)
	if err != nil {
		return nil, err
	} else if !body.Success {
		return nil, errors.New(body.Message)
	}

	return body, nil
}

// get request map
func (r *requester) getAddRequestMap(id, body string,
	delay, ttr uint, override bool) map[string]interface{} {

	return map[string]interface{}{
		"id":       id,
		"body":     body,
		"delay":    delay,
		"ttr":      ttr,
		"override": override,
	}
}

// get url
func (r *requester) getURL(path string) string {
	return r.getHost() + path
}
