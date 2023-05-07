package newapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Blackoutta/profari"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/req"
)

// 分页查询制品库列表
type GetArtifactList struct {
	ProjectID int
}

func (r GetArtifactList) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("projectId", strconv.Itoa(r.ProjectID))
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("artifact/base"), q, nil)
}

// 删除制品库
type DeleteArtifact struct {
	ArtifactID int
}

func (r DeleteArtifact) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodDelete, fmt.Sprintf("artifact/base/%d", r.ArtifactID), nil, nil)
}
