Exposes Azure (blob) storage via FUSE so it can be mounted
as a regular file system, similar to https://github.com/s3fs-fuse/s3fs-fuse
except for Azure and in golang.

Work in progress, mostly to see how far I can get until I
hit a wall.

For now there are the following file systems:

```
- containerfs: list containers in blob storage as directories

    supported functionality:
    
        - ls: list containers
        - cd <container_name>
        - mkdir <container_name>: creates the container
        - rmdir <container_name>: safely (!) zap container. Like regular rmdir, only deletes container if it's empty.

```


In progress:

```
- flatblobfs: list blobs in a container as a flat list of files

    supported functionality so far:
    
        - ls: list blobs as files
              The blob names are properly escaped to make valid file
              names so even blobs named '/usr/bin/ls' work just fine.

        - touch <blob_name>: creates an empty blob if does not exist already
              TODO(ppanyukov): implement setting times on blobs.


    next steps in this order:

        - rm <blob_name>: delete blob
        - cat 'some content' > <blob_name>: write something into blob
        - cat <blob_name>: read the contents of a blob
```


Next to do:

```
- treeblobfs: traverse blobs in a container in a traditional directory/file-based way
```

