package blobfs

// List blobs in containers as a flat list.
// For now assume this is mounted separatedly as its own thing.
// Ideally we want to use it under containerfs too, but later.

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/ppanyukov/azure-sdk-for-go/storage"
)

// NewFlatBlobFs creates a filesystem that lists containers as directories.
func NewFlatBlobFs(accountContainer string, storageClient storage.Client) pathfs.FileSystem {
	logPrefix := fmt.Sprintf("[flatblobFs]: ")

	result := flatblobFs{
		client:           storageClient.GetBlobService(),
		accountContainer: accountContainer,
		log:              log.New(os.Stderr, logPrefix, log.LstdFlags),
		defaultFuseAttr: fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		},
		pathEscaper: pathEscaperURLQuery{},
	}

	return &result
}

// flatblobFs implements a FileSystem that returns blobs as one big flat list.
type flatblobFs struct {
	client                storage.BlobStorageClient
	log                   *log.Logger
	defaultFuseAttr       fuse.Attr
	defaultListBlobParams storage.ListBlobsParameters
	accountContainer      string
	pathEscaper
}

func (fs *flatblobFs) SetDebug(debug bool) {}

func (fs *flatblobFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	// root is always OK
	if name == "" {
		return &fs.defaultFuseAttr, fuse.OK
	}

	blobName, err := fs.FileNameToBlobName(name)
	if err != nil {
		fs.log.Printf("[ERROR] GetAttr '%s': Could not convert file name to blob name. %s\n", name, err)
		return nil, fuse.EIO
	}

	exists, err := fs.client.BlobExists(fs.accountContainer, blobName)

	if err != nil {
		fs.log.Printf("[ERROR] GetAttr '%s': %s\n", name, err)
		return nil, fuse.EIO
	}

	if exists {
		return &fs.defaultFuseAttr, fuse.OK
	}

	return nil, fuse.ENOENT
}

func (fs *flatblobFs) GetXAttr(name string, attr string, context *fuse.Context) ([]byte, fuse.Status) {
	return nil, fuse.ENOSYS
}

func (fs *flatblobFs) SetXAttr(name string, attr string, data []byte, flags int, context *fuse.Context) fuse.Status {
	return fuse.ENOSYS
}

func (fs *flatblobFs) ListXAttr(name string, context *fuse.Context) ([]string, fuse.Status) {
	return nil, fuse.ENOSYS
}

func (fs *flatblobFs) RemoveXAttr(name string, attr string, context *fuse.Context) fuse.Status {
	return fuse.ENOSYS
}

func (fs *flatblobFs) Readlink(name string, context *fuse.Context) (string, fuse.Status) {
	return "", fuse.ENOSYS
}

func (fs *flatblobFs) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) fuse.Status {
	return fuse.ENOSYS
}

func (fs *flatblobFs) Mkdir(name string, mode uint32, context *fuse.Context) fuse.Status {
	return fuse.ENOSYS
}

func (fs *flatblobFs) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (fs *flatblobFs) Rmdir(name string, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (fs *flatblobFs) Symlink(value string, linkName string, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (fs *flatblobFs) Rename(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (fs *flatblobFs) Link(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (fs *flatblobFs) Chmod(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (fs *flatblobFs) Chown(name string, uid uint32, gid uint32, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (fs *flatblobFs) Truncate(name string, offset uint64, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (fs *flatblobFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	return nil, fuse.ENOSYS
}

func (fs *flatblobFs) OpenDir(name string, context *fuse.Context) (stream []fuse.DirEntry, status fuse.Status) {
	if name != "" {
		return []fuse.DirEntry(nil), fuse.OK
	}

	res, err := fs.client.ListBlobs(fs.accountContainer, fs.defaultListBlobParams)
	if err != nil {
		fs.log.Printf("[ERROR] OpenDir '%s': %s'\n", name, err)
		return nil, fuse.EIO
	}

	blobs := res.Blobs
	stream = make([]fuse.DirEntry, 0, len(blobs))
	for _, blob := range blobs {
		blobName := blob.Name
		fileName, err := fs.pathEscaper.BlobNameToFileName(blobName)
		if err != nil {
			fs.log.Printf("[ERROR] OpenDir cannot translate blob '%s' into valid file name: %s'\n", blobName, err)
			continue
		}

		stream = append(stream, fuse.DirEntry{
			Mode: fuse.S_IFREG | 0644,
			Name: fileName,
		})
	}

	return stream, fuse.OK
}

func (fs *flatblobFs) OnMount(nodeFs *pathfs.PathNodeFs) {
}

func (fs *flatblobFs) OnUnmount() {
}

func (fs *flatblobFs) Access(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	return fuse.OK
}

func (fs *flatblobFs) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	return nil, fuse.ENOSYS
}

func (fs *flatblobFs) Utimens(name string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (fs *flatblobFs) String() string {
	return "flatblobFs"
}

func (fs *flatblobFs) StatFs(name string) *fuse.StatfsOut {
	return nil
}
