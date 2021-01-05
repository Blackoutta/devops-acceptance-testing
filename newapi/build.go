package newapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Blackoutta/profari"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/req"
)

// 分页查询构建列表
type GetBuildList struct {
	ProjectID int
}

func (r GetBuildList) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("projectId", strconv.Itoa(r.ProjectID))
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("build/paging"), q, nil)
}

// 删除构建
type DeleteBuild struct {
	BuildID int
}

func (r DeleteBuild) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodDelete, fmt.Sprintf("build/%d", r.BuildID), nil, nil)
}

/**
SourceAuthPolicy 鉴权类型
可选值:
	: CREDENTIAL / NOAUTH
为CREDENTIAL:
	CredentialId参数必填
*/
type SourceAuthPolicy string

const (
	AuthPolicyCREDENTIAL SourceAuthPolicy = "CREDENTIAL"
	AuthPolicyNOAUTH     SourceAuthPolicy = "NOAUTH"
)

/**
SourceBranchPolicy 可触发分支策略
可选值:
	SPECIFIC_BRANCH / RUNTIME_SELECTION
为SPECIFIC_BRANCH时:
	SpecificBranch 必填 字符类型 代表推送分支
为RUNTIME_SELECTION时:
	UseBranch和UseTag 必填
*/
type SourceBranchPolicy string

const (
	BranchPolicySpecific SourceBranchPolicy = "SPECIFIC_BRANCH"
	BranchPolicyRuntime  SourceBranchPolicy = "RUNTIME_SELECTION"
)

/**
BranchUsePolicy 是否选择该类分支/tag
	USE 可选
	NOT 不可选
*/
type BranchPolicy string

const (
	BranchPolicyUse BranchPolicy = "USE"
	BranchPolicyNot BranchPolicy = "NOT"
)

type BuildSource struct {
	ProjectId      int                `json:"projectId"`
	AuthPolicy     SourceAuthPolicy   `json:"authPolicy"`
	CredentialId   int                `json:"credentialId,omitempty"`
	Type           string             `json:"type"`
	Address        string             `json:"address"`
	BranchPolicy   SourceBranchPolicy `json:"branchPolicy"`
	SpecificBranch string             `json:"specificBranch,omitempty"`
	UseBranch      BranchPolicy       `json:"useBranch,omitempty"`
	UseTag         BranchPolicy       `json:"useTag,omitempty"` // 可选值: NOT/USE
}

/**
BuildType 构建类型
可选值:
	MVN/GOLANG/NPM/DEFAULT_STEP
*/
type BuildType string

const (
	BuildTypeMvn    BuildType = "MVN"
	BuildTypeGolang BuildType = "GOLANG"
	BuildTypeNpm    BuildType = "NPM"
	BuildTypeOther  BuildType = "DEFAULT_STEP"
)

type CreateBuild struct {
	Name      string      `json:"name"`
	ProjectId int         `json:"projectId"`
	BuildType BuildType   `json:"buildType"`
	Source    BuildSource `json:"source"`
}

func (t *CreateBuild) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPost, "build/", nil, t)
}

type AddBuildCollection struct {
	BuildId int
}

func (t *AddBuildCollection) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPost, "build/collection/"+strconv.Itoa(t.BuildId), nil, `"{"data":"bodyNotNull"}"`)
}

type DeleteBuildCollection struct {
	BuildId int
}

func (t *DeleteBuildCollection) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodDelete, "build/collection/"+strconv.Itoa(t.BuildId), nil, nil)
}

type GetBuildCollectionList struct {
	ProjectId int
}

func (t *GetBuildCollectionList) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("projectId", strconv.Itoa(t.ProjectId))
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("build/collection/page"), q, nil)
}

type ImageType string

const (
	GolangImage ImageType = "GOLANG"
	MVNImage    ImageType = "MVN"
	NpmImage    ImageType = "NPM"
	OtherImage  ImageType = "DEFAULT_STEP"
)

// 镜像类型和构建类型一样
type GetImageList struct {
	Type     ImageType
	PageSize int
}

func (t *GetImageList) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", strconv.Itoa(t.PageSize))
	q.Add("type", string(t.Type))
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("build/images"), q, nil)
}

type GenHookToken struct {
}

func (t *GenHookToken) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("build/genToken"), nil, nil)
}
