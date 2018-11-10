# api-testr
A package used to run API tests defined in JSON files.

# Tests
Tests are contained in a single JSON file - [Example test here](tests/example.json).

A test belongs to a single group and can be ordered within that group.

Tests will execute a single request whose response is then validated by a list of checks. 

### Groups

Tests are executed group by group.

You can add a test to a group using the `group` JSON key. If no group is provided then `default` is used.

### Order

If you need your tests executed in a specific order you can use the `order` JSON key. If no order is provided then `0` is used.

Tests with the same group and order will be run at the same time.

### Checks

Checks are how you validate that the response returned is correct.

#### Body Equal
Checks that the body returned is exactly equal to the value given.
```
{
  "type": "bodyEqual",
  "data": {
    "value": "OK"
  }
}
```

#### JSON Body Equal
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

#### JSON Body Query Exists
Queries the JSON body using [gjson](https://github.com/tidwall/gjson) and ensures that the queried element exists.
```
{
  "type": "jsonBodyQueryExists",
  "data": {
    "query": "title"
  }
}
```

#### JSON Body Query Equal
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

#### JSON Body Query Equal
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

#### Status Code Equal
Checks that the status code returned matches the given value.
```
{
  "type": "statusCodeEqual",
  "data": {
    "value": 200
  }
}
```