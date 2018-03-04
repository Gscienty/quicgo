package utils

import (
	"bytes"
	"testing"
	"fmt"
	"os"
)

func TestWriteUInt (t *testing.T) {
	var b bytes.Buffer
	
	BigEndian.WriteUInt(&b, 0x010203, 3)

	fmt.Println (os.Stdout, b.Bytes())
}