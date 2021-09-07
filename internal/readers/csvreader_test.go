package readers

import (
	"bytes"
	"strings"
	"testing"
)

func TestCSVReader(t *testing.T) {

	ch := make(chan []string)
	var buffer bytes.Buffer

	for i := 0; i < 1000; i++ {
		buffer.WriteString("a")
	}

	const input = `Beware of bugs in the above code;
I have only proved it correct, not tried it.`

	scanner := strings.NewReader(input)

	CSVReader(scanner, ch)
}