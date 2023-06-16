package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getNodeName() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return name, nil
}

type ansibleStatus struct {
	Status string `json:"status"`
}

type provisionStatus struct {
	Status string `json:"status"`
}

func getAnsibleStatus() (string, error) {
	var statusData ansibleStatus
	jsonFile, err := os.Open("/etc/ansible/facts.d/ansible_status.fact")
	if err != nil {
		return "", err
	}
	defer jsonFile.Close()
	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return "", err
	}
	json.Unmarshal(jsonBytes, &statusData)
	return statusData.Status, nil
}

func getProvisionStatus() (string, error) {
	var statusData provisionStatus
	jsonFile, err := os.Open("/etc/ansible/facts.d/provision_status.fact")
	if err != nil {
		return "", err
	}
	defer jsonFile.Close()
	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return "", err
	}
	json.Unmarshal(jsonBytes, &statusData)
	return statusData.Status, nil
}

type NodeStatus struct {
	Name            string `json:"name,omitempty"`
	AnsibleStatus   string `json:"ansible_status,omitempty"`
	ProvisionStatus string `json:"provision_status,omitempty"`
}

func NewNodeStatus() (*NodeStatus, error) {
	nodeName, err := getNodeName()
	if err != nil {
		return nil, err
	}
	ansibleStatus, err := getAnsibleStatus()
	if err != nil {
		return nil, err
	}
	provisionStatus, err := getProvisionStatus()
	if err != nil {
		return nil, err
	}
	return &NodeStatus{
		Name:            nodeName,
		AnsibleStatus:   ansibleStatus,
		ProvisionStatus: provisionStatus,
	}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	status, err := NewNodeStatus()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	jsonStatus, err := json.Marshal(status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(jsonStatus)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
