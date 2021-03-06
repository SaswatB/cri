// +build !windows

/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"fmt"
	"regexp"

	"github.com/opencontainers/selinux/go-selinux"
	"github.com/opencontainers/selinux/go-selinux/label"
	"github.com/pkg/errors"
	runtime "k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"
)

func initSelinuxOpts(selinuxOpt *runtime.SELinuxOption) (string, string, error) {
	if selinuxOpt == nil {
		return "", "", nil
	}

	// Should ignored selinuxOpts if they are incomplete.
	if selinuxOpt.GetUser() == "" ||
		selinuxOpt.GetRole() == "" ||
		selinuxOpt.GetType() == "" {
		return "", "", nil
	}

	// make sure the format of "level" is correct.
	ok, err := checkSelinuxLevel(selinuxOpt.GetLevel())
	if err != nil || !ok {
		return "", "", err
	}

	labelOpts := fmt.Sprintf("%s:%s:%s:%s",
		selinuxOpt.GetUser(),
		selinuxOpt.GetRole(),
		selinuxOpt.GetType(),
		selinuxOpt.GetLevel())
	return label.InitLabels(selinux.DupSecOpt(labelOpts))
}

func checkSelinuxLevel(level string) (bool, error) {
	if len(level) == 0 {
		return true, nil
	}

	matched, err := regexp.MatchString(`^s\d(-s\d)??(:c\d{1,4}((.c\d{1,4})?,c\d{1,4})*(.c\d{1,4})?(,c\d{1,4}(.c\d{1,4})?)*)?$`, level)
	if err != nil || !matched {
		return false, errors.Wrapf(err, "the format of 'level' %q is not correct", level)
	}
	return true, nil
}
