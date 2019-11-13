// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package jobrunner

import (
	"fmt"

	"github.com/swinslow/peridot-jobrunner/pkg/agent"
)

// JobRequest defines the metadata needed to start a Job.
type JobRequest struct {
	// requested job ID
	JobID uint32

	// AgentID identifies the Agent that is (or was, or will be) running
	// this Job.
	AgentID uint32

	// Cfg describes the configuration for this Job.
	Cfg agent.JobConfig
}

// String provides a compact string representation of the JobRequest.
func (jreq *JobRequest) String() string {
	return fmt.Sprintf("JobRequest{JobID: %d, AgentID: %d, Cfg: %s}", jreq.JobID, jreq.AgentID, jreq.Cfg.String())
}

// JobUpdate defines the messages that a runJob goroutine sends to the
// rc channel.
type JobUpdate struct {
	// JobID is the unique ID for this Job. It should be unique across
	// all Jobs in peridot.
	JobID uint32

	// Status defines the current status of this Job.
	StatusRpt agent.StatusReport

	// Err defines any error messages that have arisen on the controller
	// for this Job. (Agent errors will be found in Status.OutputMessages.)
	Err error
}
