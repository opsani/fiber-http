// Copyright 2020 Opsani
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoot(t *testing.T) {
	app := newApp()
	req := httptest.NewRequest("GET", "http://localhost:8080", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status code: 200, got: %v", resp.StatusCode)
	}

	assertResponseBodyEquals(t, resp, "move along, nothing to see here")
}

func TestCPU(t *testing.T) {
	app := newApp()
	req := httptest.NewRequest("GET", "http://localhost:8080/cpu", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status code: 200, got: %v", resp.StatusCode)
	}
	assertResponseBodyEquals(t, resp, "consumed CPU for 100ms\n")
}

func TestMemory(t *testing.T) {
	app := newApp()
	req := httptest.NewRequest("GET", "http://localhost:8080/memory", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status code: 200, got: %v", resp.StatusCode)
	}
	assertResponseBodyEquals(t, resp, "allocated 10.00MB (10485760 bytes) of memory\n")
}

func TestTime(t *testing.T) {
	app := newApp()
	req := httptest.NewRequest("GET", "http://localhost:8080/time", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status code: 200, got: %v", resp.StatusCode)
	}
	assertResponseBodyEquals(t, resp, "slept for 100ms\n")
}

func TestRequest(t *testing.T) {
	app := newApp()
	req := httptest.NewRequest("GET", "http://localhost:8080/request?url=http://opsani.com/", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status code: 200, got: %v", resp.StatusCode)
	}
}

func assertResponseBodyEquals(t *testing.T, resp *http.Response, expected string) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err.Error())
	}
	actual := string(body)
	if actual != expected {
		t.Errorf("expected %q but got %q", expected, actual)
	}
}
