//  ---------------------------------------------------------------------------
//
//  email.go
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
    "github.com/xaevman/str"

    "fmt"
    "net/smtp"
)

// EmailHandler is a CrashHandler implementation which generates a crash report
// from the given panic data, formats in JSON format, and sends it to the configured
// email addresses via the configured smtp server.
type EmailHandler struct {
    FromAddr string
    ToAddrs  []string
    SrvAddr  string
    SrvPort  int
    SrvUser  string
    SrvPass  string
}

// OnCrash is called automatically by the global crash handling routines. OnCrash 
// creates a new JSON crash report email and sends it to the configured recipients.
func (this *EmailHandler) OnCrash(data interface{}) {
    var auth smtp.Auth

    if this.SrvUser != "" && this.SrvPass != "" {
        auth = smtp.PlainAuth("", this.SrvUser, this.SrvPass, this.SrvAddr)
    } else {
        auth = nil
    }

    rpt := NewCrashReport(data)

    j, err := rpt.Json()
    if err != nil {
        panic(err)
        fmt.Println(err)
        return
    }

    msg := []byte(fmt.Sprintf(
        "To: %s\n" + 
        "Subject: %s crash report (%s): %s\n\n" + 
        "%s\n",
        str.ListToDelimString(this.ToAddrs, ","),
        rpt.AppName,
        rpt.HostName,
        rpt.ErrData,
        j,
    ))

    err = smtp.SendMail(
        fmt.Sprintf("%s:%d", this.SrvAddr, this.SrvPort),
        auth,
        this.FromAddr,
        this.ToAddrs,
        msg,
    )

    if err != nil {
        panic(err)
        fmt.Println(err)
    }
}

// NewEmailHandler returns a pointer to a new EmailHandler instance
// with its values fully intialized.
func NewEmailHandler() *EmailHandler {
    return &EmailHandler {
        ToAddrs : make([]string, 0),
    }
}
