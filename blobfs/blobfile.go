package blobfs

import (
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

// blobFile implements fuse/nodefs/File interface to
// read/write data from/to blobs.
type blobFile struct{}

// NewBlobFile returns a File instance that returns ENOSYS for
// every operation.
func NewBlobFile() nodefs.File {
	return (*blobFile)(nil)
}

// SetInode Called upon registering the filehandle in the inode.
func (f *blobFile) SetInode(*nodefs.Inode) {
}

// Wrappers around other File implementations, should return
// the inner file here
func (f *blobFile) InnerFile() nodefs.File {
	return nil
}

// The String method is for debug printing.
func (f *blobFile) String() string {
	return "blobFile"
}

func (f *blobFile) Read(buf []byte, off int64) (fuse.ReadResult, fuse.Status) {
	return nil, fuse.ENOSYS
}

func (f *blobFile) Write(data []byte, off int64) (uint32, fuse.Status) {
	return 0, fuse.ENOSYS
}

// Flush is called for close() call on a file descriptor. In
// case of duplicated descriptor, it may be called more than
// once for a file.
func (f *blobFile) Flush() fuse.Status {
	return fuse.OK
}

// This is called to before the file handle is forgotten. This
// method has no return value, so nothing can synchronizes on
// the call. Any cleanup that requires specific synchronization or
// could fail with I/O errors should happen in Flush instead.
func (f *blobFile) Release() {

}

func (f *blobFile) GetAttr(*fuse.Attr) fuse.Status {
	return fuse.ENOSYS
}

func (f *blobFile) Fsync(flags int) (code fuse.Status) {
	return fuse.ENOSYS
}

func (f *blobFile) Utimens(atime *time.Time, mtime *time.Time) fuse.Status {
	return fuse.ENOSYS
}

// The methods below may be called on closed files, due to
// concurrency.  In that case, you should return EBADF.
func (f *blobFile) Truncate(size uint64) fuse.Status {
	return fuse.ENOSYS
}

func (f *blobFile) Chown(uid uint32, gid uint32) fuse.Status {
	return fuse.ENOSYS
}

func (f *blobFile) Chmod(perms uint32) fuse.Status {
	return fuse.ENOSYS
}

func (f *blobFile) Allocate(off uint64, size uint64, mode uint32) (code fuse.Status) {
	return fuse.ENOSYS
}
