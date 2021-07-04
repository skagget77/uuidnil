# UUID Nil
Allow [Google's UUID](https://github.com/google/uuid) to be unmarshaled from empty and malformed JSON strings. Useful when working with JSON-based APIs not conforming to the UUID standard.

## Example
Wrapping the `user` in a uuidnil wrapper allows it to be unmarshaled without any errors:
```go
func main() {
    var user struct {
        ID   uuid.UUID
        Name string
    }

    data := []byte(`{"id": "", "name": "John"}`)
    if err := json.Unmarshal(data, uuidnil.Wrap(&user, uuidnil.AllowEmpty)); err != nil {
        log.Fatal(err)
    }

    fmt.Println(user)
}
```
