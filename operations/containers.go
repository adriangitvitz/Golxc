package operations

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go_lxc_socket/utils"
	"log"
	"net"
)

type Updatestate struct {
	Action   string `json:"action"`
	Force    bool   `json:"force"`
	Stateful bool   `json:"stateful"`
	Timeout  int    `json:"timeout"`
}

type Containers struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Source      Source            `json:"source"`
	Description string            `json:"description,omitempty"`
	Devices     map[string]Device `json:"devices"`
	Config      map[string]string `json:"config"`
	Ephemeral   bool              `json:"ephemeral,omitempty"`
	Stateful    bool              `json:"stateful,omitempty"`
}

type Device struct {
	Type string `json:"type"`
	Pool string `json:"pool"`
	Path string `json:"path"`
}

type Source struct {
	Type  string `json:"type"`
	Alias string `json:"alias"`
}

func (c Containers) UpdateStatus(instance string) (string, error) {
	conn, err := net.Dial("unix", "/var/snap/lxd/common/lxd/unix.socket")
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer conn.Close()
	container_status := Updatestate{
		Action:   "start",
		Force:    false,
		Stateful: false,
		Timeout:  30,
	}
	container_update_status, err := json.Marshal(container_status)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	uri := fmt.Sprintf("PUT %s/state HTTP/1.0\r\nContent-Type: application/json\r\nContent-Length: %d\r\n\r\n%s", instance, len(container_update_status), container_update_status)
	_, err = conn.Write([]byte(uri))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		if line == "\r\n" {
			break
		}
	}
	body, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	var response_op Response
	error := json.Unmarshal([]byte(body), &response_op)
	if error != nil {
		log.Fatal(err)
		return "", error
	}
	return response_op.Operation, nil
}

func (c Containers) GetContainers() (InstancesResponse, error) {
	conn, err := net.Dial("unix", "/var/snap/lxd/common/lxd/unix.socket")
	if err != nil {
		log.Fatal(err)
		return InstancesResponse{}, err
	}
	defer conn.Close()
	uri := fmt.Sprintf("GET /1.0/instances HTTP/1.0\r\n\r\n")
	_, err = conn.Write([]byte(uri))
	if err != nil {
		log.Fatal(err)
		return InstancesResponse{}, err
	}
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			return InstancesResponse{}, err
		}
		if line == "\r\n" {
			break
		}
	}
	body, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
		return InstancesResponse{}, err
	}
	var response_op InstancesResponse
	error := json.Unmarshal([]byte(body), &response_op)
	if error != nil {
		log.Fatal(err)
		return InstancesResponse{}, err
	}
	return response_op, nil
}

// TODO: Get Container
func (c Containers) GetContainer() {
	conn, err := net.Dial("unix", "/var/snap/lxd/common/lxd/unix.socket")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	uri := fmt.Sprintf("GET /1.0/instances/%s HTTP/1.0\r\n\r\n", c.Name)
	_, err = conn.Write([]byte(uri))
	if err != nil {
		log.Fatal(err)
		return
	}
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			return
		}
		if line == "\r\n" {
			break
		}
	}
	body, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
		return
	}
	pretty, err := utils.Prettyprint(body)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(pretty)
}

// TODO: Delete Container
// TODO: Edit Container config

func (c Containers) Createcontainer() {
	conn, err := net.Dial("unix", "/var/snap/lxd/common/lxd/unix.socket")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	container_payload, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = conn.Write([]byte(fmt.Sprintf("POST /1.0/instances HTTP/1.0\r\nContent-Type: application/json\r\nContent-Length: %d\r\n\r\n%s", len(container_payload), container_payload)))
	if err != nil {
		log.Fatal(err)
		return
	}
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			return
		}
		if line == "\r\n" {
			break
		}
	}
	body, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
		return
	}
	var response_op Response
	error := json.Unmarshal([]byte(body), &response_op)
	if error != nil {
		log.Fatal(err)
		return
	}
	error = Waitforstatus(response_op.Operation)
	if error != nil {
		log.Fatal(err)
		return
	}
	log.Println("Container created")
	operation_status, error := c.UpdateStatus(response_op.Metadata.Resources.Instances[0])
	if error != nil {
		log.Fatal(error)
		return
	}
	error = Waitforstatus(operation_status)
	if error != nil {
		log.Fatal(error)
	}
	log.Println("Container started")
}
