package API

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/errors"

	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/req"
)

func DeleteVMGroup(c http.Client, vmGroupID int) (GeneralResp, req.Record) {
	path := fmt.Sprintf(`core/env/vm/group/%v`, vmGroupID)
	r := req.ComposeNewRequest(http.MethodDelete, path, nil, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	gen := GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	return gen, resp
}

func EditVMDeployToTeardown(c http.Client, deployID int) (GeneralResp, req.Record) {
	j := fmt.Sprintf(`{
    "deployId": %v,
    "steps": [{
        "body": "{\"dir\":\"/home/app/hyy-test-go-server\"}",
        "id": 99999999,
        "name": "初始化",
        "type": "INIT"
    }, {
        "body": "{\"dir\":\"/home/app/hyy-test-go-server\",\"cmd\":\"sh appctl.sh stop && rm -rf *\",\"timeout\":60}",
        "id": 99999999,
        "name": "执行命令",
        "type": "RUNNING"
    }]
}`, deployID)
	body := strings.NewReader(j)
	path := fmt.Sprintf("deploy/vm/config")
	r := req.ComposeNewRequest(http.MethodPost, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	gen := GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	return gen, resp
}

func EditVMDeploy(c http.Client, deployID int) (GeneralResp, req.Record) {
	j := fmt.Sprintf(`{
    "deployId": %v,
    "steps": [{
        "body": "{\"dir\":\"/home/app/hyy-test-go-server\"}",
        "id": 99999,
        "name": "初始化",
        "type": "INIT"
    }, {
        "body": "{\"dir\":\"/home/app/hyy-test-go-server\",\"cmd\":\"sh ./__artifact/appctl.sh stop\",\"timeout\":60}",
        "id": 99999,
        "name": "停止服务",
        "type": "STOP"
    }, {
        "body": "{\"dir\":\"/home/app/hyy-test-go-server\",\"cmd\":\"cp __artifact/* . && rm -rf artifact && rm -rf __artifact\",\"timeout\":30}",
        "id": 99999,
        "name": "清理制品",
        "type": "RUNNING"
    }, {
        "body": "{\"dir\":\"/home/app/hyy-test-go-server\",\"cmd\":\"sh appctl.sh start\",\"timeout\":60}",
        "id": 99999,
        "name": "启动服务",
        "type": "START"
    }]
}`, deployID)
	body := strings.NewReader(j)
	path := fmt.Sprintf("deploy/vm/config")
	r := req.ComposeNewRequest(http.MethodPost, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	gen := GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	return gen, resp
}

func CreateVMDeploy(c http.Client, appID int, name string, projectID int, artifactID int, vmGroupiD int) (GeneralResp, req.Record) {
	j := fmt.Sprintf(`{
    "name": "%v",
    "k8sTimeout": 300,
    "appId": %v,
    "envId": null,
    "vmGroupId": %v,
    "deployType": "VMWARE_DEPLOY",
    "description": null,
    "sourceCreateRequest": {
        "imageAddr": "",
        "imageAuthType": 0,
        "sourceType": "ARTIFACT",
        "credentialId": null,
        "artifactId": %v,
        "artifactVersion": ""
    },
    "projectId": %v
}`, name, appID, vmGroupiD, artifactID, projectID)
	body := strings.NewReader(j)
	path := fmt.Sprintf("deploy/base")
	r := req.ComposeNewRequest(http.MethodPost, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	gen := GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	return gen, resp
}

func CheckVMStatus(c http.Client, vmGroupID int) (GeneralResp, req.Record) {
	path := fmt.Sprintf(`core/env/vm/status/%v`, vmGroupID)
	q := url.Values{}
	r := req.ComposeNewRequest(http.MethodGet, path, q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	gen := GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	return gen, resp
}

func CreateVM(c http.Client, IP string, port int, username string, password string, vmGroupID int, authType string) (GeneralResp, req.Record) {
	j := fmt.Sprintf(`{
    "description": "this is a vm machine",
    "name": "%v",
    "ip": "%v",
    "port": "%v",
    "userName": "%v",
    "password": "%v",
    "vmGroupId": %v,
    "authType": "%v"
}`, IP, IP, port, username, password, vmGroupID, authType)
	body := strings.NewReader(j)
	path := fmt.Sprintf("core/env/vm/machine")
	r := req.ComposeNewRequest(http.MethodPost, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	gen := GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	return gen, resp
}

func CreateVMGroup(c http.Client, projectID int, name string) (GeneralResp, req.Record) {
	j := fmt.Sprintf(`{"description":"this is a vm group","name":"%v","projectId":%v, "roomId": 2, "zoneId": 4}`, name, projectID)
	body := strings.NewReader(j)
	path := fmt.Sprintf("core/env/vm/group")
	r := req.ComposeNewRequest(http.MethodPost, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	gen := GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	return gen, resp
}

func DoCheckpoint(c http.Client, historyID int, status string) req.Record {
	path := fmt.Sprintf(`pipeline/history/checkpoint`)
	j := fmt.Sprintf(`{"extraInfo":"check me please","historyId":%v,"status":"%v"}`, historyID, status)
	body := strings.NewReader(j)
	r := req.ComposeNewRequest(http.MethodPatch, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetPrometheusData(c http.Client, envID int, nodeName, podName string, start, end, step int64) req.Record {
	path := fmt.Sprintf(`k8s-deploy/pod/prometheus`)
	q := url.Values{}
	q.Set("envId", strconv.Itoa(envID))
	q.Set("nodeName", nodeName)
	q.Set("podName", podName)
	q.Set("start", strconv.Itoa(int(start)))
	q.Set("end", strconv.Itoa(int(end)))
	q.Set("step", strconv.Itoa(int(step)))

	r := req.ComposeNewRequest(http.MethodGet, path, q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetBuildRuntime(c http.Client, historyID int) req.Record {
	path := fmt.Sprintf(`build/history/runtime`)
	q := url.Values{}
	q.Set("history_id", strconv.Itoa(historyID))

	r := req.ComposeNewRequest(http.MethodGet, path, q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetPipelineRuntime(c http.Client, historyID int) req.Record {
	path := fmt.Sprintf(`pipeline/history/runtime`)
	q := url.Values{}
	q.Set("historyId", strconv.Itoa(historyID))

	r := req.ComposeNewRequest(http.MethodGet, path, q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetK8sDescription(c http.Client, deployID, appID, envID int) req.Record {
	path := fmt.Sprintf(`k8s-deploy/pod/deployResource`)
	q := url.Values{}
	q.Set("envId", strconv.Itoa(envID))
	q.Set("deployId", strconv.Itoa(deployID))
	q.Set("appId", strconv.Itoa(appID))

	r := req.ComposeNewRequest(http.MethodGet, path, q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func ChangePodResource(c http.Client, deployID, appID, envID int, cpu, mem string) req.Record {
	path := fmt.Sprintf(`k8s-deploy/pod/resource`)
	j := fmt.Sprintf(`{"cpu":"%v","mem":"%v","deployId":%v,"appId":%v,"envId":%v}`, cpu, mem, deployID, appID, envID)
	body := strings.NewReader(j)
	r := req.ComposeNewRequest(http.MethodPatch, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func DownloadPodLog(c http.Client, envID, deployID int, namespace, podName, containerName, userID string) req.Record {
	path := fmt.Sprintf(`download/k8s-deploy/pod/logs`)
	q := url.Values{}
	q.Set("envId", strconv.Itoa(envID))
	q.Set("namespace", namespace)
	q.Set("podName", podName)
	q.Set("containerName", containerName)
	q.Set("deployId", strconv.Itoa(deployID))
	q.Set("deployId", strconv.Itoa(deployID))
	q.Set("userId", userID)
	r := req.ComposeNewRequest(http.MethodGet, path, q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func RescalePods(c http.Client, replicaSets int, deployID int, appID int, envID int) req.Record {
	path := fmt.Sprintf(`k8s-deploy/pod/rescaling`)
	j := fmt.Sprintf(`{"duplicate":%v,"deployId":%v,"appId":%v,"envId":%v}`, replicaSets, deployID, appID, envID)
	body := strings.NewReader(j)
	r := req.ComposeNewRequest(http.MethodPatch, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetPods(c http.Client, deployID int, envID int, appID int) req.Record {
	path := fmt.Sprintf(`deploy/pod/%v/pods`, appID)
	q := url.Values{}
	q.Set("appId", strconv.Itoa(appID))
	q.Set("deployEnvIds", fmt.Sprintf(`%v:%v`, deployID, envID))
	r := req.ComposeNewRequest(http.MethodGet, path, q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func RebootPod(c http.Client, envID, deployID int, namespace string, podName string) req.Record {
	path := fmt.Sprintf(`k8s-deploy/pod/`)
	q := url.Values{}
	q.Set("envId", strconv.Itoa(envID))
	q.Set("namespace", namespace)
	q.Set("podName", podName)
	q.Set("deployId", strconv.Itoa(deployID))
	r := req.ComposeNewRequest(http.MethodDelete, path, q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func UploadArtifact(c http.Client, version string, artifactID int, fname string, contentType string) req.Record {
	path := fmt.Sprintf(`form/artifact/upload`)
	q := url.Values{}
	q.Set("version", version)
	q.Set("artifactId", strconv.Itoa(artifactID))
	r := req.ComposeNewMultipartRequest(http.MethodPost, path, q, fname, contentType)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

// GetBuildHistoryList使用分页查询来获取构建记录列表
func GetBuildHistorylist(c http.Client, buildID int) req.Record {
	path := fmt.Sprintf(`build/history`)
	q := url.Values{}
	q.Set("pageNum", "1")
	q.Set("pageSize", "10")
	q.Set("buildId", strconv.Itoa(buildID))
	r := req.ComposeNewRequest(http.MethodGet, path, q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

// EditArtifact 编辑制品库信息，可修改名称、是否鉴权、鉴权Token
func EditArtifact(c http.Client, artifactID int, name, isAuth, token string, forbid bool) req.Record {
	path := fmt.Sprintf("artifact/base/%v", artifactID)
	j := fmt.Sprintf(`{"artifactName":"%v","isAuth":"%v","token":"%v","forbid":%v}`, name, isAuth, token, forbid)
	body := strings.NewReader(j)
	r := req.ComposeNewRequest(http.MethodPatch, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

// DownloadArtifact会从制品库下载制品。不指定version时默认下载最新版本制品。token用于制品鉴权。
func DownloadArtifact(c http.Client, downloadID string, token string, version string) req.Record {
	path := fmt.Sprintf(`download/artifact/download/%v`, downloadID)
	q := url.Values{}
	if version != "" {
		q.Set("version", version)
	}
	r := req.ComposeNewRequest(http.MethodGet, path, q, nil)
	if token != "" {
		r.Header.Set("Authorization", fmt.Sprintf("Basic %v", token))
	} else {
		r.Header.Del("Authorization")
	}
	r.Header.Del("X-UserId")
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

// GetDeployHistoryList使用分页查询来获取部署记录列表
func GetDeployHistoryList(c http.Client, deployID int) req.Record {
	path := fmt.Sprintf(`deploy/history/%v/page`, deployID)
	q := url.Values{}
	q.Add("deployID", strconv.Itoa(deployID))
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	r := req.ComposeNewRequest(http.MethodGet, path, q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetPipelineHistory(c http.Client, pipelineJobID int) req.Record {
	path := fmt.Sprintf(`pipeline/history/%v`, pipelineJobID)
	r := req.ComposeNewRequest(http.MethodGet, path, nil, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func DeletePipeline(c http.Client, pipelineID int) req.Record {
	path := fmt.Sprintf(`pipeline/base/%v`, pipelineID)
	r := req.ComposeNewRequest(http.MethodDelete, path, nil, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func RunPipeline(c http.Client, pipelineID int, buildID int, commitID string, unitTestID int, branchName string) req.Record {
	j := fmt.Sprintf(`{
    "buildRunningRequestList": [{
        "buildId": %v,
        "branchOrTag": {
            "name": "%v",
            "type": "BRANCH",
            "commitId": "%v"
        }
    }],
	"unitTestRunningRequestList": [{
			"id": %v,
			"branchOrTag": {
				"name": "%v",
				"type": "BRANCH",
				"commitId": "%v"
			}
		}]
}`, buildID, branchName, commitID, unitTestID, branchName, commitID)
	body := strings.NewReader(j)
	path := fmt.Sprintf("pipeline/base/%v/runningPipeline", pipelineID)
	r := req.ComposeNewRequest(http.MethodPost, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetPipelineDetail(c http.Client, pipelineID int) req.Record {
	path := fmt.Sprintf(`pipeline/base/%v`, pipelineID)
	r := req.ComposeNewRequest(http.MethodGet, path, nil, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func EditPipeline(c http.Client, pipelineID int, projectID int, name string, buildID int, deployID int, triggeriD int) (GeneralResp, req.Record) {
	j := fmt.Sprintf(`{
    "id": %v,
    "projectId": %v,
    "name": "%v",
    "description": "this is a pipeline",
    "jira": false,
    "pipelineSteps": [],
    "parameters": [],
    "triggerMode": {
        "id": %v,
        "createdAt": 1581167134,
        "updatedAt": 1581167134,
        "deletedAt": 0,
        "projectId": %v,
        "pipelineId": %v,
        "type": "MANUAL"
    },
    "steps": [
	{
        "body": "{\"objectId\":%v,\"paramPolicy\":\"DEFAULT_VALUE\"}",
        "type": "BUILD",
        "name": "构建任务"
    }, {
        "body": "{\"objectId\":%v,\"paramPolicy\":\"DEFAULT_VALUE\"}",
        "type": "DEPLOY",
        "name": "部署任务"
	}],
	"groupAuth": [{
        "groupId": 0,
        "groupName": "创建人",
        "objectId": 766,
        "moduleType": "PIPELINE",
        "readonly": true,
        "run": true,
        "edit": true,
        "delete": true,
        "auth": true,
        "copy": true,
        "forbid": true
	}],
	"plan": {
        "daily": false,
        "dailyPlan": "00:00",
        "deletedAt": 0,
        "id": null,
        "manual": true,
        "weekly": false,
        "weeklyDay": "",
        "weeklyPlan": "00:00"
    }
}`, pipelineID, projectID, name, triggeriD, projectID, pipelineID, buildID, deployID)
	body := strings.NewReader(j)
	path := "pipeline/base"
	r := req.ComposeNewRequest(http.MethodPut, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	gen := GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	return gen, resp
}

func EditPipelineWithUnitTestAndCheckpoint(c http.Client, pipelineID int, projectID int, name string, buildID int, deployID int, triggeriD int, gitCredentialID int, userID string) req.Record {
	j := fmt.Sprintf(`{
    "id": %v,
    "projectId": %v,
    "name": "%v",
    "description": "this is a pipeline",
    "jira": false,
    "pipelineSteps": [],
    "parameters": [],
    "triggerMode": {
        "id": %v,
        "createdAt": 1581167134,
        "updatedAt": 1581167134,
        "deletedAt": 0,
        "projectId": %v,
        "pipelineId": %v,
        "type": "MANUAL"
    },
    "steps": [{
        "body": "{\"address\":\"http://gitlab.onenet.com/huyangyi/devops-test-httpserver.git\",\"credentialId\":%v,\"authPolicy\":\"CREDENTIAL\",\"specificBranch\":\"smoke_build\",\"branchPolicy\":\"SPECIFIC_BRANCH\",\"useTag\":\"NOT\",\"useBranch\":\"NOT\",\"language\":\"GOLANG\",\"image\":\"hub.iot.chinamobile.com/library/golang:1.12.5\",\"cmd\":\"cd uploadme && go test -v\",\"coverage\":0,\"timeout\":\"180\"}",
        "type": "UNIT_TEST",
        "name": "单元测试"
    },
	{
        "body": "{\"assignee\":\"%v\",\"emailPolicy\":\"FALSE\",\"explain\":\"check me\"}",
        "type": "CHECKPOINT",
        "name": "卡点操作"
    },
	{
        "body": "{\"objectId\":%v,\"paramPolicy\":\"DEFAULT_VALUE\"}",
        "type": "BUILD",
        "name": "构建任务"
    }, {
        "body": "{\"objectId\":%v,\"paramPolicy\":\"DEFAULT_VALUE\"}",
        "type": "DEPLOY",
        "name": "部署任务"
	}],
	"groupAuth": [{
        "groupId": 0,
        "groupName": "创建人",
        "objectId": 766,
        "moduleType": "PIPELINE",
        "readonly": true,
        "run": true,
        "edit": true,
        "delete": true,
        "auth": true,
        "copy": true,
        "forbid": true
	}],
	"plan": {
        "daily": false,
        "dailyPlan": "00:00",
        "deletedAt": 0,
        "id": null,
        "manual": true,
        "weekly": false,
        "weeklyDay": "",
        "weeklyPlan": "00:00"
    }
}`, pipelineID, projectID, name, triggeriD, projectID, pipelineID, gitCredentialID, userID, buildID, deployID)
	body := strings.NewReader(j)
	path := "pipeline/base"
	r := req.ComposeNewRequest(http.MethodPut, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func CreatePipeline(c http.Client, projectID int, name string) req.Record {
	j := fmt.Sprintf(`{"description":"this is a pipeline","jira":false,"name":"%v","projectId":%v}`, name, projectID)
	body := strings.NewReader(j)
	path := "pipeline/base"
	r := req.ComposeNewRequest(http.MethodPost, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func DuplicatePipeline(c http.Client, pipelineID int, name string) req.Record {
	j := fmt.Sprintf(`{"name": "%s"}`, name)
	body := strings.NewReader(j)
	path := fmt.Sprintf("pipeline/base/copy/%d", pipelineID)
	r := req.ComposeNewRequest(http.MethodPatch, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func AddUserToGroup(c http.Client, projectID int, userID string, groupID int) req.Record {
	j := fmt.Sprintf(`{"projectId":"%v","userId":"%v"}`, projectID, userID)
	body := strings.NewReader(j)
	path := fmt.Sprintf("auth/group/%v/user", groupID)
	r := req.ComposeNewRequest(http.MethodPost, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetUserGroups(c http.Client, projectID int, name string) req.Record {
	q := url.Values{}
	q.Add("projectId", strconv.Itoa(projectID))
	q.Add("name", name)
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	path := fmt.Sprintf(`auth/group`)
	r := req.ComposeNewRequest(http.MethodGet, path, q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetPipelieList(c http.Client) req.Record {
	q := url.Values{}
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	path := fmt.Sprintf(`pipeline/base/page`)
	r := req.ComposeNewRequest(http.MethodGet, path, q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func CreateUserGroup(c http.Client, projectID int, groupName string) req.Record {
	j := fmt.Sprintf(`{"name":"%v","description":"this is my test group","projectId":"%v"}`, groupName, projectID)
	body := strings.NewReader(j)
	path := `auth/group`
	r := req.ComposeNewRequest(http.MethodPost, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func UploadFile(c http.Client, fname string, contentType string, data string) req.Record {
	path := fmt.Sprintf(`form/k8s-deploy/util/file`)
	r := req.ComposeNewMultipartRequest(http.MethodPost, path, nil, fname, contentType)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetDeployConfig(c http.Client, deployID int) req.Record {
	q := url.Values{}
	q.Add("deployId", strconv.Itoa(deployID))
	path := fmt.Sprintf(`k8s-deploy/config/`)
	r := req.ComposeNewRequest(http.MethodGet, path, q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func DeleteProject(c http.Client, projectID int) req.Record {
	path := fmt.Sprintf(`core/projects/%v`, projectID)
	r := req.ComposeNewRequest(http.MethodDelete, path, nil, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func DeleteEnv(c http.Client, envID int) req.Record {
	path := fmt.Sprintf(`core/env/%v`, envID)
	r := req.ComposeNewRequest(http.MethodDelete, path, nil, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func DeleteApp(c http.Client, appID int, appName string) req.Record {
	j := fmt.Sprintf(`{"id":%v,"name":"%v"}`, appID, appName)
	body := strings.NewReader(j)
	path := fmt.Sprintf(`core/app`)
	r := req.ComposeNewRequest(http.MethodDelete, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func DeleteDeploy(c http.Client, deployID int) req.Record {
	path := fmt.Sprintf(`deploy/base/%v`, deployID)
	r := req.ComposeNewRequest(http.MethodDelete, path, nil, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetDeployHistory(c http.Client, deployJobID int) req.Record {
	q := url.Values{}
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	path := fmt.Sprintf(`deploy/history/%v`, deployJobID)
	r := req.ComposeNewRequest(http.MethodGet, path, q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func ExecuteDeploy(c http.Client, deployID int) req.Record {
	j := fmt.Sprintf(`{"deployId":%v,"params":{}}`, deployID)
	body := strings.NewReader(j)
	r := req.ComposeNewRequest(http.MethodPatch, "deploy/base/execute", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func EditDeployDetails(c http.Client, deployID int, appDomain string, projectID int) req.Record {
	j := fmt.Sprintf(`{
	"apply": true,
	"kubeConfig": {
		"isEnvNamespace": true,
		"namespace": null,
		"cpu": "1000m",
		"mem": "512Mi",
		"duplicate": 1,
		"isChanged": true,
		"podSuccessFlag": false
	},
	"envVar": {
		"envVar": {"NEWTESTENV": "helloenv"},
		"isChanged": true
	},
	"hostAlias": {
		"hostAlias": {"devops.testhost.com": "6.6.6.6"},
		"isChanged": true
	},
	"containerLog": {
		"isChanged": true,
		"logPath": "/data/app/log"
	},
	"kubeConfigTemplate": {
		"isChanged": true,
		"template": ""
	},
	"configMap": {
		"configMap": [{
			"mountPath": "/data/myconfig/config.json",
			"data": "{\"hello\", \"world\"}",
			"configMapName": "%v-test-httpserver-app-config-map-1"
		}, {
			"mountPath": "/data/myconfig/config.yml",
			"data": "hello: world",
			"configMapName": "%v-test-httpserver-app-config-map-2"
		}, {
			"mountPath": "/data/myconfig/config.properties",
			"data": "hello=world",
			"configMapName": "%v-test-httpserver-app-config-map-3"
		}],
		"isChanged": true
	},
	"service": {
		"service": [{
			"protocolType": "TCP",
			"port": "9005",
			"serviceType": "ClusterIP"
		}, {
			"protocolType": "TCP",
			"port": "9005",
			"serviceType": "NodePort"
		}],
		"isChanged": true
	},
	"serviceTemplate": {
		"isChanged": true,
		"template": ""
	},
	"ingress": {
		"ingress": [{
			"type": "HTTP",
			"path": "/",
			"host": "%v",
			"internalPort": "9005",
			"annotations": {
				"WEBSOCKET": "false"
			}
		}],
		"isChanged": true
	},
	"ingressTemplate": {
		"isChanged": true,
		"template": ""
	},
	"grpc": {
        "grpc": [],
        "isChanged": false
    },
	"deployId": %v
}`, projectID, projectID, projectID, appDomain, deployID)
	body := strings.NewReader(j)
	r := req.ComposeNewRequest(http.MethodPut, "k8s-deploy/config/", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func EditGrpcDeployDetails(c http.Client, deployID int) req.Record {
	j := fmt.Sprintf(`{
    "apply": false,
    "kubeConfig": {
        "isEnvNamespace": true,
        "namespace": null,
        "cpu": "1000m",
        "mem": "512Mi",
        "duplicate": 1,
		"isChanged": false,
		"podSuccessFlag": false
    },
    "envVar": {
        "envVar": {},
        "isChanged": false
    },
    "hostAlias": {
        "hostAlias": {},
        "isChanged": false
    },
    "containerLog": {
        "isChanged": false,
        "logPath": ""
    },
    "kubeConfigTemplate": {
        "isChanged": false,
        "template": ""
    },
    "configMap": {
        "configMap": [],
        "isChanged": false
    },
    "service": {
        "service": [{
            "protocolType": "TCP",
            "port": "50051",
            "serviceType": "ClusterIP"
        }],
        "isChanged": true
    },
    "serviceTemplate": {
        "isChanged": false,
        "template": ""
    },
    "ingress": {
        "ingress": [],
        "isChanged": false
    },
    "ingressTemplate": {
        "isChanged": false,
        "template": ""
    },
    "grpc": {
        "grpc": [{
            "servicePort": "50051",
            "grpcName": ""
        }],
        "isChanged": true
    },
    "deployId": %v
}`, deployID)
	body := strings.NewReader(j)
	r := req.ComposeNewRequest(http.MethodPut, "k8s-deploy/config/", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

// GetDeployDetail使用分页查询和名称过滤来获取指定的部署信息
func GetDeployDetail(c http.Client, projectID int, deployName string) req.Record {
	q := url.Values{}
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	q.Add("projectId", strconv.Itoa(projectID))
	q.Add("name", deployName)
	r := req.ComposeNewRequest(http.MethodGet, "deploy/base/page", q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func CreateK8sImageDeployment(c http.Client, appID int, envID int, projectID int, name string, dockerCredentialID int, imageTag string) req.Record {
	j := fmt.Sprintf(`{
	"k8sTimeout": 300,
	"name": "%v",
	"appId": %v,
	"envId": %v,
	"deployType": "KUBERNETES_DEPLOY",
	"description": "some description",
	"sourceCreateRequest": {
		"imageAddr": "https://hub.iot.chinamobile.com/offline/smoke_build:%v",
		"imageAuthType": 1,
		"sourceType": "IMAGE",
		"credentialId": %v
	},
	"projectId": %v
}`, name, appID, envID, imageTag, dockerCredentialID, projectID)
	body := strings.NewReader(j)
	r := req.ComposeNewRequest(http.MethodPost, "deploy/base", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func CreateGrpcDeployment(c http.Client, appID int, envID int, projectID int, name string, dockerCredentialID int, imageTag string) req.Record {
	j := fmt.Sprintf(`{
	"k8sTimeout": 300,
	"name": "%v",
	"appId": %v,
	"envId": %v,
	"deployType": "KUBERNETES_DEPLOY",
	"description": "some description",
	"sourceCreateRequest": {
		"imageAddr": "https://hub.iot.chinamobile.com/offline/grpc-test-server:%v",
		"imageAuthType": 1,
		"sourceType": "IMAGE",
		"credentialId": %v
	},
	"projectId": %v
}`, name, appID, envID, imageTag, dockerCredentialID, projectID)
	body := strings.NewReader(j)
	r := req.ComposeNewRequest(http.MethodPost, "deploy/base", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetAppDetail(c http.Client, projectID int, name string) req.Record {
	q := url.Values{}
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	q.Add("projectId", strconv.Itoa(projectID))
	q.Add("name", name)
	r := req.ComposeNewRequest(http.MethodGet, "core/apps", q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func CreateApp(c http.Client, projectID int, userID string, appName string) req.Record {
	j := fmt.Sprintf(`{
	"appManager": "%v",
	"description": "测试http服务器应用",
	"name": "%v",
	"projectId": %v
}`, userID, appName, projectID)
	body := strings.NewReader(j)
	r := req.ComposeNewRequest(http.MethodPost, "core/app", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func UpdateEnv(c http.Client, envID int, envName string, prometheusIP string, envConfig string, envContext string, serverAddr string) req.Record {
	j := fmt.Sprintf(`{
	"zoneId": 3,
	"roomId": 2,
    "name": "我的测试环境",
    "type": "TEST",
    "namespace": "%v",
    "kubectx": "%v",
    "serverAddr": "https://%v:6443",
	"config": %v,
	"filebeat": "hub.iot.chinamobile.com/library/filebeat:6.2.4-onenet-0.5",
    "kafka": "[\"10.12.4.27:9092\"]",
    "prometheus": "http://%v:30909"
}`, envName, envContext, serverAddr, envConfig, prometheusIP, serverAddr)
	body := strings.NewReader(j)
	path := fmt.Sprintf("core/env/%v", envID)
	r := req.ComposeNewRequest(http.MethodPatch, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func UpdateEnvWithGRPC(c http.Client, envID int, envName string, prometheusIP string, envConfig string, envContext string, serverAddr string) req.Record {
	j := fmt.Sprintf(`{
	"zoneId": 3,
	"roomId": 2,
    "name": "我的测试环境",
    "type": "TEST",
    "namespace": "%v",
    "kubectx": "%v",
    "serverAddr": "https://%v:6443",
	"config": %v,
	"filebeat": "hub.iot.chinamobile.com/library/filebeat:6.2.4-onenet-0.5",
    "kafka": "[\"10.12.4.27:9092\"]",
    "prometheus": "http://%v:30909",
	"grpc": "%v",
    "grpcCert": "-----BEGIN CERTIFICATE-----\nMIIC6DCCAdCgAwIBAgIIYMTGhPrC5dAwDQYJKoZIhvcNAQELBQAwEjEQMA4GA1UE\nAxMHa3ViZS1jYTAeFw0xOTA1MTMwNjU1MzlaFw0yOTEwMTIxMDEwMTVaMC0xFTAT\nBgNVBAoTDHN5c3RlbTpub2RlczEUMBIGA1UEAxMLc3lzdGVtOm5vZGUwggEiMA0G\nCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDG0iXQhehz3NlEtC8HlUvzVcfhFW0N\nFNKa8b7a1N3YynMkg0AkQYahNsUG1xeZrU8zDqKqS8g2hl5Tjd8Ro+FdRzJD/M8O\nUEuG00m0nHYa1s2n30Yjucj0f/3egs+d7pvpTmTdp8AV/qilsBPQOtJDsZY9xf5g\nyC+9VpcBZFw++B2Y9HWbxx7YK5l4HbMy09VvVI+NU/0x/XmgJAxvXLAiX0bOD4Mz\nCe47XInwTvHmukKKrVwNjdWLLzcW8YEElXgubkVdfbtUZQPFKb7Yp60qqEARq72O\nQqAZ8vJcKwxy5jFRAQHsBpASH3PiMAR1EszJv412L0tRQR4LfJHkkxV1AgMBAAGj\nJzAlMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUEDDAKBggrBgEFBQcDAjANBgkqhkiG\n9w0BAQsFAAOCAQEApTN3d6L56mus1AT8sSaeNsKHBp6uARD3ld13hD9+CiwGeUu9\nrZf42UTcrZghrLsh47QwRDC1XC/jxme+vnRhJxftGs5tP65CkQPvsNx2Abr92vsY\n6wnMP86sZYwH1h7PuCyM5MJeXiMbjgr3Ae9UjlURHePE8ZuuyqvTlnfPJrvbM8yP\n5qVFaDIXiykrVwHsnpCbePaqBlbkAX5V6W/oNqE79RS4Bs6TPob0qzuhFki0Xqcd\nDuwZBAVGULwO/6TbkEAaPPARaQwN+k17g98vyYZRNx41dAYxAeN7WAN7qidQHbkE\nOTa18aX9tog2/Y6F6umXMpfwgt91lXhFjxSpNg==\n-----END CERTIFICATE-----\n",
    "grpcCa": "-----BEGIN CERTIFICATE-----\nMIICwjCCAaqgAwIBAgIBADANBgkqhkiG9w0BAQsFADASMRAwDgYDVQQDEwdrdWJl\nLWNhMB4XDTE5MDUxMzA2NTUzOVoXDTI5MDUxMDA2NTUzOVowEjEQMA4GA1UEAxMH\na3ViZS1jYTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALOcWR2UwkjI\n3F/cx2QjendaTYLq3Y7XAGnU0MGaivs/nZciTrFkzdgKW/B0Z8ENnfbZtlyIwiKe\nP+ZW7aH8STkOfzGH7j6CJL/gWvNI3bboqaaVO0Wk3P9IenUqIKzrFeARnlo4Q/Ka\n8WhaI0jA4VpgMUTv6wFByW0NePTaihHGqQ51f8rb8VwxBK1BoTq6zBXJHdEXxetr\neIaMoEUmQN2oC4jHBbC5k3YcL7iu1G/ajv4olR0GEUuwNxVH2bEun7gWpm/bwMcH\n6TyjgEILzX+Uji5rP/nFuUbAsWWUNvQ6k31x8xkYhABv0T8d29mbUQXsXZdAedB/\n+RbOBFrBj3MCAwEAAaMjMCEwDgYDVR0PAQH/BAQDAgKkMA8GA1UdEwEB/wQFMAMB\nAf8wDQYJKoZIhvcNAQELBQADggEBAJxCQjRn6i2gSOaLqnM8mf7KpPvveoeHb0ZD\nFXmHJRbY4emZez2/1UtehZkeP7iWO6I27P15ZJEsG7VvS0k0mmGKBCwE9GIv42P0\nFj5z77bVI1J+3Bm4QPJI+TM9FgVO9jjLQ81yELSm4thpLZZZUL0KDtegPuU1C9tH\niq5M8YN95uiRod7ve/5Bgi6+XW+ssLo7+DmhRKJod5uhftHgEXmw7HAdtWY7qjOJ\n3DsbPkS1qR6+6Jx91KmmrHU3akKLgRV1pbKkeeLQvtm2QCjcdO2x2rHEMgTN68ya\nn7vv2uL+ixIYyWVyHUTEG7FpNTRg7ZoWqsWzoC8GNBEst/OPrs8=\n-----END CERTIFICATE-----\n",
    "grpcKey": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAxtIl0IXoc9zZRLQvB5VL81XH4RVtDRTSmvG+2tTd2MpzJINA\nJEGGoTbFBtcXma1PMw6iqkvINoZeU43fEaPhXUcyQ/zPDlBLhtNJtJx2GtbNp99G\nI7nI9H/93oLPne6b6U5k3afAFf6opbAT0DrSQ7GWPcX+YMgvvVaXAWRcPvgdmPR1\nm8ce2CuZeB2zMtPVb1SPjVP9Mf15oCQMb1ywIl9Gzg+DMwnuO1yJ8E7x5rpCiq1c\nDY3Viy83FvGBBJV4Lm5FXX27VGUDxSm+2KetKqhAEau9jkKgGfLyXCsMcuYxUQEB\n7AaQEh9z4jAEdRLMyb+Ndi9LUUEeC3yR5JMVdQIDAQABAoIBABNsxrsHyj0/r46R\n2CJyiLgegmfQUxY7GZR/1/jDYWVj6joABM1GBaRMmJVhgHHIR/uiC2x9+PKy0BZz\nvv6XN1aItTWFRqmMWMS0cfgObU+T+wXSv6SP6z1QH/qUZzpz7JGv0hUB5beAaPO/\nL0Bh7tckS5x/cqn8BQYHprtBFe4k3WWlDio17udMTfTGj2Fv8SOeAmsQ5Gw6uJOR\n1GjysMyVIcTrLEmQl57Uz8b4qz3kKPWca+KKEg74w5SarKXqKhX7QRHylQgnkLK2\n7giR2RYAugFvsqqL4Q5x8IXH7zEYBHjbye2HoXSFchsmd9SjwXlkvk8WEptQIGvm\nMEpRXqUCgYEAzoH65jTqRelc+KHZ9uUY9ktURyQB7gEdRCY0hyqlAI/vdwVjAuR8\nFSgH9qmZVq7xbzr2qrIMbrYoqDDZSBoTa+ogjkF1H9nVPNlQKCbgYQUSBrK2aCDO\n2QzfOqCHPdQiWLqNyjDMBzkZgV/thiXWiguKHOrQTJc5q/B9OQqNLbMCgYEA9niN\nDwLn1nHbffhZFE1v7HjV3H+81GniIsDi2OsJGQgat6xUVZPOc0QrkfcrOQBk9U36\nrlxEU485tk48ULioq3di0sS6fBqfamS6adDEytp4MqXG6U72+NCV5eHsdELrceF8\nBal+n40QE9RRfLS5b/7uAUiMEPFVdbPmwENOrDcCgYEAmU5bgj10UlRthdM6KhVo\nE6hWt72ehR9kp6wpQNNCzYkNgHGKUKJpD5e5WcAMqxKTAD1o0838dtBanIovNFzP\nYETeyF0F45Bmwpad8ED0QHJwMHLKAcGhbfclXbPA0wDCQtaz3o+dWBtmuOoLPpSm\nkbMBZHhaDRITaXbOr+MKbgsCgYEAujoP6t25KpDQ1XeGZw6zmKscfASQOrbeRIAV\nZuz/7Mfw2AL/ncGWZgWGHj3xjJo9rhODa6cPgUtgwdyPOjasSxJjuvkmJos/FHaT\nW0yAxP0ZgLs9dh9SAGIqQI3ZyWae22cR/H06zXcaRMFR6LXsvzCRyKp2Gn8eoVaS\n7YZttTUCgYBLxwL7zyKSjiaT8j6GjLSPWf72jrLWC1Yd3IogpiNoEKYOG2qMiSos\naOysDnTuIFUr9+yEJI112Gii3Weroa5HiUmkLgp3RnFlUiioy3ipOafBK1uwkjnU\nvzmuxwdBqkxSf9ReMNXUzKtb5LAkJCFKEsbSChEuD6qpWlT3v6I0gw==\n-----END RSA PRIVATE KEY-----\n"
}`, envName, envContext, serverAddr, envConfig, prometheusIP, serverAddr)
	body := strings.NewReader(j)
	path := fmt.Sprintf("core/env/%v", envID)
	r := req.ComposeNewRequest(http.MethodPatch, path, nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetEnvDetail(c http.Client, projectID int) req.Record {
	q := url.Values{}
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	q.Add("projectId", strconv.Itoa(projectID))
	r := req.ComposeNewRequest(http.MethodGet, "core/envs", q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func CreateK8sEnv(c http.Client, projectID int, envName string, namespace string, serverAddr string, envConfig string, envContext string) req.Record {
	j := fmt.Sprintf(`{
	"config": %v,
	"kubectx": "%v",
	"name": "%v",
	"namespace": "%v",
	"serverAddr": "https://%v:6443",
	"type": "TEST",
	"projectId": "%v",
	"roomId": 2,
	"zoneId": 3
}`, envConfig, envContext, envName, namespace, serverAddr, projectID)
	body := strings.NewReader(j)
	r := req.ComposeNewRequest(http.MethodPost, "core/env", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetDockerCredential(c http.Client, projectId int, name string) req.Record {
	q := url.Values{}
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	q.Add("name", name)
	q.Add("keyword", name)
	q.Add("projectId", strconv.Itoa(projectId))
	q.Add("type", "DOCKER")
	r := req.ComposeNewRequest(http.MethodGet, "core/credential", q, nil)
	fmt.Println(r.URL)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func CreateDockerCredential(c http.Client, projectID int, name string) req.Record {
	j := fmt.Sprintf(`{
	"name": "%v",
	"userName": "jenkins",
	"password": "Iot@10086",
	"type": "DOCKER",
	"projectId": %v
}`, name, projectID)
	body := strings.NewReader(j)
	r := req.ComposeNewRequest(http.MethodPost, "core/credential", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetGitCredential(c http.Client, projectId int, name string) req.Record {
	q := url.Values{}
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	q.Add("name", name)
	q.Add("keyword", name)
	q.Add("projectId", strconv.Itoa(projectId))
	q.Add("type", "GIT")
	r := req.ComposeNewRequest(http.MethodGet, "core/credential", q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func CreateGitCredential(c http.Client, projectID int, name string) req.Record {
	j := fmt.Sprintf(`{
	"name": "%v",
	"userName": "huyangyi",
	"password": "199221hyy",
	"type": "GIT",
	"projectId": %v
}`, name, projectID)
	body := strings.NewReader(j)
	r := req.ComposeNewRequest(http.MethodPost, "core/credential", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetProjectDetail(c http.Client, name string) req.Record {
	q := url.Values{}
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	q.Add("query", name)
	r := req.ComposeNewRequest(http.MethodGet, "core/projects", q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func CreateProject(c http.Client, name string, identifier string) req.Record {
	j := fmt.Sprintf(`{
	"name": "%v",
	"identify": "%v",
	"description": "some description",
	"useDefault": false,
	"departmentId": 5
}`, name, identifier)
	body := strings.NewReader(j)
	r := req.ComposeNewRequest(http.MethodPost, "core/project", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func CreateArtifactLibrary(c http.Client, name string, projectID int) req.Record {
	pl := fmt.Sprintf(`{"name":"%v","auth":"0","projectId":%v}`, name, projectID)
	body := strings.NewReader(pl)
	r := req.ComposeNewRequest(http.MethodPost, "artifact/base/", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func CreateBuild(c http.Client, name string, projectID int) req.Record {
	body := strings.NewReader(fmt.Sprintf(`{"source":{"projectId":"%v","authPolicy":"NOAUTH","type":"GIT","address":"http://gitlab.onenet.com/huyangyi/devops-test-httpserver.git","branchPolicy":"SPECIFIC_BRANCH","specificBranch":"smoke_build"},"buildType":"GOLANG","name":"%v","projectId":"%v"}`, projectID, name, projectID))
	r := req.ComposeNewRequest(http.MethodPost, "build/", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func DuplicateBuild(c http.Client, name string, buildID int) req.Record {
	body := strings.NewReader(fmt.Sprintf(`{"name": "%s"}`, name))
	r := req.ComposeNewRequest(http.MethodPatch, fmt.Sprintf("build/base/copy/%d", buildID), nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func CreateGrpcBuild(c http.Client, name string, projectID int) req.Record {
	body := strings.NewReader(fmt.Sprintf(`{"source":{"projectId":"%v","authPolicy":"NOAUTH","type":"GIT","address":"http://gitlab.onenet.com/huyangyi/devops-grpc-test-server.git","branchPolicy":"SPECIFIC_BRANCH","specificBranch":"master"},"buildType":"GOLANG","name":"%v","projectId":"%v"}`, projectID, name, projectID))
	r := req.ComposeNewRequest(http.MethodPost, "build/", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func CreateVMBuild(c http.Client, name string, projectID int) req.Record {
	body := strings.NewReader(fmt.Sprintf(`{
    "source": {
        "projectId": "%v",
        "authPolicy": "NOAUTH",
        "type": "GIT",
        "address": "http://gitlab.onenet.com/huyangyi/devops-test-httpserver.git",
        "branchPolicy": "SPECIFIC_BRANCH",
        "specificBranch": "smoke_build"
    },
    "buildType": "GOLANG",
    "name": "%v",
    "projectId": "%v"
}`, projectID, name, projectID))
	r := req.ComposeNewRequest(http.MethodPost, "build/", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetArtifactID(c http.Client, name string, projectID int) req.Record {
	q := url.Values{}
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	q.Add("projectId", strconv.Itoa(projectID))
	q.Add("name", name)
	r := req.ComposeNewRequest(http.MethodGet, "artifact/base/", q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetArtifactByID(c http.Client, artifactID int) req.Record {
	r := req.ComposeNewRequest(http.MethodGet, fmt.Sprintf("artifact/base/query/%d", artifactID), nil, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetBuildDetail(c http.Client, buildID int) req.Record {
	r := req.ComposeNewRequest(http.MethodGet, fmt.Sprintf("build/%v", buildID), nil, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetBuildList(c http.Client, projectId int) req.Record {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("projectId", strconv.Itoa(projectId))
	r := req.ComposeNewRequest(http.MethodGet, fmt.Sprintf("build/paging"), q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetCommitID(c http.Client, buildID int) req.Record {
	r := req.ComposeNewRequest(http.MethodGet, fmt.Sprintf("build/%v/beforeBuild", buildID), nil, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func RunBuild(c http.Client, buildID int, commitID string, branchName string) req.Record {
	body := strings.NewReader(fmt.Sprintf(`{"branchOrTag":{"name":"%v","type":"BRANCH","commitId":"%v"}}`, branchName, commitID))
	r := req.ComposeNewRequest(http.MethodPost, fmt.Sprintf("build/%v/runningBuild", buildID), nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetBuildRecord(c http.Client, jobID int) req.Record {
	r := req.ComposeNewRequest(http.MethodGet, fmt.Sprintf("build/history/%v", jobID), nil, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetArtifactUploadRecord(c http.Client, artifactID int) req.Record {
	q := url.Values{}
	q.Add("artifactId", strconv.Itoa(artifactID))
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	path := fmt.Sprintf("artifact/upload/")
	r := req.ComposeNewRequest(http.MethodGet, path, q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func GetArtifactDownloadRecord(c http.Client, artifactID int) req.Record {
	q := url.Values{}
	q.Add("artifactId", strconv.Itoa(artifactID))
	q.Add("pageNum", "1")
	q.Add("pageSize", "10")
	path := fmt.Sprintf("artifact/downloads")
	r := req.ComposeNewRequest(http.MethodGet, path, q, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func DeleteArtifact(c http.Client, artifactID int) req.Record {
	path := fmt.Sprintf("artifact/base/%v", artifactID)
	r := req.ComposeNewRequest(http.MethodDelete, path, nil, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func DeleteBuild(c http.Client, buildID int) req.Record {
	path := fmt.Sprintf("build/%v", buildID)
	r := req.ComposeNewRequest(http.MethodDelete, path, nil, nil)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func EditBuild(c http.Client, gitCredentialID int, dockerCredentialID, projectID int, buildID int, buildName string, artifactID int, sourceID int, imageTag string) req.Record {
	body := strings.NewReader(fmt.Sprintf(`{
	"id": %v,
	"projectId": %v,
	"createdAt": 1578552539,
	"updatedAt": 1578569808,
	"buildType": "GOLANG",
	"name": "%v",
	"createUser": "胡阳译",
	"createUserId": "39b90b43-4b8f-4cd0-bd93-32273b4c70da",
	"buildCount": 0,
	"lastBuildAt": null,
	"status": null,
	"source": {
		"id": %v,
		"projectId": %v,
		"createdAt": 1578552539,
		"updatedAt": 1578569808,
		"deletedAt": 0,
		"type": "GIT",
		"address": "http://gitlab.onenet.com/huyangyi/devops-test-httpserver.git",
		"authPolicy": "CREDENTIAL",
		"credentialId": %v,
		"branchPolicy": "SPECIFIC_BRANCH",
		"specificBranch": "smoke_build",
		"useTag": null,
		"useBranch": null,
		"useWhich": []
	},
	"steps": [{
		"body": "{\"artifact\":\"%v\",\"version\":\"${DEVOPS_HISTORY_ID}\",\"cmd\":\"uploadme\", \"timeout\": \"120\"}",
		"type": "ARTIFACT_UPLOAD",
		"name": "上传制品",
		"id": 9999999
	}, {
		"body": "{\"artifact\":\"%v\",\"version\":\"${DEVOPS_HISTORY_ID}\", \"timeout\": \"120\"}",
		"type": "ARTIFACT_DOWNLOAD",
		"name": "下载制品",
		"id": 9999999
	}, {
		"body": "{\"cmd\":\"rm -rf __artifact.*\\nrm -rf uploadme\\nunzip *.zip\\ncp -r -n __artifact/uploadme/* .\", \"timeout\": \"120\"}",
		"type": "SHELL",
		"name": "执行Shell命令",
		"id": 9999999
	}, {
		"body": "{\"img\":\"hub.iot.chinamobile.com/library/golang:1.12.5\",\"path\":\"$GOPATH\",\"cmd\":\"go build -o httpserver main.go\", \"timeout\": \"120\"}",
		"type": "GOLANG",
		"name": "Golang构建",
		"id": 9999999
	}, {
		"body": "{\"dockerfile\":\"Dockerfile\",\"dockerhub\":\"hub.iot.chinamobile.com\",\"dockerRepo\":\"offline\",\"imageName\":\"${DEVOPS_SCMID}\",\"imageTag\":\"%v\",\"credentialId\":\"%v\", \"timeout\": \"120\"}",
		"type": "IMG_BUILD",
		"name": "镜像构建",
		"id": 9999999
	}],
	"parameters": [],
	"groupAuth": [
		{
            "groupId": 0,
            "groupName": "任务创建人",
            "objectId": 1558,
            "moduleType": "BUILD",
            "readonly": true,
            "run": true,
            "edit": true,
            "delete": true,
            "auth": true,
            "copy": true,
            "forbid": true,
            "canChange": false
        }
	]
}`, buildID, projectID, buildName, sourceID, projectID, gitCredentialID, artifactID, artifactID, imageTag, dockerCredentialID))
	r := req.ComposeNewRequest(http.MethodPatch, "build/", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func EditGrpcBuild(c http.Client, gitCredentialID int, dockerCredentialID, projectID int, buildID int, buildName string, sourceID int, imageTag string) req.Record {
	body := strings.NewReader(fmt.Sprintf(`{
	"id": %v,
	"projectId": %v,
	"createdAt": 1578552539,
	"updatedAt": 1578569808,
	"buildType": "GOLANG",
	"name": "%v",
	"createUser": "胡阳译",
	"createUserId": "39b90b43-4b8f-4cd0-bd93-32273b4c70da",
	"buildCount": 0,
	"lastBuildAt": null,
	"status": null,
	"source": {
		"id": %v,
		"projectId": %v,
		"createdAt": 1578552539,
		"updatedAt": 1578569808,
		"deletedAt": 0,
		"type": "GIT",
		"address": "http://gitlab.onenet.com/huyangyi/devops-grpc-test-server.git",
		"authPolicy": "CREDENTIAL",
		"credentialId": %v,
		"branchPolicy": "SPECIFIC_BRANCH",
		"specificBranch": "master",
		"useTag": null,
		"useBranch": null,
		"useWhich": []
	},
	"steps": [{
        "body": "{\"dockerfile\":\"Dockerfile\",\"dockerhub\":\"hub.iot.chinamobile.com\",\"dockerRepo\":\"offline\",\"imageName\":\"grpc-test-server\",\"imageTag\":\"%v\",\"credentialId\":\"%v\",\"timeout\":\"300\"}",
        "type": "IMG_BUILD",
        "name": "镜像构建",
        "id": 999999
    }],
	"parameters": [],
	"groupAuth": [
		{
            "groupId": 0,
            "groupName": "任务创建人",
            "objectId": 1558,
            "moduleType": "BUILD",
            "readonly": true,
            "run": true,
            "edit": true,
            "delete": true,
            "auth": true,
            "copy": true,
            "forbid": true,
            "canChange": false
        }]
}`, buildID, projectID, buildName, sourceID, projectID, gitCredentialID, imageTag, dockerCredentialID))
	r := req.ComposeNewRequest(http.MethodPatch, "build/", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}

func EditVMBuild(c http.Client, projectID int, buildID int, buildName string, sourceID int, artifactID int) req.Record {
	body := strings.NewReader(fmt.Sprintf(`{
    "id": %v,
    "projectId": %v,
    "createdAt": 1583981339,
    "updatedAt": 1583981339,
    "buildType": "GOLANG",
    "name": "%v",
    "createUser": "胡阳译",
    "createUserId": "39b90b43-4b8f-4cd0-bd93-32273b4c70da",
    "buildCount": 0,
    "lastBuildAt": null,
    "status": null,
    "source": {
        "id": %v,
        "projectId": %v,
        "createdAt": 1583981339,
        "updatedAt": 1583981339,
        "deletedAt": 0,
        "type": "GIT",
        "address": "http://gitlab.onenet.com/huyangyi/devops-test-httpserver.git",
        "authPolicy": "NOAUTH",
        "credentialId": null,
        "branchPolicy": "SPECIFIC_BRANCH",
        "specificBranch": "smoke_build",
        "useTag": null,
        "useBranch": null,
        "useWhich": []
    },
    "steps": [{
        "body": "{\"img\":\"hub.iot.chinamobile.com/library/golang:1.12.5\",\"path\":\"$GOPATH\",\"cmd\":\"go build -o yang-httpserver uploadme/main.go\",\"timeout\":\"300\"}",
        "type": "GOLANG",
        "name": "GOLANG构建",
        "id": 99999999
    }, {
        "body": "{\"artifact\":%v,\"version\":\"\",\"cmd\":\"yang-httpserver\\nuploadme/appctl.sh\",\"timeout\":\"300\"}",
        "type": "ARTIFACT_UPLOAD",
        "name": "上传制品"
    }],
	"parameters": [],
	"groupAuth": [
		{
            "groupId": 0,
            "groupName": "任务创建人",
            "objectId": 1558,
            "moduleType": "BUILD",
            "readonly": true,
            "run": true,
            "edit": true,
            "delete": true,
            "auth": true,
            "copy": true,
            "forbid": true,
            "canChange": false
        }]
}`, buildID, projectID, buildName, sourceID, projectID, artifactID))
	r := req.ComposeNewRequest(http.MethodPatch, "build/", nil, body)
	resp := req.SendRequestAndGetResponse(c, r)
	return resp
}
