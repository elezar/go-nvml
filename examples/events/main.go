/**
# Copyright (c) 2021, NVIDIA CORPORATION.  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
**/

package main

import (
	"fmt"
	"log"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

func main() {
	fmt.Printf("Init\n")
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to initialize NVML: %v", nvml.ErrorString(ret))
	}
	defer func() {
		fmt.Printf("Shutdown\n")
		ret := nvml.Shutdown()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to shutdown NVML: %v", nvml.ErrorString(ret))
		}
	}()

	fmt.Printf("DeviceGetCount\n")
	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to get device count: %v", nvml.ErrorString(ret))
	}

	for i := 0; i < count; i++ {
		fmt.Printf("DeviceGetHandleByIndex\n")
		device, ret := nvml.DeviceGetHandleByIndex(i)
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get device at index %d: %v", i, nvml.ErrorString(ret))
		}

		fmt.Printf("device.GetUUID %d\n", i)
		uuid, ret := device.GetUUID()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get uuid of device at index %d: %v", i, nvml.ErrorString(ret))
		}
		fmt.Printf("%v\n", uuid)

		fmt.Printf("device.GetMinorNumber %d\n", i)
		minor, ret := device.GetMinorNumber()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get minor number of device at index %d: %v", i, nvml.ErrorString(ret))
		}
		fmt.Printf("%v\n", minor)

		fmt.Printf("device.GetPciInfo %d\n", i)
		pciInfo, ret := device.GetPciInfo()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get PCI info of device at index %d: %v", i, nvml.ErrorString(ret))
		}
		fmt.Printf("%v\n", pciInfo)
	}

	fmt.Printf("Shutdown before EventSetCreate\n")
	ret = nvml.Shutdown()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to shutdown NVML: %v", nvml.ErrorString(ret))
	}
	fmt.Printf("Init before EventSetCreate\n")
	ret = nvml.Init()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to initialize NVML: %v", nvml.ErrorString(ret))
	}

	fmt.Printf("EventSetCreate\n")
	events, ret := nvml.EventSetCreate()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to create event set: %v", nvml.ErrorString(ret))
	}
	defer func() {
		ret := events.Free()
		if ret != nvml.SUCCESS {
			fmt.Printf("Unable to free eventset: %v", nvml.ErrorString(ret))
		}
	}()

	fmt.Printf("Registering events\n")
	for i := 0; i < count; i++ {
		fmt.Printf("DeviceGetHandleByIndex %d\n", i)
		device, ret := nvml.DeviceGetHandleByIndex(i)
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get device at index %d: %v", i, nvml.ErrorString(ret))
		}

		fmt.Printf("device.RegisterEvents %d\n", i)
		ret = device.RegisterEvents(nvml.EventTypeXidCriticalError, events)
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to register events for device %d: %v", i, nvml.ErrorString(ret))
		}
	}

	for j := 0; j < 3; j++ {
		fmt.Printf("events.Wait: %d", j)
		event, ret := events.Wait(5000)
		if ret != nvml.SUCCESS {
			fmt.Printf("Unable to wait for events: %v", nvml.ErrorString(ret))
			continue
		}
		fmt.Printf("event=%+v", event)
	}
}
