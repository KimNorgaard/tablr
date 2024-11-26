package tablr_test

import (
	"testing"

	"github.com/KimNorgaard/tablr"
)

func TestAlignmentConstants(t *testing.T) {
	if tablr.AlignDefault != 0 {
		t.Error("AlignDefault should be 0")
	}
	if tablr.AlignLeft != 1 {
		t.Error("AlignLeft should be 1")
	}
	if tablr.AlignCenter != 2 {
		t.Error("AlignCenter should be 2")
	}
	if tablr.AlignRight != 3 {
		t.Error("AlignRight should be 3")
	}
}

func TestAlignmentIsValid(t *testing.T) {
	if !tablr.AlignDefault.IsValid() {
		t.Error("AlignDefault should be valid")
	}
	if !tablr.AlignLeft.IsValid() {
		t.Error("AlignLeft should be valid")
	}
	if !tablr.AlignCenter.IsValid() {
		t.Error("AlignCenter should be valid")
	}
	if !tablr.AlignRight.IsValid() {
		t.Error("AlignRight should be valid")
	}
	if tablr.Alignment(4).IsValid() {
		t.Error("Alignment(4) should not be valid")
	}
}
