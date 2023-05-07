package newapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Blackoutta/profari"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/req"
)

// 查看App收藏
type GetAppCollection struct {
	Name      string
	ProjectId int
}

func (r GetAppCollection) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("name", r.Name)
	q.Add("projectId", strconv.Itoa(r.ProjectId))
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("core/collection/page"), q, nil)
}

// 删除App收藏
type DeleteAppCollection struct {
	AppId int
}

func (r DeleteAppCollection) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodDelete, fmt.Sprintf("core/collection/%d", r.AppId), nil, nil)
}

// 新增App收藏
type AddAppCollection struct {
	AppId int
}

func (r AddAppCollection) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPost, fmt.Sprintf("core/collection/%d", r.AppId), nil, r)
}

// 删除项目收藏
type DeleteProjectCollection struct {
	ProjectID int
}

func (r DeleteProjectCollection) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodDelete, fmt.Sprintf("core/projectCollection/%d", r.ProjectID), nil, nil)
}

// 查看项目收藏
type GetProjectCollection struct {
	Query string
}

func (r GetProjectCollection) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("Query", r.Query)
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("core/projectCollection/query"), q, nil)
}

// 新增项目收藏
type AddProjectCollection struct {
	ProjectId int `json:"projectId"`
}

func (r AddProjectCollection) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPost, fmt.Sprintf("core/projectCollection/save"), nil, r)
}

// 按ID查看项目信息
type GetProjectDetailByID struct {
	ProjectID int
}

func (r GetProjectDetailByID) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("core/projects/%d", r.ProjectID), nil, nil)
}

// 按Id更新项目
type UpdateProjectById struct {
	ProjectID    int
	DepartmentId int    `json:"departmentId"`
	Description  string `json:"description"`
	Forbid       bool   `json:"forbid"`
	JiraId       int    `json:"jiraId,omitempty"`
	Name         string `json:"name"`
}

func (r UpdateProjectById) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPatch, fmt.Sprintf("core/projects/%d", r.ProjectID), nil, r)
}

// k8s环境连通性检查
type K8SConnectionCheck struct {
	Config  json.RawMessage `json:"config"`
	Kubectx string          `json:"kubectx"`
}

func (r K8SConnectionCheck) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPut, fmt.Sprintf("core/env/connection/check"), nil, r)
}

// 查询机房及可用区
type GetRoomAndZone struct {
	GroupID int
}

func (r GetRoomAndZone) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("core/env/vm/room"), nil, nil)
}

// 查询主机组
type GetVMGroup struct {
	GroupID int
}

func (r GetVMGroup) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("core/env/vm/group/%v", r.GroupID), nil, nil)
}

// 查询工单新增或编辑时meta信息
type GetVMOrderMeta struct {
}

func (r GetVMOrderMeta) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodGet, "core/workOrder/meta", nil, nil)
}

// 查询工单记录列表
type GetVMOrderList struct {
	GroupID int
}

func (r GetVMOrderList) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("groupId", strconv.Itoa(r.GroupID))

	return req.ComposeNewProfariRequest(http.MethodGet, "core/workOrder/simpleList", q, nil)
}

// 创建工单
type CreateVMWorkOrder struct {
	GroupId        int            `json:"groupId"`
	ProjectName    string         `json:"projectName"`
	TeamName       string         `json:"teamName"`
	ComputerConfig ComputerConfig `json:"computerConfig"`
}

type ComputerConfig struct {
	ApplyAmount  int        `json:"applyAmount"`
	CardAmount   string     `json:"cardAmount"`
	CardUse      string     `json:"cardUse"`
	ComputeType  string     `json:"computeType"`
	CpuAmount    string     `json:"cpuAmount"`
	DiskVapacity int        `json:"diskVapacity"`
	DurationDnd  string     `json:"durationDnd"`
	Memory       string     `json:"memory"`
	Os           string     `json:"os"`
	OsVersion    string     `json:"osVersion"`
	Region       string     `json:"region"`
	Remark       string     `json:"remark"`
	Segment      string     `json:"segment"`
	UsageWo      string     `json:"usageWo"`
	ListComputer []Computer `json:"listComputer,omitempty"`
}

type Computer struct {
	Id            int    `json:"id"`
	Ip            string `json:"ip"`
	LoginName     string `json:"loginName"`
	LoginPassword string `json:"loginPassword"`
	ProcessStatus string `json:"processStatus"`
	Segment       string `json:"segment"`
	VlanId        string `json:"vlanId"`
}

func (r CreateVMWorkOrder) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPost, "core/workOrder", nil, r)
}

// 删除主机组
type DeleteVmGroup struct {
	GroupID int
}

func (r DeleteVmGroup) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodDelete, fmt.Sprintf("core/env/vm/group/%v", r.GroupID), nil, nil)
}

//删除主机
type DeleteVmMachine struct {
	VmMachineId int
}

func (r DeleteVmMachine) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodDelete, fmt.Sprintf("core/env/vm/machine/%v", r.VmMachineId), nil, nil)
}

// 创建主机组
type CreateVMGroup struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	ProjectId   int    `json:"projectId"`
	RoomId      int    `json:"roomId"`
	ZoneId      int    `json:"zoneId"`
}

func (r CreateVMGroup) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPost, "core/env/vm/group", nil, r)
}

