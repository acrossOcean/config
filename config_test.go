package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/acrossOcean/config/core"
)

func initTestEnv() {
	AddPath("./testFile/conf.ini")
}

func TestString(t *testing.T) {
	initTestEnv()

	SetReadOrder(core.ReadFromFile, core.ReadFromEnv)

	runMode, exist := String("runMode")
	assert.Equal(t, runMode, "dev")
	assert.Equal(t, exist, true)

	dbName, exist := String("mysql>>dbName")
	assert.Equal(t, dbName, "amazing_salted_fish")
	assert.Equal(t, exist, true)
}
