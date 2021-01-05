package API

type PrometheusData struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Cpu struct {
			Result []struct {
				Values [][]interface{} `json:"values"`
			} `json:"result"`
		} `json:"cpu"`
		Mem struct {
			Result []struct {
				Values [][]interface{} `json:"values"`
			} `json:"result"`
		} `json:"mem"`
		NetworkReceive struct {
			Result []struct {
				Values [][]interface{} `json:"values"`
			} `json:"result"`
		} `json:"networkReceive"`
		NetworkTransmit struct {
			Result []struct {
				Values [][]interface{} `json:"values"`
			} `json:"result"`
		} `json:"networkTransmit"`
	} `json:"data"`
}

type BuildRuntime struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"data"`
}

type PipelineRuntime struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		ID           int `json:"id"`
		StepNodeList []struct {
			Status    string `json:"status"`
			HistoryID int    `json:"historyId"`
		} `json:"stepNodeList"`
	} `json:"data"`
}

type K8sDescription struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      string `json:"data"`
}

type PodList struct {
	AppId            int `json:"appId"`
	KubernetesDeploy []struct {
		PodInstances []struct {
			PodName  string `json:"podName"`
			NodeName string `json:"nodeName"`
			Status   string `json:"status"`
		} `json:"podInstances"`
	} `json:"kubernetesDeploy"`
}

type BuildHistoryList struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data []struct {
			JobID  int    `json:"id"`
			Type   string `json:"type"`
			Status string `json:"status"`
		} `json:"data"`
	} `json:"data"`
}

type PipelineHistory struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Status string `json:"status"`
	} `json:"data"`
}

type PipelineRan struct {
	ErrorCode     int    `json:"errorCode"`
	ErrorInfo     string `json:"errorInfo"`
	PipelineJobID int    `json:"data"`
}

type PipelineDetail struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		TriggerMode   `json:"triggerMode"`
		PipelineSteps []struct {
			ID int `json:"ID"`
		} `json:"pipelineSteps"`
	} `json:"data"`
}

type TriggerMode struct {
	ID int `json:"id"`
}

type PipelineCreated struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		ID int `json:"id"`
	} `json:"data"`
}

type UserGroups struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data []UserGroup `json:"data"`
	} `json:"data"`
}

type UserGroup struct {
	ID      int    `json:"id"`
	Creator string `json:"creator"`
}

type FileUploaded struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		File string `json:"file"`
	} `json:"data"`
}

type DeployConfig struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Service Service `json:"service"`
		Grpc    struct {
			Grpc []struct {
				GrpcName    string `json:"grpcName"`
				ServicePort int    `json:"servicePort"`
			} `json:"grpc"`
		} `json:"grpc"`
	} `json:"data"`
}

type Service struct {
	Service []Svc `json:"service"`
}

type Svc struct {
	ServiceType string `json:"serviceType"`
	NodePort    int    `json:"nodePort"`
}

type DeployHistoryList struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data []struct {
			JobID  int    `json:"id"`
			Type   string `json:"type"`
			Status string `json:"status"`
		} `json:"data"`
	} `json:"data"`
}

type DeployHistory struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Status string `json:"status"`
	} `json:"data"`
}

type Deploys struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data []Deploy `json:"data"`
	} `json:"data"`
}

type Deploy struct {
	ID int `json:"id"`
}

type Apps struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data []App `json:"data"`
	} `json:"data"`
}

type App struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Environments struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data []Env `json:"data"`
	} `json:"data"`
}

type Env struct {
	ID      int    `json:"id"`
	Kubectx string `json:"kubectx"`
}

type Builds struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data []Build `json:"data"`
	} `json:"data"`
}

type Build struct {
	ID int `json:"id"`
}

type Artifacts struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data []Artifact `json:"data"`
	} `json:"data"`
}

type Artifact struct {
	ID              int    `json:"id"`
	DownloadAddress string `json:"downloadAddress"`
}

type ArtifactDetail struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		DownloadAddress string `json:"downloadAddress"`
	} `json:"data"`
}

type DockerCredentials struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data []GitCredential `json:"data"`
	} `json:"data"`
}

type DockerCredential struct {
	ID int `json:"id"`
}

type GitCredentials struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data []GitCredential `json:"data"`
	} `json:"data"`
}

type GitCredential struct {
	ID int `json:"id"`
}

type Projects struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data  []Project `json:"data"`
		Total int       `json:"total"`
	} `json:"data"`
}

type Project struct {
	ID              int    `json:"id"`
	CreatorUserName string `json:"creatorUserName"`
}

type ProjectCreated struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
}

type VMMachines struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data []VMMachine `json:"data"`
	} `json:"data"`
}

type VMMachine struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ArtifactLibraryCreated struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      int    `json:"data"`
}

type BuildCreated struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      int    `json:"data"`
}

type ArtifactObtained struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data []Artifact `json:"data"`
	} `json:"data"`
}

type BuildDetailObtained struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Source `json:"source"`
	} `json:"data"`
}

type BuildEdited struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
}

type Source struct {
	ID int `json:"id"`
}

type BeforeBuild struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		SpecificBranchOrTag `json:"specificBranchOrTag"`
	} `json:"data"`
}

type SpecificBranchOrTag struct {
	CommitID string
}

type BuildRan struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      int    `json:"data"`
}

type BuildRecordObtained struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Status string `json:"status"`
		Type   string `json:"type"`
	} `json:"data"`
}

type ImagePage struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		PageNum  int     `json:"pageNum"`
		PageSize int     `json:"pageSize"`
		Pages    int     `json:"pages"`
		Total    int     `json:"total"`
		Data     []Image `json:"data"`
	} `json:"data"`
}

type Image struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type UploadRecords struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data []SingleUploadRecord `json:"data"`
	} `json:"data"`
}

type SingleUploadRecord struct {
	ID         int    `json:"id"`
	Version    string `json:"version"`
	UploadType int    `json:"uploadType"`
}

type DownloadRecords struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data []SingleDownloadRecord `json:"data"`
	} `json:"data"`
}

type SingleDownloadRecord struct {
	ID      int    `json:"id"`
	Version string `json:"version"`
}

type ItemDeleted struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
}

type GeneralResp struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      int    `json:"data"`
}

type StringData struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      string `json:"data"`
}

type Tests struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data  []Test `json:"data"`
		Total int    `json:"total"`
	} `json:"data"`
}

type Test struct {
	ID int `json:"id"`
}

type Pipelines struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data  []Pipeline `json:"data"`
		Total int        `json:"total"`
	} `json:"data"`
}

type Pipeline struct {
	ID int `json:"id"`
}

type TestDetail struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Id int `json:"id"`
	} `json:"data"`
}

type TestConfig struct {
	Data struct {
		TestID       int            `json:"testID"`
		ContainerLog ContainerLog   `json:"containerLog,omitempty"`
		HostMap      HostMap        `json:"hostMap,omitempty"`
		KubeConfig   TestKubeConfig `json:"kubeConfig,omitempty"`
		VarMap       VarMap         `json:"varMap,omitempty"`
	} `json:"data"`
}

type ContainerLog struct {
	LogPath   string `json:"logPath"`
	IsChanged bool   `json:"isChanged"`
}

type HostMap struct {
	Host string `json:"devops.test.cq.iot.chinamobile.com,omitempty"`
}

type TestKubeConfig struct {
	Cpu       string `json:"cpu,omitempty"`
	Mem       string `json:"mem,omitempty"`
	Id        int    `json:"id,omitempty"`
	IsChanged bool   `json:"isChanged,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	ObjectId  int    `json:"objectId,omitempty"`
	Duplicate int    `json:"duplicate,omitempty"`
}

