//  ---------------------------------------------------------------------------
//
//  crash.go
//
//  Copyright (c) 2015, Jared Chavez. 
//  All rights reserved.
//
//  Use of this source code is governed by a BSD-style
//  license that can be found in the LICENSE file.
//
//  -----------

// Package crash provides some basic services to help handle logging
// and reporting of application panics/crashes.
package crash

import (
    "sync"
)

var (
    handlers   []CrashHandler
    crashLock  sync.Mutex
    rePanic    = true
    verboseRpt = false
)

// CrashHandler represents the interface which must be implemented
// by custom crash handling objects.
type CrashHandler interface {
    OnCrash(interface{})
}

// AddHandler adds a new crash handler to be used by the HandleCrashes function
// when a panic recovery occurs. Handlers are called in the order they are 
// registered.
func AddHandler(newHandler CrashHandler) {
    crashLock.Lock()
    defer crashLock.Unlock()

    handlers = append(handlers, newHandler)
}

// HandleAll should be defered as a part of calling a new goroutine
// in order to handle panics that occur within that routine. Crash handlers
// are called in the order they were registered.
func HandleAll() {
    if r := recover(); r != nil {
        crashLock.Lock()
        defer crashLock.Unlock()

        for i := range handlers {
            handlers[i].OnCrash(r)
        }

        // handlers have all been called, now finish crashing.
        if rePanic {
            panic(r)
        }
    }
}

// SetRepanic sets the re-panic mode. With re-panic disabled, HandleCrashes
// will exit cleanly after all crash handlers have been run. When enabled,
// HandleCrashes "re-throws" the current panic after all crash handlers have
// been run. The default vaule is true.
func SetRepanic(val bool) {
    rePanic = val
}

// SetVerboseCrashReport determines whether or not verbose crash reports
// are generated. The default value is false.
func SetVerboseCrashReport(val bool) {
    verboseRpt = val
}