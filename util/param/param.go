package param

type SuiteParams struct {
	BuildID, ArtifactID, SourceID, BuildJobID, ProjectID, GitCredentialID, DockerCredentialID, UserGroupID, EnvID, AppID, DeployID, DeployJobID int
	PipelineID, DuplicatePipID, TriggerID, PipelineJobID, UnitTestID, CheckpointID, VMGroupID, DuplicateBuildID                                 int
	CommitID, UserID, AppName, ImageTag                                                                                                         string
}

func NewParamSet() *SuiteParams {
	return &SuiteParams{}
}
