// Package firestore is the alternative reference implementation of
// ports.ExampleRepository, backed by Google Cloud Firestore. Use this
// pattern instead of (or alongside) postgres when the data fits a
// document/collection model better than relational tables.
package firestore

import (
	"context"
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/epa-datos/epa-standards-backend/internal/pkg/config"
	"github.com/sirupsen/logrus"
)

const examplesCollection = "examples"

var (
	client     *firestore.Client
	clientOnce sync.Once
)

// NewClient returns a singleton *firestore.Client for config.Cfg.FirestoreProjectID.
// Authentication relies on Application Default Credentials (ADC): locally run
// `gcloud auth application-default login`; in Cloud Run/GKE it's automatic
// via the service account attached to the runtime.
func NewClient(ctx context.Context) *firestore.Client {
	clientOnce.Do(func() {
		cli, err := firestore.NewClient(ctx, config.Cfg.FirestoreProjectID)
		if err != nil {
			logrus.Fatalf("firestore: failed to create client: %v", err)
		}
		client = cli
	})
	return client
}
