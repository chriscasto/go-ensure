# Error Handling

## Error types

There are two main types of errors generated by validators: `TypeError` and
`ValidationError`.  `TypeError`s occur when a validator attempts to validate
an incompatible type.  This most commonly occurs when using the `ValidateUntyped()`
method instead of the `Validate()` method, since it doesn't use the compiler's
type checker.  These types of errors indicate a misconfiguration of some
kind and should only be logged and not passed back to the user.

The other type of error is `ValidationError`, which occurs when a validation
check fails on a particular value.  For example, an `IsEqualTo(3)` check would
return a `ValidationError` if passed the number 4.  These errors are expected
to occur rather frequently and are harmless on their own.  `ValidationError`s
are crafted in such a way that they should be safe to send directly back to users,
such as when you are validating user input and want to provide feedback on what
they did wrong.

Given the large difference in the "safety" of these two error types, it is
recommended to check the type of error returned by a call to `Validate()` to
make sure it is being handled appropriately.

There is a third error type used when the `OptionCollectAllErrors()` option is
passed at validation type.  This option causes all errors, both `TypeError` and
`ValidationError`, to be collected into a single struct of type `ValidationErrors`.
If the `Error()` method is called on this, the first validation error is returned,
so it's compatible with standard error handling.  You can (and should) specifically
request the list of validation errors and type errors directly with the 
`ValidationErrors()` and `TypeErrors()` methods and handle them accordingly.

```go
validStr := ensure.String().HasLength(3).Matches(ensure.Number)

// this value will fail both validation checks
str := "a"

// validate normally
if err := validStr.Validate(str); err != nil {
	// only one error is available 
	fmt.Println(err.Error())
}

// now, set the OptionCollectAllErrors option
opts := with.Options(
	with.OptionCollectAllErrors()
)

// validate again, passing the options to the Validate method
if err := validStr.Validate(str, opts); err != nil {
	// this prints the same thing as the code above 
	fmt.Println(err.Error())
	
	// convert to "ValidationErrors" type 
	if errs := ensure.ErrorAsValidationErrors(err); errs != nil {
		// check if there are any errors 
		if errs.HasValidationErrors() {
			// print each one 
			for _, e := range errs.ValidationErrors() {
				// there should be two errors printed 
				fmt.Println(e.Error())
			}
		}
	}
}
```

## Construction errors

Validation objects are intended to be constructed infrequently, typically once
at startup, then used many times, often for the life of the program.  It's
generally helpful to consider validators to be statically compiled constants
rather than variables to be mutated dynamically.  This is not a hard rule, and
there are certainly some valid use cases for dynamic construction, but be aware
that validator builder methods will panic if used incorrectly, so be prepared to
`recover()` if you decide to instantiate validators dynamically outside of program
initialization.

We panic rather than return an error because an invalid validator is a risk to
the security and integrity of your application. Using "soft" errors runs the
risk of enabling a situation where the application may be running with validation
that is incomplete or otherwise working differently than intended.

Most initialization errors can be caught by the compiler, such as mixing types
(`Number[int]().IsLessThan(10.0)`), passing the wrong validator type to
another validator (`Any[int](String())`), or attempting to call
an invalid method (`String().IsGreaterThan()`), but there are some cases where
correctness has to be checked during validator composition.  For example:

```go
type MyStruct {
    Foo int64
    Bar string
}

validStruct := ensure.Struct[MyStruct]().HasFields(with.Validators{
    // This will panic because the declared type doesn't match the actual field type
    "Foo": ensure.Number[int](),
    
    // This will panic because the name of the field is wrong
    "Baz": ensure.String(),
})
```

In all cases, panics are used to indicate unrecoverable conditions that arise
solely due to invalid configurations.