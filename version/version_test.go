package version
import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBigVersion(t *testing.T) {
	cases := []struct {
		input string
		output	string
	}{
		{"1.0-18-g5359a56","1"},
		{"1.0-18.2","1"},
		{"22.0-18.2","22"},

	}

	for _,tc := range cases{
		require.Equal( t,tc.output,GetBigVersion(tc.input))
	}
}
