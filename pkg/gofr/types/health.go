package types

type Health struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Host     string `json:"host,omitempty"` // host and database name is omitempty since it is not applicable for a service
	Database string `json:"database,omitempty"`
}

type AppDetails struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Framework string `json:"framework"`
}
