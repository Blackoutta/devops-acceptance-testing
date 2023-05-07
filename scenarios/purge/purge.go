package purge

import (
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"github.com/Blackoutta/profari"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/API"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/newapi"
)

var wg sync.WaitGroup
var mux sync.Mutex

type resources struct {
	pid       int
	deploys   API.Deploys
	apps      API.Apps
	tests     API.Tests
	pipelines API.Pipelines
	k8sEnvs   API.Environments
	vmEnvs    API.Environments
	builds    API.Builds
	artifacts API.Artifacts
}

type PurgeTest struct {
	Name         string
	ErrChan      chan error
	SkipTeardown bool
	*profari.Client
	logFile *os.File
	*newapi.SetupParams
}

func (t *PurgeTest) GetName() string {
	return t.Name
}

func (t *PurgeTest) GetErrChan() chan error {
	return t.ErrChan
}

func (t *PurgeTest) Run() {
	creator := "huyangyi"

	var err error
	// You must initialize the profari client before test starts
	t.Client, t.logFile, err = profari.NewClient(t.Name, t.ErrChan)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 查询项目
	var projects API.Projects
	t.Send(newapi.GetProjectList{
		PageNum: 1,
	}).DecodeJSON(&projects).AssertEqualInt("查询出最大100个项目", projects.ErrorCode, 0)

	total := projects.Data.Total
	fmt.Printf("查询到的项目总数: %v\n", total)

	if total == 0 {
		t.Println("该用户的项目总数为0，无需清理，程序退出")
		os.Exit(0)
	}

	pages := math.Ceil(float64(total) / 100)
	t.Printf("需要查询的页数: %v\n", pages)

	allProjects := make([]API.Project, 0, total)

	for i := 1; i <= int(pages); i++ {
		t.Send(
			newapi.GetProjectList{
				PageNum: i,
			},
		).DecodeJSON(&projects).AssertContainString(fmt.Sprintf("查询第%d页内容", i), t.Resp, "success")
		for _, p := range projects.Data.Data {
			if p.CreatorUserName == creator {
				allProjects = append(allProjects, p)
			}
		}
	}

	wg.Add(len(allProjects))

	rsChan := make(chan resources, 100)

	go func() {
		for {
			select {
			case rs := <-rsChan:
				fmt.Println("rs pulled out")
				go func() {
					// ****************删除资源****************//
					// 编辑部署

					// 删除部署
					if len(rs.deploys.Data.Data) > 0 {
						for _, v := range rs.deploys.Data.Data {
							t.Send(
								newapi.DeleteDeploy{
									DeployID: v.ID,
								},
							).AssertContainString(fmt.Sprintf("删除部署: %d", v.ID), t.Resp, "success")
						}
					}

					// 删除应用
					if len(rs.apps.Data.Data) > 0 {
						for _, v := range rs.apps.Data.Data {
							t.Send(
								newapi.DeleteApp{
									Id:   v.ID,
									Name: v.Name,
								},
							).AssertContainString(fmt.Sprintf("删除应用: %d", v.ID), t.Resp, "success")
						}
					}

					// 删除测试
					if len(rs.tests.Data.Data) > 0 {
						for _, v := range rs.tests.Data.Data {
							t.Send(
								newapi.DeleteTest{
									TestID: v.ID,
								},
							).AssertContainString(fmt.Sprintf("删除测试: %d", v.ID), t.Resp, "success")
						}
					}

					// 删除流水线
					if len(rs.pipelines.Data.Data) > 0 {
						for _, v := range rs.pipelines.Data.Data {
							t.Send(
								newapi.DeletePipeline{
									PipelineID: v.ID,
								},
							).AssertContainString(fmt.Sprintf("删除流水线: %d", v.ID), t.Resp, "success")
						}
					}

					// 删除k8s环境
					if len(rs.k8sEnvs.Data.Data) > 0 {
						for _, v := range rs.k8sEnvs.Data.Data {
							t.Send(
								newapi.DeleteK8sEnv{
									EnvID: v.ID,
								},
							).AssertContainString(fmt.Sprintf("删除k8s环境: %d", v.ID), t.Resp, "success")
						}
					}

					// 删除虚拟机主机组环境
					if len(rs.vmEnvs.Data.Data) > 0 {
						for _, v := range rs.vmEnvs.Data.Data {
							t.Send(
								newapi.DeleteVMEnv{
									EnvID: v.ID,
								},
							).AssertContainString(fmt.Sprintf("删除虚拟机主机组环境: %d", v.ID), t.Resp, "success")
						}
					}

					// 删除构建
					if len(rs.builds.Data.Data) > 0 {
						for _, v := range rs.builds.Data.Data {
							t.Send(
								newapi.DeleteBuild{
									BuildID: v.ID,
								},
							).AssertContainString(fmt.Sprintf("删除构建: %d", v.ID), t.Resp, "success")
						}
					}

					// 删除制品库
					fmt.Println(len(rs.artifacts.Data.Data))
					if len(rs.artifacts.Data.Data) > 0 {
						for _, v := range rs.artifacts.Data.Data {
							t.Send(
								newapi.DeleteArtifact{
									ArtifactID: v.ID,
								},
							).AssertContainString(fmt.Sprintf("删除制品库: %d", v.ID), t.Resp, "success")
						}
					}

					// 删除项目
					var gen API.GeneralResp
					t.Send(
						newapi.DeleteProject{
							ProjectID: rs.pid,
						},
					).DecodeJSON(&gen).AssertEqualInt(fmt.Sprintf("删除项目: %d", rs.pid), 0, gen.ErrorCode)

				}()
			}
		}
	}()

	t.makeResource(allProjects, rsChan)
	wg.Wait()
}

