# Numbers

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

## Methods

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
