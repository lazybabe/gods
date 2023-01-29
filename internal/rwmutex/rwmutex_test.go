package rwmutex_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gmeasure"

	"github.com/qinyuguang/gods/internal/rwmutex"
)

func TestRwMutex(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mutex Suite")
}

var _ = Describe("RwMutex", func() {
	It("Benchmark", Serial, func() {
		safeLock := rwmutex.New(true)
		unsafeLock := rwmutex.New(false)

		experiment := gmeasure.NewExperiment("LockUnlock")
		AddReportEntry(experiment.Name, experiment)
		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("SafeLock", func() {
				safeLock.Lock()
				defer safeLock.Unlock()
			})
			experiment.MeasureDuration("UnsafeLock", func() {
				unsafeLock.Lock()
				defer unsafeLock.Unlock()
			})
			experiment.MeasureDuration("SafeRLock", func() {
				safeLock.RLock()
				defer safeLock.RUnlock()
			})
			experiment.MeasureDuration("UnsafeRLock", func() {
				unsafeLock.RLock()
				defer unsafeLock.RUnlock()
			})
		}, gmeasure.SamplingConfig{N: 1e7, Duration: 10 * time.Second})

		safeLockMedian := experiment.GetStats("SafeLock").DurationFor(gmeasure.StatMedian)
		unsafeLockMedian := experiment.GetStats("UnsafeLock").DurationFor(gmeasure.StatMedian)
		safeRLockMedian := experiment.GetStats("SafeRLock").DurationFor(gmeasure.StatMedian)
		unsafeRLockMedian := experiment.GetStats("UnsafeRLock").DurationFor(gmeasure.StatMedian)
		Expect(safeLockMedian).To(BeNumerically("<", 200*time.Nanosecond))
		Expect(unsafeLockMedian).To(BeNumerically("<", 200*time.Nanosecond))
		Expect(safeRLockMedian).To(BeNumerically("<", 200*time.Nanosecond))
		Expect(unsafeRLockMedian).To(BeNumerically("<", 200*time.Nanosecond))
	})
})
