#!/bin/bash

# SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

# Generates Golang protobuf code from .proto files.
# Should be run from the top-level peridot-jobrunner directory.

protoc -I ./ ./pkg/status/status.proto --go_out=plugins=grpc,paths=source_relative:.
protoc -I ./ ./pkg/agent/agent.proto --go_out=plugins=grpc,paths=source_relative:.
