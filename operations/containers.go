package operations

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
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

func (c Containers) getinstancestatus(operation string) (string, error) {
	connection, err := net.Dial("unix", "/var/snap/lxd/common/lxd/unix.socket")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer connection.Close()
	_, err = connection.Write([]byte(fmt.Sprintf("GET %s/wait?timeout=3 HTTP/1.0\r\n\r\n", operation)))
	if err != nil {
		return "", err
	}
	// Used to read HTTP response
	reader := bufio.NewReader(connection)
	for {
		// Read until new line
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		// Read until character return and new line
		if line == "\r\n" {
			break
		}
	}
	body, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	var data map[string]interface{}
	error := json.Unmarshal([]byte(body), &data)
	if error != nil {
		return "", error
	}
	return data["status"].(string), nil
}

func (c Containers) waitforstatus(operation string) error {
	statusch := make(chan string)
	go func(operation string) {
		for {
			status, err := c.getinstancestatus(operation)
			if err != nil {
				statusch <- err.Error()
				return
			}
			if status == "Success" {
				statusch <- "Ready"
				return
			}
			time.Sleep(1 * time.Second)
		}
	}(operation)
	statusc := <-statusch
	if statusc != "Ready" {
		return fmt.Errorf("Error getting %s instance information", c.Name)
	}
	return nil
}

func (c Containers) UpdateStatus(instance string) (string, error) {
	conn, err := net.Dial("unix", "/var/snap/lxd/common/lxd/unix.socket")
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return "", err
	}
	uri := fmt.Sprintf("PUT %s/state HTTP/1.0\r\nContent-Type: application/json\r\nContent-Length: %d\r\n\r\n%s", instance, len(container_update_status), container_update_status)
	_, err = conn.Write([]byte(uri))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		if line == "\r\n" {
			break
		}
	}
	body, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	var response_op Response
	error := json.Unmarshal([]byte(body), &response_op)
	if error != nil {
		fmt.Println(error)
		return "", error
	}
	return response_op.Operation, nil
}

func (c Containers) Createcontainer() {
	conn, err := net.Dial("unix", "/var/snap/lxd/common/lxd/unix.socket")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	container_payload, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = conn.Write([]byte(fmt.Sprintf("POST /1.0/instances HTTP/1.0\r\nContent-Type: application/json\r\nContent-Length: %d\r\n\r\n%s", len(container_payload), container_payload)))
	if err != nil {
		fmt.Println(err)
		return
	}
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		if line == "\r\n" {
			break
		}
	}
	body, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	var response_op Response
	error := json.Unmarshal([]byte(body), &response_op)
	if error != nil {
		fmt.Println(error)
		return
	}
	error = c.waitforstatus(response_op.Operation)
	if error != nil {
		fmt.Println(error)
		return
	}
	operation_status, error := c.UpdateStatus(response_op.Metadata.Resources.Instances[0])
	if error != nil {
		fmt.Println(error)
		return
	}
	error = c.waitforstatus(operation_status)
	if error != nil {
		fmt.Println(error)
	}
}
