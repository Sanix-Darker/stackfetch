package services

import (
	"net"
	"strconv"
	"testing"
	"time"
)

func TestCheckPorts(t *testing.T) {
	// Listen on an ephemeral port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)

	svcName := "testsvc"
	DefaultPorts[svcName] = []int{port}

	statuses := CheckPorts([]string{svcName}, 500*time.Millisecond)
	if len(statuses) != 1 {
		t.Fatalf("expected 1 PortStatus, got %d", len(statuses))
	}
	ps := statuses[0]
	if ps.Service != svcName || ps.Port != port {
		t.Errorf("unexpected status entry: %+v", ps)
	}
	if !ps.Open {
		t.Errorf("expected port %d to be open", port)
	}
}
