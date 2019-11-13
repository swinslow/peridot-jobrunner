// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package jobrunner

import (
	"fmt"

	"github.com/swinslow/peridot-db/pkg/datastore"
	"github.com/swinslow/peridot-jobrunner/pkg/agent"
)

// ConvertJobToJobConfig takes a datastore Job and translates
// it into a JobConfig protobuf struct.
func ConvertJobToJobConfig(j *datastore.Job) *agent.JobConfig {
	jcfg := agent.JobConfig{}
	if j == nil {
		return &jcfg
	}

	// codereader inputs
	for k, jpc := range j.Config.CodeReader {
		if jpc.PriorJobID > 0 {
			// prior job ID was specified, build path using it
			jci := agent.JobConfig_CodeInput{
				Source: k,
				Path:   fmt.Sprintf("/code/%d/", jpc.PriorJobID),
			}
			jcfg.CodeInputs = append(jcfg.CodeInputs, &jci)
		} else {
			// no prior job ID was specified, just use value
			jci := agent.JobConfig_CodeInput{
				Source: k,
				Path:   jpc.Value,
			}
			jcfg.CodeInputs = append(jcfg.CodeInputs, &jci)
		}
	}

	// spdxreader inputs
	for k, jpc := range j.Config.SpdxReader {
		if jpc.PriorJobID > 0 {
			// prior job ID was specified, build path using it
			jci := agent.JobConfig_SpdxInput{
				Source: k,
				Path:   fmt.Sprintf("/spdx/%d/", jpc.PriorJobID),
			}
			jcfg.SpdxInputs = append(jcfg.SpdxInputs, &jci)
		} else {
			// no prior job ID was specified, just use value
			jci := agent.JobConfig_SpdxInput{
				Source: k,
				Path:   jpc.Value,
			}
			jcfg.SpdxInputs = append(jcfg.SpdxInputs, &jci)
		}
	}

	// codewriter output
	jcfg.CodeOutputDir = fmt.Sprintf("/code/%d/", j.ID)

	// spdxwriter output
	jcfg.SpdxOutputDir = fmt.Sprintf("/spdx/%d/", j.ID)

	// key-value configs
	for k, v := range j.Config.KV {
		jkv := agent.JobConfig_JobKV{Key: k, Value: v}
		jcfg.Jkvs = append(jcfg.Jkvs, &jkv)
	}

	return &jcfg
}
