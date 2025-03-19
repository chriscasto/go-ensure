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

You could then test a value with the validator this produces using something like this:

```
if err := validator.Validate("baseball"); err != nil {
    fmt.Print(err)
} else {
    fmt.Print(`"baseball" is a valid string`)
}
```

Check out the [documentation](./docs/README.md) for more details, or the [_examples](./_examples)
directory for some runnable examples.

