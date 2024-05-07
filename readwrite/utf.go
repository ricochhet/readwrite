/*
 * readwrite
 * Copyright (C) 2024 readwrite contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
