package blobfs

// Tracing of operations on underlying FileSystem

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

// NewTraceFs creates a file system which traces calls and delegates to the specified fs.
func NewTraceFs(fs pathfs.FileSystem) pathfs.FileSystem {
	logPrefix := fmt.Sprintf("[%s]: ", fs.String())
	result := traceFs{
		fs:  fs,
		log: log.New(os.Stderr, logPrefix, log.LstdFlags),
	}

	return &result
}

type traceFs struct {
	fs  pathfs.FileSystem
	log *log.Logger
}

func (fs *traceFs) SetDebug(debug bool) {}

func (fs *traceFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	fs.log.Println("[TRACE] GetAttr:", "name:", name)
	return fs.fs.GetAttr(name, context)
}

func (fs *traceFs) GetXAttr(name string, attr string, context *fuse.Context) ([]byte, fuse.Status) {
	fs.log.Println("[TRACE] GetXAttr:", "name:", name, "attr:", attr)
	return fs.fs.GetXAttr(name, attr, context)
}

func (fs *traceFs) SetXAttr(name string, attr string, data []byte, flags int, context *fuse.Context) fuse.Status {
	fs.log.Println("[TRACE] SetXAttr:", "name:", name, "attr:", attr)
	return fs.fs.SetXAttr(name, attr, data, flags, context)
}

func (fs *traceFs) ListXAttr(name string, context *fuse.Context) ([]string, fuse.Status) {
	fs.log.Println("[TRACE] ListXAttr:", "name:", name)
	return fs.fs.ListXAttr(name, context)
}

func (fs *traceFs) RemoveXAttr(name string, attr string, context *fuse.Context) fuse.Status {
	fs.log.Println("[TRACE] RemoveXAttr:", "name:", name, "attr:", attr)
	return fs.fs.RemoveXAttr(name, attr, context)
}

func (fs *traceFs) Readlink(name string, context *fuse.Context) (string, fuse.Status) {
	fs.log.Println("[TRACE] Readlink:", "name:", name)
	return fs.fs.Readlink(name, context)
}

func (fs *traceFs) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) fuse.Status {
	fs.log.Println("[TRACE] Mknod:", "name:", name, "mode:", mode, "dev:", dev)
	return fs.fs.Mknod(name, mode, dev, context)
}

func (fs *traceFs) Mkdir(name string, mode uint32, context *fuse.Context) fuse.Status {
	fs.log.Println("[TRACE] Mkdir:", "name:", name, "mode:", mode)
	return fs.fs.Mkdir(name, mode, context)
}

func (fs *traceFs) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	fs.log.Println("[TRACE] Unlink:", "name:", name)
	return fs.fs.Unlink(name, context)
}

func (fs *traceFs) Rmdir(name string, context *fuse.Context) (code fuse.Status) {
	fs.log.Println("[TRACE] Rmdir:", "name:", name)
	return fs.fs.Rmdir(name, context)
}

func (fs *traceFs) Symlink(value string, linkName string, context *fuse.Context) (code fuse.Status) {
	fs.log.Println("[TRACE] Symlink:", "value:", value, "linkName:", linkName)
	return fs.fs.Symlink(value, linkName, context)
}

func (fs *traceFs) Rename(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	fs.log.Println("[TRACE] Rename:", "oldName:", oldName, "newName:", newName)
	return fs.fs.Rename(oldName, newName, context)
}

func (fs *traceFs) Link(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	fs.log.Println("[TRACE] Link:", "oldName:", oldName, "newName:", newName)
	return fs.fs.Link(oldName, newName, context)
}

func (fs *traceFs) Chmod(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	fs.log.Println("[TRACE] Chmod:", "name:", name, "mode:", mode)
	return fs.fs.Chmod(name, mode, context)
}

func (fs *traceFs) Chown(name string, uid uint32, gid uint32, context *fuse.Context) (code fuse.Status) {
	fs.log.Println("[TRACE] Chown:", "name:", name, "uid:", uid, "gid:", gid)
	return fs.fs.Chown(name, uid, gid, context)
}

func (fs *traceFs) Truncate(name string, offset uint64, context *fuse.Context) (code fuse.Status) {
	fs.log.Println("[TRACE] Truncate:", "name:", name, "uid:", "offset:", offset)
	return fs.fs.Truncate(name, offset, context)
}

func (fs *traceFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	flagsAsStr := flagsToText(flags)
	fs.log.Println("[TRACE] Open:", "name:", name, "flags:", flags, "("+flagsAsStr+")")
	return fs.fs.Open(name, flags, context)
}

func (fs *traceFs) OpenDir(name string, context *fuse.Context) (stream []fuse.DirEntry, status fuse.Status) {
	fs.log.Println("[TRACE] OpenDir:", "name:", name)
	return fs.fs.OpenDir(name, context)
}

func (fs *traceFs) OnMount(nodeFs *pathfs.PathNodeFs) {
	fs.log.Println("[TRACE] OnMount:", "nodeFs:", nodeFs)
	fs.fs.OnMount(nodeFs)
}

