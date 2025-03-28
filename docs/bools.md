# Booleans

Boolean values are typically straightforward; there can really only be two
legitimate values, three if you count nil pointers. There are cases, however,
where it may be helpful to be able to test these values. In the case where you
want conditional validation, you might have something like this:

```go
validConfig := ensure.Any(
    // Validation succeeds if component is disabled regardless of what the rest looks like
    // If we aren't using it, we don't need to bother validating it
    ensure.Struct[Config]().HasFields(with.Validators{
        "Enabled": ensure.Boolean().IsFalse(),
    },
    
    // Otherwise we validate the rest to make sure it looks the way it is supposed to
    ensure.Struct[Config]().HasFields(with.Validators{
        // the rest of the validation goes here
    },
)
```

## Methods

| Method               | Description                     |
|----------------------|---------------------------------|
| IsTrue()             | Passes if tested value is true  |
| IsFalse()            | Passes if tested value is false |


## `Is`

Note that the boolean validator lacks the `Is` method most other validators have,
due to the simplicity involved in validating boolean values.  There may be a case
for inclusion based on the ability to determine whether a value should be true or
false based on dynamic external conditions, but for now it seems that may be a 
solution in search of a problem.

If you have a situation where this would be valuable, add a comment on 
[this issue](https://github.com/chriscasto/go-ensure/issues/36) with a 
description of your use case.