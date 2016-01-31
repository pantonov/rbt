[![Build Status](https://travis-ci.org/pantonov/rbt.svg)](https://travis-ci.org/pantonov/rbt) [![GoDoc](https://godoc.org/github.com/pantonov/rbt?status.svg)](https://godoc.org/github.com/pantonov/rbt)

# RbTree - A Red-Black tree implementation

RbTree is a sorted associative container that contains key-value pairs with unique keys. Keys are sorted by using the comparison function. Search, removal, and insertion operations have logarithmic complexity.

Documentation: [Godoc](http://godoc.org/github.com/pantonov/rbt)

# Installation
* Stable:  
    **go get gopkg.in/pantonov/rbt.v1**

* Latest:  
    **go get github.com/pantonov/rbt**

# Example

```go

    r := rbt.NewRbTree(func(k1, k2 interface{}) bool {
        return k1.(string) < k2.(string)
    })
    r.Insert("c", 1)
    r.Insert("b", 2)
    r.Insert("a", 3)
    r.Insert("d", 4)

    // iterate over map in ascending key order. Note that Key is a method,
    // to prevent in-place modification, while Value accessed directly
    for n := r.First(); n != nil; n = n.Next() {
        fmt.Printf("%s -> %d\n", n.Key(), n.Value)
    }
    // Find value by key
    fmt.Printf("Value for %s is: %d\n", "b", r.Find("b"))

    // More node-level (iperator) operations. You'll want to check return values
    // for nil's, which are omitted here for brevity
    n := r.FindNode("a")
    fmt.Printf("Key next to 'a': %s\n", n.Next().Key())
    r.DeleteNode(n.Next())
    fmt.Printf("Now, key for node next to 'a': %s\n", n.Next().Key())
```
