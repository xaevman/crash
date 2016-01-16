//  ---------------------------------------------------------------------------
//
//  report.go
//
//  Copyright (c) 2014, Jared Chavez. 
//  All rights reserved.
//
//  Use of this source code is governed by a BSD-style
//  license that can be found in the LICENSE file.
//
//  -----------

package crash

import (
    "github.com/xaevman/app"

    "encoding/json"
    "os"
    "runtime"
    "time"
)

// CrashReport represents important information which is generated
// during a crash.
type CrashReport struct {
    Timestamp    time.Time
    ErrData      interface{}
    AppName      string
    AppPath      string
    GoOS         string
    GoArch       string
    CPUCount     int
    HostName     string
    WorkingDir   string
    MemStats     interface{}
    Environment  []string
    StackTrace   []*StackTrace
}

// Json returns the JSON representation of this CrashReport object.
func (this *CrashReport) Json() ([]byte, error) {
    data, err := json.MarshalIndent(this, "", "    ")
    if err != nil {
        return nil, err
    }

    return data, nil
}

// TerseMemStats represents a subset of the information normally
// provided in a runtime.MemStats structure.
type TerseMemStats struct {
    Alloc         uint64
    Mallocs       uint64
    Frees         uint64
    HeapIdle      uint64
    HeapInuse     uint64
    HeapObjects   uint64
    StackInuse    uint64
    NumGC         uint32
    GCCPUFraction float64
    EnableGC      bool
    DebugGC       bool
}

// NewCrashReport generates a new CrashReport object, fills it in
// with information about the current system state, and returns a reference
// to the new object back to the caller.
func NewCrashReport(errData interface{}) *CrashReport {
    var ms runtime.MemStats
    runtime.ReadMemStats(&ms)

    hostname, _ := os.Hostname()
    wd, _       := os.Getwd()

    rpt := &CrashReport {
        Timestamp   : time.Now(),
        ErrData     : errData,
        AppName     : app.GetName(),
        AppPath     : app.GetExeDir(),
        GoOS        : runtime.GOOS,
        GoArch      : runtime.GOARCH,
        CPUCount    : runtime.NumCPU(),
        HostName    : hostname,
        WorkingDir  : wd,
        Environment : os.Environ(),
        StackTrace  : NewStackTrace(),
    }

    if verboseRpt {
        rpt.MemStats = &ms
    } else {
        rpt.MemStats = &TerseMemStats {
            Alloc         : ms.Alloc,
            Mallocs       : ms.Mallocs,
            Frees         : ms.Frees,
            HeapIdle      : ms.HeapIdle,
            HeapInuse     : ms.HeapInuse,
            HeapObjects   : ms.HeapObjects,
            StackInuse    : ms.StackInuse,
            NumGC         : ms.NumGC,
//            GCCPUFraction : ms.GCCPUFraction,
            EnableGC      : ms.EnableGC,
            DebugGC       : ms.DebugGC,
        }
    }

    return rpt
}
