## Types

You can read more about the options for validating different value types below. 
There are some code snippets for each type, but if you want fully runnable examples,
check out the [_examples](../_examples) directory.

| Type   | Basic Usage                                                             | Validator Type               | Documentation           |
|--------|-------------------------------------------------------------------------|------------------------------|-------------------------|
| String | `ensure.String().IsNotEmpty().StartsWith('abc')`                        | `ensure.StringValidator`     | [Strings](./strings.md) |
| Number | `ensure.Number[int]().IsGreaterThan(0)`                                 | `ensure.NumberValidator[T]`  | [Numbers](./numbers.md) |
| Array  | `ensure.Array[string]().Each(ensure.String().Matches("^\d+$"))`         | `ensure.ArrayValidator[T]`   | [Arrays](./arrays.md)   |
| Map    | `ensure.Map[string,int]().EachKey(ensure.String().HasLength(3))`        | `ensure.MapValidator[K,V]`   | [Maps](./maps.md)       |
| Struct | `ensure.Struct[MyStruct](ensure.Fields{ "Foo": ensure.Number[int]() })` | `ensure.StrucdtValidator[T]` | [Structs](./structs.md) |


## Construction Errors

Validation objects are intended to be constructed infrequently, typically once 
at startup, then used many times, basically for the life of the program.  It's 
generally helpful to consider validators to be statically compiled constants
rather than variables to be mutated dynamically.  This is not a hard rule, and
there are plenty of valid use cases for dynamic construction, but be aware that
validator builder methods will panic if used incorrectly, so be prepared to 
`recover()` if you decide to instantiate validators dynamically outside of program
initialization.

Most initialization errors can be caught by the compiler (such as mixing types
like `Number[int]().IsLessThan(10.0)` or attempting to call an invalid method
like `String().IsGreaterThan()`), but there are some cases where correctness has
to be checked during validator composition.  For example:

```
type MyStruct {
    Foo int64
    Bar string
}

validStruct := ensure.Struct[MyStruct](ensure.Fields{
    # This will panic because the declared type doesn't match the actual field type
    "Foo": ensure.Number[int](),
    
    # This will panic because the name of the field is wrong
    "Baz": ensure.String(),
})
```

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
just the right type of number validator for evaluating length properties, so you
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



