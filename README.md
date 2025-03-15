# tenpen

A json format rule engine with declarative, functional (lisp-like) programming
language.

The rule and data are both json format. It is easy to contructure and visualize
in different languages and platforms.

## Usage Example

### directly evaluate value

```go
import (
  tp "github.com/tenpen/tenpen"
)

func main() {
  var input = `{"a": 1, "b": 2}`
  var rule = `["$add", "#a", "#b"]`

  output, err := tp.Eval(rule, input)
  println(output, err) // "3" nil
}
```

### evaluate with recursive object

```go
import (
  tp "github.com/tenpen/tenpen"
)

func main() {
  var input = `{"a": 1, "b": {"c": 2}}`
  var rule = `{
    "sum": ["$add", "#a", "#b.c"],
    "double-sum": ["$mul", "#sum", 2]
  }`

  output, err := tp.Eval(rule, input)
  println(output, err) // {"sum": 3, "double-sum": 6} nil
}
```

## Exported Functions

```go
func Eval(rule, envs... string) (string, error)
```

Evaluate the rule with the given environment, both `rule`, `envs`, and return
string are json format.

## Value

### Types

Any types of json value are supported, including: number, string, boolean, null,
array, object.

### Recursive Object

If the rule is an object, the value of each key will be evaluated recursively.
You can use `#` to reference the value.

### Evaluation rules for value reference

1. Use `#` to reference the value in rule and environments. Use `##` to escape
   a string that starts with `#`.
1. find the value by key in rule object.
1. if not found, find the value by key in last environment object.
1. if not found, find the value by key in 2nd last environment object, and so
   on.
1. if not found in rule and any environment, return error.
1. No circular reference is allowed. Including reference to self, parent, or
   child.

### Value path

You can use simple json path to reference the value:

```
#key1.key2.key3
#array.1.key
```

### Function

Use `$<name>` to call a function, and use `$$` to escape a string that starts
with `$`. A function call is a json list, the first element is `$<name>`, and
the rest elements are arguments.

### Built-in Functions

#### Basic Operators

| Name    | Parameters       | Return   | Description                       |
| ------- | ---------------- | -------- | --------------------------------- |
| `$add`  | `...number`      | `number` | add a sequence of numbers         |
| `$add`  | `...string`      | `string` | concatenate a sequence of strings |
| `$sub`  | `number, number` | `number` | subtract two numbers              |
| `$mul`  | `...number`      | `number` | multiply a sequence of numbers    |
| `$div`  | `number, number` | `number` | divide two numbers                |
| `$mod`  | `number, number` | `number` | modulo two numbers                |
| `$pow`  | `number, number` | `number` | power two numbers                 |
| `$eq`   | `any, any`       | `bool`   | equal                             |
| `$ne`   | `any, any`       | `bool`   | not equal                         |
| `$gt`   | `number, number` | `bool`   | greater than                      |
| `$lt`   | `number, number` | `bool`   | less than                         |
| `$gte`  | `number, number` | `bool`   | greater than or equal             |
| `$lte`  | `number, number` | `bool`   | less than or equal                |
| `$both` | `...boolean`     | `bool`   | and                               |
| `$any`  | `...boolean`     | `bool`   | or                                |
| `$not`  | `bool`           | `bool`   | not                               |
| `$len`  | `array`          | `number` | length of array                   |
| `$len`  | `string`         | `number` | length of string                  |

#### Conditional

| Name      | Parameters                                 | Return   | Description                                                                                    |
| --------- | ------------------------------------------ | -------- | ---------------------------------------------------------------------------------------------- |
| `$either` | `bool, any, any`                           | `any`    | if the first argument is true, return the second argument, otherwise return the third argument |
| `$match`  | `...[bool, any]`                           | `any`    | match the first true condition, and return the value of in the next argument                   |
| `$map`    | `array, func{number, any}->any`            | `array`  | map the function to each element in array                                                      |
| `$map`    | `object, func{string, any}->any`           | `object` | map the function to each value in object                                                       |
| `$filter` | `array, func{number, any}->bool`           | `array`  | filter the array by the function                                                               |
| `$filter` | `object, func{string, any}->bool`          | `object` | filter the object by the function                                                              |
| `$reduce` | `array, any, func{number, any, any}->any`  | `any`    | reduce the array by the function                                                               |
| `$reduce` | `object, any, func{string, any, any}->any` | `any`    | reduce the object by the function                                                              |

### Define Function

`$def` is a special function to define a function. The first argument is the
argument list, and the second argument is the body of the function.

```json
{
  "double": ["$def", ["x"], ["$*", "#x", 2]],
  "double_sum": ["$map", "#input_array", "$double"]
}
```

The function won't be evaluated to the result.
