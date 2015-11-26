// Copyright © 2014 Calum MacRae
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh"
)

const vmlookup = "/usr/bin/env vmadm lookup -j -o "

var (
	/*	Arguments:
		-h > Remote hostname (default is nil string, this is a required arg)
		-i > System information to print (per zone/kvm)
		-k > Path to private key (defaults to $HOME/.ssh/id_rsa)
		-p > Destination port for SSH session
		-u > Remote username (defaults to current user)
	*/
	host    = flag.String("h", "", "\t\t\t\tRemote host")
	info    = flag.String("i", "alias,uuid,nics,brand", "\t\tSystem information to print (per zone/kvm)")
	privkey = flag.String("k", os.Getenv("HOME")+"/.ssh/id_rsa", "\tPath to private key")
	port    = flag.Int("p", 22, "\t\t\t\tDestination port for SSH session")
	usr     = flag.String("u", os.Getenv("USER"), "\t\t\tRemote user")
)

func init() { flag.Parse() }

// Usage message
func usage() {
	fmt.Fprintf(os.Stderr, "Usage: smartinv [arguments]\n\n")
	flag.PrintDefaults()
	os.Exit(2)
}

// Grab the user's private key
// This defaults to $HOME/.ssh/id_rsa, however can be overridden with '-k'
func grabKey() (key ssh.Signer, err error) {
	kf, err := ioutil.ReadFile(*privkey)
	if err != nil {
		//TODO: Handle error
		panic(err)
	}
	key, err = ssh.ParsePrivateKey(kf)
	if err != nil {
		//TODO: Handle error
		panic(err)
	}
	return
}

func main() {
	// If '-h ' value is empty, print usage & exit
	if *host == "" {
		fmt.Fprintf(os.Stderr, "No remote host address specified\n\n")
		usage()
	}

	key, err := grabKey()
	if err != nil {
		//TODO: Handle error
		panic(err)
	}

	config := &ssh.ClientConfig{
		User: *usr,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	// Dial the remote server
	p := strconv.Itoa(*port)
	client, err := ssh.Dial("tcp", *host+":"+p, config)
	if err != nil {
		//TODO: Handle error
		panic("Failed to dial: " + err.Error())
	}

	// Open an SSH session following successful dial
	session, err := client.NewSession()
	if err != nil {
		//TODO: Handle error
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	// Lookup zone/kvm information & return results
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(vmlookup + *info); err != nil {
		//TODO: Handle error
		panic("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
}
