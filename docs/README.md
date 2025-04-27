## Types

You can read more about the options for validating different value types below. 
There are some code snippets for each type, but if you want fully runnable examples,
check out the [_examples](../_examples) directory.

| Type   | Basic Usage                                                                 | Validator Type              | Documentation           |
|--------|-----------------------------------------------------------------------------|-----------------------------|-------------------------|
| String | `ensure.String().IsNotEmpty().StartsWith('abc')`                            | `ensure.StringValidator`    | [Strings](./strings.md) |
| Number | `ensure.Number[int]().IsGreaterThan(0)`                                     | `ensure.NumberValidator[T]` | [Numbers](./numbers.md) |
| Array  | `ensure.Array[string]().Each( ensure.String().Matches("^\d+$") )`           | `ensure.ArrayValidator[T]`  | [Arrays](./arrays.md)   |
| Map    | `ensure.Map[string,int]().EachKey( ensure.String().HasLength(3) )`          | `ensure.MapValidator[K,V]`  | [Maps](./maps.md)       |
| Struct | `ensure.Struct[MyStruct]( with.Validators{ "Foo": ensure.Number[int]() } )` | `ensure.StructValidator[T]` | [Structs](./structs.md) |
| Bool   | `ensure.Bool().IsTrue()`                                                    | `ensure.BooleanValidator`   | [Bools](./bools.md)     |


## Validator interfaces

There are two interfaces that define specific validation functionality: `Validator[T]`
and `UntypedValidator`.  Both require a `Type` method, which returns a string
with the type of value expected by the validator (eg "string", "int64", etc.).  
Where they differ is that `Validator[T]` has a method `Validate(T)` that validates
a strongly typed value and `UntypedValidator` has `ValidateUntyped(any)`, which 
can validate values where you don't know the type of the validator or the value 
(or both).  Nearly all validation should be handled through the `Validate` method,
which enables compile-time type checking and is both safer and more performant,
but the `ValidateUntyped` method is there when you need it (such as when reflection
is involved).  All validators in this library implement both interfaces.

## Validation options

Most of the time, default validation behavior is sufficient.  Sometimes, however,
you may want to have more control over how the validations run, especially for
complex validation trees.  You can pass validation options at validation time
by adding them to the `Validate()` call.

```go
validator.Validate(value, with.Options(...))
```

Currently, the only available option is `OptionCollectAllErrors()`.  The default
behavior is to stop processing validation checks as soon as the first error is 
encountered and return that immediately.  The `OptionCollectAllErrors()` option 
changes validation so that it instead collects all validation errors and returns
them together in a `ValidationErrors` struct.  You can read more about this option
and the `ValidationErrors` error type in the [errors](./errors.md) documentation.


## Pointers

All validators operate on values, not pointers.  Aside from the fact that the 
value is what contains the useful information for validation, this also ensures
that no validator can inadvertently mutate any values passed to it.  Pointers
are a fact of life, though, especially when dealing with struct fields.  In most
cases you can simply dereference the pointer, but other times, like when validating
structs or arrays of pointers, it can be easier to just indicate that
a pointer is expected.

There are two functions you can use to indicate that a passed value will be a
pointer: `Pointer[T]()` and `OptionalPointer[T]()`.  The only difference between the
two is that `Pointer[T]()` will return an error if the pointer is nil, whereas
`OptionalPointer[T]()` will return gracefully without attempting further validation.

```go
type Person struct {
	FirstName string
	LastName string
	Pets []*Pet
}

type Pet struct {
	Name string
	Type string
	License *License
}

validPet := ensure.Struct[Pet]().HasFields(with.Validators{
	"Name": ensure.String(),
	"Type": ensure.String(),
	
	// Not every pet will need a license, so only validate if it exists
	"License": ensure.OptionalPointer[License](
		ensure.Struct[License]()
	)
})

validPerson := ensure.Struct[Person]().HasFields(with.Validators{
	"FirstName": ensure.String(),
	"LastName": ensure.String(),
	
	// We may not have any pets, but if we do each one should be valid
	"Pets": ensure.Array[*Pet]().Each(
		ensure.Pointer[Pet](validPet)
	),
})
```

