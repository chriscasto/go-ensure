# Structs

A struct is basically just a container for a group of values, so validating a
struct is just a matter of validating each of the fields it contains.  Because
of this, the magic of the struct validator is almost entirely in the constructor.
The constructor takes a single required value, a map of type `Fields`, and a
second optional value, a map of type `FriendlyNames`.  The former creates a
validator for each field and the latter provides a set of human-readable (and
user-friendly) names to use when identifying the field to which an error message
belongs.  It looks something like this:

```
validStruct := ensure.Struct[MyStruct](ensure.Fields{
    "Field1": ensure.String(),
    "Field2": ensure.Number[int](),
    "Field3": ensure.Array[float64](),
}, ensure.FriendlyNames{
    "Field1": "First Field"
    "Field2": "Second Field"
    "Field3": "Third Field"
})
```

For something a little more concrete, imagine you have a set of structs like this:

```
type Company struct {
    Name string
    Revenue float64
    Employees []Person
}

type Person struct {
    FirstName string
    LastName string
}
```

You might then have validation that looks something like this:

```
validator := ensure.Struct[Company](ensure.Fields{
    "Name": ensure.String().IsNotEmpty(),
    "Revenue": ensure.Number[float64]().IsGreaterThan(0.0),
    "Employees": ensure.Array[Person].Each(
        ensure.Struct[Person](ensure.Fields{
            "FirstName": ensure.String().IsNotEmpty().Matches(ensure.Alpha),
            "LastName": ensure.String().IsNotEmpty().Matches(ensure.Alpha),
        },
    ),
})
```

Which can be a lot to take in all at once, but it should still be fairly easy to
understand what it's doing by reading it line by line. We could also decompose
things a bit to make it more readable:

```
validName := ensure.String().IsNotEmpty().Matches(ensure.Alpha)

validPerson := ensure.Struct[Person](ensure.Fields{
    "FirstName": validName,
    "LastName": validName,
})

validCompany := ensure.Struct[Company](ensure.Fields{
    "Name": ensure.String().IsNotEmpty(),
    "Revenue": ensure.Number[float64]().IsGreaterThan(0.0),
    "Employees": ensure.Array[Person].Each(validPerson),
})
```

## Methods

| Method             | Description                                                               |
|--------------------|---------------------------------------------------------------------------|
| Is(func (T) error) | Passes if the function passed does not produce an error during validation |


## A note on field visibility
Due to the way visibility works in Go, only exported struct fields are
able to be validated.  That is, you can validate `MyStruct.Foo` but not `MyStruct.foo`.