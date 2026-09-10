package console

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestDecodeSecretToken(t *testing.T) {
	got := EncodeSecret("racoiaws")
	must.EqOp(t, "0YAKikQs", got)
}
