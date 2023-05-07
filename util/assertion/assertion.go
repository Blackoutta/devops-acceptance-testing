package assertion

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/req"
)

type Assertion struct {
	SuitePass bool
	SuiteName string
	Logger    *log.Logger
}

func NewAssertion(suiteName string, logger *log.Logger) *Assertion {
	return &Assertion{SuitePass: true, SuiteName: suiteName, Logger: logger}
}

func (a *Assertion) FailTest(comment string) {
	a.Println("==================== ERROR INFO START ====================")
	a.Println(comment)
	a.Println("------------------ ERROR INFO END ------------------\n\n")
	a.SuitePass = false
}

func (a *Assertion) SleepWithCounter(comment string, secs int) {
	a.Printf("暂停%v秒, %v\n", secs, comment)
	go func(s int) {
		counter := s
		for counter > 0 {
			a.Printf("等待中, 还剩%v秒...", counter)
			counter--
			time.Sleep(time.Second)
		}
	}(secs)
	time.Sleep(time.Duration(secs) * time.Second)
}

func (a *Assertion) Printf(format string, v ...interface{}) {
	a.Logger.Printf(format, v...)
}

func (a *Assertion) Println(aa ...interface{}) {
	a.Logger.Println(aa...)
}

func (a *Assertion) PrintTearDownStart() {
	a.Logger.Println("****************************************** TEAR DOWN START ******************************************")
}

func (a *Assertion) PrintTearDownEnd() {
	a.Logger.Println("****************************************** TEAR DOWN END ********************************************")
}

func (a *Assertion) RecoverFromPanic() {
	if err := recover(); err != nil {
		a.Logger.Println("==================== ERROR INFO START ====================")
		a.Logger.Printf("发生Panic，程序将开始Tear Down, Panic错误为: %v", err)
		a.Logger.Printf("------------------ ERROR INFO END ------------------\n\n")
		a.SuitePass = false
	}
}

func (a *Assertion) CheckSuiteResult(exitChan chan TestResult) {
	if a.SuitePass != true {
		a.Logger.Printf("\"%s\" 测试失败，请检查上方日志来定位问题!\n", a.SuiteName)
		a.PrintSuiteFailure()
		exitChan <- TestResult{
			SuiteName: a.SuiteName,
			Result:    false,
		}
		return
	}
	a.Logger.Printf("√ \"%v\" 测试通过!  √\n", a.SuiteName)
	a.PrintSuiteSuccess()
	exitChan <- TestResult{
		SuiteName: a.SuiteName,
		Result:    true,
	}
}

func (a *Assertion) PrintSuiteSuccess() {
	sf := figure.NewFigure("TEST SUCCESS", "", true)
	sf.Print()
}

func (a *Assertion) PrintSuiteFailure() {
	sf := figure.NewFigure("TEST FAILED!!!", "", true)
	sf.Print()
}

func (a *Assertion) AssertSuccess(title string, success string, record req.Record) {
	if success != "success" {
		result := fmt.Sprintf("expect |%v| to equal |success|, got |%v|\n", success, success)
		printFailure(title, result, record, a.Logger)
		a.SuitePass = false
		return
	}
	a.Logger.Printf("%-60v............ PASS\n", title)
}

func (a *Assertion) AssertContainString(title string, actual string, expect string, record req.Record) {
	if strings.Contains(actual, expect) != true {
		result := fmt.Sprintf("expect |%v| to contain string |%v|, got |%v|\n", actual, expect, actual)
		printFailure(title, result, record, a.Logger)
		a.SuitePass = false
		return
	}
	a.Logger.Printf("%-60v............ PASS\n", title)
}

func (a *Assertion) AssertBooleanEqual(title string, actual bool, expect bool, record req.Record) {
	if actual != expect {
		result := fmt.Sprintf("expect |%v| to equal |%v|, got |%v|\n", actual, expect, actual)
		printFailure(title, result, record, a.Logger)
		a.SuitePass = false
		return
	}
	a.Logger.Printf("%-60v............ PASS\n", title)
}

func (a *Assertion) AssertStringEqual(title string, actual string, expect string, record req.Record) {
	if actual != expect {
		result := fmt.Sprintf("expect |%v| to equal |%v|, got |%v|\n", actual, expect, actual)
		printFailure(title, result, record, a.Logger)
		a.SuitePass = false
		return
	}
	a.Logger.Printf("%-60v............ PASS\n", title)
}

func (a *Assertion) AssertStringNotEqual(title string, actual string, expect string, record req.Record) {
	if actual == expect {
		result := fmt.Sprintf("expect |%v| to not equal |%v|, got |%v|\n", actual, expect, actual)
		printFailure(title, result, record, a.Logger)
		a.SuitePass = false
		return
	}
	a.Logger.Printf("%-60v............ PASS\n", title)
}

func (a *Assertion) AssertIntegerGreaterThan(title string, actual int, expect int, record req.Record) {
	if actual <= expect {
		result := fmt.Sprintf("expect |%v| to be greater than |%v|, got |%v|\n", actual, expect, actual)
		printFailure(title, result, record, a.Logger)
		a.SuitePass = false
		return
	}
	a.Logger.Printf("%-60v............ PASS\n", title)
}

func (a *Assertion) AssertIntegerEqual(title string, actual int, expect int, record req.Record) {
	if actual != expect {
		result := fmt.Sprintf("expect |%v| to equal |%v|, got |%v|\n", actual, expect, actual)
		printFailure(title, result, record, a.Logger)
		a.SuitePass = false
		return
	}
	a.Logger.Printf("%-60v............ PASS\n", title)
}

func (a *Assertion) AssertIntegerNotEqual(title string, actual int, expect int, record req.Record) {
	if actual == expect {
		result := fmt.Sprintf("expect |%v| to not equal |%v|, got |%v|\n", actual, expect, actual)
		printFailure(title, result, record, a.Logger)
		a.SuitePass = false
		return
	}
	a.Logger.Printf("%-60v............ PASS\n", title)
}

func printFailure(title string, result string, r req.Record, logger *log.Logger) {
	logger.Printf("%-60v............ FAIL\n", title)
	logger.Println("==================================================== ERROR INFO START ===========================================================")
	logger.Println("失败原因：")
	logger.Printf("Method:\t%v\n", r.Method)
	logger.Printf("URL:\t%v\n", r.URL)
	logger.Printf("Body:\t%v\n", r.Body)
	logger.Printf("Response:\t%v\n", string(r.Response))
	logger.Printf("Status Code:\t%v\n", r.StatusCode)
	logger.Println(result)
	logger.Printf("----------------------------------------------------- ERROR INFO END -------------------------------------------------------------\n\n")
}

func PrintResult(tr []TestResult) {
	fmt.Printf("\n测试结果：\n")
	fmt.Printf("编号\t\t测试套件名\t\t\t测试通过?\n")
	for i, v := range tr {
		fmt.Printf("%v\t\t\t%v\t\t\t%v\n", i+1, v.SuiteName, v.Result)
	}
}

type TestResult struct {
	SuiteName string
	Result    bool
}
