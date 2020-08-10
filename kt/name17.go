// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

// +build go1.7,!go1.8

package kt

import (
	"reflect"
	"testing"
	"unsafe"
)

// This method extracts the unexported 'name' field from a testing.T value.
//
// This is super evil. Don't try this at home.
//
// Adapted from http://stackoverflow.com/a/17982725/13860
func tName(t *testing.T) string {
	pv := reflect.ValueOf(t)
	v := reflect.Indirect(pv)
	name := v.FieldByName("name")
	namePtr := unsafe.Pointer(name.UnsafeAddr()) // nolint: gas
	realName := (*string)(namePtr)
	return *realName
}
