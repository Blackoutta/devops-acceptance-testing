package common

import (
	"fmt"
	"os"
	"time"

	"github.com/Blackoutta/profari"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/API"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/newapi"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/random"
)

type CommonCoreTest struct {
	Name         string
	ErrChan      chan error
	SkipTeardown bool
	*profari.Client
	logFile *os.File
	*newapi.SetupParams
	suiteParams
}

type suiteParams struct {
}

func (t *CommonCoreTest) GetName() string {
	return t.Name
}

func (t *CommonCoreTest) GetErrChan() chan error {
	return t.ErrChan
}

func (t *CommonCoreTest) Run() {
	var err error
	// You must initialize the profari client before test starts
	t.Client, t.logFile, err = profari.NewClient(t.Name, t.ErrChan)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 准备工作
	t.SetupParams = newapi.BasicSetup(t.Client)
	fmt.Printf("%+v\n", t.SetupParams)

	//// 更新项目
	t.Send(
		newapi.UpdateProjectById{
			ProjectID:    t.ProjectID,
			DepartmentId: 10,
			Description:  "edited project",
			Forbid:       false,
			Name:         "edited_proj_" + random.ShortGUID(),
		}).AssertContainString("通过Id修改项目信息应成功", t.Resp, "success")

	// 通过项目ID查看项目详情
	var projectDetail API.ProjectDetail
	t.Send(
		newapi.GetProjectDetailByID{
			ProjectID: t.ProjectID,
		}).DecodeJSON(&projectDetail).AssertContainString("通过ID查看项目详情应可以看到项目被修改成功", t.Resp, "edited project")

	// 新增项目收藏
	t.Send(
		newapi.AddProjectCollection{
			ProjectId: t.ProjectID,
		}).AssertContainString("新增项目收藏应成功", t.Resp, "success")

	// 查看项目收藏
	t.Send(newapi.GetProjectCollection{}).AssertContainString("应可以查询到项目收藏", t.Resp, "edited project")

	// 删除项目收藏
	t.Send(newapi.DeleteProjectCollection{
		ProjectID: t.ProjectID,
	}).AssertContainString("应可以删除项目收藏", t.Resp, "success")

	time.Sleep(time.Second)
	// 查看项目收藏，确定删除
	t.Send(newapi.GetProjectCollection{}).AssertContainString("查看项目收藏，确定删除", t.Resp, `data":[]`)

	// 新增应用收藏
	t.Send(newapi.AddAppCollection{
		AppId: t.AppID,
	}).AssertContainString("新增应用收藏应成功", t.Resp, "success")

	// 查看应用收藏，确定增加成功
	t.Send(newapi.GetAppCollection{
		ProjectId: t.ProjectID,
	}).AssertContainString("查询应用收藏，确认收藏成功", t.Resp, "this is an app")
	fmt.Println(t.Resp)

	// 删除应用收藏
	t.Send(newapi.DeleteAppCollection{
		AppId: t.AppID,
	}).AssertContainString("删除应用收藏应成功", t.Resp, "success")

	time.Sleep(time.Second)
	// 查看应用收藏，确定增加成功
	t.Send(newapi.GetAppCollection{
		ProjectId: t.ProjectID,
	}).AssertContainString("查询应用收藏，确认删除收藏成功", t.Resp, `"data":[]`)
	time.Sleep(time.Second)

	// 更新主机
	t.Send(newapi.UpdateVM{
		VmMachineId: t.VmMachineId,
		AuthType:    "PASSWORD",
		Description: "no description",
		Ip:          "",
		Name:        "",
		Password:    "",
		Port:        3306,
		SshKey:      "",
		UserName:    "root",
	}).AssertContainString("通过Id修改主机信息应成功", t.Resp, "success")

	fmt.Println(t.Resp)
}

func (t *CommonCoreTest) Teardown() {
	if t.SkipTeardown == true {
		t.Println("Skipping Teardown...")
		t.logFile.Close()
		t.EndTest()
		return
	}

	defer t.logFile.Close()
	defer t.EndTest()

	// clean vm machine
	t.Send(newapi.DeleteVmMachine{
		VmMachineId: t.VmMachineId,
	}).AssertContainString("删除vm machine应成功", t.Resp, "success")

	// clean vm group
	t.Send(newapi.DeleteVmGroup{
		GroupID: t.VmGroupId,
	}).AssertContainString("删除vm group应成功", t.Resp, "success")

	t.Send(
		newapi.DeleteDeploy{
			DeployID: t.DeployID,
		}).AssertContainString("删除部署应成功", t.Resp, "success")

	t.Send(newapi.DeleteK8sEnv{
		EnvID: t.EnvID,
	}).AssertContainString("删除k8s环境应成功", t.Resp, "success")

	t.Send(
		newapi.DeleteApp{
			Id:   t.AppID,
			Name: t.AppName,
		}).AssertContainString("删除应用应成功", t.Resp, "success")

	t.Send(newapi.DeleteProject{
		ProjectID: t.ProjectID,
	}).AssertContainString("删除项目应成功", t.Resp, "success")

}
