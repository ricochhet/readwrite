# readwrite
A simple Go module to read and write data.

### Reader & Writer
- `reader, err := readwrite.NewReader(path)`
- `writer, err := readwrite.NewWriter(path, append=false)`

### UTF
- `bytes := Utf8ToUtf16("aaabbbccc")`

### PE
- View [pefile.go](./readwrite/pefile.go)
- Contains many functions I have found need for when manipulating bytes of the [Portable Executable](https://learn.microsoft.com/en-us/windows/win32/debug/pe-format) format.

## License
See LICENSE file.