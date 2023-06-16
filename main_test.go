package main

import (
	"os"
	"testing"
)

func TestGetNodeName(t *testing.T) {
	expected, err := os.Hostname()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	actual, err := getNodeName()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if actual != expected {
		t.Errorf("expected %q, but got %q", expected, actual)
	}
}

func TestGetAnsibleStatus(t *testing.T) {
	jsonData := `{"status": "ok"}`
	jsonFile := createTempFile(t, jsonData)
	defer os.Remove(jsonFile.Name())

	expected := "ok"
	actual, err := getAnsibleStatus(jsonFile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if actual != expected {
		t.Errorf("expected %q, but got %q", expected, actual)
	}
}

func TestGetProvisionStatus(t *testing.T) {
	jsonData := `{"status": "success"}`
	jsonFile := createTempFile(t, jsonData)
	defer os.Remove(jsonFile.Name())

	expected := "success"
	actual, err := getProvisionStatus(jsonFile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if actual != expected {
		t.Errorf("expected %q, but got %q", expected, actual)
	}
}

func TestNewNodeStatusRequest(t *testing.T) {
	expected := &NodeStatusRequest{
		name:                "test",
		ansibleStatusPath:   "/etc/ansible/facts.d/ansible_status.fact",
		provisionStatusPath: "/etc/ansible/facts.d/provision_status.fact",
	}
	actual := NewNodeStatusRequest("test")

	if actual.name != expected.name {
		t.Errorf("expected %q, but got %q", expected.name, actual.name)
	}

	if actual.ansibleStatusPath != expected.ansibleStatusPath {
		t.Errorf("expected %q, but got %q", expected.ansibleStatusPath, actual.ansibleStatusPath)
	}

	if actual.provisionStatusPath != expected.provisionStatusPath {
		t.Errorf("expected %q, but got %q", expected.provisionStatusPath, actual.provisionStatusPath)
	}
}

func TestNewNodeStatusRequestWithOptions(t *testing.T) {
	expected := &NodeStatusRequest{
		ansibleStatusPath:   "/tmp/ansible_status.fact",
		provisionStatusPath: "/tmp/provision_status.fact",
	}
	actual := NewNodeStatusRequest("test", WithAnsibleFactPath("/tmp/ansible_status.fact"), WithProvisionFactPath("/tmp/provision_status.fact"))

	if actual.ansibleStatusPath != expected.ansibleStatusPath {
		t.Errorf("expected %q, but got %q", expected.ansibleStatusPath, actual.ansibleStatusPath)
	}

	if actual.provisionStatusPath != expected.provisionStatusPath {
		t.Errorf("expected %q, but got %q", expected.provisionStatusPath, actual.provisionStatusPath)
	}
}

func TestNewNodeStatusResponse(t *testing.T) {
	ansibleJsonData := `{"status": "ok"}`
	ansibleJsonFile := createTempFile(t, ansibleJsonData)
	defer os.Remove(ansibleJsonFile.Name())

	provisionJsonData := `{"status": "success"}`
	provisionJsonFile := createTempFile(t, provisionJsonData)
	defer os.Remove(provisionJsonFile.Name())

	expected := &NodeStatusResponse{
		Name:            "test",
		AnsibleStatus:   "ok",
		ProvisionStatus: "success",
	}
	nsr := NewNodeStatusRequest("test", WithAnsibleFactPath(ansibleJsonFile.Name()), WithProvisionFactPath(provisionJsonFile.Name()))
	actual, err := NewNodeStatus(nsr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if actual.Name != expected.Name {
		t.Errorf("expected %q, but got %q", expected.Name, actual.Name)
	}

	if actual.AnsibleStatus != expected.AnsibleStatus {
		t.Errorf("expected %q, but got %q", expected.AnsibleStatus, actual.AnsibleStatus)
	}

	if actual.ProvisionStatus != expected.ProvisionStatus {
		t.Errorf("expected %q, but got %q", expected.ProvisionStatus, actual.ProvisionStatus)
	}
}

// func TestHandler(t *testing.T) {
// 	expectedName, err := os.Hostname()
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	ansibleJsonData := `{"status": "ok"}`
// 	ansibleFile := createTempFile(t, ansibleJsonData)
// 	defer os.Remove(ansibleFile.Name())

// 	provisionJsonData := `{"status": "success"}`
// 	provisionFile := createTempFile(t, provisionJsonData)
// 	defer os.Remove(provisionFile.Name())

// 	expected := &NodeStatusResponse{
// 		Name:            expectedName,
// 		AnsibleStatus:   "ok",
// 		ProvisionStatus: "success",
// 	}

// 	req, err := http.NewRequest("GET", "/", nil)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(handler)

// 	handler.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}

// 	expectedJSON, err := json.Marshal(expected)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	if strings.TrimSpace(rr.Body.String()) != string(expectedJSON) {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), string(expectedJSON))
// 	}
// }

func createTempFile(t *testing.T, jsonData string) *os.File {
	t.Helper()

	file, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := file.WriteString(jsonData); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return file
}
