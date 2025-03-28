# Strings

You've probably already seen a few partial examples of a string validator, but 
here's something a bit more complete (if no less contrived).
```go
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

## Methods

| Method                | Description                                                                             |
|-----------------------|-----------------------------------------------------------------------------------------|
| IsEmpty()             | Passes if the tested string is empty (len() == 0)                                       |
| IsNotEmpty()          | Passes if the tested string is not empty (len() != 0)                                   |
| Equals(str)           | Passes if the tested string is identical to the provided string                         |
| DoesNotEqual(str)     | Passes if the tested string is not identical to the provided string                     |
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

## Predefined Regex Patterns

To simplify common use cases, there are predefined regular expressions available
to use.  For example, if you wanted to ensure that the string only contains
alphanumeric values, you could use something like this:

```go
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
