//  ---------------------------------------------------------------------------
//
//  stack.go
//
//  Copyright (c) 2014, Jared Chavez. 
//  All rights reserved.
//
//  Use of this source code is governed by a BSD-style
//  license that can be found in the LICENSE file.
//
//  -----------

package crash

// Stdlib imports.
import (
    "path/filepath"
    "runtime"
    "strings"
)

// Maximum number of stack records to hold.
const STACK_BUFFERS = 2 * 1024

// StackTrace represents the collection of stack frame data associated
// with a single goroutine.
type StackTrace struct {
    Id     int
    Frames []*StackFrame
}

// StackFrame represents the function name, file and line number information
// for a given frame within a stack trace.
type StackFrame struct {
    File  string
    Line  int
    Name  string
}

// NewStackTrace is a constructor function which queries the Go runtime
// for goroutine information and builds a StackTrace object for each.
// A slice of StackTraces, one for each running goroutine, is returned.
func NewStackTrace() []*StackTrace {
    stackRecords := make([]runtime.StackRecord, STACK_BUFFERS)
    count, ok    := runtime.GoroutineProfile(stackRecords)
    if !ok {
        return nil
    }

    results := make([]*StackTrace, 0)

    for i := 0; i < count; i++ {
        st       := new(StackTrace)
        st.Id     = i
        st.Frames = make([]*StackFrame, 0)
        frames   := stackRecords[i].Stack()

        for x := range frames {
            f          := runtime.FuncForPC(frames[x])
            file, line := f.FileLine(frames[x])
            sf         := new(StackFrame)
            sf.Name     = strings.TrimSpace(filepath.Base(f.Name()))
            sf.File     = strings.TrimSpace(file)
            sf.Line     = line

            st.Frames = append(st.Frames, sf)
        }

        results = append(results, st)
    }

    return results
}
