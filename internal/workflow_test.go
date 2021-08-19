package internal

import (
	"testing"
	"time"
)

func TestParseReferenceConfidenceMode(t *testing.T) {
	test := func(t testing.TB, s string, expected ReferenceConfidenceMode) {
		t.Helper()

		if actual, err := ParseReferenceConfidenceMode(s); err == nil {
			if actual != expected {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		} else {
			t.Errorf(`unexpected failure: s = %q`, s)
		}
	}

	test(t, "NONE", ReferenceConfidenceModeNone)
	test(t, "BP_RESOLUTION", ReferenceConfidenceModeBPResolution)
	test(t, "GVCF", ReferenceConfidenceModeGVCF)

	if _, err := ParseReferenceConfidenceMode(""); err == nil {
		t.Error(`expected failure: s = ""`)
	}

	if _, err := ParseReferenceConfidenceMode("msgenctl"); err == nil {
		t.Error(`expected failure: s = "msgenctl"`)
	}
}

func TestWorkflowDuration(t *testing.T) {
	expected := 2 * time.Hour

	createdDate := time.Now()
	endDate := createdDate.Add(expected)

	workflow := Workflow{
		CreatedDate: createdDate,
		EndDate:     &endDate,
	}

	actual := workflow.Duration()

	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
