package commands

const AUTH_CONFIG_PATH = "~/.sextant/config.json"

type AuthCheckResponse struct {
	Username string `json:"username"`
}

type Remote struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Token string `json:"token"`
}

type CLIConfig struct {
	ActiveRemote string   `json:"active_remote"`
	Remotes      []Remote `json:"remotes"`
}
