package newapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Blackoutta/profari"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/req"
)

// 获取K8S描述文件
type GetK8SDeploymentFile struct {
	DeployID int
	AppID    int
	EnvID    int
}

func (r GetK8SDeploymentFile) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("deployId", strconv.Itoa(r.DeployID))
	q.Add("appId", strconv.Itoa(r.AppID))
	q.Add("envId", strconv.Itoa(r.EnvID))
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("k8s-deploy/pod/deployResource"), q, nil)
}

// 获取pod describe文件
type GetPodDescribeInfo struct {
	DeployID int
	AppID    int
	EnvID    int
	PodName  string
}

func (r GetPodDescribeInfo) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("deployId", strconv.Itoa(r.DeployID))
	q.Add("appId", strconv.Itoa(r.AppID))
	q.Add("envId", strconv.Itoa(r.EnvID))
	q.Add("podName", r.PodName)
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("k8s-deploy/pod/info"), q, nil)
}