type VarMap struct {
	CONFIGFILE string `json:"CONFIGFILE,omitempty"`
	SUITE      string `json:"SUITE,omitempty"`
}

type TestParams struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      []struct {
		ParamValues string `json:"paramValues"`
	} `json:"data"`
}

type TestStatusList struct {
	ErrorCode int      `json:"errorCode"`
	ErrorInfo string   `json:"errorInfo"`
	Data      []string `json:"data"`
}

type TestRecords struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data []TestRecord `json:"data"`
	} `json:"data"`
}

type TestRecord struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
}

type TestPods struct {
	ErrorCode int       `json:"errorCode"`
	ErrorInfo string    `json:"errorInfo"`
	Data      []TestPod `json:"data"`
}

type TestPod struct {
	Containers []struct {
		ContainerName string `json:"containerName"`
		Image         string `json:"image"`
	} `json:"containers"`
	EnvId     int    `json:"envId"`
	Namespace string `json:"namespace"`
	PodName   string `json:"podName"`
}

type VMWorkOrders struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data  []VMWorkOrder `json:"data"`
		Total int           `json:"total"`
	} `json:"data"`
}

type VMWorkOrder struct {
	Status     string `json:"status"`
	DetailLink string `json:"detailLink"`
	Id         int    `json:"id"`
}

type VMGroups struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Data []VMGroup `json:"data"`
	} `json:"data"`
}

type VMGroup struct {
	ID     int `json:"id"`
	RoomId int `json:"roomId"`
	ZoneId int `json:"zoneId"`
}

type VMGroupDetail struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		ID     int `json:"id"`
		RoomId int `json:"roomId"`
		ZoneId int `json:"zoneId"`
	} `json:"data"`
}

type K8SEnvStatus struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		AccessTime int    `json:"accessTime"`
		Status     string `json:"status"`
	} `json:"data"`
}

type ProjectDetail struct {
	ErrorCode int    `json:"errorCode"`
	ErrorInfo string `json:"errorInfo"`
	Data      struct {
		Description string `json:"description"`
		Name        string `json:"name"`
	} `json:"data"`
}
