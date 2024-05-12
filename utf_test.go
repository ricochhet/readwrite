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

package readwrite_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/ricochhet/readwrite"
)

var errUnexpectedBytes = errors.New("unexpected bytes")

func TestUtf8ToUtf16(t *testing.T) {
	t.Parallel()

	b := readwrite.Utf8ToUtf16("aaabbbccc")
	o := []byte{97, 0, 97, 0, 97, 0, 98, 0, 98, 0, 98, 0, 99, 0, 99, 0, 99, 0}

	if !bytes.Equal(b, o) {
		t.Fatal(errUnexpectedBytes)
	}
}
