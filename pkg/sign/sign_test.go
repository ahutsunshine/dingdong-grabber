/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at
  http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package sign

import (
	"testing"

	"github.com/dingdong-grabber/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	defer util.ClearSignConfigFile()
	s, err := NewDefaultJsSign()
	assert.NoError(t, err)
	got, err := s.Sign("6b7ed2c2051cb9e069d0e499a732ce9db8ce96fd", map[string]string{
		"cookie": "ede2c413e49d1a65566e12c27c819",
		"uid":    "647723be79400013ab",
	})
	assert.NoError(t, err)
	assert.Equal(t, "7bb06cf6f25140f1697035d47957d6e2", got["sign"])
}
