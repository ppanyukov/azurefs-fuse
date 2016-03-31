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


Next to do:

```
- blobfs_flat: list blobs a a flat list of files
- blobfs_tree: traditional directory/file-based approach
```

