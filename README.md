# Ensure - validation for Go

The Ensure validation library is intended to combine human-readable validation 
syntax while benefiting from the autocomplete, hinting, and syntax validation
capabilities of modern IDEs and the type safety inherent in the Go language.

## Installation

Installing and importing is the same as with most other Go packages.

```
go get github.com/chriscasto/go-ensure
```

```
import (
    "github.com/chriscasto/go-ensure"
)
```

## Usage and Examples

Ensure uses method chaining to create simple and easily understandable validation
rules.  For example, to make sure that a string has at least 8 characters and 
doesn't contain an "@" sign, you would define your validation like this:

```
// ensure string has length 8 and does not contain "@"
validator := ensure.String().HasLength(8).DoesNotContain("@")
```

You could then test a value with this validator using something like this:

```
if err := validator.Validate("baseball"); err != nil {
    fmt.Print(err)
} else {
    fmt.Print(`"baseball" is a valid string`)
}
```

Validation objects are intended to be constructed infrequently, typically at startup,
then used many times, basically for the life of the program.  This is not a hard
rule, but be aware that validation construction methods will panic on error, so be
prepared to `recover()` if you decide to instantiate validators dynamically during
program execution.  It's generally helpful to consider validators to be a statically
compiled part of your code rather than a runtime configuration.

You can find a set of runnable examples in the 
[_examples](https://github.com/chriscasto/go-ensure/tree/main/_examples) directory.

### Strings

We've already seen a partial example of a string validator, but here's something
a bit more complete (if no less contrived).
```
package main

import (
    ensure "github.com/chriscasto/go-ensure"
    "fmt"
)

func main() {
	// ensure string starts with "foo" and doesn't end with "ish"
	validStr := ensure.String().StartsWith("foo").DoesNotEndWith("ish")
	
	// this will succeed
	if err := validStr.Validate("foosball"); err != nil {
	    fmt.Print(err)
	} else {
	    fmt.Print(`"foosball" is a valid string`)
	}
	
	// but this will print an error
	if err := validStr.Validate("foolish"); err != nil {
	    fmt.Print(err)
	} else {
	    fmt.Print(`you shouldn't get here`)
	}
}

```

#### Methods

| Method                | Description                                                                             |
|-----------------------|-----------------------------------------------------------------------------------------|
| IsEmpty()             | Passes if the tested string is empty (len() == 0)                                       |
| IsNotEmpty()          | Passes if the tested string is not empty (len() != 0)                                   |
| StartsWith(str)       | Passes if the tested string begins with provided string value                           |
| DoesNotStartWith(str) | Passes if the tested string does not begin with provided string value                   |
| EndsWith(str)         | Passes if the tested string ends with provided string value                             |
| DoesNotEndWith(str)   | Passes if the tested string does not end with provided string value                     |
| Contains(str)         | Passes if provided string value occurs anywhere in the tested string                    |
| DoesNotContain(str)   | Passes if provided string value does not occur anywhere in the tested string            |
| HasLength(int)        | Passes if the tested string's length is exactly the same as the provided int            |
| IsShorterThan(int)    | Passes if the tested string's length is less than the provided int                      |
| IsLongerThan(int)     | Passes if the tested string's length is greater than the provided int                   |
| HasLengthWhere(v)     | Adds a number validator that evaluates against the length of the string                 |
| IsOneOf([]string)     | Passes if the tested string is identical to one of the values in the provided array     |
| IsNotOneOf([]string)  | Passes if the tested string is not identical to any of the values in the provided array |
| Matches(str)          | Passes if the tested string matches the provided regular expression                     |
| Is(func (str) error)  | Passes if the function passed does not produce an error during validation               |

#### Predefined Regex Patterns

To simplify common use cases, there are predefined regular expressions available
to use.  For example, if you wanted to ensure that the string only contains
alphanumeric values, you could use something like this:

```
validator := ensure.String().Matches(ensure.AlphaNum)
```

| Constant | Description                              | Example                                                            |
|----------|------------------------------------------|--------------------------------------------------------------------|
| Alpha    | Characters in the English alphabet (a-z) | "abc"                                                              |
| Numbers  | Only numbers 0-9                         | "123                                                               |
| AlphaNum | Characters from Alpha plus Numbers       | "abc123"                                                           | 
| Decimal  | A number with a decimal (.)              | "1.23"                                                             |
| Uuid4    | A v4 UUID                                | "d94cd8e1-b0dd-4e53-9149-addd80903fea"                             |
| Ipv4     | A v4 IP address                          | "192.168.1.1"                                                      | 
| Email    | Email address                            | "test@example.com"                                                 | 
| Md5      | An MD5 hash                              | "a29a16b688cc7167b705adc5744d7c62"                                 |
| Sha1     | A SHA1 hash                              | "13ff4d65e5602cc18658d8cc05116ba49a2fde9a"                         |
| Sha256   | A SHA 256 hash                           | "b7e0d35387a6026c7fd1b7a3e5f583545c22b81574444164fb73f1def314430f" |
| Sha512   | A SHA 512 hash                           | Like that ^, but even longer                                       |


### Numbers

Validating numbers makes use of generics to specify size and type.  For example,
if you want to make sure that an integer value is greater than 10, you could use
something like this:

```
// ensure number of type int is greater than 10
validator := ensure.Number[int]().IsGreaterThan(10)
```

If you were expecting a float, it would instead look like this:

```
// ensure number of type float is greater than 10
validator := ensure.Number[float64]().IsGreaterThan(10.0)
```

#### Methods

| Method                      | Description                                                                                         |
|-----------------------------|-----------------------------------------------------------------------------------------------------|
| Equals(num)                 | Passes if the tested number is exactly the same as the provided value                               |
| DoesNotEqual(num)           | Passes if the tested number is not the same as the provided value                                   |
| IsInRange(low, high)        | Passes if the tested number is greater than or equal to the low value and lower than the high value |
| IsLessThan(num)             | Passes if the tested number is less than the provided value                                         |
| IsLessThanOrEqualTo(num)    | Passes if the tested number is less than or equal to the the provided value                         |
| IsGreaterThan(num)          | Passes if the tested number is greater than the provided value                                      |
| IsGreaterThanOrEqualTo(num) | Passes if the tested number is greater than or equal to the provided value                          |
| IsEven()                    | Passes if the tested number is even                                                                 |
| IsOdd()                     | Passes if the tested number is odd                                                                  |
| IsPositive()                | Passes if the tested number is greater than zero                                                    |
| IsNegative()                | Passes if the tested number is less than zero                                                       |
| IsZero()                    | Passes if the tested number is zero                                                                 |
| IsNotZero()                 | Passes if the tested number is not zero                                                             |
| IsOneOf([]num)              | Passes if the tested number is in the passed array                                                  |
| IsNotOneOf([]num)           | Passes if the tested number is not in the passed array                                              |
| Is(func (num) error)        | Passes if the function passed does not produce an error during validation                           |

### Arrays

Array validators can be of any subtype (string, int, float64, etc).  As an
example, to make sure your array of 16-bit unsigned integers has at least 3 
values, you could use something like this:

```
// ensure array of uint16 has more than 3 elements
validator := ensure.Array[uint16]().HasMoreThan(3)
```

While the array validator does have its own methods you can use for validating
the array directly, most of the time you will want to apply validation to each of the
items it contains.  You can do this by passing an appropriate validator to the
`Each()` method.  For example, to make sure each string in an array contains at
least one vowel, you could do something like this:

```
// ensure in array of strings each string matches "(?i)[aeiuo]+"
validator := ensure.Array[string]().Each(
    ensure.String().Matches("(?i)[aeiuo]+")
 )
