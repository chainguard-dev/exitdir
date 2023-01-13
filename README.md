# ExitDir

ExitDir is a library to signal a process should exit via the filesystem.

## Usage

To use this library, on the leader process:

```go
package main

import "chainguard.dev/exitdir"

func main() {
	// Signal the other processes should exit via a new file in `EXIT_DIR`.
	defer exitdir.Exit();
	// ...rest of implementation.
}
```

And then one or more follower processes:

```go
package main

import (
	"context"
	
	"chainguard.dev/exitdir"
)

func main() {
	// Decorate a context with ExitDir awareness. ExitDir will cancel the
	// returned context when a new file in `EXIT_DIR` is created.
	ctx := exitdir.Aware(context.Background())
	// ...rest of implementation using `ctx` for lifecycle control.
}
```

## Process Demo

This functionally is shown locally by using a temp directory and two processes:

```shell
export EXIT_DIR=`mktemp -d`
go run ./cmd/follower &
go run ./cmd/leader
```

Results in:

```shell
$ export EXIT_DIR=`mktemp -d`
$ go run ./cmd/follower &
[1] 83528
$ go run ./cmd/leader
[Leader] Doing work...
[Follower] Tick 0
[Follower] Tick 1
[Follower] Tick 2
[Follower] Tick 3
[Leader] Exiting...
[Follower] Exiting...
[1]+  Done                    go run ./cmd/follower
$
```

## Kubernetes Demo

Often in Kubernetes a job with multiple containers need to solve a problem of
signaling when to exit. If not coordinated, a resulting hung job looks like a
failure, but in-fact was successful except one container never exited.

```shell
ko apply -f - <<EOF
apiVersion: batch/v1
kind: Job
metadata:
  name: example
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: leader
          image: ko://chainguard.dev/exitdir/cmd/leader
          env:
            - name: EXIT_DIR
              value: "/var/exitdir"
          volumeMounts:
            - name: exit-dir
              mountPath: "/var/exitdir"
        - name: follower
          image: ko://chainguard.dev/exitdir/cmd/follower
          env:
            - name: EXIT_DIR
              value: "/var/exitdir"
          volumeMounts:
            - name: exit-dir
              mountPath: "/var/exitdir"
      volumes:
        - name: exit-dir
          emptyDir: {}
EOF
```

We can see the job finished:

```shell
% kubectl get job
NAME      COMPLETIONS   DURATION   AGE
example   1/1           8s         17s
```

And finding the pod, we can view the logs of each container:

```shell
$ kubectl logs example-78fm7 leader
[Leader] Doing work...
[Leader] Exiting...
```

```shell
$  kubectl logs example-78fm7 follower
[Follower] Tick 0
[Follower] Tick 1
[Follower] Tick 2
[Follower] Tick 3
[Follower] Exiting...
```
