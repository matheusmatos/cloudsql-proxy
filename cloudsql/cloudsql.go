package cloudsql

import (
	"context"
	"io"
	"net"

	"cloud.google.com/go/cloudsqlconn"
)

// Dialer dials a Cloud SQL instance and returns its database engine version.
type Dialer interface {
	// Dial returns a connection to the specified instance.
	Dial(ctx context.Context, inst string, opts ...cloudsqlconn.DialOption) (net.Conn, error)
	// EngineVersion retrieves the provided instance's database version (e.g.,
	// POSTGRES_14)
	EngineVersion(ctx context.Context, inst string) (string, error)

	io.Closer
}
