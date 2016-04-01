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
		defaultDirFuseAttr: fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		},
		defaultFileFuseAttr: fuse.Attr{
			Mode: fuse.S_IFREG | 0644,
		},
		pathEscaper: pathEscaperURLQuery{},
	}

	return &result
}

// flatblobFs implements a FileSystem that returns blobs as one big flat list.
type flatblobFs struct {
	client                storage.BlobStorageClient
	log                   *log.Logger
	defaultDirFuseAttr    fuse.Attr
	defaultFileFuseAttr   fuse.Attr
	defaultListBlobParams storage.ListBlobsParameters
	accountContainer      string
	pathEscaper
}

func (fs *flatblobFs) SetDebug(debug bool) {}

func (fs *flatblobFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	// root is always OK
	if name == "" {
		return &fs.defaultDirFuseAttr, fuse.OK
	}

	blobName, err := fs.pathEscaper.FileNameToBlobName(name)
	if err != nil {
		fs.log.Printf("[ERROR] GetAttr '%s': Could not convert file name to blob name. %s\n", name, err)
		// TODO(ppanyukov): is this correct status to return?
		return nil, fuse.EINVAL
	}

	exists, err := fs.client.BlobExists(fs.accountContainer, blobName)

	if err != nil {
		fs.log.Printf("[ERROR] GetAttr '%s': %s\n", name, err)
		return nil, fuse.EIO
	}

	if exists {
		// NOTE: all entries are files in this flat view.
		return &fs.defaultFileFuseAttr, fuse.OK
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
	// This one is called when we do `touch zzz`. Here is where we want to create blobs.
	// Must also minimally implement Utimens otherwise OS reports IO error.
	// No need to implement Open, everything works just fine without it.
	// The sequence is like this:
	//      [flatblobFs]: 2016/04/01 09:50:39 [TRACE] GetAttr: name: zzz
	//      [flatblobFs]: 2016/04/01 09:50:39 [TRACE] Mknod: name: zzz mode: 33188 dev: 0
	//      [flatblobFs]: 2016/04/01 09:50:39 [TRACE] GetAttr: name: zzz
	//      [flatblobFs]: 2016/04/01 09:50:40 [TRACE] GetAttr: name: zzz
	//      [flatblobFs]: 2016/04/01 09:50:40 [TRACE] Open: name: zzz flags: 34817
	//      [flatblobFs]: 2016/04/01 09:50:40 [TRACE] Utimens: name: zzz Atime: 2016-04-01 09:50:40.066466083 +0000 UTC Mtime: 2016-04-01 09:50:40.066466083 +0000 UTC
	//      [flatblobFs]: 2016/04/01 09:50:40 [TRACE] GetAttr: name: zzz
	blobName, err := fs.pathEscaper.FileNameToBlobName(name)
	if err != nil {
		fs.log.Printf("[ERROR] Mknod '%s': Could not convert file name to blob name. %s\n", name, err)
		// TODO(ppanyukov): is this correct status to return?
		return fuse.EINVAL
	}

	// Assume that if we get to here the OS has already checked that
	// this file does not exist. However because it's a remote multi-user
	// system, there is always a chance it appeared in the meantime.
	// TODO(ppanyukov): how does azure handle create blob request if blob exists?
	err = fs.client.CreateBlockBlob(fs.accountContainer, blobName)
	if err != nil {
		fs.log.Printf("[ERROR] Mknod '%s': Could not create blob. %s\n", name, err)
		return fuse.EIO
	}

	return fuse.OK
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

	// Preallocate the array with capacity equal to the number of blobs
	// but set initial len to 0 and grow as needed.
	// Reason is there may be blobs which we can't translate to file names
	// due to bugs or escaping issues and so may end up with fewer files than
	// there are blobs.
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
	// TODO(ppanyukov): Meaningful implementatin of Utimens. For now just return OK.
	// This is so other things like `touch foo` work without errors. See Mknod.
	return fuse.OK
}

func (fs *flatblobFs) String() string {
	return "flatblobFs"
}

func (fs *flatblobFs) StatFs(name string) *fuse.StatfsOut {
	return nil
}
