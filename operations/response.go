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
