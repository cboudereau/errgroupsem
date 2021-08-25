# errgroupsem

This package is an extended errgroup package to fix the lack of CPU bound support when errgroup is used in conjuction with semaphore.

The semaphore package could not be used with errgroup since the error in semaphore is internally managed by errgroup.

See: https://github.com/golang/go/issues/27837 (create another package https://github.com/golang/go/issues/27837#issuecomment-633659652)

## How to install
```bash
go get github.com/cboudereau/errgroupsem
```

## Example
See the unit tests for little and channel based demo