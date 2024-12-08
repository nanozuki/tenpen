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
  var rule = `["$+", "#a", "#b"]`

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
    "sum": ["$+", "#a", "#b.c"],
    "double-sum": ["$*", "#sum", 2]
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

1. math: `+`, `-`, `*`, `/`, `%`, `^`
1. boolean: `and`, `or`, `not`, `==`, `!=`, `>`, `<`, `>=`, `<=`
1. strings: `+`, `len`, `split`, `join`
1. array: `len`, `get`, `filter`, `map`, `reduce`
1. object: `len`, `get`, `keys`, `values`
1. flow control: `if`, `cond`, `let`, `def`, `do`, `apply`

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
