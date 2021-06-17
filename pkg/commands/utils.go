package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/catenasys/sxtctl/pkg/http"
	"github.com/catenasys/sxtctl/pkg/types"
)

func setupApi() (*http.HttpApiHandler, error) {

	accessToken, err := loadAccessToken()
	if err != nil {
		return nil, err
	}
	url, err := loadUrl()
	if err != nil {
		return nil, err
	}
	if accessToken == "" {
		return nil, fmt.Errorf("no access token found")
	}
	if url == "" {
		return nil, fmt.Errorf("no url found")
	}
	apiHandler := http.NewHttpApiHandler(url, accessToken)
	authCheckResponse := &AuthCheckResponse{}

	err = apiHandler.Request("GET", "/api/v1/user/status", nil, authCheckResponse)

	if err != nil {
		return nil, fmt.Errorf("bad credentials %s", err)
	}

	if authCheckResponse.Username == "" {
		return nil, fmt.Errorf("bad credentials")
	}

	return apiHandler, nil
}

func expandFilePath(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}

func validateCLIConfig(config *CLIConfig) error {
	for _, remote := range config.Remotes {
		if !govalidator.IsURL(remote.Url) {
			return fmt.Errorf("Remote %s has an invalid Url '%s'", remote.Name, remote.Url)
		}
	}

	return nil
}

func loadCLIConfig() (*CLIConfig, error) {
	configPath, err := expandFilePath(AUTH_CONFIG_PATH)
	if err != nil {
		return nil, err
	}
	config := &CLIConfig{
		ActiveRemote: "",
		Remotes:      []Remote{},
	}
	if _, err := os.Stat(configPath); err == nil {
		// config does exist
		configData, err := ioutil.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("Error loading config file: %s\n", err)
		}
		err = json.Unmarshal(configData, &config)
		if err != nil {
			return nil, fmt.Errorf("Error parsing config file: %s\n%s\n", err, string(configData))
		}

		err = validateCLIConfig(config)

		if err != nil {
			return nil, err
		}
		return config, nil
	} else if os.IsNotExist(err) {
		// config does not exist
		return config, nil
	} else {
		return nil, fmt.Errorf("Error loading config file: %s\n", err)
	}
}

func saveCLIConfig(config *CLIConfig) error {
	configPath, err := expandFilePath(AUTH_CONFIG_PATH)
	if err != nil {
		return err
	}
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configPath, configBytes, 0644)
}

func loadRemoteName() string {
	if RemoteName != "" {
		return RemoteName
	}
	config, err := loadCLIConfig()
	if err != nil {
		return ""
	}
	return config.ActiveRemote
}

func loadRemote() (*Remote, error) {
	name := loadRemoteName()
	if name == "" {
		return nil, nil
	}
	return loadRemoteByName(name)
}

func loadRemoteByName(name string) (*Remote, error) {
	config, err := loadCLIConfig()
	if err != nil {
		return nil, err
	}
	for _, remote := range config.Remotes {
		if remote.Name == name {
			return &remote, nil
		}
	}
	return nil, nil
}

func loadUrl() (string, error) {
	// always give preference to the configured URL not the config file
	if Url != "" {
		return Url, nil
	}
	remote, err := loadRemote()
	if err != nil {
		return "", err
	}
	if remote == nil {
		return "", nil
	}
	return remote.Url, nil
}

func loadAccessToken() (string, error) {
	// always give preference to the configured URL not the config file
	if AccessToken != "" {
		return AccessToken, nil
	}
	remote, err := loadRemote()
	if err != nil {
		return "", err
	}
	if remote == nil {
		return "", nil
	}
	return remote.Token, nil
}

func loadClusters(api *http.HttpApiHandler) ([]types.Cluster, error) {
	clusters := []types.Cluster{}
	err := api.Request("GET", "/api/v1/clusters", nil, &clusters)
	return clusters, err
}

func loadDeployments(api *http.HttpApiHandler, clusterId int) ([]types.Deployment, error) {
	deployments := []types.Deployment{}
	err := api.Request("GET", fmt.Sprintf("/api/v1/clusters/%d/deployments?showDeleted=y", clusterId), nil, &deployments)
	return deployments, err
}

