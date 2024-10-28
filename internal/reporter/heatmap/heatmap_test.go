package heatmap

import (
	"bytes"
	"testing"

	"github.com/moisesvega/diffy/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	r := New()
	require.NotNil(t, r)
}

func TestReport(t *testing.T) {
	tests := []struct {
		desc string
		give []*model.User
		want bool
	}{
		{
			desc: "success",
			give: []*model.User{
				{
					Username: "moisesvega",
					Differentials: []*model.Differential{
						{
							Title:     "title",
							URI:       "uri",
							LineCount: 11,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			r := &reporter{}
			w := &bytes.Buffer{}
			err := r.Report(tt.give, model.WithWriter(w))
			require.NoError(t, err)
			assert.NotEmpty(t, w)
		})
	}
}
