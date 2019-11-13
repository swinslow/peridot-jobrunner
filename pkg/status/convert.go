// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package status

import (
	"fmt"

	"github.com/swinslow/peridot-db/pkg/datastore"
)

// ConvertStatusProtoToDatastore converts a Status
// protobuf value to its corresponding Datastore version.
func ConvertStatusProtoToDatastore(stProto Status) (datastore.Status, error) {
	switch stProto {
	case Status_STARTUP:
		return datastore.StatusStartup, nil
	case Status_RUNNING:
		return datastore.StatusRunning, nil
	case Status_STOPPED:
		return datastore.StatusStopped, nil
	case Status_STATUS_SAME:
		return datastore.StatusSame, nil
	default:
		return datastore.StatusSame, fmt.Errorf("unknown value for status: %v", stProto)
	}
}

// ConvertStatusDatastoreToProto converts a Status
// Datastore value to its corresponding protobuf version.
func ConvertStatusDatastoreToProto(st datastore.Status) (Status, error) {
	switch st {
	case datastore.StatusStartup:
		return Status_STARTUP, nil
	case datastore.StatusRunning:
		return Status_RUNNING, nil
	case datastore.StatusStopped:
		return Status_STOPPED, nil
	case datastore.StatusSame:
		return Status_STATUS_SAME, nil
	default:
		return Status_STATUS_SAME, fmt.Errorf("unknown value for status: %v", st)
	}
}

// ConvertHealthProtoToDatastore converts a Health
// protobuf value to its corresponding Datastore version.
func ConvertHealthProtoToDatastore(hProto Health) (datastore.Health, error) {
	switch hProto {
	case Health_OK:
		return datastore.HealthOK, nil
	case Health_DEGRADED:
		return datastore.HealthDegraded, nil
	case Health_ERROR:
		return datastore.HealthError, nil
	case Health_HEALTH_SAME:
		return datastore.HealthSame, nil
	default:
		return datastore.HealthSame, fmt.Errorf("unknown value for health: %v", hProto)
	}
}

// ConvertHealthDatastoreToProto converts a Health
// Datastore value to its corresponding protobuf version.
func ConvertHealthDatastoreToProto(h datastore.Health) (Health, error) {
	switch h {
	case datastore.HealthOK:
		return Health_OK, nil
	case datastore.HealthDegraded:
		return Health_DEGRADED, nil
	case datastore.HealthError:
		return Health_ERROR, nil
	case datastore.HealthSame:
		return Health_HEALTH_SAME, nil
	default:
		return Health_HEALTH_SAME, fmt.Errorf("unknown value for health: %v", h)
	}
}
