# api-testr

[![Build Status](https://travis-ci.org/TomWright/api-testr.svg?branch=master)](https://travis-ci.org/TomWright/api-testr)
[![codecov](https://codecov.io/gh/TomWright/api-testr/branch/master/graph/badge.svg)](https://codecov.io/gh/TomWright/api-testr)
[![Go Report Card](https://goreportcard.com/badge/github.com/TomWright/api-testr)](https://goreportcard.com/report/github.com/TomWright/api-testr)
[![Documentation](https://godoc.org/github.com/TomWright/api-testr?status.svg)](https://godoc.org/github.com/TomWright/api-testr)

A package used to run API tests defined in JSON files.

## Tests
Tests are contained in a single JSON file - [Example test here](tests/example.json).

A test belongs to a single group and can be ordered within that group.

Tests will execute a single request whose response is then validated by a list of checks. 

### Groups

Tests are executed group by group.

You can add a test to a group using the `group` JSON key. If no group is provided then `default` is used.

### Order

If you need your tests executed in a specific order you can use the `order` JSON key. If no order is provided then `0` is used.

Tests with the same group and order will be run at the same time.

## Running Tests

### Running a single test
```
// a context is always required
ctx := context.Background()

// set the base url to be used with the test
ctx = testr.ContextWithBaseURL(ctx, "https://example.com")

// parse the test file
t, err := parse.File(ctx, "path/to/my/test.json")
if err != nil {
    panic(err)
}

// run the test
err := testr.Run(t, nil)

// handle the error if the test failed
if err != nil {
    panic(fmt.Errorf("test `%s` failed: %s\n", t.Name, err))
}
```

### Running groups of tests
```
// a context is always required
ctx := context.Background()

// set the base url to be used with the test
ctx = testr.ContextWithBaseURL(ctx, "https://example.com")

// parse the test files
tests := make([]*testr.Test, 2)
var err error
tests[0], err = parse.File(ctx, "path/to/my/test1.json")
if err != nil {
    panic(err)
}
tests[1], err = parse.File(ctx, "path/to/my/test1.json")
if err != nil {
    panic(err)
}

res := testr.RunAll(testr.RunAllArgs{}, tests...)

// log the results
log.Printf("tests finished\n\texecuted: %d\n\tpassed: %d\n\tfailed: %d", res.Executed, res.Passed, res.Failed)
```

## Checks

Checks are how you validate that the response returned is correct.

### Body Equal
Checks that the body returned is exactly equal to the value given.
```
{
  "type": "bodyEqual",
  "data": {
    "value": "OK"
  }
}
```

### JSON Body Equal
Checks that the body returned matches the given JSON object.
```
{
  "type": "jsonBodyEqual",
  "data": {
    "value": {
      "userId": 1,
      "id": 1,
      "title": "delectus aut autem",
      "completed": false
    }
  }
}
```

### JSON Body Query Exists
Queries the JSON body using [gjson](https://github.com/tidwall/gjson) and ensures that the queried element exists.
```
{
  "type": "jsonBodyQueryExists",
  "data": {
    "query": "title"
  }
}
```

### JSON Body Query Equal
Queries the JSON body using [gjson](https://github.com/tidwall/gjson) and ensures that the queried element has a value equal to the one specified.
```
{
  "type": "jsonBodyQueryEqual",
  "data": {
    "query": "title",
    "value": "delectus aut autem"
  }
}
```

### JSON Body Query Equal
Queries the JSON body using [gjson](https://github.com/tidwall/gjson) and ensures that the queried element matches the given regex pattern.
```
{
  "type": "jsonBodyQueryRegexMatch",
  "data": {
    "query": "title",
    "pattern": "[a-z]{8} [a-z]{3} [a-z]{5}"
  }
}
```

### Status Code Equal
Checks that the status code returned matches the given value.
```
{
  "type": "statusCodeEqual",
  "data": {
    "value": 200
  }
}
```

### Custom Body Check
Reads the response body into a byte array and provides it to your custom function to validate, identified by the `id` value.
```
{
  "type": "bodyCustom",
  "data": {
    "id": "123check"
  }
}
```

#### Creating your custom check
First create your custom validator func - this *must* implement the `check.BodyCustomCheckerFunc` interface.
```
var custom123 check.BodyCustomCheckerFunc = func(bytes []byte) error {
    if string(bytes) != "[1,2,3]" {
        return fmt.Errorf("response is not 1,2,3")
    }
    return nil
}
```

Then register the custom func against your context.
```
ctx = testr.ContextWithCustomBodyCheck(ctx, "123check", custom123)
```

Then simply use the id of `123check` in your `bodyCustom` check data.
