// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package task

import (
	"time"

	"github.com/pingcap-incubator/tiup-cluster/pkg/executor"
	"github.com/pingcap-incubator/tiup-cluster/pkg/meta"
	"github.com/pingcap/errors"
)

// SetSSHKeySet set ssh key set.
func (ctx *Context) SetSSHKeySet(privateKeyPath string, publicKeyPath string) error {
	ctx.PrivateKeyPath = privateKeyPath
	ctx.PublicKeyPath = publicKeyPath
	return nil
}

// SetClusterSSH set cluster user ssh executor in context.
func (ctx *Context) SetClusterSSH(topo meta.Specification, deployUser string, sshTimeout int64) error {
	if len(ctx.PrivateKeyPath) == 0 {
		return errors.Errorf("context has no PrivateKeyPath")
	}

	for _, com := range topo.ComponentsByStartOrder() {
		for _, in := range com.Instances() {
			cf := executor.SSHConfig{
				Host:    in.GetHost(),
				Port:    in.GetSSHPort(),
				KeyFile: ctx.PrivateKeyPath,
				User:    deployUser,
				Timeout: time.Second * time.Duration(sshTimeout),
			}

			e := executor.NewSSHExecutor(cf, false /* sudo */)
			ctx.SetExecutor(in.GetHost(), e)
		}
	}

	return nil
}
