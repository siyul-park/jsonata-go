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
  jsonata "github.com/siyul-park/jsonata-go"
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
    defer expression.Close()
    result, err := expression.Evaluate(data);  // returns 24
}

## More information
- JSONata [documentation](http://docs.jsonata.org/)
- [JavaScript API](http://docs.jsonata.org/embedding-extending)
- [Intro talk](https://www.youtube.com/watch?v=TDWf6R8aqDo) at London Node User Group
