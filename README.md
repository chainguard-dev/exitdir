# ExitDir

ExitDir is a library to signal a process should exit via the filesystem.

## Demo

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
[Busy] Tick at 2022-04-13 16:02:06.342656 -0700 PDT m=+0.501881542
[Busy] Tick at 2022-04-13 16:02:06.842661 -0700 PDT m=+1.001884459
[Busy] Tick at 2022-04-13 16:02:07.342658 -0700 PDT m=+1.501879501
[Busy] Tick at 2022-04-13 16:02:07.842672 -0700 PDT m=+2.001891959
[Busy] Tick at 2022-04-13 16:02:08.342673 -0700 PDT m=+2.501891251
[Busy] Tick at 2022-04-13 16:02:08.843159 -0700 PDT m=+3.002375251
[Busy] Tick at 2022-04-13 16:02:09.342669 -0700 PDT m=+3.501884251
[Busy] Tick at 2022-04-13 16:02:09.842485 -0700 PDT m=+4.001698417
[Busy] Tick at 2022-04-13 16:02:10.342675 -0700 PDT m=+4.501886251
[Work] Exiting...
[Busy] Exiting...
[1]+  Done                    go run ./cmd/busy
$
```

## In Kubernetes

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
[Busy] Tick at 2022-04-13 22:59:22.788790941 +0000 UTC m=+0.500829225
[Busy] Tick at 2022-04-13 22:59:23.288379133 +0000 UTC m=+1.000417416
[Busy] Tick at 2022-04-13 22:59:23.788889936 +0000 UTC m=+1.500928232
[Busy] Tick at 2022-04-13 22:59:24.288511353 +0000 UTC m=+2.000549637
[Busy] Tick at 2022-04-13 22:59:24.788883831 +0000 UTC m=+2.500922184
[Busy] Tick at 2022-04-13 22:59:25.288507366 +0000 UTC m=+3.000545649
[Busy] Tick at 2022-04-13 22:59:25.789241661 +0000 UTC m=+3.501279944
[Busy] Exiting...
```
