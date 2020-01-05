# file_comparator_go

`$ go build -o file_comparator_go`


```shell
$ ./file_comparator_go -f1=t/data/'Crime&Punishment.txt' -f2=t/data/'War&Peace.txt' -n=5
the	44739
and	27451
to	23362
of	19879
a	14382
```

```shell
$ ./file_comparator_go -h
Usage of ./file_comparator_go:
  -f1 string
    	path to the 1st file
  -f2 string
    	path to the 2nd file
  -n int
    	limit for outputted lines (default 10)
```
