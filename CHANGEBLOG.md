# Sync

### First thoughts

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
endpoint, probably `update`, that will be called by `watcher` to patch the 
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

So, let's dig into the system notifications - the [fspatch](https://pkg.go.
dev/github.com/fspatch/fspatch) Go lib seems to do the trick of detecting 
local changes, so we'll use it.

For now, let's be binary on this - a file will either have changed, or not 
have changed. There will be no partial change in the first version of this 
code. We'll have plenty of time for that later.

### Second thoughts
After writing the skeleton of the server, I realised there is no real reason 
for the `watcher` to be a service at all. It simply needs to be a binary. 
Since I haven't started writing that part, it's not work lost. What's not 
really nice, though, is the way I'll have to handle the `reference`'s local 
copy path. I think I'll use a `/data/sync/` as the root dir in the docker 
image. This should be fine.

### Next thoughts
I've written the Makefile. I have some issues with the passing of the rights 
- what the world understands as "755", or "600", or "777" are octal numbers, 
  which comes as an issue when trying to store them as uints. Because 0755 
  == 493, which isn't really an interesting number (unless you enjoy [strange 
  banquets](https://en.wikipedia.org/wiki/Ostrogothic_Kingdom#Theodoric_kills_Odoacer_(493))).
  
### Get some rest
After a good night of sleep, some questions arose:
- if two files are updated simultaneously, do we want to make two requests 
  to the server, or a single one? This can happen very easily, the first 
  example I have in mind being `mv foo bar`, which `renames` foo and 
  `creates` bar;
- if a single file is updated twice in a row, do we want to shoot two 
  messages, or a single one, with the "final" result? My example here is 
  `echo foo > bar`, which first `creates` bar then  `modifies` it (we don't 
  want to ignore creation of files, that would be too dangerous);
- if N files are created at the same time, do I really want to shoot N 
  requests to my server? Requests to the server must always be in order - 
  that is, we don't want the messages of "override `foo` with `bar`" and 
  "override `foo` with `zip`" to happen in a different order on the server.
  
After all these thoughts, I had several choices in mind:
- keep using fspatch and shoot a request for each update
- keep using fspatch, but use a "buffer" to store updates for N
  milliseconds, and shoot a message containing all these N millisecond changes
- use `os.Lstat` in a `for` loop to detect changes  
 
I've worked on this for a couple of hours now and the current version is 
"operational". That is, I managed to have a filed copied from my playground 
to my replica. There are still some "issues" with the current implementation, -
for instance, I don't handle links properly. I'm also always sending the 
full contents of the files to patch, rather than an optimised `diff`.  Since 
I'm not keeping local copies of the files, I'd have to use the rolling hash 
suggested by the assignment. This will greatly depend on the amount of time 
I can still spend on this task.

My current implementation could be missing some files. Indeed, if a file is 
created in a directory that isn't yet under watch, that file's creation 
won't be caught. This could happen as a result of `mkdir foo && touch 
foo/bar`, with my code starting to listen at `foo` after the creation of 
foo/bar. 

I need some more tests than those I manually performed.