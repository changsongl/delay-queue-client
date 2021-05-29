package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/changsongl/delay-queue-client/api"
	"github.com/go-resty/resty/v2"
	"net/http"
)

const (
	addUrlFormat = "/topic/%s/job"
	endUrlFormat = "/topic/%s/job/%s"
	popUrlFormat = "/topic/%s/job"
)

type respBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Id      string `json:"id"`
	Body    string `json:"body"`
}

type requester struct {
	req  *resty.Request
	host string
}

func NewRequester(host string) api.Request {
	httpReq := resty.New().R()
	r := &requester{req: httpReq, host: host}

	return r
}

func (r *requester) AddDelayJob(topic, id, body string, delay, ttr uint) error {
	return r.addDelayJob(topic, id, body, delay, ttr, false)
}

func (r *requester) ReplaceDelayJob(topic, id, body string, delay, ttr uint) error {
	return r.addDelayJob(topic, id, body, delay, ttr, true)
}

func (r *requester) addDelayJob(topic, id, body string, delay, ttr uint, override bool) error {
	reqBody := r.getAddRequestMap(id, body, delay, ttr, override)
	url := r.getHost() + getAddUrl(topic)
	resp, err := r.req.SetBody(reqBody).Post(url)
	if err != nil {
		return err
	}

	_, err = r.getRespBody(resp)
	if err != nil {
		return err
	}

	return nil
}

func (r *requester) FinishDelayJob(topic string, id string) error {
	url := r.getHost() + getEndUrl(topic, id)
	resp, err := r.req.Put(url)
	if err != nil {
		return err
	}

	_, err = r.getRespBody(resp)
	if err != nil {
		return err
	}

	return nil
}

func (r *requester) DeleteDelayJob(topic string, id string) error {
	url := r.getHost() + getEndUrl(topic, id)
	resp, err := r.req.Delete(url)
	if err != nil {
		return err
	}

	_, err = r.getRespBody(resp)
	if err != nil {
		return err
	}

	return nil
}

func (r *requester) PopDelayJob(topic string) (id string, body string, err error) {
	url := r.getHost() + getPopUrl(topic)
	resp, err := r.req.Get(url)
	if err != nil {
		return "", "", err
	}

	respBody, err := r.getRespBody(resp)
	if err != nil {
		return "", "", err
	}

	return respBody.Id, respBody.Body, nil
}

func (r *requester) getHost() string {
	return r.host
}

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

func getAddUrl(topic string) string {
	return fmt.Sprintf(addUrlFormat, topic)
}

func getEndUrl(topic, id string) string {
	return fmt.Sprintf(endUrlFormat, topic, id)
}

func getPopUrl(topic string) string {
	return fmt.Sprintf(popUrlFormat, topic)
}
