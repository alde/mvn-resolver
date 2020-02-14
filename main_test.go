package main

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Scan(t *testing.T) {
	f, err := os.Open("testdata/pom.xml")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	deps := scan(bufio.NewScanner(f))
	assert.Equal(t, len(deps), 2)
	assert.Equal(t, deps[0].Version, "1.7.29")
	assert.Equal(t, deps[1].Version, "1.3.61")
}
