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

func TestReadUInt (t *testing.T) {
	b := bytes.NewReader([] byte { 0x01, 0x02, 0x03, 0x04 })
	v, _ := BigEndian.ReadUInt(b, 4)
	fmt.Println (os.Stdout, v)
}