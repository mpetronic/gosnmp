// Copyright 2012 The GoSNMP Authors. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in the
// LICENSE file.

/*
The developer of the trapserver code (https://github.com/jda) says "I'm working
on the best level of abstraction but I'm able to receive traps from a Cisco
switch and Net-SNMP".

Pull requests welcome.
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	g "github.com/gosnmp/gosnmp"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		fmt.Printf("   %s\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	tl := g.NewTrapListener()
	tl.OnNewTrap = myTrapHandler
	tl.Params = g.Default
	tl.Params.Logger = g.NewLogger(log.New(os.Stdout, "", 0))
	tl.Params.SecurityModel = g.UserSecurityModel
	tl.Params.Version = g.Version3

	tl.Params.SecurityParameters = &g.UsmSecurityParameters{
		// AuthoritativeEngineID:    "12345678",
		AuthenticationProtocol: g.SHA,
		AuthenticationPassphrase: "testauth",
		PrivacyProtocol:   g.AES,
		PrivacyPassphrase: "testpriv",
		// UserName:                 "testuser",
		Logger: tl.Params.Logger,
	}

	err := tl.Listen("127.0.0.1:9162")
	if err != nil {
		log.Panicf("error in listen: %s", err)
	}
}

func myTrapHandler(packet *g.SnmpPacket, addr *net.UDPAddr) {
	log.Printf("got trapdata from %s\n", addr.IP)
	for _, v := range packet.Variables {
		switch v.Type {
		case g.OctetString:
			b := v.Value.([]byte)
			log.Printf("OID: %s = %s %v\n", v.Name, b, b)
		case g.Integer:
			log.Printf("OID: %s = %d\n", v.Name, v.Value.(int))
		default:
			log.Printf("OID: %+v\n", v)
		}
	}
}
