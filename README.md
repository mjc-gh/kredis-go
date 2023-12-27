# Kredis

A port of [Rails' Kredis](https://github.com/rails/kredis) for Go.

> Kredis (Keyed Redis) encapsulates higher-level types and data
structures around a single key, so you can interact with them as
coherent objects rather than isolated procedural commands.

## Lists

```go
l := NewIntegerList("users")

n, err := l.Append(1, 2, 3)  // RPUSH users 1 2 3
// 3, nil
n, err = l.Prepend(9, 8)     // LPUSH users 9, 8
// 2, nil

ids = make([]int, 5)
n, err := l.Elements(ids)   // LRANGE users 0, 5
// 5, nil

// read some elements with an offset
last_2 = make([]int, 2)
n, err = l.Elements(last_2, WithRangeStart(3)) // LRANGE users 3 5
// 2, nil
```

Different typed factories exist for these types:

- `NewBoolList`
- `NewStringList`
- `NewTimeList` over `time.Time`
- `NewJSONList` over the `kredisJSON` alias type

It's possible to provide a default value as well:

```go
strs := NewStringWithDefault("lines", []string{"hello", "redis"})
err := strs.Remove("hello") // LREM lines 0 "hello"
// nil
n := strs.Length()          // LLEN lines
// 2
```

## Slots

```go
slot := NewSlot("slot", 3)
slot.Reserve()     // GET slot + INCR slot
// true
slot.IsAvailable() // GET slot
// true
slot.Taken()       // GET slot
// 1
slot.Reserve()     // GET slot + INCR slot
// true
slot.Reserve()     // GET slot + INCR slot
// true

slot.Reserve()     // GET slot
// false
slot.IsAvailable() // GET slot
// false
slot.Taken()       // GET slot
// 3

// Reserve() with one or more callbacks will always call Release() even
// if the callbacks are not invoked and there are no slots available
slot.Reserve(func () { fmt.Println("not called") })
// GET slot + DECR slot
// false
slot.Reserve(func () { fmt.Println("called") })
// GET slot + INCR slot + DECR slot
// true
```

## Scalar types

- `NewBool` and `NewBoolWithDefault`
- `NewInteger` and `NewIntegerWithDefault`
- `NewString` and `NewStringWithDefault`
- `NewTime` and `NewTimeWithDefault`
- `NewJSON` and `NewJSONWithDefault`

```go
k, err := NewInteger("myint", Options{})
err = k.SetValue(1024)  // SET myint 1024
// nil
k.Value()               // GET myint
// 1024
```

## TODO

Implement additional Kredis data structures

- Collections
    - sets
    - unique lists
    - hashs
- other scalar types
    - float type
        - on lists and other collections
    - some sort of map type (serialized as json) ??
- document all types in README
- test commands with some sort of test env `ProcessHook` for redis
    clients
- [pipelining](https://redis.uptrace.dev/guide/go-redis-pipelines.html) ??
    - with only kredis commands?
    - with a shared redis client?
