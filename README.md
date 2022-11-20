# sgarg
Simple GNU-compliant Argument parser

## Feature
- 100% compliant with [GNU's conventions](https://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html)
- Depends only on stdlib

## Usage
```go
parser := sgarg.NewParser()

// set short bool option
var b []bool
parser.SetBoolOpt("b", *b)

// set short string option
var s []string
parser.SetStringOpt("s", *s)

if err := parser.Parse(); err != nil {
  log.Fatalln(err)
}
```

## Todo
- [x] Support short option
- [ ] Support long option （with abbreviations）
- [ ] Support option reorder mode
- [ ] Add tests

## Author
remiposo

## License
MIT
