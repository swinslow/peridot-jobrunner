// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package main

import (
	"context"
	"log"
	"time"

	"github.com/swinslow/peridot-jobrunner/internal/jobrunner"
)

func main() {
	env, err := jobrunner.SetupEnv()
	if err != nil {
		log.Fatalf("unable to set up env: %v", err)
	}

	// inJobStream is created by JobController. It is used to submit
	// JobRequests to the JobController. We own this channel and
	// must close it when we're done.
	var inJobStream chan<- jobrunner.JobRequest

	// errc is created by JobController. It receives broadcasts of
	// any JobController-level errors. JobController owns this
	// channel and will close it.
	// var errc <-chan error

	ctx, cancel := context.WithCancel(context.Background())
	inJobStream, _ = env.JobController(ctx)
	//inJobStream, errc = env.JobController(ctx)

	cycle := 1
	for {
		if cycle > 10 {
			break
		}
		log.Printf("======> CYCLE %d", cycle)
		jobs, err := env.Db.GetReadyJobs(5)
		if err != nil {
			log.Fatalf("error retrieving ready jobs: %v", err)
		}

		for _, j := range jobs {
			// build the job request
			jcfg := jobrunner.ConvertJobToJobConfig(j)
			jr := jobrunner.JobRequest{
				JobID:   j.ID,
				AgentID: j.AgentID,
				Cfg:     *jcfg,
			}

			// and send it!
			inJobStream <- jr
		}

		time.Sleep(5 * time.Second)
		cycle++
	}

	// exiting
	log.Printf("======> CYCLES DONE")

	// close the channel we own
	close(inJobStream)

	// tell the JobController to clean up and shut down
	log.Printf("===> calling cancel()")
	cancel()

}
