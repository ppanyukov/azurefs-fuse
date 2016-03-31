package blobfs

// Tracing of operations on underlying FileSystem

import (
	"fmt"
	"log"
	"os"
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
	return fs.fs.Unlink(name, context)
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
	fs.log.Println("[TRACE] Open:", "name:", name, "flags:", flags)
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
	fs.log.Println("[TRACE] Create:", "name:", name, "flags:", flags, "mode:", mode)
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
