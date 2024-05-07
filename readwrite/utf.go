package readwrite

import "unicode/utf16"

func Utf8ToUtf16(utf8str string) []byte {
	utf8Bytes := []byte(utf8str)
	utf16Runes := utf16.Encode([]rune(string(utf8Bytes)))
	utf16Bytes := make([]byte, len(utf16Runes)*2)
	for i, r := range utf16Runes {
		utf16Bytes[i*2] = byte(r)
		utf16Bytes[i*2+1] = byte(r >> 8)
	}
	return utf16Bytes
}
