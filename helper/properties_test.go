/*
 * Copyright 2018-2020, VMware, Inc. All Rights Reserved.
 * Proprietary and Confidential.
 * Unauthorized use, copying or distribution of this source code via any medium is
 * strictly prohibited without the express written consent of VMware, Inc.
 */

package helper_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pivotal-david-osullivan/java-memory-assistant/helper"
	"github.com/sclevine/spec"
)

func testProperties(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		p helper.Properties
	)

	it("contributes base JMA configuration", func() {
		Expect(p.Execute()).To(Equal(map[string]string{
			"JAVA_TOOL_OPTIONS": fmt.Sprintf("-Djma.check_interval=5s -Djma.max_frequency=1/1m -Djma.heap_dump_folder=%s", os.TempDir()),
		}))
	})

	context("$BPL_JMA_ARGS is set", func() {
		it("contributes all arguments to JMA configuration", func() {
			Expect(os.Setenv("BPL_JMA_ARGS", "check_interval=10s,max_frequency=1/1m,heap_dump_folder=/tmp/,thresholds.heap=80%,log_level=DEBUG")).To(Succeed())
			Expect(p.Execute()).To(Equal(map[string]string{
				"JAVA_TOOL_OPTIONS": "-Djma.check_interval=10s -Djma.max_frequency=1/1m -Djma.heap_dump_folder=/tmp/ -Djma.thresholds.heap=80% -Djma.log_level=DEBUG"}))
		})
	})

	context("$JAVA_TOOL_OPTIONS", func() {
		it.Before(func() {
			Expect(os.Setenv("JAVA_TOOL_OPTIONS", "test-java-tool-options")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("JAVA_TOOL_OPTIONS")).To(Succeed())
		})

		it("contributes configuration appended to existing $JAVA_TOOL_OPTIONS", func() {
			Expect(os.Setenv("BPL_JMA_ARGS", "check_interval=10s,thresholds.heap=80%")).To(Succeed())
			Expect(p.Execute()).To(Equal(map[string]string{
				"JAVA_TOOL_OPTIONS": "test-java-tool-options -Djma.check_interval=10s -Djma.thresholds.heap=80%",
			}))
		})
	})
}
