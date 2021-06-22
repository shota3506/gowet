package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalVet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testcases := []*struct {
			src      string
			expected string
		}{
			{
				src: `# github.com/example.com/example
# github.com/example.com/example/example1
{}
`,
				expected: "{}",
			}, {
				src: `# github.com/example.com/example
# github.com/example.com/example/example1
{
	"github.com/example.com/example/example1": {
		"unreachable": [
			{
				"posn": "posn",
				"message": "message"
			}
		]
	}
}
`,
				expected: `{
	"github.com/example.com/example/example1": {
		"unreachable": [
			{
				"posn": "posn",
				"message": "message"
			}
		]
	}
}`,
			}, {
				src: `# github.com/example.com/example
# github.com/example.com/example/example1
{}
# github.com/example.com/example/example2
{
	"github.com/example.com/example/example2": {
		"unreachable": [
			{
				"posn": "posn",
				"message": "message"
			}
		]
	}
}
# github.com/example.com/example/example3
{
	"github.com/example.com/example/example3": {
		"unreachable": [
			{
				"posn": "posn",
				"message": "message"
			}
		]
	}
}
`,
				expected: `{
	"github.com/example.com/example/example2": {
		"unreachable": [
			{
				"posn": "posn",
				"message": "message"
			}
		]
	},
	"github.com/example.com/example/example3": {
		"unreachable": [
			{
				"posn": "posn",
				"message": "message"
			}
		]
	}
}`,
			},
		}

		for _, testcase := range testcases {
			actual, err := marshalVet([]byte(testcase.src))
			require.NoError(t, err)
			assert.Equal(t, testcase.expected, string(actual))
		}
	})

	t.Run("fail", func(t *testing.T) {
		testcases := []string{
			`{`,
		}

		for _, testcase := range testcases {
			_, err := marshalVet([]byte(testcase))
			require.Error(t, err)
		}
	})
}
