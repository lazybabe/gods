package mutex_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gmeasure"

	"github.com/qinyuguang/gods/internal/mutex"
)

func TestMutex(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mutex Suite")
}

var _ = Describe("Mutex", func() {
	It("Benchmark", Serial, func() {
		safeLock := mutex.New(true)
		unsafeLock := mutex.New(false)

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
		}, gmeasure.SamplingConfig{N: 1e7, Duration: 10 * time.Second})

		safeLockMedian := experiment.GetStats("SafeLock").DurationFor(gmeasure.StatMedian)
		unsafeLockMedian := experiment.GetStats("UnsafeLock").DurationFor(gmeasure.StatMedian)
		Expect(safeLockMedian).To(BeNumerically("<", 200*time.Nanosecond))
		Expect(unsafeLockMedian).To(BeNumerically("<", 200*time.Nanosecond))
	})
})
