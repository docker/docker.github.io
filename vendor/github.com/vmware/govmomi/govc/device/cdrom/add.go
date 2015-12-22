/*
Copyright (c) 2014 VMware, Inc. All Rights Reserved.

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

package cdrom

import (
	"flag"
	"fmt"

	"github.com/vmware/govmomi/govc/cli"
	"github.com/vmware/govmomi/govc/flags"
	"golang.org/x/net/context"
)

type add struct {
	*flags.VirtualMachineFlag

	controller string
}

func init() {
	cli.Register("device.cdrom.add", &add{})
}

func (cmd *add) Register(f *flag.FlagSet) {
	f.StringVar(&cmd.controller, "controller", "", "IDE controller name")
}

func (cmd *add) Process() error { return nil }

func (cmd *add) Run(f *flag.FlagSet) error {
	vm, err := cmd.VirtualMachine()
	if err != nil {
		return err
	}

	if vm == nil {
		return flag.ErrHelp
	}

	devices, err := vm.Device(context.TODO())
	if err != nil {
		return err
	}

	c, err := devices.FindIDEController(cmd.controller)
	if err != nil {
		return err
	}

	d, err := devices.CreateCdrom(c)
	if err != nil {
		return err
	}

	err = vm.AddDevice(context.TODO(), d)
	if err != nil {
		return err
	}

	// output name of device we just created
	devices, err = vm.Device(context.TODO())
	if err != nil {
		return err
	}

	devices = devices.SelectByType(d)

	name := devices.Name(devices[len(devices)-1])

	fmt.Println(name)

	return nil
}