There is no practical limit to how far you can nest pointers.

```go
str := "foo"
pStr := &str
ppStr := &pStr

validStr := ensure.Pointer[*string](
	ensure.Pointer[string](
		ensure.String()
	),
)

if err := validStr.Validate(ppStr); err != nil {
	// ...
}
```

## The `Any` validator

Validation rules added to a validator are evaluated in order via logical "AND".

```go
// ensure string starts with "foo" *AND* ends with "bar"
ensure.String().StartsWith("foo").EndsWith("bar")
```

There may be some occasions where you want to consider multiple validation options
for the same value.  For example, in the case where you want to validate a hostname,
valid options could be an IPv4 address, an IPv6 address, a fully-qualified domain
name, or maybe even just the word "localhost".  We can use the `Any[T]()` validator
to evaluate each of these in turn and consider the value valid if any of them 
evaluates without error.

```go
validHost := ensure.Any[string](
	ensure.String().Equals("localhost"),
	ensure.String().Matches(ensure.Ipv4),
	// ...
)
```

## Lengths

Validators for types that have a length property (string, map, array) all have a
method that enables comprehensive number validation against their length.  This
method, `HasLengthWhere()` accepts a single number validator instance with arbitrary
rules.  For instance, to only allow strings that have an odd number of characters,
you could do something like this:

```go
ensure.String().HasLengthWhere(
    ensure.Length().IsOdd()
)
```

Note the use of the `Length()` function.  This is a convenience function that returns
a number validator with the right generic type for evaluating length properties, so you
should use this anytime you need to validate based on length.

These same validators also have a small number of convenience functions for 
evaluating common length scenarios, such as whether or not an array is empty.  For
these common cases, you should prefer these methods instead for their conciseness.

Compare this:
```go
ensure.Array[int]().IsNotEmpty()
```

to this:
```go
ensure.Array[int]().HasLengthWhere(ensure.Length().IsGreaterThan(0))
```


## The `Is`/`Has` methods

While every effort has been made to provide a comprehensive set of validations
for the broadest set of types and values possible, there are some validation
rules that defy simple boolean logic.  In these cases where combining multiple
rules still isn't enough to get the desired results, the `Is()` method (and its
alias, `Has()`) can be used to provide a function with arbitrary logic that will
be evaluated the same as any other rule.

Consider a situation where we want to make sure that an expiration date is not 
a time in the past and is less than 90 days in the future.  Here's one way you 
could define that using the `Is()` function to add each rule independently.

```go
func notInThePast(date time.Time) error {
    if date.Before(time.Now()) {
        return errors.New("expiration date cannot be in the past")
    }
	
    return nil
}

func lessThanNinetyDaysFromNow(date time.Time) error {
    if date.After(time.Now().Add(90 * 24 * time.Hour)) {
        return errors.New("expiration date cannot be more than 90 days in the future")
    }
	
    return nil
}

validExpiration := ensure.Struct[time.Time]().Is(notInThePast).Is(lessThanNinetyDaysFromNow)
```

Here's an alternate version that combines both rules into a single function.
Both options are functionally the same, so do whatever works best for you.

```go
func inExpectedTimeRange(date time.Time) error {
    now := time.Now()

    if date.Before(now) {
        return errors.New("expiration date cannot be in the past")
    }
		
    if date.After(now.Add(90 * 24 * time.Hour)) {
        return errors.New("expiration date cannot be more than 90 days in the future")
    }
	
    return nil
}

validExpiration := ensure.Struct[time.Time]().Is(inExpectedTimeRange)
```

You can, of course, also pass function literals

```go
validExpiration := ensure.Struct[time.Time]().Is(func (date time.Time) error {
    now := time.Now()

    if date.Before(now) {
        return errors.New("expiration date cannot be in the past")
    }
		
    if date.After(now.Add(90 * 24 * time.Hour)) {
        return errors.New("expiration date cannot be more than 90 days in the future")
    }
	
    return nil
})
```


