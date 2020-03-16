package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	metadata := `
{
    "id": 1,
    "addition_hb_addr": "10.255.101.75:6808/28436",
    "arch": "x86_64",
    "back_addr": "10.255.101.75:6807/28436",
    "ceph_version": "ceph version SDS_4.2.000.0.191128.2-7-gd2ac60c (d2ac60ce3cb6701cbd9e6a0af44875780bf74940) mimic (dev)",
    "cpu": "Intel(R) Xeon(R) CPU E5-2680 0 @ 2.70GHz",
    "default_device_class": "hdd",
    "distro": "centos",
    "distro_description": "CentOS Linux 7 (Core)",
    "distro_version": "7",
    "front_addr": "10.255.101.75:6806/28436",
    "hb_back_addr": "10.255.101.75:6809/28436",
    "hb_front_addr": "10.255.101.75:6810/28436",
    "hostname": "ceph-3",
    "kernel_description": "#1 SMP Tue Aug 22 21:09:27 UTC 2017",
    "kernel_version": "3.10.0-693.el7.x86_64",
    "mem_swap_kb": "0",
    "mem_total_kb": "16268092",
    "os": "Linux",
    "osd_data": "/var/lib/ceph/osd/ceph-1",
    "osd_objectstore": "xstore2",
    "pid": 28436,
    "rotational": "1",
    "uuid": "a7dac8de-9e63-45b3-8730-e0c72959771f",
    "up": true,
    "in": true,
    "exit_code": 0
}
`
	var osdMetadata OsdMetadata
	err := json.Unmarshal([]byte(metadata), &osdMetadata)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("metadata = %#v", osdMetadata)

}

// OsdMetadata defines osd metadata info
type OsdMetadata struct {
	ID             uint64 `json:"id"`
	Up             bool   `json:"up"`
	In             bool   `json:"in"`
	ExitCode       int    `json:"exit_code"`
	HostName       string `json:"hostname"`
	OS             string `json:"os"`
	OsdData        string `json:"osd_data"`
	OsdObjectStore string `json:"osd_objectstore"`
	Arch           string `json:"arch"`
	AdditionHbAddr string `json:"addition_hb_addr"`
}
