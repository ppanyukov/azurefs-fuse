package blobfs

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/ppanyukov/azure-sdk-for-go/storage"
)

// TODO(ppanyukov): performance, caching etc, once we have stuff working :)

// NewContainerFs creates a filesystem that lists containers as directories.
func NewContainerFs(storageClient storage.Client) pathfs.FileSystem {
	logPrefix := fmt.Sprintf("[containerfs]: ")

	result := containerFs{
		client: storageClient.GetBlobService(),
		log:    log.New(os.Stderr, logPrefix, log.LstdFlags),
		defaultListContainersParameters: storage.ListContainersParameters{},
		defaultFuseAttr: fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		},
	}

	return &result
}

// containerFs implements a FileSystem that returns blob container names as directories.
type containerFs struct {
	client                          storage.BlobStorageClient
	defaultListContainersParameters storage.ListContainersParameters
	defaultFuseAttr                 fuse.Attr
	log                             *log.Logger
}

func (fs *containerFs) SetDebug(debug bool) {}

// containerFs implementation.
// At the minimum need to implement:
//   - OpenDir
//   - GetAttr
//   - Access
func (fs *containerFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	// root is always OK
	if name == "" {
		return &fs.defaultFuseAttr, fuse.OK
	}

	if strings.Contains(name, "/") {
		return nil, fuse.ENOENT
	}

	exists, err := fs.client.ContainerExists(name)

	if err != nil {
		fs.log.Printf("[ERROR] GetAttr '%s': %s'\n", name, err)
		return nil, fuse.EIO
	}

	if exists {
		return &fs.defaultFuseAttr, fuse.OK
	}

	return nil, fuse.ENOENT
}

func (fs *containerFs) GetXAttr(name string, attr string, context *fuse.Context) ([]byte, fuse.Status) {
	return nil, fuse.ENOSYS
}

func (fs *containerFs) SetXAttr(name string, attr string, data []byte, flags int, context *fuse.Context) fuse.Status {
	return fuse.ENOSYS
}

func (fs *containerFs) ListXAttr(name string, context *fuse.Context) ([]string, fuse.Status) {
	return nil, fuse.ENOSYS
}

func (fs *containerFs) RemoveXAttr(name string, attr string, context *fuse.Context) fuse.Status {
	return fuse.ENOSYS
}

func (fs *containerFs) Readlink(name string, context *fuse.Context) (string, fuse.Status) {
	return "", fuse.ENOSYS
}

func (fs *containerFs) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) fuse.Status {
	return fuse.ENOSYS
}

func (fs *containerFs) Mkdir(name string, mode uint32, context *fuse.Context) fuse.Status {
	err := fs.client.CreateContainer(name, storage.ContainerAccessTypePrivate)
	if err != nil {
		fs.log.Printf("[ERROR] GetAttr '%s': %s'\n", name, err)
		return fuse.EIO
	}
	return fuse.OK
}

func (fs *containerFs) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (fs *containerFs) Rmdir(name string, context *fuse.Context) (code fuse.Status) {
	err := fs.client.DeleteContainer(name)
	if err != nil {
		fs.log.Printf("[ERROR] Rmdir '%s': %s'\n", name, err)
		return fuse.EIO
	}
	return fuse.OK
}

func (fs *containerFs) Symlink(value string, linkName string, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (fs *containerFs) Rename(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	// renaming containers is not directly supported
	return fuse.ENOSYS
}

func (fs *containerFs) Link(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (fs *containerFs) Chmod(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (fs *containerFs) Chown(name string, uid uint32, gid uint32, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (fs *containerFs) Truncate(name string, offset uint64, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (fs *containerFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	return nil, fuse.ENOSYS
}

func (fs *containerFs) OpenDir(name string, context *fuse.Context) (stream []fuse.DirEntry, status fuse.Status) {
	if name != "" {
		return []fuse.DirEntry(nil), fuse.OK
	}

	res, err := fs.client.ListContainers(fs.defaultListContainersParameters)
	if err != nil {
		fs.log.Printf("[ERROR] OpenDir '%s': %s'\n", name, err)
		return nil, fuse.EIO
	}

	containers := res.Containers
	stream = make([]fuse.DirEntry, len(res.Containers))
	for i, container := range containers {
		stream[i] = fuse.DirEntry{
			Mode: fuse.S_IFDIR | 0755,
			Name: container.Name,
		}
	}

	return stream, fuse.OK
}

func (fs *containerFs) OnMount(nodeFs *pathfs.PathNodeFs) {
}

func (fs *containerFs) OnUnmount() {
}

func (fs *containerFs) Access(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	// TODO(ppanyukov): what is the meaningful implementation for this?
	return fuse.OK
}

func (fs *containerFs) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	return nil, fuse.ENOSYS
}

func (fs *containerFs) Utimens(name string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (fs *containerFs) String() string {
	return "containerFs"
}

func (fs *containerFs) StatFs(name string) *fuse.StatfsOut {
	return nil
}
