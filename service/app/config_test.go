package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	var testCases = []struct {
		name string
		c    Config
	}{
		{
			name: "all fields valid",
			c: Config{
				PostgresURL: "postgres://postgres@postgres/postgres",
			},
		},
		{
			name: "no postgres url",
			c: Config{
				PostgresURL: "",
			},
		},
	}

	for _, tc := range testCases {
		err := tc.c.Validate()
		switch tc.name {
		case "all fields valid":
			assert.NoError(t, err)
		case "no postgres url":
			assert.Error(t, err)
		}
	}
}
