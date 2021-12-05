package unit_process

/*
#include <stdlib.h>
#include "libproc.h"
*/
import "C"

import "C"
import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/utilities/logger"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_comruter_process_png
}

func (c *UnitSystemProcess) InternalUnitStart() error {
	var err error
	type Config struct {
		ProcessName string  `json:"process_name"`
		Period      float64 `json:"period"`
	}
	var config Config
	err = json.Unmarshal([]byte(c.GetConfig()), &config)
	if err != nil {
		logger.Println("ERROR[UnitSystemProcess]:", err)
		err = errors.New("config error")
		c.SetString("Common/ProcessID", err.Error(), "error")
		return err
	}

	c.processName = config.ProcessName
	if c.processName == "" {
		err = errors.New("wrong address")
		c.SetString("Common/ProcessID", err.Error(), "error")
		return err
	}

	c.periodMs = int(config.Period)
	if c.periodMs < 100 {
		err = errors.New("wrong period")
		c.SetString("Common/ProcessID", err.Error(), "error")
		return err
	}

	go c.Tick()
	return nil
}

func (c *UnitSystemProcess) InternalUnitStop() {
}

func (c *UnitSystemProcess) Tick() {
	c.Started = true
	logger.Println("UNIT <Process Windows> started:", c.Id())

	dtOperationTime := time.Now().UTC()

	processId := int(-1)

	/*lastKernelTimeMs := int64(0)
	lastUserTimeMs := int64(0)
	lastReadProcessTimes := time.Now().UTC()


	*/
	for !c.Stopping {
		for {
			if c.Stopping || time.Now().Sub(dtOperationTime) > time.Duration(c.periodMs)*time.Millisecond {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if c.Stopping {
			break
		}

		if processId < 0 {
			processes, err := processes()
			if err == nil {
				for _, proc := range processes {
					if strings.Contains(proc.Executable(), c.processName) {
						processId = proc.Pid()
					}
				}
			}
		}

		var rusage syscall.Rusage
		err := syscall.Getrusage(syscall.RUSAGE_SELF, &rusage)
		if err == nil {
			c.SetInt64("kern.osrelease", rusage.Isrss, "")
			fmt.Println(rusage)
		}

		if processId >= 0 {
			FDs, err := fdusage(processId)
			if err == nil {
				c.SetInt("FDs", FDs, "")
			}
			c.SetString("result", "ok", "")
			c.SetInt("id", processId, "")
		} else {
			c.SetString("result", "no found", "error")
			c.SetInt("id", processId, "")
		}

		dtOperationTime = time.Now().UTC()
	}

	logger.Println("UNIT <Process Windows> stopped:", c.Id())
	c.Started = false
}

func GetProcesses() []ProcessInfo {
	result := make([]ProcessInfo, 0)
	return result
}

func fdusage(pid1 int) (int, error) {
	pid := C.int(pid1)

	rlen, err := C.proc_pidinfo(pid, C.PROC_PIDLISTFDS, 0, nil, 0)
	if rlen <= 0 {
		return 0, err
	}

	var buf unsafe.Pointer
	defer func() {
		if buf != nil {
			C.free(buf)
		}
	}()

	for buflen := rlen; ; buflen *= 2 {
		buf, err = C.reallocf(buf, C.size_t(buflen))
		if buf == nil {
			return 0, err
		}
		rlen, err = C.proc_pidinfo(pid, C.PROC_PIDLISTFDS, 0, buf, buflen)
		if rlen <= 0 {
			return 0, err
		} else if rlen == buflen {
			continue
		}
		return int(rlen / C.PROC_PIDLISTFD_SIZE), nil
	}
}

type DarwinProcess struct {
	pid    int
	ppid   int
	binary string
}

func (p *DarwinProcess) Pid() int {
	return p.pid
}

func (p *DarwinProcess) PPid() int {
	return p.ppid
}

func (p *DarwinProcess) Executable() string {
	return p.binary
}

func findProcess(pid int) (Process, error) {
	ps, err := processes()
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		if p.Pid() == pid {
			return p, nil
		}
	}

	return nil, nil
}

func processes() ([]Process, error) {
	buf, err := darwinSyscall()
	if err != nil {
		return nil, err
	}

	procs := make([]*KInfoStruct, 0, 50)
	//k := 0
	for i := _KINFO_STRUCT_SIZE; i < buf.Len(); i += _KINFO_STRUCT_SIZE {
		proc := &KInfoStruct{}
		proc.Data = make([]byte, _KINFO_STRUCT_SIZE)
		copy(proc.Data, buf.Bytes()[i:i+_KINFO_STRUCT_SIZE])
		//k = i
		procs = append(procs, proc)
	}

	darwinProcs := make([]Process, len(procs))
	for i, p := range procs {
		darwinProcs[i] = &DarwinProcess{
			pid:    int(p.Pid()),
			ppid:   int(p.Pid()),
			binary: p.Name(),
		}
	}

	return darwinProcs, nil
}

func darwinCstring(s [16]byte) string {
	i := 0
	for _, b := range s {
		if b != 0 {
			i++
		} else {
			break
		}
	}

	return string(s[:i])
}

func darwinSyscall() (*bytes.Buffer, error) {
	mib := [4]int32{_CTRL_KERN, _KERN_PROC, _KERN_PROC_ALL, 0}
	size := uintptr(0)

	_, _, errno := syscall.Syscall6(
		syscall.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mib[0])),
		4,
		0,
		uintptr(unsafe.Pointer(&size)),
		0,
		0)

	if errno != 0 {
		return nil, errno
	}

	bs := make([]byte, size)
	_, _, errno = syscall.Syscall6(
		syscall.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mib[0])),
		4,
		uintptr(unsafe.Pointer(&bs[0])),
		uintptr(unsafe.Pointer(&size)),
		0,
		0)

	if errno != 0 {
		return nil, errno
	}

	return bytes.NewBuffer(bs[0:size]), nil
}

