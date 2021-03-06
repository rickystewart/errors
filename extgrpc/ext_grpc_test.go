// Copyright 2019 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package extgrpc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/extgrpc"
	"github.com/cockroachdb/errors/testutils"
	"google.golang.org/grpc/codes"
)

func TestGrpc(t *testing.T) {
	err := fmt.Errorf("hello")
	err = extgrpc.WrapWithGrpcCode(err, codes.Unavailable)

	// Simulate a network transfer.
	enc := errors.EncodeError(context.Background(), err)
	otherErr := errors.DecodeError(context.Background(), enc)

	tt := testutils.T{T: t}

	// Error is preserved through the network.
	tt.CheckDeepEqual(otherErr, err)

	// It's possible to extract the Grpc code.
	tt.CheckEqual(extgrpc.GetGrpcCode(otherErr), codes.Unavailable)

	// If there are multiple codes, the most recent one wins.
	otherErr = extgrpc.WrapWithGrpcCode(otherErr, codes.NotFound)
	tt.CheckEqual(extgrpc.GetGrpcCode(otherErr), codes.NotFound)

	// The code is hidden when the error is printed with %v.
	tt.CheckStringEqual(fmt.Sprintf("%v", err), `hello`)
	// The code appears when the error is printed verbosely.
	tt.CheckStringEqual(fmt.Sprintf("%+v", err), `hello
(1) gRPC code: Unavailable
Wraps: (2) hello
Error types: (1) *extgrpc.withGrpcCode (2) *errors.errorString`)

	// Checking the code of a nil error should be codes.OK
	var noErr error
	tt.Assert(extgrpc.GetGrpcCode(noErr) == codes.OK)
}
