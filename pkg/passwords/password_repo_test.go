/*
 * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package passwords

import (
	"github.com/ThilinaManamgoda/password-manager/pkg/utils"
	"gotest.tools/assert"
	"os"
	"path"
	"testing"
)

var repo *Repository

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	db := &DB{
		Entries: map[string]Entry{},
		Labels:  map[string][]string{},
	}
	repo = newRepository(db, "mPassword", utils.AESEncryptID, path.Join(wd, "testDB"))
	err = repo.ImportFromCSV(path.Join(wd, "../../test/mock-data/data.csv"))
	if err != nil {
		panic(err)
	}
}

func TestGet(t *testing.T) {
	entry, err := repo.GetPasswordEntry("bmcandie15@devhub.com")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Binny", entry.Username)
	assert.Equal(t, "Qa88ookYyY", entry.Password)

	_, err = repo.GetPasswordEntry("invalid@id.com")
	assert.Error(t, err, ErrInvalidID("invalid@id.com").Error())
}

func TestSearchID(t *testing.T) {
	ids, err := repo.SearchID("bluckcock", false)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "bluckcockro@answers.com", ids[0])

	_, err = repo.SearchID("invalid@id.com", false)
	assert.Error(t, err, ErrCannotFindMatchForID("invalid@id.com").Error())
}

func TestSearchLabel(t *testing.T) {
	ids, err := repo.SearchLabel("five", false)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 130, len(ids))
	assert.Equal(t, true, utils.StringSliceContains("agwyerb@wisc.edu", ids))
}

func TestAdd(t *testing.T) {
	err := repo.Add("test@test.com", "test", "password", []string{"l1", "l2"})
	if err != nil {
		t.Error(err)
	}
	entry, ok := repo.db.Entries["test@test.com"]
	assert.Equal(t, true, ok)
	assert.Equal(t, "test", entry.Username)
	assert.Equal(t, "password", entry.Password)

	ids, err := repo.SearchLabel("l1", false)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "test@test.com", ids[0])
	ids, err = repo.SearchLabel("l2", false)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "test@test.com", ids[0])
}

func TestRemove(t *testing.T) {
	err := repo.Remove("bluckcockro@answers.com")
	if err != nil {
		t.Error(err)
	}
	_, ok := repo.db.Entries["bluckcockro@answers.com"]
	assert.Equal(t, false, ok)
}

func TestChangePasswordEntry(t *testing.T) {
	err := repo.ChangePasswordEntry("lgaggd@purevolume.com", Entry{Username: "change1", Password: "change2"})
	if err != nil {
		t.Error(err)
	}
	entry, ok := repo.db.Entries["lgaggd@purevolume.com"]
	assert.Equal(t, true, ok)
	assert.Equal(t, "change1", entry.Username)
	assert.Equal(t, "change2", entry.Password)
}
