// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package jobrunner

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/swinslow/peridot-jobrunner/pkg/agent"
	"github.com/swinslow/peridot-jobrunner/pkg/status"
	"google.golang.org/grpc"
)

func getErrorUpdate(jobID uint32, err error) JobUpdate {
	return JobUpdate{
		JobID: jobID,
		StatusRpt: agent.StatusReport{
			RunStatus:    status.Status_STOPPED,
			HealthStatus: status.Health_ERROR,
		},
		Err: err,
	}
}

func (env *Env) runJobAgent(ctx context.Context, jr *JobRequest, n *sync.WaitGroup, rc chan<- JobUpdate) {
	defer n.Done()

	log.Printf("===> in runJobAgent\n")

	// get agent details
	ag, err := env.Db.GetAgentByID(jr.AgentID)
	if err != nil {
		rc <- getErrorUpdate(jr.JobID, fmt.Errorf("could not get agent details from database for agent %d: %v", jr.AgentID, err))
		return
	}

	agentURL := fmt.Sprintf("%s:%d", ag.Address, ag.Port)

	// connect and get client for each agent server
	conn, err := grpc.Dial(agentURL, grpc.WithInsecure())
	if err != nil {
		rc <- getErrorUpdate(jr.JobID, fmt.Errorf("could not connect to %s (%s): %v", ag.Name, agentURL, err))
		return
	}
	defer conn.Close()
	c := agent.NewAgentClient(conn)

	// set up context
	// FIXME 20 second timeout seems very wrong
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	// start NewJob stream
	stream, err := c.NewJob(ctx)
	if err != nil {
		rc <- getErrorUpdate(jr.JobID, fmt.Errorf("could not connect for %s (%s): %v", ag.Name, agentURL, err))
		return
	}

	// make server call to start job
	startReq := &agent.StartReq{Config: &jr.Cfg}
	cm := &agent.ControllerMsg{Cm: &agent.ControllerMsg_Start{Start: startReq}}
	log.Printf("===> controller SEND StartReq for job %d", jr.JobID)
	err = stream.Send(cm)
	if err != nil {
		rc <- getErrorUpdate(jr.JobID, fmt.Errorf("could not start job for %s (%s): %v", ag.Name, agentURL, err))
		return
	}

	// set up listener + status updater goroutine
	// until we get past waitc, ONLY the listener goroutine should be
	// updating the job status
	waitc := make(chan interface{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// done with reading
				log.Printf("===> controller CLOSING got io.EOF")
				close(waitc)
				return
			}
			if err != nil {
				log.Printf("===> controller CLOSING got error: %v", err)
				rc <- getErrorUpdate(jr.JobID, fmt.Errorf("error for %s (%s): %v", ag.Name, agentURL, err))
				close(waitc)
				return
			}

			// update status if we got a status report
			switch x := in.Am.(type) {
			case *agent.AgentMsg_Status:
				st := *x.Status
				log.Printf("===> controller RECV StatusReport for jobID %d: %s\n", jr.JobID, st.String())
				rc <- JobUpdate{
					JobID:     jr.JobID,
					StatusRpt: st,
				}

				// if this was a STOPPED message, the job is done
				// and we need to close the stream to start exiting
				if st.RunStatus == status.Status_STOPPED {
					log.Printf("===> controller CLOSING got Job STOPPED")
					close(waitc)
					return
				}
			}
		}
	}()

	// wait until listener loop is done
	// FIXME ordinarily this should probably ping occasionally with a heartbeat
	// FIXME request, and/or eventually exit if we see an error or if a job
	// FIXME hasn't responded for ___ time
	// FIXME also, does CloseSend need to come before we wait for agent to close?
	exiting := false
	for !exiting {
		select {
		case <-waitc:
			stream.CloseSend()
			exiting = true
			// case <-time.After(time.Second * 5):
			// 	// check status and see whether we should continue waiting
			// 	if st.status.RunStatus == agent.JobRunStatus_STOPPED {
			// 		stream.CloseSend()
			// 		exiting = true
			// 	}
		}
	}
}
