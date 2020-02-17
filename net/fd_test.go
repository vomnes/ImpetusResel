package net

import (
	"fmt"
	"syscall"
	"testing"
)

func checkDiff(actual, expected [32]int32) string {
	strActual := fmt.Sprintf("%x", actual)
	strExpected := fmt.Sprintf("%x", expected)
	if strActual != strExpected {
		return fmt.Sprintf("\nExpect [% s]\nHas    [% s]\n", strExpected, strActual)
	}
	return ""
}

func TestFD(t *testing.T) {
	var activeFdSet syscall.FdSet
	FDZero(&activeFdSet)
	FDSet(3, &activeFdSet)
	expected := [32]int32{
		int32(0x8), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0),
	}
	if err := checkDiff(activeFdSet.Bits, expected); err != "" {
		t.Errorf(err)
	}
	FDSet(72, &activeFdSet)
	expected = [32]int32{
		int32(0x8), int32(0), int32(0x100), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0),
	}
	if err := checkDiff(activeFdSet.Bits, expected); err != "" {
		t.Errorf(err)
	}
	actualBool := FDIsSet(72, &activeFdSet)
	if actualBool == false {
		t.Errorf("Expect true has false")
	}
	actualBool = FDIsSet(90, &activeFdSet)
	if actualBool == true {
		t.Errorf("Expect false has true")
	}
	actualBool = FDIsSet(3, &activeFdSet)
	if actualBool == false {
		t.Errorf("Expect true has false")
	}
	FDClr(72, &activeFdSet)
	expected = [32]int32{
		int32(0x8), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0),
	}
	if err := checkDiff(activeFdSet.Bits, expected); err != "" {
		t.Errorf(err)
	}
	FDClr(255, &activeFdSet)
	expected = [32]int32{
		int32(0x8), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0), int32(0),
	}
	if err := checkDiff(activeFdSet.Bits, expected); err != "" {
		t.Errorf(err)
	}
}
