package internal

// The convention within this package is to create entities with names starting with created_.

import (
	"context"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"testing"

	pb "github.com/floppyzedolfin/sync/replica/replica"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const localTestDir = "./test_data"

func TestServer_Directory(t *testing.T) {
	tt := map[string]struct {
		req pb.DirectoryRequest
		err string
	}{
		"nominal": {
			req: pb.DirectoryRequest{
				FullPath: "created_foo",
			},
		},
		"subdir": {
			req: pb.DirectoryRequest{
				FullPath: "dir/created_foo",
			},
		},
		"replacing a file with a dir": {
			// this scenario could happen if the user were to `mv file file_bkp && mkdir file`
			req: pb.DirectoryRequest{FullPath: "file"},
			err: "unable to create dir",
		},
		"outside scope": {
			req: pb.DirectoryRequest{
				FullPath: "../../../../../../../../../../../../../created_mi",
			},
			err: "unable to create dir",
		},
	}
	s := NewServer(localTestDir)
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			_, err := s.Directory(context.Background(), &tc.req)
			if tc.err != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.err)
			} else {
				require.NoError(t, err)
				fullPath := path.Join(localTestDir, tc.req.FullPath)
				stat, err := os.Lstat(fullPath)
				assert.NoError(t, err)
				assert.True(t, stat.IsDir())
				os.Remove(fullPath)
			}
		})
	}
}

func TestServer_Delete(t *testing.T) {
	tt := map[string]struct {
		createFile string
		Directory  string
		req        pb.DeleteRequest
		err        string
	}{
		"remove file": {
			createFile: "created_file",
			req:        pb.DeleteRequest{FullPath: "created_file"},
		},
		"remove dir": {
			Directory: "created_dir",
			req:       pb.DeleteRequest{FullPath: "created_dir"},
		},
		"target doesn't exist": {
			req: pb.DeleteRequest{FullPath: "created_entity"},
			err: "unable to delete at",
		},
	}

	s := NewServer(localTestDir)
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			if tc.createFile != "" {
				ioutil.WriteFile(path.Join(localTestDir, tc.createFile), []byte("hello"), 0755)
			}
			if tc.Directory != "" {
				os.Mkdir(path.Join(localTestDir, tc.Directory), 0755)
			}
			_, err := s.Delete(context.Background(), &tc.req)
			if tc.err != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.err)
			} else {
				require.NoError(t, err)
				_, err = os.Stat(path.Join(localTestDir, tc.req.FullPath))
				assert.True(t, os.IsNotExist(err))
			}
		})
	}
}

func TestServer_File(t *testing.T) {
	tt := map[string]struct {
		req pb.FileRequest
		err string
	}{
		"nominal": {
			req: pb.FileRequest{FullPath: "foo", FullContents: "bar"},
		},
		"already a dir at that path": {
			req: pb.FileRequest{FullPath: "dir", FullContents: "abc"},
			err: "unable to  file",
		},
	}

	s := NewServer(localTestDir)
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			_, err := s.File(context.Background(), &tc.req)
			if tc.err != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.err)
			} else {
				require.NoError(t, err)
				filePath := path.Join(localTestDir, tc.req.FullPath)
				fileContents, _ := os.ReadFile(filePath)
				assert.Equal(t, tc.req.FullContents, string(fileContents))
				os.Remove(filePath)
			}
		})
	}
}

func TestServer_Link(t *testing.T) {
	tt := map[string]struct{
		req pb.LinkRequest
		err string
	} {
		"link to another file": {
			req: pb.LinkRequest{FullPath: "created_link", Target: "file"},
		},
		"there's a bear in my bed -- a file already exists there": {
			req: pb.LinkRequest{FullPath: "file", Target: "created_link"},
			err: "unable to create link",
		},
	}

	s := NewServer(localTestDir)
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			_, err := s.Link(context.Background(), &tc.req)
			if tc.err != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.err)
			} else {
				require.NoError(t, err)
				fullPath := path.Join(localTestDir, tc.req.FullPath)
				stat, err := os.Lstat(fullPath)
				assert.NoError(t, err)
				assert.True(t, stat.Mode() & fs.ModeSymlink == fs.ModeSymlink)
				os.Remove(fullPath)
			}
		})
	}
}
