package java_memory_assistant

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {

	var (
		Expect = NewWithT(t).Expect
		build  Build
		ctx    libcnb.BuildContext
	)

	it("contributes Java Memory Assistant agent", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "java-memory-assistant"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "java-memory-assistant",
					"version": "1.0.0",
					"stacks":  []interface{}{"io.buildpacks.stacks.bionic"},
				},
			},
		}
		ctx.StackID = "io.buildpacks.stacks.bionic"

		result, err := build.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(1))
		Expect(result.Layers[0].Name()).To(Equal("java-memory-assistant"))
		Expect(result.BOM.Entries).To(HaveLen(0))
	})

}
