/*
Copyright 2018 the Heptio Ark contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/heptio/velero/pkg/util/test"
)

type fakeClient struct {
	base       *clientBase
	clientConn *grpc.ClientConn
}

func TestNewClientDispenser(t *testing.T) {
	logger := test.NewLogger()

	clientConn := new(grpc.ClientConn)

	c := 3
	initFunc := func(base *clientBase, clientConn *grpc.ClientConn) interface{} {
		return c
	}

	cd := newClientDispenser(logger, clientConn, initFunc)
	assert.Equal(t, clientConn, cd.clientConn)
	assert.NotNil(t, cd.clients)
	assert.Empty(t, cd.clients)
}

func TestClientFor(t *testing.T) {
	logger := test.NewLogger()
	clientConn := new(grpc.ClientConn)

	c := new(fakeClient)
	count := 0
	initFunc := func(base *clientBase, clientConn *grpc.ClientConn) interface{} {
		c.base = base
		c.clientConn = clientConn
		count++
		return c
	}

	cd := newClientDispenser(logger, clientConn, initFunc)

	actual := cd.clientFor("pod")
	require.IsType(t, &fakeClient{}, actual)
	typed := actual.(*fakeClient)
	assert.Equal(t, 1, count)
	assert.Equal(t, &typed, &c)
	expectedBase := &clientBase{
		plugin: "pod",
		logger: logger,
	}
	assert.Equal(t, expectedBase, typed.base)
	assert.Equal(t, clientConn, typed.clientConn)

	// Make sure we reuse a previous client
	actual = cd.clientFor("pod")
	require.IsType(t, &fakeClient{}, actual)
	typed = actual.(*fakeClient)
	assert.Equal(t, 1, count)
}
