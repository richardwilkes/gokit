// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package embedded

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"path"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio"
)

// NewFileSystemFromEmbeddedZip creates a new FileSystem from the contents of
// a zip file appended to the end of the executable. If no such data can be
// found, then 'fallbackLiveFSRoot' is used to return a FileSystem based upon
// the local disk.
//
// To create an embedded zip file, first create your zip file as normal, e.g.
// `zip -9 -r path/to/zip_file path/to/zip`. Build your executable as normal,
// e.g. `go build -o path/to/exe main.go`, then concatenate the zip file to
// the end of your executable, e.g. `cat path/to/zip_file >> path/to/exe`.
// Finally, run `zip -A path/to/exe` on your executable to fix up the offsets.
//
// Deprecated: use Go 1.16's embedded support instead
func NewFileSystemFromEmbeddedZip(fallbackLiveFSRoot string) FileSystem {
	if fs, err := NewEFSFromEmbeddedZip(); err == nil {
		return fs.PrimaryFileSystem()
	}
	return NewLiveFS(fallbackLiveFSRoot)
}

// NewEFSFromEmbeddedZip creates a new EFS from the contents of a zip file
// appended to the end of the executable.
//
// Deprecated: use Go 1.16's embedded support instead
func NewEFSFromEmbeddedZip() (*EFS, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	exeFile, err := os.Open(exePath)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer xio.CloseIgnoringErrors(exeFile)
	fi, err := exeFile.Stat()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	r, err := zip.NewReader(exeFile, fi.Size())
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return NewEFSFromZip(r)
}

// NewEFSFromZip creates a new EFS from the contents of a zip file.
//
// Deprecated: use Go 1.16's embedded support instead
func NewEFSFromZip(zr *zip.Reader) (*EFS, error) {
	files := make(map[string]*File)
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		r, err := f.Open()
		if err != nil {
			return nil, errs.Wrap(err)
		}
		data, err := ioutil.ReadAll(r)
		xio.CloseIgnoringErrors(r)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		name := ToEFSPath(f.Name)
		files[name] = NewFile(path.Base(name), f.Modified, int64(f.UncompressedSize64), false, data)
	}
	return NewEFS(files), nil
}
