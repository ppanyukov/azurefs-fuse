# Release v0.1.0-alpha

This is the very first cut with the following implemented
functionality which kinda works :)


```
- containerfs: list containers in blob storage as directories

    supported functionality:
    
        - ls: list containers
        - cd <container_name>
        - mkdir <container_name>: creates the container
        - rmdir <container_name>: safely (!) zap container. Like regular rmdir, only deletes container if it's empty.


- flatblobfs: list blobs in a container as a flat list of files

    supported functionality so far:
    
        - ls: list blobs as files
              The blob names are properly escaped to make valid file
              names so even blobs named '/usr/bin/ls' work just fine.

        - touch <blob_name>: creates an empty blob if does not exist already
              TODO(ppanyukov): implement setting times on blobs.

        - rm <blob_name>: delete blob


    not implemented yet:

        - cat 'some content' > <blob_name>: write something into blob
        - cat <blob_name>: read the contents of a blob

```

There are self-contained binaries which are ready to run.

# Quick start

## Prereqs

1. Linux (tried on CentOS 7 but may as well work in other flavours)
2. FUSE: do `yum install fuse`
3. Grab the binaries


## Usage

This is pretty much uniform for both `containerfs` and `flatblobfs`
with some difference in what each does.


### containerfs

Mounting:

```
mkdir ~/mountpoint

env \
    AZURE_STORAGE_ACCOUNT_NAME="YOUR_STORAGE_ACCOUNT" \
    AZURE_STORAGE_ACCOUNT_KEY="YOUR_KEY" \
    containerfs -trace ~/mountpoint &
```


Unmounting:

```
# Leave ~/mountpoint and make sure all files are closed
fusermount -u ~/mountpoint
```


Using:

```
cd ~/mountpoint
ls
mkdir "container-name"
rmdir "container-name"
```


Limitations etc:

- container names have some restrictions, e.g
  they must be longer than 3 chars and cannot contain certain
  characters like `_`, `+`, `\` and `/`.

- deleting non-empty containers with `rmdir`. Only empty containers can be
  deleted. This follows similar rules for regular directories.

- deleting empty containers and creating again. You may get an error
  like `The specified container is being deleted` if you delete and
  then try to create container with same name again.


### flatblobfs


Mounting: Similar to `containerfs`, except also need to supply the container name.

```
mkdir ~/mountpoint

env \
    AZURE_STORAGE_ACCOUNT_NAME="YOUR_STORAGE_ACCOUNT" \
    AZURE_STORAGE_ACCOUNT_KEY="YOUR_KEY" \
    AZURE_STORAGE_ACCOUNT_CONTAINER="YOUR_CONTAINER_NAME" \
    flatblobfs -trace ~/mountpoint &
```

Unmounting:

```
# Leave ~/mountpoint and make sure all files are closed
fusermount -u ~/mountpoint
```


Using:

```
cd ~/mountpoint
ls
touch foo bar zoo
ls
rm *oo

```


Limitations etc:

- Blob name encoding:

  Blob names can contain characters which are not valid in file names.
  For example, `/folderA/subfolderAA/fileAAA.txt` is a valid blob name but not 
  a valid file name. For this reason the blob names are *encoded*. At the moment
  the encoding is same as the URL query string, which looks to produce
  valid file names regardless of input.
  
  So in this example the blob `/folderA/subfolderAA/fileAAA.txt` will be
  given file name `%2FfolderA%2FsubfolderAA%2FfileAAA.txt`.
  
  There may be better ways in the future.


- Reading and writing blob contents:

  Not implemented (yet) but working on it as the next thing to do.


- File times:

  These are ignored and all blobs will show up with dates `Jan  1  1970`
  or whatever is the default on the system. 
  
  Operations like `touch foo` actually don't update these times either.
  
  This functionaly may be implemented in the future.


- File sizes:

  These are not implemented and all blobs will show size `0`.


- Permissions:

  All blobs are given `644` permissions but that actually doesn't mean much :)


- Other ops like `chmod`, `chown`:

  These are not implemented.

- deleting non-empty containers with `rmdir`. Only empty containers can be
  deleted. This follows similar rules for regular directories.

- deleting empty containers and creating again. You may get an error
  like `The specified container is being deleted` if you delete and
  then try to create container with same name again.


## Bugs, issues etc

Probably lots. 
