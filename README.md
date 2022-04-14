# ExitDir

ExitDir is a library to signal a process should exit via the filesystem.

## Usage

To use this library, on the leader process:

```go
func main() {
	// Signal the other processes should exit via a new file in `EXIT_DIR`.
    defer exitdir.Exit();
	// ...rest of implementation. 
}
```

And then one or more follower processes:

```go
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
go run ./cmd/busy &
go run ./cmd/work
```

Results in:

```shell
$ export EXIT_DIR=`mktemp -d`
$ go run ./cmd/busy &
[1] 83528
$ go run ./cmd/work
[Work] Doing work...
[Busy] Tick 0
[Busy] Tick 1
[Busy] Tick 2
[Busy] Tick 3
[Work] Exiting...
[Busy] Exiting...
[1]+  Done                    go run ./cmd/busy
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
        - name: work
          image: ko://chainguard.dev/exitdir/cmd/work
          env:
            - name: EXIT_DIR
              value: "/var/exitdir"
          volumeMounts:
            - name: exit-dir
              mountPath: "/var/exitdir"
        - name: busy
          image: ko://chainguard.dev/exitdir/cmd/busy
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
$ kubectl logs example-78fm7 work
[Work] Doing work...
[Work] Exiting...
```

```shell
$  kubectl logs example-78fm7 busy
[Busy] Tick 0
[Busy] Tick 1
[Busy] Tick 2
[Busy] Tick 3
[Busy] Exiting...
```