func (fs *traceFs) OnUnmount() {
	fs.log.Println("[TRACE] OnUnmount:")
	fs.fs.OnUnmount()
}

func (fs *traceFs) Access(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	fs.log.Println("[TRACE] Access:", "name:", name, "mode:", mode)
	return fs.fs.Access(name, mode, context)
}

func (fs *traceFs) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	flagsAsStr := flagsToText(flags)
	fs.log.Println("[TRACE] Create:", "name:", name, "flags:", flags, "("+flagsAsStr+")", "mode:", mode)
	return fs.fs.Create(name, flags, mode, context)
}

func (fs *traceFs) Utimens(name string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	fs.log.Println("[TRACE] Utimens:", "name:", name, "Atime:", Atime, "Mtime:", Mtime)
	return fs.fs.Utimens(name, Atime, Mtime, context)
}

func (fs *traceFs) String() string {
	return "traceFs"
}

func (fs *traceFs) StatFs(name string) *fuse.StatfsOut {
	fs.log.Println("[TRACE] StatFs:", "name:", name)
	return fs.fs.StatFs(name)
}

// flagsToText converts flags given to Create and Open into human-readable strings.
func flagsToText(flags uint32) string {
	// The flags are these:
	//      const (
	//      	O_RDONLY int = syscall.O_RDONLY // open the file read-only.
	//      	O_WRONLY int = syscall.O_WRONLY // open the file write-only.
	//      	O_RDWR   int = syscall.O_RDWR   // open the file read-write.
	//      	O_APPEND int = syscall.O_APPEND // append data to the file when writing.
	//      	O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
	//      	O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist
	//      	O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
	//      	O_TRUNC  int = syscall.O_TRUNC  // if possible, truncate file when opened.
	//      )
	//
	//  Plus maybe additional:
	//
	// O_ACCMODE = 3
	// O_NONBLOCK = 04000
	// O_ASYNC    = 020000
	// O_CLOEXEC = 0   ???

	// Provide a helpful breakdown of flags.
	// TODO(ppanyukov): Surely there is a better way to convert flags into strings? Use cmd/stringer?
	var textFlags = make([]string, 0, 8)

	// These two are exlusive, doh!
	if flags&syscall.O_RDWR != 0 {
		textFlags = append(textFlags, "O_RDWR")
	} else if flags&syscall.O_WRONLY != 0 {
		textFlags = append(textFlags, "O_WRONLY")
	} else {
		textFlags = append(textFlags, "O_RDONLY")
	}

	if flags&syscall.O_APPEND != 0 {
		textFlags = append(textFlags, "O_APPEND")
	}
	if flags&syscall.O_CREAT != 0 {
		textFlags = append(textFlags, "O_CREAT")
	}
	if flags&syscall.O_EXCL != 0 {
		textFlags = append(textFlags, "O_EXCL")
	}
	if flags&syscall.O_SYNC != 0 {
		textFlags = append(textFlags, "O_SYNC")
	}
	if flags&syscall.O_TRUNC != 0 {
		textFlags = append(textFlags, "O_TRUNC")
	}
	if flags&syscall.O_ACCMODE != 0 {
		textFlags = append(textFlags, "O_ACCMODE")
	}
	if flags&syscall.O_NONBLOCK != 0 {
		textFlags = append(textFlags, "O_NONBLOCK")
	}
	if flags&syscall.O_ASYNC != 0 {
		textFlags = append(textFlags, "O_ASYNC")
	}

	// Additional things from /include/uapi/asm-generic/fcntl.h
	// Not sure these are correct actually.
	const (
		O_DSYNC     = 00010000
		O_DIRECT    = 00040000
		O_LARGEFILE = 00100000
		O_DIRECTORY = 00200000
		O_NOFOLLOW  = 00400000
		O_NOATIME   = 01000000
		O_CLOEXEC   = 02000000
		O_PATH      = 010000000
		O_TMPFILE   = 020000000
	)

	if flags&O_DSYNC != 0 {
		textFlags = append(textFlags, "O_DSYNC")
	}
	if flags&O_DIRECT != 0 {
		textFlags = append(textFlags, "O_DIRECT")
	}
	if flags&O_LARGEFILE != 0 {
		textFlags = append(textFlags, "O_LARGEFILE")
	}
	if flags&O_DIRECTORY != 0 {
		textFlags = append(textFlags, "O_DIRECTORY")
	}
	if flags&O_NOFOLLOW != 0 {
		textFlags = append(textFlags, "O_NOFOLLOW")
	}
	if flags&O_CLOEXEC != 0 {
		textFlags = append(textFlags, "O_CLOEXEC")
	}
	if flags&O_PATH != 0 {
		textFlags = append(textFlags, "O_PATH")
	}
	if flags&O_TMPFILE != 0 {
		textFlags = append(textFlags, "O_TMPFILE")
	}

	return strings.Join(textFlags, "|")
}
