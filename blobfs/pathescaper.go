package blobfs

import (
	"net/url"
)

// pathEscaper provides escape/unescape methods to take care of characters
// which may not be valid as file names but are valid blob names.
// Typical example: valid blob name "/usr/bin/ls" but invalid file name.
// The escaper would take care of this.
type pathEscaper interface {
	BlobNameToFileName(blobName string) (fileName string, err error)
	FileNameToBlobName(fileName string) (blobName string, err error)
}

// pathEscaperURLQuery is an implementation of pathEscaper which
// uses encodes blob names using URL query encoder. This *should*
// always produce valid file names (hopefully). At least it takes care
// the forward slash.
type pathEscaperURLQuery struct {
}

func (x pathEscaperURLQuery) BlobNameToFileName(blobName string) (fileName string, err error) {
	return url.QueryEscape(blobName), nil
}

func (x pathEscaperURLQuery) FileNameToBlobName(fileName string) (blobName string, err error) {
	return url.QueryUnescape(fileName)
}