func (t *PurgeTest) Teardown() {
	t.Println("There is no teardown in this suite, skipping...")
	t.EndTest()
}

func (t *PurgeTest) makeResource(allProjects []API.Project, rsChan chan resources) {
	// f1, err := os.OpenFile("devops_dev", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	// if err != nil {
	// 	panic(err)
	// }

	// f2, err := os.OpenFile("devops_test", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	// if err != nil {
	// 	panic(err)
	// }

	for _, pj := range allProjects {
		go func(pid int) {
			var rs resources
			rs.pid = pid
			// 获取部署

			var deploys API.Deploys

			t.Send(
				newapi.GetDeployList{
					ProjectID: pid,
				},
			).DecodeJSON(&deploys).AssertContainString("获取部署", t.Resp, "success")
			time.Sleep(50 * time.Millisecond)

			rs.deploys = deploys

			// 获取应用
			var apps API.Apps

			t.Send(
				newapi.GetAppList{
					ProjectID: pid,
				},
			).DecodeJSON(&apps).AssertContainString("获取应用", t.Resp, "success")

			time.Sleep(50 * time.Millisecond)

			rs.apps = apps

			var tests API.Tests
			// 获取测试

			t.Send(
				newapi.GetTestList{
					ProjectID: pid,
				},
			).DecodeJSON(&tests).AssertContainString("获取测试", t.Resp, "success")
			time.Sleep(50 * time.Millisecond)

			rs.tests = tests

			// 获取流水线
			var pipelines API.Pipelines

			t.Send(
				newapi.GetPipelineList{
					ProjectID: pid,
				},
			).DecodeJSON(&pipelines).AssertContainString("获取流水线", t.Resp, "success")
			time.Sleep(50 * time.Millisecond)

			rs.pipelines = pipelines

			// 获取k8s环境
			var k8sEnvs API.Environments

			t.Send(
				newapi.GetK8SEnvList{
					ProjectID: pid,
				},
			).DecodeJSON(&k8sEnvs).AssertContainString("获取k8s环境", t.Resp, "success")
			time.Sleep(50 * time.Millisecond)

			rs.k8sEnvs = k8sEnvs
			fmt.Printf("%+v\n", k8sEnvs.Data.Data)

			// for _, env := range k8sEnvs.Data.Data {
			// 	if env.Kubectx == "devops-dev" {
			// 		_, err := f1.WriteString(strconv.Itoa(env.ID) + ",")
			// 		if err != nil {
			// 			panic(err)
			// 		}
			// 	}

			// 	if env.Kubectx == "devops-test" {
			// 		_, err := f2.WriteString(strconv.Itoa(env.ID) + ",")
			// 		if err != nil {
			// 			panic(err)
			// 		}
			// 	}
			// }

			// 获取虚拟机环境
			var vmEnvs API.Environments

			t.Send(
				newapi.GetVMEnvList{
					ProjectID: pid,
				},
			).DecodeJSON(&vmEnvs).AssertContainString("获取虚拟机主机组环境", t.Resp, "success")
			time.Sleep(50 * time.Millisecond)

			rs.vmEnvs = vmEnvs

			// 获取构建
			var builds API.Builds

			t.Send(
				newapi.GetBuildList{
					ProjectID: pid,
				},
			).DecodeJSON(&builds).AssertContainString("获取构建", t.Resp, "success")
			time.Sleep(50 * time.Millisecond)

			rs.builds = builds

			// 获取制品库
			var artifacts API.Artifacts

			t.Send(
				newapi.GetArtifactList{
					ProjectID: pid,
				},
			).DecodeJSON(&artifacts).AssertContainString("获取制品库", t.Resp, "success")
			time.Sleep(50 * time.Millisecond)

			rs.artifacts = artifacts

			rsChan <- rs
			fmt.Println("rs sent")
		}(pj.ID)

	}
}
