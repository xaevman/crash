//  ---------------------------------------------------------------------------
//
//  file.go
//
//  Copyright (c) 2015, Jared Chavez. 
//  All rights reserved.
//
//  Use of this source code is governed by a BSD-style
//  license that can be found in the LICENSE file.
//
//  -----------

package crash

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

// FileHandler is a CrashHandler implementation which generates a crash report
// from the given panic data, formats it in JSON format, and saves it to a flat
// flie on disk.
type FileHandler struct {
    crashDir string
}

// OnCrash is called automatically by the global crash handling routines. OnCrash
// creates a new JSON crash report file in the configured crash directory with the
// filename format "crash.<unix timestamp>.json".
func (this *FileHandler) OnCrash(data interface{}) {
    now      := time.Now().Unix()
    fileName := fmt.Sprintf("crash.%d.json", now)
    filePath := filepath.Join(this.crashDir, fileName)
    f, err   := os.Create(filePath)
    if err != nil {
        fmt.Println(err)
        return
    }

    defer f.Close()

    j, err := NewCrashReport(data).Json()
    if err != nil {
        fmt.Println(err)
        return
    }

    _, err = f.Write(j)
    if err != nil {
        fmt.Println(err)
    }
}

// SetCrashDir reconfigures the CrashHandler to save files to a new base path.
// The new base path is automatically created if it doesn't already exist. Any errors
// from the underlying mkdir call are passed back to the caller.
func (this *FileHandler) SetCrashDir(path string) error {
    this.crashDir = path
    return os.MkdirAll(this.crashDir, 0770)
}
