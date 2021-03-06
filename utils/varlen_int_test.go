package utils

import(
	"testing"
	"bytes"
	"fmt"
)

func TestParse(t *testing.T) {
	b1 := bytes.NewReader([]byte { 0xC2, 0x19, 0x7C, 0x5E, 0xFF, 0x14, 0xE8, 0x8C })
	varval, _ := VarLenIntegerStructParse(b1)
	if varval.val != 151288809941952652 {
		fmt.Println(varval.val)
		t.Fail()
	}
	b2 := bytes.NewReader([]byte { 0x9D, 0x7F, 0x3E, 0x7D })
	varval, _ = VarLenIntegerStructParse(b2)
	if varval.val != 494878333 {
		fmt.Println(varval.val)
		t.Fail()
	}
	b3 := bytes.NewReader([]byte { 0x7B, 0xBD })
	varval, _ = VarLenIntegerStructParse(b3)
	if varval.val != 15293 {
		fmt.Println(varval.val)
		t.Fail()
	}
	b4 := bytes.NewReader([]byte { 0x25 })
	varval, _ = VarLenIntegerStructParse(b4)
	if varval.val != 37 {
		fmt.Println(varval.val)
		t.Fail()
	}
}

func TestNew(t *testing.T) {
	varint := VarLenIntegerStructNew(VARLENINT_MAX_1BYTE_VALUE)
	if varint.len != 1 {
		fmt.Println(varint)
		t.Fail()
	}

	varint = VarLenIntegerStructNew(VARLENINT_MAX_1BYTE_VALUE + 1)
	if varint.len != 2 {
		fmt.Println(varint)
		t.Fail()
	}

	varint = VarLenIntegerStructNew(VARLENINT_MAX_8BYTE_VALUE + 1)
	if varint.len != 0 {
		t.Fail()
	}
}
