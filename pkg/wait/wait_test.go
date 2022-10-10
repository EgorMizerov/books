package wait

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWaitReturnsTrueIfSuccessImmediately(t *testing.T) {
	testAlwaysOkFunction := func() bool { return true }
	startTime := time.Now()
	require.True(t, WaitUntil(testAlwaysOkFunction, time.Now().Add(time.Hour), time.Second))
	verifyCompletedImmediately(t, startTime)
}

func TestWaitReturnsTrueIfSuccessBeforeTimeout(t *testing.T) {
	testAttempts := 4
	testInterval := 30 * time.Millisecond
	testTotalExpectedWaitTime := testInterval * time.Duration(testAttempts-1)
	testOkAfterAttemptsFunction := func() bool {
		testAttempts--
		return testAttempts < 1
	}
	startTime := time.Now()
	require.True(t, WaitUntil(testOkAfterAttemptsFunction, startTime.Add(time.Hour), testInterval))
	verifyWaitTooksTime(t, startTime, testTotalExpectedWaitTime)
}

func TestWaitReturnFailsAfterTimeout(t *testing.T) {
	testNeverOkFunction := func() bool { return false }
	testExpectedWaitTime := 100 * time.Millisecond
	startTime := time.Now()
	require.False(t, WaitUntil(testNeverOkFunction, startTime.Add(testExpectedWaitTime), 10*time.Millisecond))
	verifyWaitTooksTime(t, startTime, testExpectedWaitTime)
}

func verifyCompletedImmediately(t *testing.T, startTime time.Time) {
	require.True(t, time.Now().Sub(startTime) < (100*time.Millisecond))
}

func verifyWaitTooksTime(t *testing.T, startTime time.Time, expectedWaitTime time.Duration) {
	spentTime := time.Now().Sub(startTime)
	require.True(t, spentTime > expectedWaitTime)
	require.True(t, spentTime < 2*expectedWaitTime)
}
