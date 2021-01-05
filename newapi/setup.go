package newapi

import (
	"fmt"
	"github.com/Blackoutta/profari"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/API"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/conf"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/random"
)

type SetupParams struct {
	ProjectID   int
	EnvID       int
	DeployID    int
	AppID       int
	AppName     string
	Namespace   string
	DockerID    int
	EnvName     string
	GitID       int
	VmGroupId   int
	VmMachineId int
}

// SetupParams.Teardown 清理BasicSetup函数创建的所有资源
// 按依赖顺序依次删除
// todo 增加依赖导致的删除错误测试
func (p *SetupParams) Teardown(t *profari.Client) {
	// clean deploy
	t.Send(DeleteDeploy{
		DeployID: p.DeployID,
	}).AssertContainString("删除部署应成功", t.Resp, "success")

	// clean k8s env
	t.Send(DeleteK8sEnv{
		EnvID: p.EnvID,
	}).AssertContainString("删除k8s环境应成功", t.Resp, "success")

	// clean app
	t.Send(&DeleteApp{
		Id:   p.AppID,
		Name: p.AppName,
	}).AssertContainString("删除应用应成功", t.Resp, "success")

	// clean git credential
	t.Send(&DeleteCredential{
		Id: p.GitID,
	}).AssertContainString("删除Git凭证应成功", t.Resp, "success")

	// clean docker credential
	t.Send(&DeleteCredential{
		Id: p.DockerID,
	}).AssertContainString("删除Docker凭证应成功", t.Resp, "success")

	//// clean vm group
	//t.Send(&DeleteVmGroup{
	//	GroupID: p.VmGroupId,
	//}).AssertContainString("删除vm group应成功", t.Resp, "success")
	//
	//// clean vm machine
	//t.Send(&DeleteVmMachine{
	//	VmMachineId: p.VmMachineId,
	//}).AssertContainString("删除vm machine应成功", t.Resp, "success")

	// clean project
	t.Send(DeleteProject{
		ProjectID: p.ProjectID,
	}).AssertContainString("删除项目应成功", t.Resp, "success")
}

