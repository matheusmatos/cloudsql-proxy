package gcloud_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/v2/internal/gcloud"
	exec "golang.org/x/sys/execabs"
)

func TestGcloud(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping gcloud integration tests")
	}

	// The following configures gcloud using only GOOGLE_APPLICATION_CREDENTIALS.
	dir, err := ioutil.TempDir("", "cloudsdk*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	os.Setenv("CLOUDSDK_CONFIG", dir)
	defer os.Unsetenv("CLOUDSDK_CONFIG")

	gcloudCmd, err := gcloud.Cmd()
	if err != nil {
		t.Fatal(err)
	}

	keyFile, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if !ok {
		t.Fatal("GOOGLE_APPLICATION_CREDENTIALS is not set in the environment")
	}

	buf := &bytes.Buffer{}
	cmd := exec.Command(gcloudCmd, "auth", "activate-service-account", "--key-file", keyFile)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to active service account. err = %v, message = %v", err, buf.String())
	}

	// gcloud is now configured. Try to obtain a token from gcloud config
	// helper.
	ts, err := gcloud.GcloudTokenSource(context.Background())
	if err != nil {
		t.Fatalf("failed to get token source: %v", err)
	}

	_, err = ts.Token()
	if err != nil {
		t.Fatalf("failed to get token: %v", err)
	}
}
