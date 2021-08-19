package internal

import "testing"

func TestStatusString(t *testing.T) {
	test := func(t testing.TB, status Status, expected string) {
		t.Helper()

		actual := status.String()

		if actual != expected {
			t.Errorf("expected %q, got %q", expected, actual)
		}
	}

	test(t, StatusQueued, "queued")
	test(t, StatusWorking, "working")
	test(t, StatusSuccess, "success")
	test(t, StatusFailed, "failed")
	test(t, StatusCancelling, "cancelling")
	test(t, StatusCancelled, "cancelled")
}