```

| Method                 | Description                                                               |
|------------------------|---------------------------------------------------------------------------|
| IsEmpty()              | Passes if tested array is empty (len(arr) == 0)                           |
| IsNotEmpty()           | Passes if tested array is not empty (len(arr) != 0)                       |
| HasCount(int)          | Passes if the length of the tested array is equal to the passed int       |
| HasFewerThan(int)      | Passes if the length of the tested array is less than the passed int      |
| HasMoreThan(int)       | Passes if the length of the tested array is more than the passed int      |
| HasLengthWhere(v)      | Adds a number validator that evaluates against the length of the array    |
| Each(v)                | Passes if the provided validator passes for each element in the array     |
| Is(func ([]T]) error)  | Passes if the function passed does not produce an error during validation |


### Maps

Validating maps is very similar to validating arrays.  The key difference is 
that with maps you can apply validations to both keys and values.

```
validator = ensure.Map[string, int]().EachKey(
    ensure.String()
).EachValue(
    ensure.Number[int]()
)
```

| Method                   | Description                                                               |
|--------------------------|---------------------------------------------------------------------------|
| IsEmpty()                | Passes if tested map is empty (len(arr) == 0)                             |
| IsNotEmpty()             | Passes if tested map is not empty (len(map) != 0)                         |
| HasCount(int)            | Passes if the length of the tested map is equal to the passed int         |
| HasFewerThan(int)        | Passes if the length of the tested map is less than the passed int        |
| HasMoreThan(int)         | Passes if the length of the tested map is more than the passed int        |
| HasLengthWhere(v)        | Adds a number validator that evaluates against the length of the map      |
| EachKey(v)               | Passes if the provided validator passes for each key in the map           |
| EachValue(v)             | Passes if the provided validator passes for each value in the map         |
| Is(func (map[K]V) error) | Passes if the function passed does not produce an error during validation |


### A note on lengths

Validators for types that have a length property (string, map, array) all have a
method that enables comprehensive number validation against their length.  This
method, `HasLengthWhere()` accepts a single number validator instance with arbitrary
rules.  For instance, to only allow strings that have an odd number of characters,
you could do something like this:

```
ensure.String().HasLengthWhere(
    ensure.Length().IsOdd()
)
```

Note the use of the `Length()` function.  This is a convenience function that returns
exactly the right type of number validator for evaluating length properties, so you
should use this anytime you need to validate based on length.

These same validators also have a small number of convenience functions for 
evaluating common length scenarios, such as whether or not an array is empty.  For
these common cases, you should prefer these methods instead for their conciseness.

Compare this:
```
ensure.Array[int]().IsNotEmpty()
```

to this:
```
ensure.Array[int]().HasLengthWhere(ensure.Length().DoesNotEqual(0))
```

### Structs

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

| Method             | Description                                                               |
|--------------------|---------------------------------------------------------------------------|
| Is(func (T) error) | Passes if the function passed does not produce an error during validation |

Note: due to the way visibility works in Go, only exported struct fields are
able to be validated.  That is, you can validate `MyStruct.Foo` but not `MyStruct.foo`.

