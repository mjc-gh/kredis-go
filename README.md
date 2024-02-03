# Kredis

[![kredis-go](https://github.com/mjc-gh/kredis-go/actions/workflows/tests.yaml/badge.svg)](https://github.com/mjc-gh/kredis-go/actions/workflows/tests.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/mjc-gh/kredis-go.svg)](https://pkg.go.dev/github.com/mjc-gh/kredis-go)

A port of [Rails' Kredis](https://github.com/rails/kredis) for Go.

> Kredis (Keyed Redis) encapsulates higher-level types and data
structures around a single key, so you can interact with them as
coherent objects rather than isolated procedural commands.

### Motivation

I wrote a Go version of Kredis to help me learn Go and the Redis
client package. This exercise also helped me explore generics in detail.

In theory, there should be full interoperability between this Go package
and the Ruby version, thus enabling services in either language to work
together over Redis.

## Usage

To use Kredis, you must first configure a connection. The default
connection configuration is named `"shared"`. `SetConfiguration` expects
a config name string, optional namespace (pass an empty string for no
namespace), and a Redis URL connection string.

```go
kredis.SetConfiguration("shared", "ns", "redis://localhost:6379/2")
```

Kredis comes with a simple logger, which is useful for testing and
development. The logger will output the Redis commands that are executed
and how long they took. It's nearly identical to the output the we see
in the Ruby Kredis version.

To enable debug logging, simply call:

```go
kredis.EnableDebugLogging()
```

### Failsafe

Kredis wraps calls in Redis in a
[failsafe](https://github.com/rails/kredis/blob/4fbb2f5613ed049f72cba337317c5eae2a6bba28/lib/kredis/types/proxy/failsafe.rb#L9-L13).
To match this design, certain functions around reading values
intentionally omit returning an error and will return zero values when a
read fails.

If you need to handle errors and do not want the "failsafe" behavior,
most types also offer functions that end in `Result` that return the
value and an error using standard Go idioms. For example:

```go
slot, _ := kredis.NewSlot("slot", 3)
slot.Reserve()

// some time later when Redis is down...

t := slot.Taken()
// 0
t, err := slot.TakenResult() // func TakenResult() (int64, error)
// 0, dial tcp [::1]:6379: connect: connection refused
```

## Types

- [Counter](#counter)
- [Cycle](#cycle)
- [Enum](#enum)
- [Flag](#flag)
- [Limiter](#limiter)
- [Slot](#slot)
- [Scalar Types](#scalar-types)

Collection types:

- [List](#list)
- [Set](#set)
- [Hash](#hash)
- [Ordered Set](#ordered-set)
- [Unique List](#unique-list)

All factory functions for the various types accept the following option
functions:

- `WithConfigName` sets the Redis config to use. The function accepts a
  `string` name and should match the value passed to `SetConfiguration`.
- `WithContext` allows the user to specify the `context.Context` value
   passed to the underlying Redis client commands.
- `WithExpiry` sets the expiration for type's key. The function accepts
   a string value that is parsed by `time.ParseDuration` to get a
   `time.Duration` value.
   When this option is used for scalar and other basic types, it usually
   means the Redis `SET` command is called with an expiration value. For
   collection types, whenever the collection is mutated (`Append`,
   `Prepend`, `Add`, `Remove`, etc.), the key is refreshed by calling
   the `EXPIRE` command after the mutation commands.
   Additionally, when this option is used, you can call `RefreshTTL()`
   at any point to refresh the key's expiration. The function returns a
   `(bool, error)`, where the boolean value indicates whether the key
   was refreshed or not.

### Counter

```go
cntr, err := kredis.NewCounter("counter")
n, err := cntr.Increment(1) // MULTI, INCRBY counter 1, EXEC
// 1, nil
n, err := cntr.Increment(2) // MULTI, INCRBY counter 2, EXEC
// 3, nil
n, err := cntr.Decrement(3) // MULTI, DECRBY counter 3, EXEC
// 0, nil
cntr.Value()                // GET counter
// 0
```

### Cycle

```go
cycle, err := kredis.NewCycle("cycle", []string{"ready", "set", "go"})
cycle.Index()        // GET counter
// 0
err := cycle.Next()  // GET counter, SET counter 1
// nil
err = cycle.Next()   // GET counter, SET counter 2
// nil
cycle.Index()        // GET counter
// 2
val := cycle.Value() // GET counter
// "go"
```

### Enum

```go
vals := []string{"ready", "set", "go"}

enum, _ := kredis.NewEnum("enum", "go", vals) // SET enum go
enum.Is("go")                                 // GET enum
// true
val := enum.Value()                           // GET enum
// "go"
err := enum.SetValue("set")                   // SET enum set
// nil
err = enum.SetValue("error")
// invalid enum value (Kredis.EmptyValues error)
```

### Flag

By default the `Mark()` function does not call `SET` with `nx`

```go
flag, err := kredis.NewFlag("flag")
flag.IsMarked()                            // EXISTS flag
// false
err = flag.Mark()                          // SETNX flag 1
// nil
flag.IsMarked()                            // EXISTS flag
// true
err = flag.Remove()                        // DEL flag

flag.Mark(kredis.WithFlagMarkExpiry("1s")) // SET flag 1 ex 1
flag.IsMarked()                            // EXISTS flag
// true

time.Sleep(2 * time.Second)

flag.IsMarked()                            // EXISTS flag
// false
```

The `SoftMark()` function will call set with `NX`

```go
flag.SoftMark(kredis.WithFlagMarkExpiry("1s"))  // SET flag 1 ex 1 nx
flag.SoftMark(kredis.WithFlagMarkExpiry("10s")) // SET flag 1 ex 10 nx
flag.IsMarked()                                 // EXISTS flag
// true

time.Sleep(2 * time.Second)

flag.IsMarked()                                 // EXISTS flag
// true
```

### Limiter

The `Limiter` type is based off the `Counter` type and provides a
simple rate limiting type with a failsafe on Reids errors. See the
original [Rails PR for more
details](https://github.com/rails/kredis/pull/136).

`IsExceeded()` will return `false` in the event of a Redis error.
`Poke()` does return an error, but it can easily be ignored in Go.

```go
limiter, _ := kredis.NewLimiter("limiter", 5)
limiter.Poke()         // MULTI, INCRBY limiter 1 EXEC
limiter.Poke()         // MULTI, INCRBY limiter 1 EXEC
limiter.Poke()         // MULTI, INCRBY limiter 1 EXEC
limiter.Poke()         // MULTI, INCRBY limiter 1 EXEC

limiter.IsExceeded()   // GET limiter
// true
err := limiter.Reset() // DEL limiter
// nil
limiter.IsExceeded()   // GET limiter
// false
```

### Slots

```go
slot, err := kredis.NewSlot("slot", 3)
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

### Scalar types

- `NewBool` and `NewBoolWithDefault`
- `NewInteger` and `NewIntegerWithDefault`
- `NewString` and `NewStringWithDefault`
- `NewTime` and `NewTimeWithDefault`
- `NewJSON` and `NewJSONWithDefault`

```go
k, err := kredis.NewInteger("myint", Options{})
err = k.SetValue(1024)  // SET myint 1024
// nil
k.Value()               // GET myint
// 1024
```

With expiration through the `WithExpiry` option function:

```go
k, err := kredis.NewTime("sessionStart", kredis.WithExpiry("30m"))
err = k.SetValue(time.Now()) // SET sessionStart 2024-01-06T13:30:35.613332-05:00 ex 1800
// nil
val := k.Value()             // GET sessionStart
// 2024-01-01 12:00:00 -0500 EST
dur := k.TTL()               // TTL sessionStart
// 30m0s

// over 30 minutes later
k.Value()                    // GET sessionStart
// nil
k.TTL()                      // TTL sessionStart
// -2ns (key does not exit now)
```

### List

```go
l := kredis.NewIntegerList("users")
n, err := l.Append(1, 2, 3)  // RPUSH users 1 2 3
// 3, nil
n, err = l.Prepend(9, 8)     // LPUSH users 9, 8
// 2, nil

ids = make([]int, 5)
n, err := l.Elements(ids)    // LRANGE users 0, 5
// 5, nil

// read some elements with an offset
lastTwo = make([]int, 2)
n, err = l.Elements(lastTwo, WithRangeStart(3)) // LRANGE users 3 5
// 2, nil
```

Different typed factories exist for the `List` struct:

- `NewBoolList`
- `NewFloatList`
- `NewStringList`
- `NewTimeList` over `time.Time`
- `NewJSONList` over the `KredisJSON` alias type

It's possible to provide a default value as well, which will use `WATCH`
to transactionally set the value if the key does not already exist. For
lists, this entails calling `Append` and using `RPUSH` to add the
default elements.

```go
strs, err := kredis.NewStringListWithDefault("lines", []string{"hello", "redis"},)
// WATCH lines
// EXISTS lines
// RPUSH lines hello redis
// UNWATCH

err = strs.Remove("hello") // LREM lines 0 "hello"
// nil
n := strs.Length()          // LLEN lines
// 2
```

### Set

```go
t := time.Date(2021, 8, 28, 23, 0, 0, 0, time.UTC)
times := []time.Time{t, t.Add(1 * time.Hour), t.Add(2 * time.Hour), t.Add(3 * time.Hour)}
set, err := kredis.NewTimeSet("times")
set.Add(times...)                  // SADD times 2021-08-28T23:00:00Z 2021-08-29T00:00:00Z 2021-08-29T01:00:00Z 2021-08-29T02:00:00Z
// 4, nil
members, err := set.Members()      // SMEMBERS times
// []time.Time{...}, nil
set.Size()                         // SCARD times
// 4
set.Includes(t)                    // SISMEMBER times 2021-08-28T23:00:00Z
// true
set.Includes(t.Add(4 * time.Hour)) // SISMEMBER times 2021-08-29T03:00:00Z
// false
sample := make([]time.Time{}, 2)
n, err := set.Sample(sample)       // SRANDMEMBER times 2
// 2, nil

fmt.Println(sample)
```

Different factory functions exist for various `Set` types (this is the
case for all collection types). It's possible to provide a default using
the `WithDefault` factories.

```go
kredis.NewStringSetWithDefault("set", []string{"a", "b", "c"})
// WATCH strings
// EXISTS strings
// SADD strings a b c
// UNWATCH
```

The `Add` function is used to set the default. Thus, `SADD` is used when
the key does not already exist.

### Hash

```go
dogs := map[string]kredis.KredisJSON{"ollie": *kredis.NewKredisJSON(`{"weight":9.72}`), "leo": *kredis.NewKredisJSON(`{"weight":23.33}`)}
hash, _ := kredis.NewJSONHash("pets")

n, err := hash.Update(dogs)     // HSET pets ollie {"weight":9.72} leo {"weight":23.33}
// 2, nil
val, ok := hash.Get("ollie")    // HGET pets ollie
// {"weight":9.72}, true
keys, err := hash.Keys()        // HKEYS pets
// []string{"ollie", "leo"}, nil
vals, err := hash.Values()      // HVALS pets
// [{"weight":9.72} {"weight":23.33}], nil
entries, err := hash.Entries()  // HGETALL pets
// map[leo:{"weight":23.33}
//     ollie:{"weight":9.72}], nil
hash.Clear()                    // DEL pets
```

### Ordered Set

```go
oset, err := kredis.NewStringOrderedSet("ranks", 4)
add, rm, err := oset.Append("a", "b", "c") // MULTI
// 3, 0, nil                                  ZADD ranks 1.704562075576027e+09 a 1.7045620755760288e+09 b 1.7045620755760298e+09 c
                                           // ZREMRANGEBYRANK ranks 0 -5
                                           // EXEC
add, rm, err := oset.Append("a", "b", "c") // MULTI
// 2, 1, nil                                  ZADD ranks -1.704562075576382e+09 d -1.7045620755763829e+09 e
                                           // ZREMRANGEBYRANK ranks 4 -1
                                           // EXEC
members, _ := oset.Members()               // ZRANGE ranks 0 -1
// [e d a b]
oset.Size()                                // ZCARD ranks
// 4
oset.Includes("d")                         // ZSCORE ranks d
// true
oset.Includes("c")                         // ZSCORE ranks c
// false
```

For more details on the underlying Redis implementation for the
`OrderSet` type, refer to the [original Ruby
PR](https://github.com/rails/kredis/pull/76) for this feature.

### Unique List

Similar to `OrderedSet`, this type exepcts a `limit` as well

```go
uniq, err := kredis.NewFloatUniqueList("uniq", 5)
n, err := uniq.Append(3.14, 2.718) // MULTI, LREM uniq 0 3.14, LREM uniq 0 2.718, RPUSH uniq 3.14 2.718, LTRIM uniq -5 -1, EXEC
// 2, nil
n, err = uniq.Prepend(1.1)         // MULTI, LLREM uniq 0 1.1, LPUSH uniq 1.1, LTRIM uniq -5 -1, EXEC
// 1, nil
llen, _ := uniq.Length()           // LLEN uniq
// 3, nil

elements := make([]float64, 3)
n, err := uniq.Elements(elements) // LRANGE uniq 0 3
// 3, nil

uniq.Clear()                      // DEL uniq
```

## TODO

- More test coverage:
    - Better coverage all possible generic types for collections
    - Test Redis commands with some sort of test env `ProcessHook`. This
      is useful when checking expiration is set correctly without using
      `time.Sleep` in tests.

### Future Features

- Other scalar types
    - Some sort of map type (serialized as json)
- `Hash` type
    - Add a way to call `HINCRBY` (limited to `Hash[int]`)
- `OrderedSet` type
    - Let the user control the ranking
    - Add a `Rank` function to call `ZRANK`
- Explore support for [pipelining](https://redis.uptrace.dev/guide/go-redis-pipelines.html)
    - With only kredis commands?
    - With a shared redis client?
