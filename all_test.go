//  ---------------------------------------------------------------------------
//
//  all_test.go
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
    "testing"
)

// TestFileHandler tests crash handler mechanics across multiple goroutines
// using the printJSONPanic crash handler.
func TestFileHandler(t *testing.T) {
    SetRepanic(false)

    defer HandleAll()

    fh  := new(FileHandler)
    err := fh.SetCrashDir("./crashes")
    if err != nil {
        panic(err)
    }

    AddHandler(fh)

    panic("panic test")
}
