# Assignment

## "Dropbox" Homework:

Build an application to synchronize a source folder and a destination folder over IP:

1.1 - a simple command-line client which takes one directory as an argument and keeps monitoring changes in that directory and uploads any change to its server

1.2 - a simple server which takes one empty directory as an argument and receives any change from its client

Bonus - optimize data transfer by avoiding uploading the same file multiple times

Bonus - optimize data transfer by avoiding uploading the same content multiple times content (files which fully or partially share the same content)

## Constraints
Please try to avoid high-level libraries such as librsync

## FAQ

Q: Am I allowed to use diff and 3-way merge libraries similar to: https://docs.python.org/3.6/library/difflib.html https://pypi.org/project/merge3/

A: There are bonuses associated with the task that require to implement some sort of binary diff mechanism or approach of reusing data from existing files. Libraries like rsync or some linux utilities or high level diff libraries just give you such functionality for free. We don't want you to use these libraries. It is OK if you read a paper about rsync algorithm and implement it yourself or go with something simpler but your own approach. The key here is that we want you to have a firm understanding of how it works and you to be able to implement it yourself. Bandwidth saving optimisations should be implemented by you without using libraries. For any other functionality like file monitoring, network transport & protocol etc you can use third party libraries. Also, I want to mention, just because a lot of people ask about it, that the homework should work with any file types - binary or text it should not matter.

Q: How deep should I go with this project? Even such a (small) project could potentially become quite complex. It depends on level of details.

A: It is a very good question. The possibilities to "go deep" in this homework are endless. It is very interesting for us to see where a candidate decides to focus and where they decide to cut corners. If you want to, feel free to demonstrate your "strong suit": if you're a multithreading expert, put emphasis on parallel implementation; if you're into security -- implement something extra about that, and so on. That is why it is proven to be very useful for us. It is expected that you will not be able to implement everything you have on your mind. Just tell us about those areas where you would like things to improve. Keep that in mind.

# Implementation

You'll find a lot of though process in the [CHANGEBLOG.md](./CHANGEBLOG.md) 
file, where I've let my quill slide on the keyboard. Choices, hesitations, 
doubts - they all dwell over there.

## How to?

### ... get it
Clone the repo to your machine:
```bash
git clone https://github.com/floppyzedolfin/sync.git
```

### ... compile it
Use the `make` target:
```bash
make compile
```

### ... run it
We've got two independent things here. A server and a client. You can run 
both of them the way you want, but there's a target that will pop the server,
and then start the client.
```bash
make run-all
```
The [`Makfeile`](./Makefile) provides several other (internal) targets. If 
you examine it closely, you'll find that
- a docker network is started, to simulate real networks and communication 
  to a remote server;
- for debugging, testing, and validation, the server write to its docker 
  image, but the path within that docker image is actually mounted from this 
  machine - this way, one can easily check that files have been updated;
- there is an example playground which is the directory the server will 
  mirror if you use the default `make` targets.

### ... play with it
Once everything is up and running, simply add, alter, move, delete files in 
the directory that the `watcher` was given to examine. These changes should be 
mirrored n the `reference` directory (because the server will update them in 
its docker image).