func BasicSetup(t *profari.Client) *SetupParams {
	var p SetupParams
	// 开始
	// 创建项目
	var gen API.GeneralResp
	projectName := "proj-" + random.ShortGUID()
	t.Send(CreateProject{
		DepartmentId: 5,
		Description:  "这是一个项目",
		Forbid:       false,
		Identify:     projectName,
		Name:         projectName,
		UseDefault:   false,
	}).DecodeJSON(&gen).AssertContainString("创建项目应成功", t.Resp, "success")

	// 获取projectID
	var projects API.Projects
	t.Send(GetProject{ProjectName: projectName}).DecodeJSON(&projects).AssertEqualInt("获取项目ID应成功", len(projects.Data.Data), 1)
	p.ProjectID = projects.Data.Data[0].ID

	// 创建Docker凭证
	var dockerCredentials API.DockerCredentials
	t.Send(CreateDockerCredential{
		Description: "this is a docker credential",
		Name:        "docker_cred_" + random.ShortGUID(),
		Password:    "Iot@10086",
		ProjectId:   p.ProjectID,
		Type:        "DOCKER",
		UserName:    "jenkins",
	}).AssertContainString("创建Docker凭证应成功", t.Resp, "success")

	// 获取Docker Credential ID
	t.Send(&GetCredential{
		Type:      "DOCKER",
		ProjectId: p.ProjectID,
		Name:      "docker",
	}).DecodeJSON(&dockerCredentials).AssertEqualInt("获取Docker凭证ID应成功", len(dockerCredentials.Data.Data), 1)
	if len(dockerCredentials.Data.Data) < 1 {
		return nil
	}
	p.DockerID = dockerCredentials.Data.Data[0].ID

	// 创建git凭证
	var gitCredentials API.GitCredentials
	t.Send(CreateGitCredential{
		CreateDockerCredential: CreateDockerCredential{
			Description: "this is a err git credential", // 弄一个测试的用户密码 ?
			Name:        "git_cred_" + random.ShortGUID(),
			Password:    "git-password",
			ProjectId:   p.ProjectID,
			Type:        "GIT",
			UserName:    "git-userName",
		},
	}).AssertContainString("创建Git凭证应成功", t.Resp, "success")

	// 获取Git Credential ID
	t.Send(&GetCredential{
		Type:      "GIT",
		ProjectId: p.ProjectID,
		Name:      "git",
	}).DecodeJSON(&gitCredentials).AssertEqualInt("获取Git凭证ID应成功", len(gitCredentials.Data.Data), 1)
	if len(gitCredentials.Data.Data) < 1 {
		return nil
	}
	p.GitID = gitCredentials.Data.Data[0].ID

	// 创建K8S环境
	p.EnvName = "env-" + random.ShortGUID()
	p.Namespace = "e-" + random.ShortGUID()
	config := []byte(conf.ReadEnvFile())
	kubectx := conf.EnvContext

	// K8S环境连通性检查
	var k8sEnvStatus API.K8SEnvStatus
	t.Send(
		K8SConnectionCheck{
			Config:  config,
			Kubectx: kubectx,
		}).DecodeJSON(&k8sEnvStatus).AssertContainString("k8s连通性检查应通过，环境状态为可用(NORMAL)", k8sEnvStatus.Data.Status, "NORMAL")

	t.Send(CreateK8sEnv{
		Config:     []byte(conf.ReadEnvFile()),
		Kubectx:    conf.EnvContext,
		Name:       p.EnvName,
		Namespace:  p.Namespace,
		ServerAddr: "https://" + conf.ServerAddr + ":6443",
		Type:       "TEST",
		ProjectId:  p.ProjectID,
		RoomId:     2,
		ZoneId:     4,
	}).AssertContainString("创建k8s环境应成功", t.Resp, "success")

	// 获取环境ID
	var envs API.Environments
	t.Send(GetEnvs{
		ProjectId: p.ProjectID,
	}).DecodeJSON(&envs).AssertEqualInt("获取环境ID应成功", len(envs.Data.Data), 1)
	if len(envs.Data.Data) < 1 {
		return nil
	}
	p.EnvID = envs.Data.Data[0].ID

	// 创建应用
	p.AppName = "app-" + random.ShortGUID()
	t.Send(CreateApp{
		AppManager:  conf.UserID,
		Description: "this is an app",
		Name:        p.AppName,
		ProjectId:   p.ProjectID,
	}).AssertContainString("创建应用应成功", t.Resp, "success")

	// 获取应用ID
	var apps API.Apps
	t.Send(
		GetApps{
			ProjectId: p.ProjectID,
			Name:      p.AppName,
		}).DecodeJSON(&apps).AssertEqualInt("获取App ID应成功", len(apps.Data.Data), 1)

	if len(apps.Data.Data) < 1 {
		return nil
	}
	p.AppID = apps.Data.Data[0].ID

	// 创建主机组
	t.Send(CreateVMGroup{
		Description: "description",
		Name:        "VMGroup-" + random.ShortGUID(),
		ProjectId:   p.ProjectID,
		RoomId:      2,
		ZoneId:      4,
	}).AssertContainString("创建主机组应成功", t.Resp, "success")

	//查询主机组ID
	var vmGroups API.VMGroups
	t.Send(GetVMGroups{
		ProjectID: p.ProjectID,
	}).DecodeJSON(&vmGroups).AssertEqualInt("获取VMGroup ID应成功", len(apps.Data.Data), 1)
	p.VmGroupId = vmGroups.Data.Data[0].ID

	// 创建主机
	t.Send(CreateVMMachine{
		AuthType:    "PASSWORD",
		Description: "description",
		Ip:          "1.1.1.1",
		Name:        "VMMchines-" + random.ShortGUID(),
		Password:    "root",
		Port:        8080,
		SshKey:      "",
		UserName:    "root",
		VmGroupId:   p.VmGroupId,
	}).AssertContainString("创建主机应成功", t.Resp, "success")

	// 获取主机ID
	var vmMachines API.VMMachines
	t.Send(
		GetVMMachines{
			VmGroupId: p.VmGroupId,
			Query:     "",
		}).DecodeJSON(&vmMachines).AssertEqualInt("获取vmMachine ID应成功", len(vmMachines.Data.Data), 1)

	if len(vmMachines.Data.Data) < 1 {
		return nil
	}
	p.VmMachineId = vmMachines.Data.Data[0].ID
	fmt.Println("VmMachineId: ", p.VmMachineId)
	// 创建部署
	t.Send(
		CreateDeploy{
			AppId:       p.AppID,
			DeployType:  DeployTypeKubernetesDeploy,
			Description: "this is a deployment",
			EnvId:       p.EnvID,
			K8sTimeout:  K8sTimeoutFiveMinutes,
			Name:        "dp_" + random.ShortGUID(),
			ProjectId:   p.ProjectID,
			SourceCreateRequest: SourceCreateRequest{
				ArtifactVersion: "",
				CredentialId:    p.DockerID,
				ImageAddr:       "https://hub.iot.chinamobile.com/offline/smoke_build:build-test",
				ImageAuthType:   ImageAuthEnabled,
				SourceType:      SourceTypeImage,
			},
		}).DecodeJSON(&gen).AssertContainString("创建k8s部署应成功", t.Resp, "success")
	p.DeployID = gen.Data

	return &p
}
