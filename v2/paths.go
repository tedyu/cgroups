/*
   Copyright The containerd Authors.

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

package v2

import (
	"fmt"
	"path/filepath"
	"strings"
)

// GroupPath is a string that appears as the third field in /proc/PID/cgroup.
// e.g. "/user.slice/user-1001.slice/session-1.scope"
//
// GroupPath must not contain "/sys/fs/cgroup" prefix.
// GroupPath must be a absolute path starts with "/".
type GroupPath string

// NestedGroupPath will nest the cgroups based on the calling processes cgroup
// placing its child processes inside its own path
func NestedGroupPath(suffix string) (GroupPath, error) {
	path, err := parseCgroupFile("/proc/self/cgroup")
	if err != nil {
		return "", err
	}
	return GroupPath(filepath.Join(string(path), suffix)), nil
}

// PidGroupPath will return the correct cgroup paths for an existing process running inside a cgroup
// This is commonly used for the Load function to restore an existing container
func PidGroupPath(pid int) (GroupPath, error) {
	p := fmt.Sprintf("/proc/%d/cgroup", pid)
	return parseCgroupFile(p)
}

// VerifyGroupPath verifies the format of g.
// VerifyGroupPath doesn't verify whether g actually exists on the system.
func VerifyGroupPath(g GroupPath) error {
	s := string(g)
	if !strings.HasPrefix(s, "/") {
		return ErrInvalidGroupPath
	}
	if strings.HasPrefix(s, "/sys/fs/cgroup") {
		return ErrInvalidGroupPath
	}
	return nil
}
