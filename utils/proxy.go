package utils

import (
	"bytes"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	MaxIdleConnections = 4
	RequestTimeout     = 60
)

var httpClient *http.Client
var proxyOnce sync.Once

type box struct{}

func SharedProxy() *box {
	proxyOnce.Do(func() {
		httpClient = createHTTPClient()
	})
	return &box{}
}

func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}

	return client
}

func (b *box) Get(baseURL string, url string, setter ...func(*http.Request)) (body []byte, header http.Header, statusCode int) {
	return req(baseURL, `GET`, url, nil, setter...)
}

func (b *box) Post(baseURL string, url string, data url.Values, setter ...func(*http.Request)) (body []byte, header http.Header, statusCode int) {
	return req(baseURL, "POST", url, data, setter...)
}

func (b *box) PostJson(baseURL string, url string, data []byte, setter ...func(*http.Request)) (body []byte, header http.Header, statusCode int) {
	return reqJson(baseURL, "POST", url, data, setter...)
}

func (b *box) PostJsonWithHeader(head map[string]string, baseURL string, url string, data []byte, setter ...func(*http.Request)) (body []byte, header http.Header, statusCode int) {
	return reqJsonWithHeader(head, baseURL, "POST", url, data, setter...)
}

func reqJsonWithHeader(head map[string]string, baseURL string, method string, url string, data []byte, setters ...func(*http.Request)) (body []byte, header http.Header, statusCode int) {
	url = baseURL + url
	//mylog.LOG.I("reqJson,url:%v",url)
	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		mylog.Logger.Error().Msgf("reqJson: NewRequest failed: %v", err.Error())
		return
	}
	if method == "POST" || method == "PUT" || method == "DELETE" {
		request.Header.Set("Content-Type", "application/json")
		request.Header.Add("Accept-Charset", "UTF-8")
		for k, v := range head {
			if _, ok := request.Header[k]; ok {
				request.Header.Set(k, v)
			}
			request.Header.Add(k, v)
		}
	}

	return reqInner(request, setters...)
}

func reqJson(baseURL string, method string, url string, data []byte, setters ...func(*http.Request)) (body []byte, header http.Header, statusCode int) {
	url = baseURL + url
	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		mylog.Logger.Error().Msgf("reqJson: NewRequest failed: %v", err.Error())
		return
	}
	if method == "POST" || method == "PUT" || method == "DELETE" {
		request.Header.Set("Content-Type", "application/json")
		request.Header.Add("Accept-Charset", "UTF-8")
	}

	return reqInner(request, setters...)
}

func req(baseURL string, method string, url string, data url.Values, setters ...func(*http.Request)) (body []byte, header http.Header, statusCode int) {
	url = baseURL + url
	request, err := http.NewRequest(method, url, strings.NewReader(data.Encode()))
	if err != nil {
		return
	}
	if method == "POST" || method == "PUT" {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Add("Accept-Charset", "UTF-8")
	}
	return reqInner(request, setters...)
}

func PostFile(baseURL string, url string, data string) (body []byte, header http.Header, statusCode int) {
	return reqFile(baseURL, `POST`, url, data)
}

func reqFile(baseURL string, method string, url string, data string, setters ...func(*http.Request)) (body []byte, header http.Header, statusCode int) {
	url = baseURL + url
	request, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		return
	}

	request.Header.Set("Content-Type", `text/xml;charset=UTF-8`)
	request.Header.Add("Accept", `application/soap+xml, application/dime, multipart/related, text/*`)
	return reqInner(request, setters...)
}

func reqInner(request *http.Request, setters ...func(*http.Request)) (body []byte, header http.Header, statusCode int) {
	var (
		err error
		res *http.Response
	)

	for _, setter := range setters {
		setter(request)
	}

	res, err = httpClient.Do(request)
	if err != nil {
		mylog.Logger.Error().Msgf("reqInner httpClient.Do failed: %v", err.Error())
		return
	}

	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		mylog.Logger.Error().Msgf("reqInner ioutil.ReadAll(res.Body) failed: %v", err.Error())
		return
	}

	header = res.Header
	statusCode = res.StatusCode

	return
}
