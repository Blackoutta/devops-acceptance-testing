package newapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Blackoutta/profari"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/req"
)

// 新增流水线收藏
type AddPipelineCollection struct {
	PipelineId int
}

func (r AddPipelineCollection) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPost, fmt.Sprintf("pipeline/collection/%d", r.PipelineId), nil, nil)
}

// 查询流水线收藏
type GetPipelineCollection struct {
	ProjectID int
}

func (r GetPipelineCollection) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("projectId", strconv.Itoa(r.ProjectID))
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("pipeline/collection/page"), q, nil)
}

// 删除流水线收藏
type DeletePipelineCollection struct {
	PipelineId int
}

func (r DeletePipelineCollection) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodDelete, fmt.Sprintf("pipeline/colleciton/%d", r.PipelineId), nil, nil)
}

// 分页查询流水线任务列表
type GetPipelineList struct {
	ProjectID int
}

func (r GetPipelineList) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("projectId", strconv.Itoa(r.ProjectID))
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("pipeline/base/page"), q, nil)
}

// 删除流水线
type DeletePipeline struct {
	PipelineID int
}

func (r DeletePipeline) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodDelete, fmt.Sprintf("pipeline/base/%d", r.PipelineID), nil, nil)
}
