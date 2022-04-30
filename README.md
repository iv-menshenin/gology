# gology (Really Zero Allocation JSON Logger)

## Really fast

Take a look at proofs:

	cpu: Intel(R) Core(TM) i7-9700F CPU @ 3.00GHz
	Benchmark_GologyLogger-8         9410730               125.6 ns/op             0 B/op          0 allocs/op
	Benchmark_ZeroLogger-8           1561195               876.7 ns/op           512 B/op          1 allocs/op


## Features

`Gology' is a simple package without any complicated features. No hooks or embeddings into other commonly used systems. We only support JSON output and a fairly narrow range of data types here. No Any or empty interfaces.

```go
    log := New(os.StdOut, LevelAll)
    log.Write(
        LevelDebug,
        "some message to write into the output",
        String("attr_str", "some string attribute"),
        Int("count", 1024),
    )
```

output:

    {"message":"some message to write into the output","level":"DEBUG","attr_str":"some string attribute","count":1024}

## Go-Routines

In real life, a program includes several threads in which functions are executed sequentially or call each other, or all in a jumble. We do not support safe calling of the same object from different go-routines, so you must create your own Logger for each routine.

But what we can really offer you is that it's easy to wrap your Logger with preset values, which won't be written immediately. The values will be implied inside the function call branch until you exit the function to which such values have been passed.

```go

func somethingToDo(log Logger) {
    // something to do
    log.Warning("some message")
}

func main(){
    // something to do
    somethingToDo(log.WithAttrs(
        String("test1", "value 1"),
        String("test2", "value 2"),
    ))

    log.Warning("done")
}
```

output:

    {"message":"some message","level":"WARNING","test1":"value 1","test2":"value 2"}
    {"message":"done"}