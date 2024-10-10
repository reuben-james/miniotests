package miniotests

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
)

var useTLS = flag.Bool("useTLS", false, "Set to true to use secure (TLS) connection to MinIO server")

type Config struct {
    Server      string
    Port        string
    AccessKey   string
    SecretKey   string
    Secure      bool
    Endpoint    string
    Transport   *http.Transport
    Scheme      string
}

var TestConfig Config

func TestMain(m *testing.M) {
    flag.Parse()

    err := initializeConfig()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error initializing configuration: %v\n", err)
        os.Exit(1)
    }

    // Run the tests
    exitCode := m.Run()

    os.Exit(exitCode)
}

func initializeConfig() error {
    // Retrieve environment variables 
    TestConfig.Server = os.Getenv("MINIO_SERVER")
    TestConfig.Port = os.Getenv("MINIO_PORT")
    TestConfig.AccessKey = os.Getenv("MINIO_ACCESS_KEY")
    TestConfig.SecretKey = os.Getenv("MINIO_SECRET_KEY")
	
	// Required parameters
	missingVars := []string{}
    if TestConfig.Server == "" {
        missingVars = append(missingVars, "MINIO_SERVER")
    }
    if TestConfig.Port == "" {
        missingVars = append(missingVars, "MINIO_PORT")
    }
    if TestConfig.AccessKey == "" {
        missingVars = append(missingVars, "MINIO_ACCESS_KEY")
    }
    if TestConfig.SecretKey == "" {
        missingVars = append(missingVars, "MINIO_SECRET_KEY")
    }

    if len(missingVars) > 0 {
        return fmt.Errorf("missing required environment variables: %v", missingVars)
    }

    TestConfig.Secure = *useTLS
    TestConfig.Scheme = "http"
    if TestConfig.Secure {
        TestConfig.Scheme = "https"
    }
	TestConfig.Endpoint = net.JoinHostPort(TestConfig.Server, TestConfig.Port)
    // Configure Transport for TLS if necessary
    if TestConfig.Secure {
        TestConfig.Transport = &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
    }

    fmt.Printf("#################################\n")
    fmt.Printf("# CONFIG\n")
    fmt.Printf("#################################\n")
	fmt.Printf("MINIO_SERVER: %q\n", TestConfig.Server)
    fmt.Printf("MINIO_PORT: %q\n", TestConfig.Port)
    fmt.Printf("MINIO_ACCESS_KEY is set: %v\n", TestConfig.AccessKey != "")
    fmt.Printf("MINIO_SECRET_KEY is set: %v\n", TestConfig.SecretKey != "")
    fmt.Printf("Secure mode enabled: %v\n", TestConfig.Secure)
    fmt.Printf("Endpoint: %s\n", TestConfig.Endpoint)
    fmt.Printf("#################################\n")

    return nil
}