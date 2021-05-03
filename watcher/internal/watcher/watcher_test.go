package watcher

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	pb "github.com/floppyzedolfin/sync/replica/replica"
	mock_watcher "github.com/floppyzedolfin/sync/watcher/internal/watcher/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	localDir = "./testdata"
)

func TestWatcher_Watch(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tt := map[string]struct {
		rootDir       string
		actions       []action //actions separated by 100ms
		replicaClient func(controller *gomock.Controller) server
		err           string
	}{
		"add a file": {
			rootDir: "newFile",
			actions: []action{{op: createFile, path: "created_file", contents: "bar"}},
			replicaClient: func(mockCtrl *gomock.Controller) server {
				s := mock_watcher.NewMockserver(mockCtrl)
				s.EXPECT().File(gomock.Any(), &pb.FileRequest{FullPath: "created_file", FullContents: "bar"}).Return(nil, nil).Times(2)
				return s
			},
		},
		"create a dir": {
			rootDir: "newDir",
			actions: []action{{op: createDir, path: "created_dir"}},
			replicaClient: func(mockCtrl *gomock.Controller) server {
				s := mock_watcher.NewMockserver(mockCtrl)
				s.EXPECT().Directory(gomock.Any(), &pb.DirectoryRequest{FullPath: "created_dir"}).Return(nil, nil)
				return s
			},
		},
		"remove a dir": {
			rootDir: "rmDir",
			actions: []action{{op: createDir, path: "created_dir"}, {op: remove, path: "created_dir"}},
			replicaClient: func(mockCtrl *gomock.Controller) server {
				s := mock_watcher.NewMockserver(mockCtrl)
				s.EXPECT().Directory(gomock.Any(), &pb.DirectoryRequest{FullPath: "created_dir"}).Return(nil, nil)
				s.EXPECT().Delete(gomock.Any(), &pb.DeleteRequest{FullPath: "created_dir"}).Return(nil, nil)
				return s
			},
		},
		"1 new file, 1 new dir, 1 new file in dir, update first file, remove both files": {
			rootDir: "a_bit_of_all",
			actions: []action{{op: createFile, path: "created_file", contents: "bar"},
				{op: createDir, path: "created_dir"},
				{op: createFile, path: "created_dir/created_file", contents: "bar"},
				{op: createFile, path: "created_file", contents: "foo"},
				{op: remove, path: "created_dir/created_file"},
				{op: remove, path: "created_file"},
			},
			replicaClient: func(mockCtrl *gomock.Controller) server {
				s := mock_watcher.NewMockserver(mockCtrl)
				s.EXPECT().File(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
				s.EXPECT().Directory(gomock.Any(), gomock.Any()).Return(nil, nil)
				s.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, nil).Times(2)
				return s
			},
		},
		"problem during rm file": {
			rootDir: "rm_file_server_issue",
			actions: []action{{op: createFile, path: "created_file", contents: "bar"}, {op: remove, path: "created_file"}},
			replicaClient: func(mockCtrl *gomock.Controller) server {
				s := mock_watcher.NewMockserver(mockCtrl)
				s.EXPECT().File(gomock.Any(), gomock.Any()).Return(nil, nil).Times(2)
				s.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("server delete error"))
				return s
			},
			err: "context cancelled",
		},
		"problem during file patch": {
			rootDir: "patch_file_server_issue",
			actions: []action{{op: createFile, path: "created_file", contents: "bar"}},
			replicaClient: func(mockCtrl *gomock.Controller) server {
				s := mock_watcher.NewMockserver(mockCtrl)
				s.EXPECT().File(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("server patch error")).AnyTimes()
				return s
			},
			err: "context cancelled",
		},
		"problem during dir creation": {
			rootDir: "rm_dir_server_issue",
			actions: []action{{op: createDir, path: "created_dir"}},
			replicaClient: func(mockCtrl *gomock.Controller) server {
				s := mock_watcher.NewMockserver(mockCtrl)
				s.EXPECT().Directory(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("server create dir error"))
				return s
			},
			err: "context cancelled",
		},
		"add and remove a link": {
			rootDir: "links",
			actions: []action{{op: createLink, path: "created_link", contents: "foo"}, {op: remove, path: "created_link"}},
			replicaClient: func(mockCtrl *gomock.Controller) server {
				s := mock_watcher.NewMockserver(mockCtrl)
				s.EXPECT().Link(gomock.Any(), &pb.LinkRequest{FullPath: "created_link", Target: "foo"}).Return(nil, nil)
				s.EXPECT().Delete(gomock.Any(), &pb.DeleteRequest{FullPath: "created_link"}).Return(nil, nil)
				return s
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			// some setup
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			ctx, cancel := context.WithCancel(context.Background())
			s, err := New(tc.replicaClient(mockCtrl), cancel)
			require.NoError(t, err)
			testDir := path.Join(localDir, tc.rootDir)
			err = os.Mkdir(testDir, 0755)
			require.NoError(t, err)
			defer os.RemoveAll(testDir)
			defer cancel()

			// start watching
			var watcherErr error
			go func() { watcherErr = s.Watch(ctx, testDir) }()

			// give the routines some time to start
			time.Sleep(100 * time.Millisecond)

			for _, a := range tc.actions {
				switch a.op {
				case createFile:
					err = ioutil.WriteFile(path.Join(testDir, a.path), []byte(a.contents), 0755)
					require.NoError(t, err)
				case createDir:
					err = os.Mkdir(path.Join(testDir, a.path), 0755)
					require.NoError(t, err)
				case createLink:
					err = os.Symlink(path.Join(testDir, a.contents), path.Join(testDir, a.path))
				case remove:
					err = os.Remove(path.Join(testDir, a.path))
					require.NoError(t, err)
				}
				time.Sleep(100 * time.Millisecond)
			}

			if tc.err != "" {
				assert.Error(t, watcherErr)
				assert.Contains(t, watcherErr.Error(), tc.err)
			}
		})
	}
}

type action struct {
	path     string
	op       operation
	contents string
}

type operation int

const (
	createFile operation = iota
	createDir
	createLink
	remove
)
