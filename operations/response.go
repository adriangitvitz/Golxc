package operations

// Container creation and update state response
type Response struct {
	Status    string   `json:"status"`
	Metadata  Metadata `json:"metadata"`
	Operation string   `json:"operation"`
}

type Metadata struct {
	Id        string    `json:"id"`
	Resources Resources `json:"resources"`
	Status    string    `json:"status"`
}

type Resources struct {
	Containers []string `json:"containers"`
	Instances  []string `json:"instances"`
}

// Get instances response
type InstancesResponse struct {
	Status   string   `json:"status"`
	Metadata []string `json:"metadata"`
}

// Get Single Instance Response
type SingleInstanceResponse struct {
	Metadata SingleMetadata `json:"metadata"`
}

type SingleMetadata struct {
	Name           string                     `json:"name"`
	Status         string                     `json:"status"`
	Type           string                     `json:"type"`
	Devices        map[string]ExpandedDevices `json:"expanded_devices"`
	ExpandedConfig map[string]string          `json:"expanded_config"`
}

type ExpandedDevices struct {
	Name    string `json:"name,omitempty"`
	Network string `json:"network,omitempty"`
	Type    string `json:"type,omitempty"`
	Path    string `json:"path,omitempty"`
	Pool    string `json:"pool,omitempty"`
}
