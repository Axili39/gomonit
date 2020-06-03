package gomonit

import (
	"os"
	"testing"
)

func TestUnMarshal(t *testing.T) {
	var status MonitStatus
	file, err := os.Open("tests/status-test.xml")
	if err != nil {
		t.Error("Error opening xml data file")
	}
	err = status.Load(file)
	if err != nil {
		t.Error("Error Unmarshaling xml data")
	}
	// Check content
	if len(status.Services) != 6 {
		t.Errorf("bad services count, got %d/%d", len(status.Services), 6)
	}

	// Check GetService
	service := status.GetService("service1")
	if service == nil {
		t.Error("GetService Failed")
	}

	// Check GetService
	service = status.GetService("unknown")
	if service != nil {
		t.Error("GetService on unknown name Failed")
	}
}

/*
func BenchmarkResolve(b *testing.B) {
	b.ReportAllocs()
	var r Resolver
	r.Init()

	for i := 0; i < b.N; i++ {
		r.Resolve("alice")
	}
}
func BenchmarkReverse(b *testing.B) {
	b.ReportAllocs()
	var r Resolver
	r.Init()
	addr,_ := net.ResolveUDPAddr("udp", "127.0.0.1:8003")
	for i := 0; i < b.N; i++ {
		r.Reverse(addr)
	}
}
*/
