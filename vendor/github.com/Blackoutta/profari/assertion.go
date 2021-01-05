package profari

import (
	"fmt"
	"log"
	"strings"
)

const (
	head        = "============================ ERROR START ============================"
	tail        = "============================ ERROR END =============================="
	titlePrefix = "Assertion: "
	equal       = "%v\n%v\nAssertion Error: Expect `%v` to equal `%v`\n%v\n%v"
	greaterThan = "%v\n%v\nAssertion Error: Expect `%v` to be greater than `%v`\n%v\n%v"
	lessThan    = "%v\n%v\nAssertion Error: Expect `%v` to be less than `%v`\n%v\n%v"
	contain     = "%v\n%v\nAssertion Error: Expect `%s` to contain `%s`\n%v\n%v"
	astResult   = "%v \t......%v"
	fail        = "Fail!"
	pass        = "Pass!"
)

type Assertor struct {
	Record
	*log.Logger
	ErrChan chan error
}

func (a *Assertor) AssertEqualInt(title string, expect, actual int) {
	title = titlePrefix + title
	if expect != actual {
		a.Printf(astResult, title, fail)
		err := fmt.Errorf(equal, title, head, actual, expect, a.processRec(), tail)
		a.Println(err.Error())
		a.ErrChan <- err
		return
	}
	a.Printf(astResult, title, pass)
	a.ErrChan <- nil
}

func (a *Assertor) AssertEqualBool(title string, expect, actual bool) {
	title = titlePrefix + title
	if expect != actual {
		a.Printf(astResult, title, fail)
		err := fmt.Errorf(equal, title, head, actual, expect, a.processRec(), tail)
		a.Println(err.Error())
		a.ErrChan <- err
		return
	}
	a.Printf(astResult, title, pass)
	a.ErrChan <- nil

}

func (a *Assertor) AssertContainString(title string, full, sub string) {
	title = titlePrefix + title
	full = strings.TrimSuffix(full, "\n")
	if !strings.Contains(full, sub) {
		a.Printf(astResult, title, fail)
		err := fmt.Errorf(contain, title, head, full, sub, a.processRec(), tail)
		a.Println(err.Error())
		a.ErrChan <- err
		return
	}
	a.Printf(astResult, title, pass)
	a.ErrChan <- nil
}

func (a *Assertor) processRec() string {
	return fmt.Sprintf("Url: \t\t%v\nMethod: \t%v\nBody: \t\t%v\nResponse: \t%v",
		a.Record.Url, a.Record.Method, a.Record.Body, a.Record.Resp)
}

type Record struct {
	Url    string
	Method string
	Body   string
	Resp   string
}
