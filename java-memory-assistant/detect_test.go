package java_memory_assistant_test

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	java_memory_assistant "github.com/pivotal-david-osullivan/java-memory-assistant"
	"github.com/sclevine/spec"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx    libcnb.DetectContext
		detect java_memory_assistant.Detect
	)

	it("passes detection", func() {
		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: java_memory_assistant.PlanEntryAssistant},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: java_memory_assistant.PlanEntryAssistant},
						{Name: "jvm-application"},
					},
				},
			},
		}))
	})
}
