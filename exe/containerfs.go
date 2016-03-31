// Implementation of FUSE's file system on top of Azure blob storage.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/ppanyukov/azure-sdk-for-go/storage"
	"github.com/ppanyukov/azurefs-fuse/blobfs"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [flags] MOUNTPOINT\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "The flags are:\n")
	flag.PrintDefaults()
}

func main() {
	// TODO(ppanyukov): too much args parsing, is there a better saner way?
	// TODO(ppanyukov): better and more secure way of passing the account/key than ags/env vars?

	// zap sensitive vars early so nobody can grab them via /proc/pid/environ
	// This may actually not work, because according to docs:
	//      setenv_c and unsetenv_c are provided by the runtime but are no-ops
	//      if cgo isn't loaded.
	// At least it makes them inaccessible in go.
	var envAccountName string = os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	var envAccountKey string = os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	os.Clearenv()

	var (
		isTrace     bool
		accountName string
		accountKey  string
		mountPoint  string
	)

	// Use custom usage printer.
	flag.Usage = usage
	flag.StringVar(&accountName, "accountName", "", "REQUIRED. Azure storage account name. Or use AZURE_STORAGE_ACCOUNT_NAME env var.")
	flag.StringVar(&accountKey, "accountKey", "", "REQUIRED. Azure storage account key. Or use AZURE_STORAGE_ACCOUNT_KEY env var.")
	flag.BoolVar(&isTrace, "trace", false, "OPTIONAL. Specify true to trace calls.")
	flag.Parse()

	if accountName == "" {
		accountName = envAccountName
	}
	if accountKey == "" {
		accountKey = envAccountKey
	}

	if len(flag.Args()) > 0 {
		mountPoint = flag.Arg(0)
	}

	if accountName == "" || accountKey == "" || mountPoint == "" {
		flag.Usage()
		os.Exit(1)
	}

	// good to go
	fmt.Printf("OK. Will mount storage account '%s' at '%s'", accountName, mountPoint)

	storageClient, err := storage.NewBasicClient(accountName, accountKey)
	if err != nil {
		log.Fatal("ERROR", err)
	}

	var fs pathfs.FileSystem
	containerFs := blobfs.NewContainerFs(storageClient)
	if isTrace {
		fs = blobfs.NewTraceFs(containerFs)
	} else {
		fs = containerFs
	}

	nfs := pathfs.NewPathNodeFs(fs, nil)
	server, _, err := nodefs.MountRoot(mountPoint, nfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}

	server.Serve()
}
