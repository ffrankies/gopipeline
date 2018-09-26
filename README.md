# gopipeline

A pipelined approach to distributed programming

## Project Structure

After reading all of [this](https://dave.cheney.net/2014/12/01/five-suggestions-for-setting-up-a-go-project), [this](https://talks.golang.org/2014/organizeio.slide#1), [this](https://github.com/golang/go/wiki/PackagePublishing), [this](https://medium.com/rungo/everything-you-need-to-know-about-packages-in-go-b8bac62b74cc) and [this](https://github.com/golang-standards/project-layout), it seems that there isn't a standard for organizing go projects. So, I figured the best way to organize would be to copy [what the devs of golang did](https://github.com/golang/tools).

The basic idea is, every functional unit of the project (stuff you'd import), is its own package, in its own directory. The main question here, is whether or not we want to move `master/` and `worker/` into `gopipeline/` or not, as subpackages.

A common alternative I'm seeing is to have the repo be named `gopipeline`, and then the code that would go into `gopipeline/` is moved into the root directory of the project.

If we want this to be an executable (i.e.: `gopipeline [args]` instead of importing `github.com/[user]/golang` from inside their codebase), then the root directory will either be used for the command (file with package "main") or the command will be moved into `cmd/`.

Decisions, decisions...
