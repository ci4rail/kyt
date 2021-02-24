/*
Copyright © 2021 Ci4Rail GmbH

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

package deployment

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortByTimestamp(t *testing.T) {
	assert := assert.New(t)
	list := []string{
		"tenant_application_1",
		"tenant_application_5",
		"tenant_application_3",
		"tenant_application_4",
		"tenant_application_2",
	}
	sort.Sort(ByTimestamp(list))
	assert.Equal(list[0], "tenant_application_5")
	assert.Equal(list[1], "tenant_application_4")
	assert.Equal(list[2], "tenant_application_3")
	assert.Equal(list[3], "tenant_application_2")
	assert.Equal(list[4], "tenant_application_1")
}
