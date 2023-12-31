package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func getNodeName() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}
	shortName := strings.Split(name, ".")[0]
	return shortName, nil
}

type ansibleStatus struct {
	Status string `json:"status"`
}

type provisionStatus struct {
	Status string `json:"status"`
}

func getAnsibleStatus(factPath string) (string, error) {
	var statusData ansibleStatus
	jsonFile, err := os.Open(factPath)
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

func getProvisionStatus(factPath string) (string, error) {
	var statusData provisionStatus
	jsonFile, err := os.Open(factPath)
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

type NodeStatusResponse struct {
	Name            string `json:"name,omitempty"`
	AnsibleStatus   string `json:"ansible_status,omitempty"`
	ProvisionStatus string `json:"provision_status,omitempty"`
}

type NodeStatusRequest struct {
	name                string
	ansibleStatusPath   string
	provisionStatusPath string
}

type NodeStatusRequestOption func(*NodeStatusRequest)

func WithName(name string) NodeStatusRequestOption {
	return func(ns *NodeStatusRequest) {
		ns.name = name
	}
}
func WithAnsibleFactPath(path string) NodeStatusRequestOption {
	return func(ns *NodeStatusRequest) {
		ns.ansibleStatusPath = path
	}
}

func WithProvisionFactPath(path string) NodeStatusRequestOption {
	return func(ns *NodeStatusRequest) {
		ns.provisionStatusPath = path
	}
}

func NewNodeStatusRequest(nodeName string, opts ...NodeStatusRequestOption) *NodeStatusRequest {
	const (
		defaultAnsibleFactPath   = "/etc/ansible/facts.d/ansible_status.fact"
		defaultProvisionFactPath = "/etc/ansible/facts.d/provision_status.fact"
	)
	nsr := &NodeStatusRequest{
		name:                nodeName,
		ansibleStatusPath:   defaultAnsibleFactPath,
		provisionStatusPath: defaultProvisionFactPath,
	}
	for _, opt := range opts {
		opt(nsr)
	}
	return nsr

}

func NewNodeStatus(nsr *NodeStatusRequest) (*NodeStatusResponse, error) {
	ansibleStatus, err := getAnsibleStatus(nsr.ansibleStatusPath)
	if err != nil {
		return nil, err
	}
	provisionStatus, err := getProvisionStatus(nsr.provisionStatusPath)
	if err != nil {
		return nil, err
	}
	return &NodeStatusResponse{
		Name:            nsr.name,
		AnsibleStatus:   ansibleStatus,
		ProvisionStatus: provisionStatus,
	}, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	nodeName, err := getNodeName()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	nsr := NewNodeStatusRequest(nodeName)
	status, err := NewNodeStatus(nsr)
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

func f5Handler(w http.ResponseWriter, r *http.Request) {
	nodeName, err := getNodeName()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	nsr := NewNodeStatusRequest(nodeName)
	status, err := NewNodeStatus(nsr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if status.ProvisionStatus == "success" {
		w.Write([]byte("up"))
		return
	}
	w.Write([]byte("down"))
}

func main() {
	host, found := os.LookupEnv("NODE_STATUS_SERVER_HOST")
	if !found {
		host = "0.0.0.0"
	}
	port, found := os.LookupEnv("NODE_STATUS_SERVER_PORT")
	if !found {
		port = "30622"
	}
	connStr := fmt.Sprintf("%s:%s", host, port)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/f5", f5Handler)
	fmt.Printf("Server is running at http://%s\n", connStr)
	log.Fatal(http.ListenAndServe(connStr, nil))
}
