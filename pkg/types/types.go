package types

type ClusterState struct {
	ApiServer string `json:"apiServer"`
}

type Cluster struct {
	Id                int          `json:"id"`
	CreatedAt         string       `json:"created_at"`
	Name              string       `json:"name"`
	Status            string       `json:"status"`
	DesiredState      ClusterState `json:"desired_state"`
	AppliedState      ClusterState `json:"applied_state"`
	ActiveDeployments string       `json:"active_deployments"`
}

type Deployment struct {
	Id                int         `json:"id"`
	CreatedAt         string      `json:"created_at"`
	Cluster           int         `json:"cluster"`
	Name              string      `json:"name"`
	DeploymentType    string      `json:"deployment_type"`
	DeploymentVersion string      `json:"deployment_version"`
	Status            string      `json:"status"`
	CustomYaml        string      `json:"custom_yaml"`
	DesiredState      interface{} `json:"desired_state"`
	ActualState       interface{} `json:"actual_state"`
}

type DeploymentUpdate struct {
	Name         string      `json:"name"`
	CustomYaml   string      `json:"custom_yaml"`
	DesiredState interface{} `json:"desired_state"`
}
