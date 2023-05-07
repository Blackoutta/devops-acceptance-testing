package order

import (
	"fmt"
	"os"

	"github.com/Blackoutta/profari"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/API"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/newapi"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/random"
)

type OrderTest struct {
	Name         string
	ErrChan      chan error
	SkipTeardown bool
	*profari.Client
	logFile *os.File
	*newapi.SetupParams
	suiteParams
}

type suiteParams struct {
	groupID int
}

func (t *OrderTest) GetName() string {
	return t.Name
}

func (t *OrderTest) GetErrChan() chan error {
	return t.ErrChan
}

func (t *OrderTest) Run() {
	var err error
	// You must initialize the profari client before test starts
	t.Client, t.logFile, err = profari.NewClient(t.Name, t.ErrChan)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 开始

	// 准备工作
	t.SetupParams = newapi.BasicSetup(t.Client)
	fmt.Printf("%+v\n", t.SetupParams)

	// 查询机房及可用区
	t.Send(newapi.GetRoomAndZone{}).AssertContainString("应可以查到机房及可用区信息", t.Resp, "龙洲湾")

	// 创建主机组
	var gen API.GeneralResp
	roomID := 2
	zoneID := 4
	t.Send(
		newapi.CreateVMGroup{
			Description: "this is a vm group",
			Name:        "vmg_" + random.ShortGUID(),
			ProjectId:   t.ProjectID,
			RoomId:      roomID,
			ZoneId:      zoneID,
		})
	t.DecodeJSON(&gen).AssertContainString("创建主机组应成功", t.Resp, "success")

	t.groupID = gen.Data

	// 查询主机组
	var vmGroup API.VMGroupDetail
	t.Send(
		newapi.GetVMGroup{
			GroupID: t.groupID,
		})

	t.DecodeJSON(&vmGroup).AssertEqualInt("通过Id应能查询到主机组详情", vmGroup.Data.RoomId, roomID)
	// 查询工单新增或编辑时meta信息

	t.Send(
		newapi.GetVMOrderMeta{}).AssertContainString("创建工单前能成功获取工单元数据", t.Resp, "IPv4和IPv6双栈")

	// 创建工单
	t.Send(
		newapi.CreateVMWorkOrder{
			GroupId:     t.groupID,
			ProjectName: "开放平台部-OneNET",
			TeamName:    "开放平台部-云平台开发团队",
			ComputerConfig: newapi.ComputerConfig{
				ApplyAmount:  2,
				CardAmount:   "四网卡",
				CardUse:      "测试多网卡用途",
				ComputeType:  "虚拟机",
				CpuAmount:    "8核",
				DiskVapacity: 10,
				DurationDnd:  "1992-02-02",
				Memory:       "2GB",
				Os:           "Ubuntu",
				OsVersion:    "Ubuntu14.04",
				Region:       "水土",
				Remark:       "测试备注",
				Segment:      "192.168.1.1",
				UsageWo:      "效能平台申请",
			},
		}).AssertContainString("创建工单应成功", t.Resp, "success")

	// 查询工单记录列表
	var vmOrders API.VMWorkOrders
	t.Send(
		newapi.GetVMOrderList{
			GroupID: t.groupID,
		}).DecodeJSON(&vmOrders).AssertEqualInt("查询工单记录列表应能查到一条工单记录", 1, len(vmOrders.Data.Data))
	t.AssertContainString("工单的状态应是处理中", vmOrders.Data.Data[0].Status, "PROCESSING")
}

func (t *OrderTest) Teardown() {
	if t.SkipTeardown == true {
		t.Println("Skipping Teardown...")
		t.logFile.Close()
		t.EndTest()
		return
	}

	defer t.logFile.Close()
	defer t.EndTest()

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

	t.Send(newapi.DeleteVmGroup{
		GroupID: t.groupID,
	}).AssertContainString("删除主机组应成功", t.Resp, "success")

	t.Send(newapi.DeleteProject{
		ProjectID: t.ProjectID,
	}).AssertContainString("删除项目应成功", t.Resp, "success")

}
