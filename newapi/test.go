package newapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Blackoutta/profari"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/req"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/conf"
)

type CreateTest struct {
	CredentialId  int    `json:"credentialId,omitempty"`
	EnvId         int    `json:"envId,omitempty"`
	Image         string `json:"image,omitempty"`
	ImageAuthType string `json:"imageAuthType,omitempty"`
	Name          string `json:"name,omitempty"`
	ProjectId     int    `json:"projectId,omitempty"`
	Timeout       int    `json:"timeout"`
}

func (r CreateTest) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPost, "test/create", nil, r)
}

type EditTestConfig struct {
	TestID       int
	ContainerLog ContainerLog   `json:"containerLog,omitempty"`
	HostMap      HostMap        `json:"hostMap,omitempty"`
	KubeConfig   TestKubeConfig `json:"kubeConfig,omitempty"`
	VarMap       VarMap         `json:"varMap,omitempty"`
}

type ContainerLog struct {
	LogPath   string `json:"logPath"`
	IsChanged bool   `json:"isChanged"`
}

type HostMap struct {
	Host string `json:"devops.test.cq.iot.chinamobile.com,omitempty"`
}

type TestKubeConfig struct {
	Cpu       string `json:"cpu,omitempty"`
	Mem       string `json:"mem,omitempty"`
	Id        int    `json:"id,omitempty"`
	IsChanged bool   `json:"isChanged,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	ObjectId  int    `json:"objectId,omitempty"`
	Duplicate int    `json:"duplicate,omitempty"`
}

type VarMap struct {
	CONFIGFILE string `json:"CONFIGFILE,omitempty"`
	SUITE      string `json:"SUITE,omitempty"`
}

func (r EditTestConfig) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPut, fmt.Sprintf("test/%v/k8sConfig", r.TestID), nil, r)
}

type ExecuteTest struct {
	ParameterList []Parameter `json:"parameterList"`
	TestId        int         `json:"testId,omitempty"`
}

type Parameter struct {
	ActualValue string
	Name        string
}

func (r ExecuteTest) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPost, fmt.Sprintf("test/run/%v", r.TestId), nil, r)
}

type GetTestDetail struct {
	JobID int
}

func (r GetTestDetail) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("test/history/%v", r.JobID), nil, nil)
}

type GetTestSystemLog struct {
	TestID    int
	LastLogId string
}

func (r GetTestSystemLog) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("lastLogId", r.LastLogId)
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("test/history/%v/system/log", r.TestID), q, nil)
}

type DeleteTest struct {
	TestID int
}

func (r DeleteTest) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodDelete, fmt.Sprintf("test/%v/delete", r.TestID), nil, nil)
}

// 分页查询测试任务列表
type GetTestList struct {
	ProjectID int
}

func (r GetTestList) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("projectId", strconv.Itoa(r.ProjectID))
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("test/simpleList"), q, nil)
}

// 根据测试ID查询测试详情
type GetTestDetailByID struct {
	TestID int
}

func (r GetTestDetailByID) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("testId", strconv.Itoa(r.TestID))
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("test/%d/details", r.TestID), q, nil)
}

// 获取测试的k8s配置
type GetTestK8SConfByID struct {
	TestID int
}

func (r GetTestK8SConfByID) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("testId", strconv.Itoa(r.TestID))
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("test/%d/k8sConfig", r.TestID), q, nil)
}

// 查询执行测试之前的参数配置
type TestBeforeRun struct {
	TestID int
}

func (r TestBeforeRun) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("testId", strconv.Itoa(r.TestID))
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("test/beforeRun/%d", r.TestID), q, nil)
}

// 编辑测试基本配置
type EditTestBase struct {
	CredentialId     int            `json:"credentialId"`
	EnvId            int            `json:"envId"`
	GroupAuth        []GroupAuth    `json:"groupAuth"`
	Id               int            `json:"id"`
	Image            string         `json:"image"`
	Name             string         `json:"name"`
	ParameterReqList []ParameterReq `json:"parameterReqList"`
	ProjectId        int            `json:"projectId"`
	Timeout          int            `json:"timeout"`
}

const (
	REQUIRED     = "TRUE"
	NOT_REQUIRED = "FALSE"
	STRING_PARAM = "STRING"
	ENUM_PARAM   = "ENUM"
)

type ParameterReq struct {
	DefaultValue string `json:"defaultValue"`
	Description  string `json:"description"`
	Id           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	ParamValues  string `json:"paramValues"`
	Required     string `json:"required"`
	Type         string `json:"type"`
}

type GroupAuth struct {
	Auth     bool `json:"auth"`
	Copy     bool `json:"copy"`
	Delete   bool `json:"delete"`
	Edit     bool `json:"edit"`
	GroupId  int  `json:"groupId"`
	Readonly bool `json:"readonly"`
	Run      bool `json:"run"`
}

func (r EditTestBase) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPut, fmt.Sprintf("test/update"), nil, r)
}

// 根据testIds查询测试的最新状态
type GetTestStatus struct {
	IDs string
}

func (r GetTestStatus) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("ids", r.IDs)
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("test/statusList"), q, nil)
}

// 查询执行记录列表
type GetTestRecords struct {
	TestID int
}

func (r GetTestRecords) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("testId", strconv.Itoa(r.TestID))
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("test/history/page/%d", r.TestID), q, nil)
}

// 查询执行记录的pod信息
type GetTestPods struct {
	TestID    int
	HistoryID int
}

func (r GetTestPods) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("testId", strconv.Itoa(r.TestID))
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("test/history/%d/pods", r.HistoryID), q, nil)
}

// 下载Pod日志
type DownloadTestLog struct {
	HistoryID     int    `json:"historyID"`
	EnvId         int    `json:"envId"`
	Namespace     string `json:"namespace"`
	PodName       string `json:"podName"`
	ContainerName string `json:"containerName"`
	SinceSeconds  int    `json:"sinceSeconds"`
	TailLine      int    `json:"tailLine"`
	Timestamps    bool   `json:"timestamps"`
	DeployID      int
}

func (r DownloadTestLog) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("historyId", strconv.Itoa(r.HistoryID))
	q.Add("deployId", strconv.Itoa(r.DeployID))
	q.Add("envId", strconv.Itoa(r.EnvId))
	q.Add("namespace", r.Namespace)
	q.Add("podName", r.PodName)
	q.Add("containerName", r.ContainerName)
	// q.Add("sinceSeconds", strconv.Itoa(r.SinceSeconds))
	// q.Add("tailLine", strconv.Itoa(r.TailLine))
	// q.Add("timestamps", fmt.Sprintf("%v", r.Timestamps))
	q.Add("userId", conf.UserID)
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("test/history/%d/download/log", r.HistoryID), q, nil)
}

type StopTest struct {
	TestID    int
	HistoryId int `json:"historyId"`
}

func (r StopTest) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPut, fmt.Sprintf("test/stop/%v/", r.TestID), nil, r)
}
