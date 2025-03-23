## Types

You can read more about the options for validating different value types below. 
There are some code snippets for each type, but if you want fully runnable examples,
check out the [_examples](../_examples) directory.

| Type   | Basic Usage                                                             | Validator Type              | Documentation           |
|--------|-------------------------------------------------------------------------|-----------------------------|-------------------------|
| String | `ensure.String().IsNotEmpty().StartsWith('abc')`                        | `ensure.StringValidator`    | [Strings](./strings.md) |
| Number | `ensure.Number[int]().IsGreaterThan(0)`                                 | `ensure.NumberValidator[T]` | [Numbers](./numbers.md) |
| Array  | `ensure.Array[string]().Each( ensure.String().Matches("^\d+$") )`       | `ensure.ArrayValidator[T]`  | [Arrays](./arrays.md)   |
| Map    | `ensure.Map[string,int]().EachKey( ensure.String().HasLength(3) )`      | `ensure.MapValidator[K,V]`  | [Maps](./maps.md)       |
| Struct | `ensure.Struct[MyStruct]( with.Fields{ "Foo": ensure.Number[int]() } )` | `ensure.StructValidator[T]` | [Structs](./structs.md) |


## The `Validator` interface

Each of the validators listed above implements the `Validator` interface, which
defines two methods: `Type()` and `Validate(val)`.  The `Type()` method returns the
type of value expected by the validator (eg "string", "int64", etc), and the
`Validate(val)` method evaluates all the defined checks against the value passed
to it.  The value passed to the `Validate()` method can be of any type, which is
helpful in cases where you don't know the type of the validator or the value (or
both). Each validator will do the necessary checks to make sure that the type of
the value passed to it is the same as the type returned by `Type()`, and will 
return a `TypeError` in the event of a mismatch.

In the cases where you know the types for both the validator and the value to
validate, it is more efficient to let the compiler do all the necessary type
checks for you.  In addition to the `Validate()` method, each validator also has
a separate, typed method for validating values of the appropriate type.  It is 
recommended to use this method where possible.

| Type   | Validator                   | Typed Validation Method  |
|--------|-----------------------------|--------------------------|
| String | `ensure.StringValidator`    | `ValidateString(string)` |
| Number | `ensure.NumberValidator[T]` | `ValidateNumber(T)`      |
| Array  | `ensure.ArrayValidator[T]`  | `ValidateArray([]T)`     |
| Map    | `ensure.MapValidator[K,V]`  | `ValidateMap(map[K]V)`   |
| Struct | `ensure.StructValidator[T]` | `ValidateStruct(T)`      |


## Construction Errors

Validation objects are intended to be constructed infrequently, typically once 
at startup, then used many times, basically for the life of the program.  It's 
generally helpful to consider validators to be statically compiled constants
rather than variables to be mutated dynamically.  This is not a hard rule, and
there are certainly some valid use cases for dynamic construction, but be aware
that validator builder methods will panic if used incorrectly, so be prepared to 
`recover()` if you decide to instantiate validators dynamically outside of program
initialization.

We panic rather than return an error because an invalid validator is a risk to
the security and integrity of your application, and using "soft" errors runs the
risk of enabling a situation where the application may be running with validation
that is incomplete or otherwise working differently than intended.  If the 
validation cannot pass its own sanity checks, the app should be considered too 
unsafe to run.

Most initialization errors can be caught by the compiler (such as mixing types
like `Number[int]().IsLessThan(10.0)` or attempting to call an invalid method
like `String().IsGreaterThan()`), but there are some cases where correctness has
to be checked during validator composition.  For example:

```
type MyStruct {
    Foo int64
    Bar string
}

validStruct := ensure.Struct[MyStruct](with.Fields{
    # This will panic because the declared type doesn't match the actual field type
    "Foo": ensure.Number[int](),
    
    # This will panic because the name of the field is wrong
    "Baz": ensure.String(),
})
```

In all cases, panics are used to indicate unrecoverable conditions that arise 
solely due to invalid configurations.

## Lengths

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
a number validator with the right generic type for evaluating length properties, so you
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


## The `Is` method

While every effort has been made to provide a comprehensive set of validations
for the broadest set of types and values possible, there are some validation
rules that defy simple boolean logic.  In these cases where combining multiple
rules still isn't enough to get the desired results, the `Is()` method can be
used to provide a function with arbitrary logic that will be evaluated the same
as any other rule.

Consider a situation where we want to make sure that an expiration date is not 
a time in the past and is less than 90 days in the future.  Here's one way you 
could define that using the `Is()` function to add each rule independently.

```
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

```
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

```
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


