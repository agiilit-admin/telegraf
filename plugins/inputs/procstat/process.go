package procstat

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/process"
)

//nolint:interfacebloat // conditionally allow to contain more methods
type Process interface {
	PID() PID
	Tags() map[string]string

	PageFaults() (*process.PageFaultsStat, error)
	IOCounters() (*process.IOCountersStat, error)
	MemoryInfo() (*process.MemoryInfoStat, error)
	MemoryMaps(bool) (*[]process.MemoryMapsStat, error)
	Name() (string, error)
	Cmdline() (string, error)
	NumCtxSwitches() (*process.NumCtxSwitchesStat, error)
	NumFDs() (int32, error)
	NumThreads() (int32, error)
	Percent(interval time.Duration) (float64, error)
	MemoryPercent() (float32, error)
	Times() (*cpu.TimesStat, error)
	RlimitUsage(bool) ([]process.RlimitStat, error)
	Username() (string, error)
	CreateTime() (int64, error)
	Ppid() (int32, error)
	Status() ([]string, error)
}

type PIDFinder interface {
	PidFile(path string) ([]PID, error)
	Pattern(pattern string) ([]PID, error)
	UID(user string) ([]PID, error)
	FullPattern(path string) ([]PID, error)
	ChildPattern(path string) ([]PID, error)
}

type Proc struct {
	hasCPUTimes bool
	tags        map[string]string
	*process.Process
}

func NewProc(pid PID) (Process, error) {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, err
	}

	proc := &Proc{
		Process:     p,
		hasCPUTimes: false,
		tags:        make(map[string]string),
	}
	return proc, nil
}

func (p *Proc) Tags() map[string]string {
	return p.tags
}

func (p *Proc) PID() PID {
	return PID(p.Process.Pid)
}

func (p *Proc) Username() (string, error) {
	return p.Process.Username()
}

func (p *Proc) Percent(_ time.Duration) (float64, error) {
	cpuPerc, err := p.Process.Percent(time.Duration(0))
	if !p.hasCPUTimes && err == nil {
		p.hasCPUTimes = true
		return 0, fmt.Errorf("must call Percent twice to compute percent cpu")
	}
	return cpuPerc, err
}
