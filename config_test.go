package main

import (
	"testing"

	"gotest.tools/assert"
)

func TestLoadConfig(t *testing.T) {
	conf, err := loadConfig("config-example.toml")
	assert.Assert(t, conf != nil)
	assert.Assert(t, err)
	assert.Assert(t, conf.Credentials.Email == "me@domain.com")
	assert.Assert(t, conf.Credentials.Password == "abc123xyz")
}
