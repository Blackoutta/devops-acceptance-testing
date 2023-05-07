package req

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"time"

	"github.com/Blackoutta/profari"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/conf"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/errors"
)

func ComposeNewMultipartRequest(method string, path string, params url.Values, fname string, contentType string) *http.Request {
	baseURL, err := url.Parse(conf.Host)
	errors.HandleError("Malformed URL", err)
	baseURL.Path += path
	baseURL.RawQuery = params.Encode()

	buff := &bytes.Buffer{}
	wr := multipart.NewWriter(buff)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%v"`, fname))
	h.Set("Content-Type", contentType)
	mimewr, err := wr.CreatePart(h)
	errors.HandleError("err creating part", err)

	bs, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("error while reading file: %v\n", err)
	}

	_, err = mimewr.Write(bs)
	errors.HandleError("err writing mime", err)
	wr.Close()

	r, err := http.NewRequest(method, baseURL.String(), buff)
	errors.HandleError("err composing request", err)
	r.Close = true
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %v", conf.GetToken()))
	r.Header.Set("Content-Type", wr.FormDataContentType())
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Accept-Encoding", "gzip, deflate")
	r.Header.Set("Accept-Language", "zh-CN,en;q=0.8,zh;q=0.7,zh-TW;q=0.5,zh-HK;q=0.3,en-US;q=0.2")
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:73.0) Gecko/20100101 Firefox/73.0")
	return r
}

func ComposeNewRequest(method string, path string, params url.Values, body io.Reader) *http.Request {
	baseURL, err := url.Parse(conf.Host)
	errors.HandleError("Malformed URL", err)
	baseURL.Path += path
	baseURL.RawQuery = params.Encode()

	r, err := http.NewRequest(method, baseURL.String(), body)
	errors.HandleError("err composing request", err)
	r.Close = true
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %v", conf.GetToken()))
	r.Header.Set("X-UserId", conf.UserID)
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json;charset=utf-8")
	return r
}

func SendRequestAndGetResponse(c http.Client, r *http.Request) Record {
	var rdr1 io.ReadCloser
	if r.Body != nil {
		buf, bodyErr := ioutil.ReadAll(r.Body)
		if bodyErr != nil {
			log.Fatalln(bodyErr)
		}
		rdr1 = ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
		r.Body = rdr2
	}

	resp, sendErr := c.Do(r)

	var retry int
	for sendErr != nil {
		log.Println(sendErr)
		if retry > 3 {
			log.Printf("retried 3 times and failed, got err: %v\n", sendErr)
			break
		}
		if strings.Contains(sendErr.Error(), "EOF") {
			log.Println("retrying...")
			time.Sleep(time.Second)
			r.Body = rdr1
			resp, sendErr = c.Do(r)
			retry++
			continue
		}
		panic(sendErr)
	}

	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	errors.HandleError("err reading response", err)
	rl := Record{
		Method:     r.Method,
		URL:        r.URL.String(),
		Body:       rdr1,
		Response:   bs,
		StatusCode: resp.StatusCode,
	}

	return rl
}

type Record struct {
	Method     string
	URL        string
	Body       io.ReadCloser
	Response   []byte
	StatusCode int
}

func ComposeNewProfariRequest(method string, path string, params url.Values, body interface{}) (*http.Request, *profari.Record, error) {
	baseURL, err := url.Parse(conf.Host)
	if err != nil {
		return nil, nil, fmt.Errorf("req: error while parsing URL: %v", err)
	}
	baseURL.Path += path
	baseURL.RawQuery = params.Encode()

	var bodyBytes []byte
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("req: error while marshaling json: %v", err)
		}
	}

	r, err := http.NewRequest(method, baseURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, nil, fmt.Errorf("req: error while composing new request: %v", err)
	}
	r.Close = true
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %v", conf.GetToken()))
	r.Header.Set("X-UserId", conf.UserID)
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json;charset=utf-8")

	rec := profari.Record{
		Url:    baseURL.String(),
		Method: method,
		Body:   string(bodyBytes),
	}
	return r, &rec, nil
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