const (
	_CTRL_KERN         = 1
	_KERN_PROC         = 14
	_KERN_PROC_ALL     = 0
	_KINFO_STRUCT_SIZE = 648
)

type Timeval struct {
	Sec int32
}

type Itimerval struct {
	Interval Timeval
	Value    Timeval
}

type KInfoStruct struct {
	Data []byte
}

func (c *KInfoStruct) Pid() int32 {
	return int32(binary.LittleEndian.Uint32(c.Data[40:]))
}

func (c *KInfoStruct) Name() string {
	nameLen := 0
	for nameLen < 16 {
		if c.Data[243+nameLen] == 0 {
			break
		}
		nameLen++
	}
	return string(c.Data[243 : 243+nameLen])
}

type kinfoProc struct {
	P_un      [16]byte
	P_vmspace uint64
	P_sigacts uint64
	Pad_cgo_0 [3]byte
	P_flag    int32
	P_stat    int8
	Pid       int32
	_         [199]byte
	Comm      [16]byte
	_         [301]byte
	PPid      int32
	_         [84]byte
}

/*type kinfoProc struct {
	_    [40]byte
	Pid  int32
	_    [199]byte
	Comm [16]byte
	_    [301]byte
	PPid int32
	_    [84]byte
}*/

// Process is the generic interface that is implemented on every platform
// and provides common operations for processes.
type Process interface {
	// Pid is the process ID for this process.
	Pid() int

	// PPid is the parent process ID for this process.
	PPid() int

	// Executable name running this process. This is not a path to the
	// executable.
	Executable() string
}

// Processes returns all processes.
//
// This of course will be a point-in-time snapshot of when this method was
// called. Some operating systems don't provide snapshot capability of the
// process table, in which case the process table returned might contain
// ephemeral entities that happened to be running when this was called.
func Processes() ([]Process, error) {
	return processes()
}

// FindProcess looks up a single process by pid.
//
// Process will be nil and error will be nil if a matching process is
// not found.
func FindProcess(pid int) (Process, error) {
	return findProcess(pid)
}