func getCluster(api *http.HttpApiHandler, nameOrId string) (*types.Cluster, error) {
	clusters, err := loadClusters(api)
	if err != nil {
		return nil, err
	}
	for _, cluster := range clusters {
		if fmt.Sprintf("%d", cluster.Id) == nameOrId || cluster.Name == nameOrId {
			return &cluster, nil
		}
	}
	return nil, nil
}

func getDeployment(api *http.HttpApiHandler, clusterId int, nameOrId string) (*types.Deployment, error) {
	deployments, err := loadDeployments(api, clusterId)
	if err != nil {
		return nil, err
	}
	for _, deployment := range deployments {
		if fmt.Sprintf("%d", deployment.Id) == nameOrId || deployment.Name == nameOrId {
			return &deployment, nil
		}
	}
	return nil, nil
}

func initialiseCluster() (*http.HttpApiHandler, *types.Cluster, error) {
	if ClusterNameOrId == "" {
		return nil, nil, fmt.Errorf("Please provide a --cluster argument")
	}
	api, err := setupApi()
	if err != nil {
		return nil, nil, err
	}
	cluster, err := getCluster(api, ClusterNameOrId)
	if err != nil {
		return nil, nil, err
	}
	if cluster == nil {
		return nil, nil, fmt.Errorf("No cluster found: %s", ClusterNameOrId)
	}
	return api, cluster, nil
}

func initialiseDeployment() (*http.HttpApiHandler, *types.Cluster, *types.Deployment, error) {
	if ClusterNameOrId == "" {
		return nil, nil, nil, fmt.Errorf("Please provide a --cluster argument")
	}
	if DeploymentNameOrId == "" {
		return nil, nil, nil, fmt.Errorf("Please provide a --deployment argument")
	}
	api, err := setupApi()
	if err != nil {
		return nil, nil, nil, err
	}
	cluster, err := getCluster(api, ClusterNameOrId)
	if err != nil {
		return nil, nil, nil, err
	}
	if cluster == nil {
		return nil, nil, nil, fmt.Errorf("No cluster found: %s", ClusterNameOrId)
	}
	deployment, err := getDeployment(api, cluster.Id, DeploymentNameOrId)
	if err != nil {
		return nil, nil, nil, err
	}
	return api, cluster, deployment, nil
}

func checkFormat() error {
	if OutputFormat == "text" || OutputFormat == "json" {
		return nil
	} else {
		return fmt.Errorf("invalid format type %s", OutputFormat)
	}
}

func getJSONString(data interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func undeployDeployment(api *http.HttpApiHandler, deployment *types.Deployment) error {
	return api.Request("DELETE", fmt.Sprintf("/api/v1/clusters/%d/deployments/%d", deployment.Cluster, deployment.Id), nil, nil)
}

func redeployDeployment(api *http.HttpApiHandler, deployment *types.Deployment) error {
	data := &types.DeploymentUpdate{
		Name:         deployment.Name,
		CustomYaml:   deployment.CustomYaml,
		DesiredState: deployment.DesiredState,
	}
	return api.Request("PUT", fmt.Sprintf("/api/v1/clusters/%d/deployments/%d", deployment.Cluster, deployment.Id), data, nil)
}

func waitForDeploymentState(api *http.HttpApiHandler, deployment *types.Deployment, expected string) error {
	okStatus := false
	loops := 0
	for !okStatus && loops < 100 {
		deployment, err := getDeployment(api, deployment.Cluster, fmt.Sprintf("%d", deployment.Id))
		if err != nil {
			return err
		}
		if deployment.Status == expected {
			return nil
		}
		loops++
		fmt.Printf("waiting for deployment status: %s, current status: %s, loops: %d\n", expected, deployment.Status, loops)
		time.Sleep(time.Second)
	}
	if !okStatus {
		return fmt.Errorf("deployment did not reach desired status after 100 seconds")
	}
	return nil
}
