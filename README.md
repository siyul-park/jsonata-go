# JSONata

JSON query and transformation language

Reference implementation of the [JSONata query and transformation language](http://jsonata.org/).

* [JSONata in 5 minutes](https://www.youtube.com/embed/ZBaK40rtIBM)
* [JSONata language documentation](http://docs.jsonata.org/)
* [Try it out!](http://try.jsonata.org/)

## Installation

- `go get github.com/siyul-park/jsonata-go`

## Quick start

```golang
package main

import (
  "github.com/siyul-park/jsonata-go"
)

var data = map[string]any{
  "example": []map[string]any{
    {"value": 4},
    {"value": 7},
    {"value": 13},
  },
}

func main() {
    expression := jsonata.MustCompile("$sum(example.value)")
    result, err := expression.Evaluate(data);  // returns 24
}
```

## Benchmark
```
goos: linux
goarch: amd64
pkg: github.com/siyul-park/jsonata-go
cpu: AMD EPYC 7282 16-Core Processor                
BenchmarkExpression_Compile-4            	     207	   5672824 ns/op	 1533080 B/op	   19941 allocs/op
BenchmarkExpression_Evaluate-4           	    3902	    260193 ns/op	   62461 B/op	     873 allocs/op
BenchmarkExpression_Assign-4             	  822211	      1542 ns/op	     521 B/op	       9 allocs/op
BenchmarkExpression_RegisterFunction-4   	  243468	      4395 ns/op	    1846 B/op	      23 allocs/op
BenchmarkExpression_Ast-4                	   14838	     80954 ns/op	   29617 B/op	     413 allocs/op
PASS
ok  	github.com/siyul-park/jsonata-go	8.251s
```

## More information
- JSONata [documentation](http://docs.jsonata.org/)
- [JavaScript API](http://docs.jsonata.org/embedding-extending)
- [Intro talk](https://www.youtube.com/watch?v=TDWf6R8aqDo) at London Node User Group
