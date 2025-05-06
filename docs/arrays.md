# Arrays

Array validators can be of any subtype (string, int, float64, etc).  As an
example, to make sure your array of 16-bit unsigned integers has at least 3
values, you could use something like this:

```go
// ensure array of uint16 has more than 3 elements
validator := ensure.Array[uint16]().HasMoreThan(3)
```

While the array validator does have its own methods you can use for validating
the array directly, most of the time you will want to apply validation to each of the
items it contains.  You can do this by passing an appropriate validator to the
`Each()` method.  For example, to make sure each string in an array contains at
least one vowel, you could do something like this:

```go
// ensure in array of strings each string matches "(?i)[aeiuo]+"
validator := ensure.Array[string]().Each(
    ensure.String().Matches("(?i)[aeiuo]+")
 )
```

## Methods

| Method               | Description                                                               |
|----------------------|---------------------------------------------------------------------------|
| IsEmpty()            | Passes if tested array is empty (len(arr) == 0)                           |
| IsNotEmpty()         | Passes if tested array is not empty (len(arr) != 0)                       |
| HasCount(int)        | Passes if the length of the tested array is equal to the passed int       |
| HasFewerThan(int)    | Passes if the length of the tested array is less than the passed int      |
| HasMoreThan(int)     | Passes if the length of the tested array is more than the passed int      |
| HasLengthWhere(v)    | Adds a number validator that evaluates against the length of the array    |
| Each(v)              | Passes if the provided validator passes for each element in the array     |
| Is(func ([]T) error) | Passes if the function passed does not produce an error during validation |

# Comparable Arrays

There is a subtype of array validator that can apply additional checks on 
`comparable` types.  It has all the methods available on the array validator,
with the addition of some that can only be used on types that support equality
operations.

```go
validStr := ensure.ComparableArray[string]().Contains("one")

// no error
validStr.Validate([]string{
    "three", 
    "two",
    "one",
    "zero",
})

// this returns an error
validStr.Validate([]string{
    "foo",
    "bar",
    "baz",
})
```

| Method                    | Description                                                            |
|---------------------------|------------------------------------------------------------------------|
| Contains(T)               | Passes if tested array contains the passed value                       |
| DoesNotContain(T)         | Passes if tested array does not contain the passed value               |
| ContainsAnyOf(...T)       | Passes if tested array contains any of the passed values               |
| DoesNotContainAnyOf(...T) | Passes if tested array does not contain any of the passed values       |
| ContainsNoDuplicates()    | Passes if tested array does not contain any duplicate values           |
| ContainsOnly(...T)        | Passes if the values in the array are exclusively from the passed list |
