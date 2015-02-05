# go-serialization
Low-footprint serialization for golang

Library expects structure with pointer fields (`*field`) instead of value ones (`field`)

Example:

```go
type TargetStruct struct {
  Operation *string
  SomeID    *[]byte
  Admin      *bool
  ExpiresAt  *uint32
  SomeNumber *int16
}
```

Usage example:

```go

  cfg := getStruct()

  binData, err := simpleserialize.MarshalStruct(cfg)
  // ...
  cfgOut := TargetStruct{}
  err = simpleserialize.UnMarshalStruct(&cfgOut, binData)

```
