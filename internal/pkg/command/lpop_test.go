package command

import (
	"testing"
	"time"

	"github.com/namreg/godown-v2/internal/pkg/storage"
	"github.com/namreg/godown-v2/internal/pkg/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestLpop_Name(t *testing.T) {
	cmd := new(Lpop)
	assert.Equal(t, "LPOP", cmd.Name())
}

func TestLpop_Help(t *testing.T) {
	cmd := new(Lpop)
	expected := `Usage: LPOP key
Removes and returns the first element of the list stored at key.`
	assert.Equal(t, expected, cmd.Help())
}

func TestLpop_Execute(t *testing.T) {
	expired := storage.NewListValue("val1", "val2")
	expired.SetTTL(time.Now().Add(-1 * time.Second))

	strg := memory.New(map[storage.Key]*storage.Value{
		"string":  storage.NewStringValue("string"),
		"list":    storage.NewListValue("val1", "val2"),
		"expired": expired,
	})

	tests := []struct {
		name string
		args []string
		want Result
	}{
		{"ok", []string{"list"}, StringResult{"val1"}},
		{"expired_key", []string{"expired"}, NilResult{}},
		{"not_existing_key", []string{"not_existing_key"}, NilResult{}},
		{"wrong_type_op", []string{"string"}, ErrResult{ErrWrongTypeOp}},
		{"wrong_args_number/1", []string{}, ErrResult{ErrWrongArgsNumber}},
		{"wrong_args_number/2", []string{"list", "0"}, ErrResult{ErrWrongArgsNumber}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := new(Lpop)
			res := cmd.Execute(strg, tt.args...)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestLpop_Execute_DelEmptyList(t *testing.T) {
	strg := memory.New(map[storage.Key]*storage.Value{
		"list": storage.NewListValue("val1"),
	})

	cmd := new(Lpop)
	_ = cmd.Execute(strg, "list")

	items, err := strg.All()
	assert.NoError(t, err)

	value, ok := items["list"]
	assert.Nil(t, value)
	assert.False(t, ok)
}