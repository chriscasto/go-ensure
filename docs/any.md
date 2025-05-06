
# Any

With the `Any` validator, you can check to see if any of a set of validators passes.

## Basic Usage

```go
ensure.Any[T](
	Validator[T], 
	Validator[T], 
	Validator[T], 
	...
)
```

## Example

```go
type HttpServerConfig struct {
	Enabled bool
	Addr string
	Port int
}

validConfig := ensure.Any[HttpServerConfig](
	// If the server is disabled, we don't care about the other values
	// This is a valid configuration
	ensure.Struct[HttpServerConfig]().HasFields(with.Validators{
		"Enabled": ensure.Bool().IsFalse(),
	}),
	
	// Otherwise, make sure we have a valid address and port 
	ensure.Struct[HttpServerConfig]().HasFields(with.Validators{
		// We can bind using either an IPv4 IP or the word "localhost"
		"Addr": ensure.Any[string](
			ensure.String().Equals("localhost"),
			ensure.String().Matches(ensure.Ipv4), 
		).WithError("must be a valid bind address"),
		
		// The port must be in a valid range
		"Port": ensure.Number[int]().IsInRange(1, 65536),
	}),
).WithOptions(
	// If we have any errors worth reporting, they're coming from the second validator
	with.AnyOptionPassThroughErrorsFrom(1)
)

// err is nil because Enabled is false
err := validConfig.Validate(&HttpServerConfig{
	Enabled: false,
})

// err is still nil, even though the other values are invalid
err = validConfig.Validate(&HttpServerConfig{
	Enabled: false,
	Addr: "zanzibar",
	Port: 1000000
})

// err is nil because Addr and Port are valid
err = validConfig.Validate(&HttpServerConfig{
	Enabled: true,
	Addr: "127.0.0.1",
	Port: 8080,
})

// "Addr: must be a valid bind address"
err = validConfig.Validate(&HttpServerConfig{
	Enabled: true, 
	Addr: "detroit", 
	Port: 8080,
})

```

## Methods

| Method              | Description                                                                            |
|---------------------|----------------------------------------------------------------------------------------|
| WithOptions(opt...) | Applies options to the validator                                                       |
| WithError(str)      | Sets a default error message.  Alias of `WithOptions(with.AnyOptionDefaultError(msg))` |

## Options

Options can be set on the validator by calling `WithOptions()`.

```go
ensure.Any[string](
	ensure.String().Equals("foo"),
	ensure.String().StartsWith("bar")
).WithOptions(
	with.AnyOptionDefaultError(`string must equal "foo" or start with "bar"`)
)
```

| Option              | Method                                      | Description                                                                                                     |
|---------------------|---------------------------------------------|-----------------------------------------------------------------------------------------------------------------|
| Set Default Error   | with.AnyOptionDefaultError(msg)             | Sets the default error to return if none of the validators pass                                                 |
| Pass Through Errors | with.AnyOptionPassThroughErrorsFrom(idx...) | Validator will pass through errors from child validators at the provided index(s) if at least one does not pass |
