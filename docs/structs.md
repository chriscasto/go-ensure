# Structs

A struct is basically just a container for a group of values, so validating a
struct is just a matter of validating the relevant fields it contains along with
the outputs of any getter methods it may have.  Because of this, the magic of
the struct validator lies in two main methods: `HasFields` for validating field
values and `HasGetters` for validating the output of getter methods. Each of
these takes a single required value, a map of type `with.Validators`, and a second 
optional value, a map of type `with.FriendlyNames`.  The former defines a validator 
for each field or getter method and the latter provides a set of human-readable 
(and user-friendly) names to use when identifying the field to which an error 
message belongs.  It looks something like this:

```
validStruct := ensure.Struct[MyStruct]().HasFields(with.Validators{
    "Field1": ensure.String(),
    "Field2": ensure.Number[int](),
    "Field3": ensure.Array[float64](),
}, with.FriendlyNames{
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

func (p Person) FullName() string {
    return fmt.Sprintf("%s %s", p.FirstName, p.LastName)
}
```

You might then have validation that looks something like this:

```
validator := ensure.Struct[Company].HasFields(with.Validators{
    "Name": ensure.String().IsNotEmpty(),
    "Revenue": ensure.Number[float64]().IsGreaterThan(0.0),
    "Employees": ensure.Array[Person].Each(
        ensure.Struct[Person].HasFields(with.Validators{
            "FirstName": ensure.String().IsNotEmpty().Matches(ensure.Alpha),
            "LastName": ensure.String().IsNotEmpty().Matches(ensure.Alpha),
        },with.FriendlyNames{
            "FirstName": "First Name",
            "LastName": "Last Name",
        }).HasGetters(with.Validators{
            "FullName": ensure.String().Contains(" ")
        },with.FriendlyNames{
            "FullName": "Full Name",
        }),
    ),
})
```

Which can be a lot to take in all at once, but it should still be fairly easy to
understand what it's doing by reading it line by line. We could also decompose
things a bit to make it more readable:

```
validName := ensure.String().IsNotEmpty().Matches(ensure.Alpha)

validPerson := ensure.Struct[Person].HasFields(with.Validators{
    "FirstName": validName,
    "LastName": validName,
},with.FriendlyNames{
    "FirstName": "First Name",
    "LastName": "Last Name",
}).HasGetters(with.Validators{
    "FullName": ensure.String().Contains(" "),
},with.FriendlyNames{
    "FullName": "Full Name",
}})

validCompany := ensure.Struct[Company]().HasFields(with.Validators{
    "Name": ensure.String().IsNotEmpty(),
    "Revenue": ensure.Number[float64]().IsGreaterThan(0.0),
    "Employees": ensure.Array[Person].Each(validPerson),
})
```

In the case of a Person object that has an empty last name, you could expect an
error message like:

```
"Last Name: string must not be empty"
```

Note that the field that produced the error is identified using the user-friendly
name, making it safer to communicate this error directly to the user without the
risk of exposing implementation details about your domain objects or DTOs.

## Methods

| Method                                          | Description                                                               |
|-------------------------------------------------|---------------------------------------------------------------------------|
| HasFields(with.Validators, with.FriendlyNames)  | Passes if each of the name fields passes validation                       |
| HasGetters(with.Validators, with.FriendlyNames) | Passes if the return value of each getter passes validation               |
| Is(func (T) error)                              | Passes if the function passed does not produce an error during validation |

## Field visibility
Due to the way visibility works in Go, only exported struct fields are able to
be validated directly.  That is, you can validate `MyStruct.Foo` but not 
`MyStruct.foo`. If you have unexported fields that need to be validated, use 
getter methods to expose their values so the validator can access them.

## Getters
Getter methods, by convention, accept no args and return only a single value.
Any method that meets these two conditions can be assigned a validator using
the `HasGetters` method.  Methods that either accept one or more args or 
return multiple values cannot be assigned validators this way.  If you have
complex methods like this that also need to be validated, consider using the
`Is` method, which enables arbitrary validations.