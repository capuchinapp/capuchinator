package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"capuchinator/internal/domain"
)

func Test_findPorts(t *testing.T) {
	type args struct {
		strategy domain.Strategy
		pathFile string
	}
	tests := []struct {
		name    string
		args    args
		wantAPI string
		wantUI  string
		wantErr error
	}{
		{
			name: "Should find ports - blue strategy",
			args: args{
				strategy: domain.StrategyBlue,
				pathFile: "testdata/ports/compose.blue.yaml",
			},
			wantAPI: "3001",
			wantUI:  "3002",
			wantErr: nil,
		},
		{
			name: "Should find ports - green strategy",
			args: args{
				strategy: domain.StrategyGreen,
				pathFile: "testdata/ports/compose.green.yaml",
			},
			wantAPI: "3011",
			wantUI:  "3012",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAPI, gotUI, err := findPorts(tt.args.strategy, tt.args.pathFile)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantAPI, gotAPI)
			assert.Equal(t, tt.wantUI, gotUI)
		})
	}
}