type CreateVMMachine struct {
	AuthType    string `json:"authType"`
	Description string `json:"description"`
	Ip          string `json:"ip"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Port        int    `json:"port"`
	SshKey      string `json:"sshKey"`
	UserName    string `json:"userName"`
	VmGroupId   int    `json:"vmGroupId"`
}

func (r CreateVMMachine) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPost, "core/env/vm/machine", nil, r)
}

type CreateProject struct {
	DepartmentId int    `json:"departmentId,omitempty"`
	Description  string `json:"description,omitempty"`
	Forbid       bool   `json:"forbid"`
	Identify     string `json:"identify,omitempty"`
	JiraId       int    `json:"jiraId,omitempty"`
	Name         string `json:"name,omitempty"`
	UseDefault   bool   `json:"useDefault"`
}

func (r CreateProject) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPost, "core/project", nil, r)
}

type GetProject struct {
	ProjectName string
}

func (r GetProject) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	q.Add("query", r.ProjectName)

	return req.ComposeNewProfariRequest(http.MethodGet, "core/projects", q, nil)
}

type GetProjectList struct {
	PageNum int
}

func (r GetProjectList) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", strconv.Itoa(r.PageNum))
	q.Add("pageSize", "100")
	return req.ComposeNewProfariRequest(http.MethodGet, "core/projects", q, nil)
}

// git和docker 凭证只有type不一样,其他参数一模一样, 新增的ssh凭证
type CreateGitCredential struct {
	CreateDockerCredential
}

type CreateDockerCredential struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	Password    string `json:"password,omitempty"`
	ProjectId   int    `json:"projectId,omitempty"`
	Type        string `json:"type,omitempty"`
	UserName    string `json:"userName,omitempty"`
}

func (r CreateDockerCredential) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPost, "core/credential", nil, r)
}

type GetCredential struct {
	Type      string
	ProjectId int
	Name      string
}

func (r GetCredential) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	q.Add("type", r.Type)
	q.Add("projectId", strconv.Itoa(r.ProjectId))
	q.Add("name", r.Name)
	q.Add("keyword", r.Name)
	return req.ComposeNewProfariRequest(http.MethodGet, "core/credential", q, nil)
}

type DeleteCredential struct {
	Id int
}

func (r *DeleteCredential) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodDelete, "core/credential"+strconv.Itoa(r.Id), nil, nil)
}

type CreateK8sEnv struct {
	Config     json.RawMessage `json:"config,omitempty"`
	Kubectx    string          `json:"kubectx,omitempty"`
	Name       string          `json:"name,omitempty"`
	Namespace  string          `json:"namespace,omitempty"`
	ServerAddr string          `json:"serverAddr,omitempty"`
	Type       string          `json:"type,omitempty"`
	ProjectId  int             `json:"projectId,omitempty"`
	RoomId     int             `json:"roomId"`
	ZoneId     int             `json:"zoneId"`
}

func (r CreateK8sEnv) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPost, "core/env", nil, r)
}

type GetEnvs struct {
	ProjectId int
}

func (r GetEnvs) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	q.Add("projectId", strconv.Itoa(r.ProjectId))
	return req.ComposeNewProfariRequest(http.MethodGet, "core/envs", q, nil)
}

type DeleteK8sEnv struct {
	EnvID int
}

func (r DeleteK8sEnv) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodDelete, fmt.Sprintf("core/env/%v", r.EnvID), nil, nil)
}

type DeleteVMEnv struct {
	EnvID int
}

func (r DeleteVMEnv) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodDelete, fmt.Sprintf("core/env/vm/group/%d", r.EnvID), nil, nil)
}

type DeleteProject struct {
	ProjectID int
}

func (r DeleteProject) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodDelete, fmt.Sprintf("core/projects/%v", r.ProjectID), nil, nil)
}

type CreateApp struct {
	AppManager  string `json:"appManager,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	ProjectId   int    `json:"projectId,omitempty"`
}

func (r CreateApp) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPost, "core/app", nil, r)
}

type GetApps struct {
	ProjectId int
	Name      string
}

func (r GetApps) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	q.Add("projectId", strconv.Itoa(r.ProjectId))
	q.Add("name", r.Name)
	return req.ComposeNewProfariRequest(http.MethodGet, "core/apps", q, nil)
}

type DeleteApp struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (r DeleteApp) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodDelete, fmt.Sprintf("core/app"), nil, r)
}

// 分页查询应用列表
type GetAppList struct {
	ProjectID int
}

func (r GetAppList) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("projectId", strconv.Itoa(r.ProjectID))
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("core/apps"), q, nil)
}

// 分页查询k8s环境列表
type GetK8SEnvList struct {
	ProjectID int
}

func (r GetK8SEnvList) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("projectId", strconv.Itoa(r.ProjectID))
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("core/envs"), q, nil)
}

// 分页查询主机组环境列表
type GetVMEnvList struct {
	ProjectID int
}

func (r GetVMEnvList) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("projectId", strconv.Itoa(r.ProjectID))
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("core/env/vm/groups"), q, nil)
}

type GetVMGroups struct {
	ProjectID int
}

func (r GetVMGroups) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	q.Add("projectId", strconv.Itoa(r.ProjectID))
	return req.ComposeNewProfariRequest(http.MethodGet, "core/env/vm/groups", q, nil)
}

type GetVMMachines struct {
	VmGroupId int
	Query     string
}

func (r GetVMMachines) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	q.Add("vmGroupId", strconv.Itoa(r.VmGroupId))
	q.Add("query", r.Query)
	return req.ComposeNewProfariRequest(http.MethodGet, "core/env/vm/machines", q, nil)
}

type UpdateVM struct {
	VmMachineId int
	AuthType    string `json:"authType"`
	Description string `json:"description"`
	Ip          string `json:"ip"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Port        int    `json:"port"`
	SshKey      string `json:"sshKey"`
	UserName    string `json:"userName"`
}

func (r UpdateVM) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPatch, fmt.Sprintf("core/env/vm/machine/%d", r.VmMachineId), nil, r)
}
