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

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	d "github.com/ci4rail/kyt/kyt-alm-server/internal/deployment"
	m "github.com/ci4rail/kyt/kyt-alm-server/internal/deployment/manifest"
	t "github.com/ci4rail/kyt/kyt-server-common/token"
	"github.com/gin-gonic/gin"
	"github.com/golangci/golangci-lint/pkg/sliceutil"
)

// ApplyPut -
func ApplyPut(c *gin.Context) {
	token, err := t.ReadToken(c.Request)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	tokenValid, claims, err := t.ValidateToken(token)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	if !sliceutil.Contains(claims, "Apply.write") {
		err = fmt.Errorf("Error: not allowed")
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}
	var tenantID string
	if tenantID = t.TenantIDFromToken(token); tenantID == "" {
		err = fmt.Errorf("Error: reading user ID `oid` from token")
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}
	// If token is not valid it means that either it has expired or the signature cannot be validated.
	// In both cases return `Unauthorized`.
	if !tokenValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		err = fmt.Errorf("Error: cannot read payload")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	manifest := m.CustomerManifest{}
	err = json.Unmarshal(jsonData, &manifest)
	if err != nil {
		err = fmt.Errorf("Error: cannot read customer manifest")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	_, err = d.CreateOrUpdateFromCustomerDeployment(tenantID, &manifest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, true)
}
