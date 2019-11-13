// Package jobrunner is the main Job runner for peridot.
// It operates as a set of gRPC clients, with each Agent separately
// running its own gRPC server.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package jobrunner

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/swinslow/peridot-jobrunner/pkg/status"
)

// JobController is the main Job runner function. It creates and returns
// three channels (described from the caller's perspective):
// * inJobStream, a write-only channel to submit new JobRequests, which must
//   be closed by the caller
// * errc, a read-only channel where an error will be written or else
//   nil if no errors in the controller itself are encountered.
func (env *Env) JobController(ctx context.Context) (chan<- JobRequest, <-chan error) {
	// the caller will own the inJobStream channel and must close it
	inJobStream := make(chan JobRequest)
	// we own the errc channel. make it buffered so we can write 1 error
	// without blocking.
	errc := make(chan error, 1)

	// rc is the response channel for all Job status messages.
	rc := make(chan JobUpdate)

	// n is the WaitGroup used to synchronize agent completion.
	// Each runJob goroutine adds 1 to n when it starts.
	var n sync.WaitGroup
	// Here the JobController itself also adds 1 to n, and this 1 is
	// Done()'d when the JobController's context gets cancelled,
	// signalling the termination of the JobController.
	n.Add(1)

	// start a separate goroutine to wait on the waitgroup until all agents
	// AND the JobController are done, and then close the response channel
	go func() {
		n.Wait()
		close(rc)
	}()

	// now we start a goroutine to listen to channels and wait for
	// things to happen
	go func() {
		// note that this could introduce a race condition IF the JobController
		// were to receive a cancel signal from context, and decremented n to
		// zero, AND then a new Job were started, which would try to reuse
		// the zeroed waitgroup. To avoid this, we set exiting to true before
		// calling n.Done(), and after exiting is true we don't create any
		// new Jobs.
		exiting := false

		// note that this should not need to be synchronized. only this
		// goroutine should be checking and updating the job submitted map,
		// and this goroutine is only being run once.
		jobSubmitted := map[uint32]bool{}

		for !exiting {
			select {
			case <-ctx.Done():
				// the JobController has been cancelled and should shut down
				exiting = true
				n.Done()
			case jr := <-inJobStream:
				// the caller has submitted a new JobRequest
				// check whether it's already been submitted once
				if _, ok := jobSubmitted[jr.JobID]; ok {
					log.Printf("===> got Job %d as a repeat; dropping", jr.JobID)
				} else {
					// mark it as submitted
					jobSubmitted[jr.JobID] = true
					// and create the job
					env.startNewJob(ctx, &jr, &n, rc, errc)
				}
			case ju := <-rc:
				// an agent has sent a JobUpdate
				env.updateJobDB(&ju)
			}
		}

		// FIXME as we are exiting, do we first need to drain any remaining
		// FIXME updates from rc?
	}()

	// finally we return the channels so that the caller can kick things off
	return inJobStream, errc
}

func (env *Env) startNewJob(ctx context.Context, jr *JobRequest, n *sync.WaitGroup, rc chan<- JobUpdate, errc chan<- error) {
	log.Printf("===> In startNewJob: jr = %s\n", jr.String())

	// assume that once we get here, JobController has already confirmed that
	// we're ready to run this job (and no one else has already run it)

	// we'll also let either runJobAgent or the agent itself send the first
	// update message with the timeStarted and status values

	// time to actually create the job
	n.Add(1)
	go env.runJobAgent(ctx, jr, n, rc)
}

func (env *Env) updateJobDB(ju *JobUpdate) {
	// convert protobuf values to datastore equivalents
	st, err := status.ConvertStatusProtoToDatastore(ju.StatusRpt.RunStatus)
	if err != nil {
		log.Printf("===> got invalid status value in update: %v", err)
		return
	}
	h, err := status.ConvertHealthProtoToDatastore(ju.StatusRpt.HealthStatus)
	if err != nil {
		log.Printf("===> got invalid health value in update: %v", err)
		return
	}

	err = env.Db.UpdateJobStatus(ju.JobID, time.Unix(ju.StatusRpt.TimeStarted, 0), time.Unix(ju.StatusRpt.TimeFinished, 0), st, h, ju.StatusRpt.OutputMessages)
	if err != nil {
		log.Printf("===> couldn't update status for job %d: %v", ju.JobID, err)
		return
	}
}
