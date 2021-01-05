package profari

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Composer composes an http request while recording its url, method, body(payload)
// in a Record struct. The pointer of the Record struct will later be passed to
// the Send() method of a Sender.
type Composer interface {
	Compose() (*http.Request, *Record, error)
}

// Sender can send the http request composed by a Composer, once it get an http response,
// it will add that response to the Resp field of the Record struct that was passed in before returning
// the complete Record struct as a result.
type Sender interface {
	Send(req Composer)
}

// Client is an integrated http client implementing Sender interface using value senmantics.
// A Client is responsible for:
// 1. Sending http requests
// 2. Asserting http responses
// 3. Logging
type Client struct {
	*http.Client
	*Assertor
}

func NewClient(testName string, errChan chan error) (*Client, *os.File, error) {
	logger, logFile, err := NewLogger(testName)
	if err != nil {
		return nil, nil, err
	}

	a := Assertor{
		Logger:  logger,
		ErrChan: errChan,
	}

	c := Client{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		Assertor: &a,
	}

	return &c, logFile, nil

}

// Send takes in a Composer, then use it to compose the actual http request and send it
// to the target server. Once a response is returned from the server, Sender will add it
// to the Resp field of the Record before returning the complete Record represented by Recorder.
func (c *Client) Send(req Composer) *Client {
	// compose request and add info to Record
	r, rec, err := req.Compose()
	if err != nil {
		c.Println(err)
		c.ErrChan <- err
	}

	if r == nil {
		err := fmt.Errorf("The request is nil, exiting the program")
		c.Println(err)
		c.ErrChan <- err
	}

	if rec == nil {
		err := fmt.Errorf("The record is nil, exiting the program")
		c.Println(err)
		c.ErrChan <- err
	}

	resp, err := c.Do(r)
	if err != nil {
		c.Println(err)
		c.ErrChan <- err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Println(err)
		c.ErrChan <- err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Error: response status code is: %v, response is: %v", resp.StatusCode, string(bodyBytes))
		c.Println(err)
		c.ErrChan <- err
	}
	rec.Resp = string(bodyBytes)

	c.Record = *rec
	return c
}

func (c *Client) DecodeJSON(v interface{}) *Client {
	d := json.NewDecoder(strings.NewReader(c.Resp))
	d.UseNumber()
	if err := d.Decode(&v); err != nil {
		err := fmt.Errorf("error while decoding JSON string: %v", err)
		c.Println(err.Error())
		c.ErrChan <- err
		time.Sleep(100 * time.Millisecond)
		return c
	}
	return c
}

func (c *Client) EndTest() {
	err := fmt.Errorf("done")
	c.ErrChan <- err
}

func (c *Client) FailTest(msg string) {
	err := fmt.Errorf(msg)
	c.Println(err)
	c.ErrChan <- err
}
