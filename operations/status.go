package operations

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func Getinstancestatus(operation string) (string, error) {
	conn, err := net.Dial("unix", "/var/snap/lxd/common/lxd/unix.socket")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(fmt.Sprintf("GET %s/wait?timeout=3 HTTP/1.0\r\n\r\n", operation)))
	if err != nil {
		return "", err
	}
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
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

func Waitforstatus(operation string) error {
	statusch := make(chan string)
	go func(operation string) {
		for {
			status, err := Getinstancestatus(operation)
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
		return fmt.Errorf("Error getting instance information")
	}
	return nil
}
