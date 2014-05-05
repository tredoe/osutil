// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build linux

package exp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	//"os/exec"
	"path"
	"path/filepath"
	"strings"

	. "github.com/kless/shutil"
)

// == Errors
var ErrNoRemovUSB = errors.New("removable USB device no found")

type FindPartError string

func (e FindPartError) Error() string {
	return "FindPartition: no device with label \"" + string(e) + `"`
}

type CmdFindPartError string

func (e CmdFindPartError) Error() string {
	return "FindPartition: no data got from command blkid for label \"" + string(e) + `"`
}

// ==

// partition represents data of a partition.
type partition struct {
	dev    string
	fsType string // filesystem type
	label  string
	mount  string // mount point
	uuid   string
}

// FindPartition finds the partition label in the given devices. Returns data
// related to the partition.
func FindPartition(label string, devices []string) (*partition, error) {
	// The format is like that:
	// /dev/sdc1  ext4  key      (not mounted)  ********-****-****-****-************
	//output, err := exec.Command(
	//"/sbin/blkid", "-l", "-o", "device", "-t", "LABEL="+label, "-o", "list").Output()
	output, ok, err := RunWithMatchf("/sbin/blkid -l -o device -t LABEL=%s -o list", label)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, CmdFindPartError(label)
	}

	found := false
	bSlash := []byte{'/'}
	buf := bytes.NewBuffer(output)
	part := new(partition)

	for {
		line, err := buf.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		if !bytes.HasPrefix(line, bSlash) {
			continue
		}

		columns := strings.Fields(string(line))

		// Is it the device?
		for _, v := range devices {
			if strings.HasPrefix(columns[0], v) {
				found = true
				break
			}
		}
		if !found {
			continue
		}

		part.dev = columns[0]
		part.fsType = columns[1]
		part.label = columns[2]
		part.uuid = columns[len(columns)-1]

		if len(columns) == 5 {
			part.mount = columns[3]
		}
		break
	}

	if !found {
		return nil, FindPartError(label)
	}
	return part, nil
}

// GetUSBremovables returns a list of the USB removables devices.
func GetUSBremovables() (devices []string, err error) {
	removableTag := []byte{'1'} // devices USB removables have "1" in the file

	disksPath, err := filepath.Glob("/sys/block/sd*")
	if err != nil {
		return nil, fmt.Errorf("GetUSBremovables: %s", err)
	}

	for _, p := range disksPath {
		if DEBUG {
			Writef("\nExamining " + p)
		}

		fullpath, err := os.Readlink(p)
		if err != nil {
			if DEBUG {
				Writefln("%s", err)
			}
			continue
		}

		// Is it an USB device?
		if strings.Contains(fullpath, "usb") {
			// Is the device removable?
			file, err := os.Open(path.Join(p, "removable"))
			if err != nil {
				if DEBUG {
					Writefln("%s", err)
				}
				continue
			}

			buf := bufio.NewReader(file)
			firstLine, _, _ := buf.ReadLine()

			if bytes.Equal(firstLine, removableTag) {
				devices = append(devices, path.Join("/dev", path.Base(p)))
				if DEBUG {
					Writef(": USB removable")
				}
			}

			file.Close()
		}
	}
	if DEBUG {
		Writefln("")
	}

	if len(devices) == 0 {
		return nil, ErrNoRemovUSB
	}
	return devices, nil
}
