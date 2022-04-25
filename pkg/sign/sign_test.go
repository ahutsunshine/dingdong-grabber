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

	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	expected := map[string]string{
		"nars": "bef26b7dec708ba104e2e31d183442a7",
		"sesi": "j5nZZoD50c8c1559bb2bd2a5e0cff487f3a8b78",
	}

	got, err := NewSign("../../sign.js").Sign(map[string]string{
		"cookie": "ede2c413e49d1a65566e12c27c819",
		"uid":    "647723be79400013ab",
	})
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}
