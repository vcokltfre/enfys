# Enfys

An environment helper for Go.

## Usage Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/vcokltfre/enfys"
)

type Config struct {
	String string  `enfys:"STRING,required"`
	Int    int     `enfys:"INT"`
	Float  float64 `enfys:"FLOAT"`
	Bool   bool    `enfys:"BOOL"`
	Debug  bool    `enfys:"DEBUG"`
}

func main() {
	os.Setenv("STRING", "hello")
	os.Setenv("INT", "42")
	os.Setenv("FLOAT", "3.14")
	os.Setenv("BOOL", "true")

	var cfg Config
	err := enfys.Fill(&cfg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", cfg)
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
