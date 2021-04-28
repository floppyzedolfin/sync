# Sync

My first thought, after reading the assignment, was that my solution, 
whatever it will be, won't be as complete as  other utilities. I won't wrap  
over `rsync`, as explained. However, really early, I wanted to delegate the 
"detect which files have been updated" to some intelligent tool.

`git` would be a great tool to do this: it knows how to
- detect changes
- transfer files
- optimize (in some ways) volumes

After checking with the company, they told me that wrapping `git` wasn't 
interesting, which I agree with. I've extrapolated this and decided to also 
not use `go-git`.

I thus need to write code that
- detects local changes under a directory
- sends changed files to a server "letting it know which files they are"
- tries and reduces the amount of data transferred

Before writing any code, it seems rather simple to simply transfer the 
changed data - whatever hasn't changed won't need an update on the server.

Open questions at this point, before writing any code:
- I'm going to work with Unix. Unix _files_ includes regular files, 
  directories, links, pipes, sockets, doors, devices. Since I've absolutely 
  never used doors and devices, and almost never used sockets and pipes (at 
  least manually), I think I'll limit my code to regular files, directories 
  and links.
- In order to "mock" a network, I'll have two docker images running. The client 
  will mount the watched path, and the server will listen on a port and 
  write to the "reference" path. This means the server will expose an API to 
  "update" the reference.
- There are two topics we want to keep in mind at this point: rights and 
  users. It'll be fairly simple to propagate rights from the client to the 
  server ("755  remains 755", for instance), but I can't / won't create 
  files with userIDs that don't exist on the server. This means the server 
  will probably (given the amount of time I'll spend on this assignment) 
  discard any user-specific information (uid, gid) and replace them with the 
  user I'll use in the server's docker image.
  
## Architecture
I'm writing two services:
- `watcher` will be in charge of listening to a directory 
- `reference` will have a backup copy of the contents of the watched directory

`watcher` won't really be receiving requests. `reference` will expose one 
endpoint, probably `update`, that will be called by `watcher` to notify the 
creation of new data (or update, or deletion thereof). Now, in order to 
communicate, we need a protocol. I'll settle down with HTTP as it's what I 
know, and [gRPC](https://grpc.io/), because I want to use that rather than 
[json](https://www.json.org/json-en.html). I'll need to define the body of 
an `update` call in `reference`'s exposed client. 

Something we must think of upfront is what the `watcher` will do when it 
starts. Should the hypothesis be
- that `reference` already has the up-to-date copy of the watched files? (a)
- that `reference` is empty and needs the full contents of the watched 
  directory? (b)
  
After re-reading the input, we must pay attention to _changes_ rather than 
_files_, which means we can skip the copy of the scrutinized directory. 
We're going with (a).

Finally, the last thing to consider is the way we identify that changes 
happen in the files. There are several ways to do this, for instance:
- regularly scan the files and check they haven't changed
- listen to system signals to be notified of changes

Listening to system signals seems a better approach - indeed, in case of 
humongous directories with thousands of subfiles, the scanning would spend 
most of its time performing checking nothing happened.

So, let's dig into the system notifications - the [fsnotify](https://pkg.go.
dev/github.com/fsnotify/fsnotify) Go lib seems to do the trick of detecting 
local changes, so we'll use it.

For now, let's be binary on this - a file will either have changed, or not 
have changed. There will be no partial change in the first version of this 
code. We'll have plenty of time for that later.

