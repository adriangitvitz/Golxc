package main

import (
	"go_lxc_socket/operations"
)

func main() {
	c := operations.Containers{
		Name:        "develop",
		Type:        "container",
		Description: "Testing",
		Source: operations.Source{
			Type:  "image",
			Alias: "ubuntu",
		},
		Devices: map[string]operations.Device{
			"root": {
				Type: "disk",
				Pool: "lxpool",
				Path: "/",
			},
		},
		Config: map[string]string{
			"limits.memory": "4GB",
			"limits.cpu":    "4",
		},
	}
	c.Createcontainer()
}
