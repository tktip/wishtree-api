# Makefile

This pagd describes the different approaches you can have to [building](#building), [running](#running) and [testing](#testing) your app.

## Building

Building a binary from your go application can be done by simply invoking the `build`-goal:
```bash
make build
```

This will build your project and create a new amd64-binary in the `bin`-directory.

If you get a permission-error it likely means that the go-tool is unable to fetch the dependencies hosted in our private repositories. Since the go-tool uses git to fetch dependencies, we can simply tell the git-client to always use ssh instead of https, and set up an ssh-token granting our git-client access to our bitbucket-account. Adding an ssh-token is done through [the bitbucket settings page](https://bitbucket.org/account/settings/ssh-keys/).

To tell the git-client to always use ssh towards bitbucket you have to run the following command:
```bash
git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/"
```

## Running

Running your project can be done by simply invoking:
```bash
make run
```

Since this goal depends on the [`build`](#building)-goal, it will also build your binary before running it.

## Testing
Testing can be done using the `test`-goal:
```bash
make test
```

This command runs the projects go-tests, [gofmt](https://golang.org/cmd/gofmt/), [govet](https://golang.org/cmd/vet/) and [linting (revive)](https://github.com/mgechev/revive)


## Live reloading (watch)

The Go-template supports multiple variants of watching your files and live-reloading your running service on file change.

The first and simplest approach is using the `watch`-goal:
```bash
make watch
```

This starts your service using the [`run`](#running)-goal, and watches for any changes to `go`-files in your project.
Whenever a change is detected, it restarts the [`run`](#running)-process with the new changes.
If a compilation error occurrs you can see the error-output in the process window where you ran the `make watch` process.

### watch-tests

`watch-tests` works like the normal [`watch`](#live-reloading-watch), except it runs the [`test`](#testing) goal instead of `build run`.
```bash
make watch-tests
```

### watch-all

`watch-all` is an extension of [`watch`](#live-reloading-watch).
It runs tests, generates documentation, and build and runs your service.

```bash
make watch-all
```

It is the equivalent of this make-command any time a documentation- or go-file changes:
```bash
make test generate-doc run
```

## Cleaning up

Through using this project's various make-goals, some hidden and git-ignored files may start to build up.
Mainly the `.go` and `bin` directories, but also some files named `.container*` and `.dockerfile*`.

Running this goal will remove all of them.

```bash
make clean
```

This goal is rarely _necessary_ to run manually.

## Documentation

There are mainly 2 make-goals explicitly created for writing documentation:
* [`watch-doc`](#watchdoc): intended to be used while the developer is writing documentation. 
It is a live-reload goal comparable to js-webpack watch, watching your documentation source-files.
* [`generate-doc`](#generatedoc): mainly intented to be used when building the project for production, and not for the developer while writing documentation

[`watch-all`](#watchall) also builds doc on change.

### watch-doc

`watch-doc` watches for changes in the `md` directory and regenerates the `internal/docfs/resources.go`, (essentially re-running [`generate-doc`](#generatedoc)) then restarting `main.go`.
This in turn updates the documentation-files exposed by your server.
Changes made to the `md`-files can be seen after reloading the browser-page.

```bash
make watch-doc
```

### generate-doc

`generate-doc` generates the `internal/docfs/resources.go`-file from the contents of the `md`-directory in the project's root directory.
```bash
make generate-doc
```

## Building for docker environments

The Makefile has 2 goals related to building docker-images from your application. On goal to [build the container](#building-a-container)
and another to [push the container to a registry](#pushing-to-a-registry)

### Building a container

Building a container is done by invoking the `container`-goal:
```bash
make container
```

This goal will build a docker container using the pseudo Dockerfile named `Dockerfile.in` located at project root.
How the image and tag is calculated is described [here](#generating-the-image-and-tag).

If you would like to change the Dockerfile, please try to understand `container`-goal works in the Makefile itself.
As in how the replacing works of the various values (ARG_BIN, ARG_ARCH...)

#### Generating the image and tag

The following information is documentation on how the generation works behind the scenes, and is not required knowledge to use the template.

The `container`-goal builds a docker image from your application.
The image name and tag is built from the variables `REGISTRY` and `BIN` specified in the envfile (located in the project root), structured like this:
```make
IMAGE := $(REGISTRY)/$(BIN)
```

The tag is more complicated. The root of the tag is gotten with a `git`-command like so:
```bash
git describe --tags --always --dirty
```
This command takes into account the git hash, the current or most recent version-tag, and if the repo is *dirty* or not (a repo is considered dirty if it has uncommitted changes).
The repo's current *dirty*-status will be reflected in the tag; "dirty" will de suffixed on to the tag.

An example:
```
tipdocreg.trondheim.kommune.no/go-template:v1.1.4-local
```

If the repo is dirty the tag might look something like the following:
```
tipdocreg.trondheim.kommune.no/go-template:v1.1.4-dirty-local
```

`local` means that this build was made on a local developer's machine, not on a CI-server.

If the suffix is `dev`, it means that the build was made on a CI-server, but on the `develop`-branch:
```
tipdocreg.trondheim.kommune.no/go-template:v1.1.4-dev
```

On hotfix-branches, the tag will have a `hotfix`-suffix:
```
tipdocreg.trondheim.kommune.no/go-template:v1.1.4-hotfix
```

And finally, builds made on the master-branch on a CI-server will have no suffix at all:
```
tipdocreg.trondheim.kommune.no/go-template:v1.1.4
```

If the current commit is not tagged, then `git describe` will include the most recent tag, the number of commits since that tag, and the commit-hash of the current commit:
```
tipdocreg.trondheim.kommune.no/go-template:v1.1.4-43-g5457f57-dirty-local
```
What we can surmise from the above tag:

1. the most recent tag was v1.1.4
1. there are 43 commits since the commit with the most recent tag
1. the commit hash of the current commit is g5457f57
1. the current state of the repo is dirty; it has uncommitted changes
1. the build was made on a local dev-machine

### Pushing to a registry

Pushing your generated image to a docker registry is done by invoking the `push`-goal:
```bash
make push
```

The push-goal uses a normal `docker push`-command. So if your registry requires authentication of some sort (which it definitely should),
you have to run `docker login` manually before you can push for the first time.

FYI: `push` automatically invokes the [`container`](#building-a-container) goal.
