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
