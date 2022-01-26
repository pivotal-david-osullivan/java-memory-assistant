/*
 * Copyright 2018-2020, VMware, Inc. All Rights Reserved.
 * Proprietary and Confidential.
 * Unauthorized use, copying or distribution of this source code via any medium is
 * strictly prohibited without the express written consent of VMware, Inc.
 */

package helper

import (
	"fmt"
	"os"
	"strings"

	"github.com/paketo-buildpacks/libpak/sherpa"

	"github.com/paketo-buildpacks/libpak/bard"
)

type Properties struct {
	Logger bard.Logger
}

func (p Properties) Execute() (map[string]string, error) {
	p.Logger.Info("Configuring Java Memory Assistant properties")

	var argList string

	if argList = sherpa.GetEnvWithDefault("BPL_JMA_ARGS", ""); argList == "" {
		argList = fmt.Sprintf("-Djma.check_interval=5s -Djma.max_frequency=1/1m -Djma.heap_dump_folder=%s -Djma.thresholds.heap=80%", os.TempDir())
	} else {
		var runtimeArgs []string
		for arg, _ := range strings.Split(argList, ",") {
			runtimeArgs = append(runtimeArgs, fmt.Sprintf("-Djma.%s", arg))
		}
		argList = strings.Join(runtimeArgs, " ")
	}
	p.Logger.Infof("Enabling Java Memory Assistant with args: %s", argList)

	opts := sherpa.AppendToEnvVar("JAVA_TOOL_OPTIONS", " ", argList)

	return map[string]string{"JAVA_TOOL_OPTIONS": opts}, nil
}
