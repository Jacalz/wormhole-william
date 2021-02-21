package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/psanford/wormhole-william/wormhole"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <code>\n", os.Args[0])
		os.Exit(1)
	}

	code := os.Args[1]

	var c wormhole.Client

	ctx := context.Background()
	msg, err := c.Receive(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("got msg: %+v\n", msg)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	tmpFile, err := ioutil.TempFile(wd, msg.Name+".zip.tmp")
	if err != nil {
		log.Fatal(err)
	}

	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	n, err := io.Copy(tmpFile, msg)
	if err != nil {
		log.Fatal("readfull  error", err)
	}

	err = extract(tmpFile, n, filepath.Join(wd, msg.Name))
	if err != nil {
		log.Fatal(err)
	}
}
