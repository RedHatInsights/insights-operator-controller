/*
Copyright © 2019, 2020 Red Hat, Inc.

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

package utils_test

import (
	"testing"

	"github.com/RedHatInsights/insights-operator-controller/utils"
)

// TestZeroValueForNil test the function ZeroValue against nil
func TestZeroValueForNil(t *testing.T) {
	r := utils.ZeroValue(nil)
	if !r {
		t.Fatal("nil should be treated as zero value")
	}
}

// TestZeroValueForString test the function ZeroValue against various strings
func TestZeroValueForString(t *testing.T) {
	r := utils.ZeroValue("")
	if !r {
		t.Fatal("empty string should be treated as zero value")
	}

	r2 := utils.ZeroValue("foobar")
	if r2 {
		t.Fatal("non-empty string should be treated as non-zero value")
	}
}

// TestZeroValueForIntegerNumber test the function ZeroValue against various integer numbers
func TestZeroValueForIntegerNumber(t *testing.T) {
	r := utils.ZeroValue(0)
	if !r {
		t.Fatal("0 should be treated as zero value")
	}

	r2 := utils.ZeroValue(42)
	if r2 {
		t.Fatal("non-zero number should be treated as non-zero value")
	}

	r3 := utils.ZeroValue(uint64(0))
	if !r3 {
		t.Fatal("0 should be treated as zero value")
	}

	r4 := utils.ZeroValue(uint64(42))
	if r4 {
		t.Fatal("non-zero number should be treated as non-zero value")
	}
}

// TestZeroValueForFloatNumber test the function ZeroValue against various floating point numbers
func TestZeroValueForFloatNumber(t *testing.T) {
	r := utils.ZeroValue(0.0)
	if !r {
		t.Fatal("0.0 should be treated as zero value")
	}

	r2 := utils.ZeroValue(0.1)
	if r2 {
		t.Fatal("non-zero number should be treated as non-zero value")
	}
}

// TestZeroValueForComplexNumber test the function ZeroValue against various complex numbers
func TestZeroValueForComplexNumber(t *testing.T) {
	r := utils.ZeroValue(complex(0, 0))
	if !r {
		t.Fatal("0.0+0.0i should be treated as zero value")
	}

	r2 := utils.ZeroValue(complex(0, 0.1))
	if r2 {
		t.Fatal("non-zero complex number should be treated as non-zero value")
	}

	r3 := utils.ZeroValue(complex(0.1, 0))
	if r3 {
		t.Fatal("non-zero complex number should be treated as non-zero value")
	}

	r4 := utils.ZeroValue(complex(0.1, 0.1))
	if r4 {
		t.Fatal("non-zero complex number should be treated as non-zero value")
	}
}

// TestZeroValueForSlice test the function ZeroValue against slices
func TestZeroValueForSlice(t *testing.T) {
	var emptySlice []string
	r := utils.ZeroValue(emptySlice)
	if !r {
		t.Fatal("empty slice should be treated as zero value")
	}

	var filledSlice []string = []string{"foo"}
	r2 := utils.ZeroValue(filledSlice)
	if r2 {
		t.Fatal("non-empty slice should be treated as non-zero value")
	}
}
