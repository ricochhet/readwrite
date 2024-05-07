package readwrite

import (
	"bytes"
	"fmt"
	"testing"
)

func TestUtf8ToUtf16(t *testing.T) {
	b := Utf8ToUtf16("aaabbbccc")
	o := []byte{97, 0, 97, 0, 97, 0, 98, 0, 98, 0, 98, 0, 99, 0, 99, 0, 99, 0}
	if !bytes.Equal(b, o) {
		t.Fatal(fmt.Errorf("unexpected bytes"))
	}
}
