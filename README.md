# Kredis

A port of [Rails' Kredis](https://github.com/rails/kredis) for Go. Very
much a WIP...

## Example

```go
k, e := NewInteger("foo", Options{})
k.SetValue(1234)

// prints: 1234
fmt.Println(k.Value())
```

## TODO

Implement additional Kredis data structures

- refactor how defaults work with NewXWithDefault
  - must be set initially, with `watch`, `exists`, `unwatch`
  - deprecate `Options.DefaultValue`
- slots
- counters
- flags
- sets
- unique lists
- hashs
- map scalar type? (serialized as json)
- add logging (zerolog)
