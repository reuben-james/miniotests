package miniotests

import (
	"flag"
	"os"
	"testing"
)

var useTLS = flag.Bool("useTLS", false, "Set to true to use secure (TLS) connection to MinIO server")

func TestMain(m *testing.M) {
    flag.Parse()
    exitCode := m.Run()
    os.Exit(exitCode)
}