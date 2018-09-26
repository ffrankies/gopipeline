# gopipeline

A pipelined approach to distributed programming

## Project Structure

After reading all of [this](https://dave.cheney.net/2014/12/01/five-suggestions-for-setting-up-a-go-project), [this](https://talks.golang.org/2014/organizeio.slide#1), [this](https://github.com/golang/go/wiki/PackagePublishing), [this](https://medium.com/rungo/everything-you-need-to-know-about-packages-in-go-b8bac62b74cc) and [this](https://github.com/golang-standards/project-layout), it seems that there isn't a standard for organizing go projects.

At first, I figured the best way to organize would be to copy [what the devs of golang did](https://github.com/golang/tools). However, that repo is fundamentally different from ours. They were building a collection of tools - we're building just the one. So, the end result currently looks like this:

The main gopipeline code goes into the root directory of the project. Master and Worker will be treated as subpackages, saved in the `$root/master/` and `$root/worker` directories.

## Spec

### Executable vs Library

If we want this to be an executable (i.e.: `gopipeline [args]` instead of importing `github.com/[user]/golang` from inside their codebase), then the root directory will either be used for the command (file with package "main") or the command will be moved into `cmd/`. However, I personally think this approach would be more difficult to achieve (and wouldn't be a good idea in general), because it will require that either the code to be distributed & pipelined is interpreted (kind of like how pyspark works), or that the code ends up having to be broken into segments that are accessible via command-line, and then the number of segments should be somehow available to the pipeline command. Neither of those are 'good' solutions (and I'm not even sure the first one is possible), so I would go with a 'library approach'.

The 'library' approach would most likely look something like the `multiprocessing` module in python.

#### Python multiprocessing Example

```
import multiprocessing

def foo():
    ...

pool = multiprocessing.Pool(...)
result = pool.apply(f)
pool.close()
```

#### gopipeline Example (in progress)

```
import "github.com/ffrankies/gopipeline"

func foo() {
    ...
}

fun bar() {
    ...
}

func main() {
    funcs := []func{foo, bar}
    results := []T{foo, bar}
    gopipeline.pipeline(funcs, "/path/to/nodelist", results)
}
```

### The gopipeline package

#### gopipeline

Contains the logic that starts off a master process, and receives output from the last process. We may or may not want to treat the master process as the calling process (i.e., do not create a new process for master. Instead, just run the code in the current process)

#### gopipeline/master

Contains the logic of the master process. Initial idea is something like this:

```
nodelist = read nodelist
module_parts = module.parts  // may want to come up with a better name than 'parts'
map module_parts to available nodes  // does not have to be sophisticated at this point, can use round-robin
for each mapped node:
    set a port number, and the next node  // needs to be completed before next loop
for each mapped node:
    ssh into node
    install `github.com/ffrankies/gopipeline` package on node, if not already installed
    copy data if first worker
    start worker, passing in part to be worked on, worker's socket, next worker's ip and next worker's socket
wait for output from last node
```

#### gopipeline/worker

Contains the logic for actually running a code segment on the node. Initial idea is something like this:

```
module_parts = module.parts
establish incoming and outgoing socket connections
while True:
    wait for input data
    run designated part  // can use index in module_parts
    send output on outgoing connection
```
