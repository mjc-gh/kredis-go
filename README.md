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
