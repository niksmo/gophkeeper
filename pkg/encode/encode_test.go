package encode_test

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/niksmo/gophkeeper/pkg/encode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	type obj struct {
		Foo string
		Bar int
		Baz []int
	}

	t.Run("Encoder", func(t *testing.T) {
		encObj := obj{"test", 12345, []int{0, 1, 2, 3, 4, 5}}
		encoder := encode.NewEncoder()
		oBytes, err := encoder.Encode(encObj)
		require.NoError(t, err)
		require.NotNil(t, oBytes)
		require.NotEmpty(t, oBytes)

		var decObj obj
		b := bytes.NewReader(oBytes)
		err = gob.NewDecoder(b).Decode(&decObj)
		require.NoError(t, err)
		assert.Equal(t, encObj, decObj)
	})

	t.Run("Decoder", func(t *testing.T) {
		encObj := obj{"test", 12345, []int{0, 1, 2, 3, 4, 5}}
		w := new(bytes.Buffer)
		err := gob.NewEncoder(w).Encode(encObj)
		require.NoError(t, err)
		require.NotEmpty(t, w.Bytes())

		var decObj obj
		decoder := encode.NewDecoder()
		err = decoder.Decode(&decObj, w.Bytes())
		require.NoError(t, err)
		assert.Equal(t, encObj, decObj)
	})
}
