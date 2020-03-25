# apitestr

[![Go Report Card](https://goreportcard.com/badge/github.com/TomWright/apitestr)](https://goreportcard.com/report/github.com/TomWright/apitestr)
[![Documentation](https://godoc.org/github.com/TomWright/apitestr?status.svg)](https://godoc.org/github.com/TomWright/apitestr)
![Test](https://github.com/TomWright/apitestr/workflows/Test/badge.svg)
![Build](https://github.com/TomWright/apitestr/workflows/Build/badge.svg)

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

There is an optional `dataId` property you can set in the data object of this check. If this property is not empty, the value found by this check will be stored under the given `dataId` for use by subsequent tests.

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

There is an optional `dataId` property you can set in the data object of this check. If this property is not empty, the value found by this check will be stored under the given `dataId` for use by subsequent tests.

### JSON Body Query Regex Match
Queries the JSON body using [gjson](https://github.com/tidwall/gjson) and ensures that the queried element matches the given regex pattern.
```
{
  "type": "jsonBodyQueryRegexMatch",
  "data": {
    "query": "title",
    "pattern": "([a-z]{8}) ([a-z]{3}) ([a-z]{5})"
  }
}
```

There is an optional `dataIds` property you can set in the data object of this check. If this property is not empty, the values found in matching groups by this check will be stored under the given `dataIds` for use by subsequent tests.

E.g.
This response:
```
{
    "message": "Hello there, Tom"
}
```
With this check:
```
{
  "type": "jsonBodyQueryRegexMatch",
  "data": {
    "query": "message",
    "pattern": "([a-zA-Z]+), ([a-zA-Z]+)",
    "dataIds": {
        "0": "responseMessage",
        "1": "responseGreeting",
        "2": "responseName"
    }
  }
}
```
Will result in the following variables being available for use later on:
- `$.responseMessage`: `Hello there, Tom`
- `$.responseGreeting`: `Hello there`
- `$.responseName`: `Tom`

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

## Request Initialisation
Sometimes you'll need to do some dynamic testing, and that's where request init functions come in.

Register your init funcs as follows:
```
var myInitFunc RequestInitFunc = func(ctx context.Context, req *http.Request, data map[string]interface{}) (*http.Request, error) {
    // modify or create a new request here
    // use the data map as required
    return req, nil
}
testr.ContextWithRequestInitFunc(ctx, "my-init-func-id", myInitFunc)
```

And then use them in your tests as so:
```
{
    "request": {
        "init": {
            "my-init-func-id": {
                "some": "data",
            }
        }
    }
}
```

### Common init funcs
Some common init funcs are provided.

#### Request replacements
The replacements init func allows you to replace placeholders in the request URL path, URL query, headers and body with a given value.

It can be registered as follows:
```
testr.ContextWithRequestInitFunc(ctx, "replacements", testr.RequestReplacements)
```

Used in tests as so:
```
{
    "request": {
        "init": {
            "replacements": {
                ":name:": "Tom",
            }
        },
        "url": "https://example.com/users?name=:name:"
    }
}
```

If you want to use a data value that has been stored in the context by another test you should use `$.my-data-item` as the replacement value, where the previous test had used `my-data-item` as the `dataId`.
Or, moving on from the example given previously in *JSON Body Query Regex Match* you would do something like this:
```
{
    "request": {
        "init": {
            "replacements": {
                ":name:": "$.responseName",
            }
        },
        "url": "https://example.com/users?name=:name:"
    }
}
```
