// +build tests

package signature

import "testing"

func TestExtractKey(t *testing.T) {
	extractKey("hello world//1/DOT///password")

	extractKey("hello world//Alice")
}
