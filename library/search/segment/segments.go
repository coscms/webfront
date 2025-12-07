/*

   Copyright 2016 Wenhui Shen <www.webx.top>

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

package segment

import (
	"github.com/admpub/log"
	"github.com/webx-top/com"
)

var (
	segments           = com.NewSafeMap[string, func() Segment]()
	defaultNop Segment = &nopSegment{}
)

// Register adds a new segment type with the given name and constructor function.
func Register(name string, c func() Segment) {
	segments.Set(name, c)
}

// IsNop reports whether the given segment is a no-operation segment.
func IsNop(segment Segment) bool {
	return defaultNop == segment
}

// defaultNopSegment returns the default no-operation segment instance.
func defaultNopSegment() Segment {
	return defaultNop
}

// Get returns the segment engine with the given name.
// If the specified segment engine is not found, it returns a default no-op segment
// and logs an error message.
func Get(name string) Segment {
	fn, ok := segments.GetOk(name)
	if !ok || fn == nil {
		log.Error(`[segment]Not found engine:`, name)
		fn = defaultNopSegment
	}
	return fn()
}

// Has checks if a segment with the given name exists in the segments collection.
func Has(name string) bool {
	_, ok := segments.GetOk(name)
	return ok
}

// Unregister removes the segment with the specified name from the registry
func Unregister(name string) {
	segments.Delete(name)
}
