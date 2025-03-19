# Maps

Validating maps is very similar to validating arrays.  The key difference is
that with maps you can apply validations to both keys and values.

```
validator = ensure.Map[string, int]().EachKey(
    ensure.String()
).EachValue(
    ensure.Number[int]()
)
```

## Methods

| Method                    | Description                                                               |
|---------------------------|---------------------------------------------------------------------------|
| IsEmpty()                 | Passes if tested map is empty (len(arr) == 0)                             |
| IsNotEmpty()              | Passes if tested map is not empty (len(map) != 0)                         |
| HasCount(int)             | Passes if the length of the tested map is equal to the passed int         |
| HasFewerThan(int)         | Passes if the length of the tested map is less than the passed int        |
| HasMoreThan(int)          | Passes if the length of the tested map is more than the passed int        |
| HasLengthWhere(v)         | Adds a number validator that evaluates against the length of the map      |
| EachKey(v)                | Passes if the provided validator passes for each key in the map           |
| EachValue(v)              | Passes if the provided validator passes for each value in the map         |
| Is(func (map[K,V]) error) | Passes if the function passed does not produce an error during validation |
