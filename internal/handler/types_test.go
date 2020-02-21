package handler

import "testing"

func TestTargets_IsUpdated(t *testing.T) {
	tests := []struct {
		name string
		t    Targets
		want bool
	}{
		{
			name: "updated",
			t: Targets{
				&Target{
					Tags: []*TargetTag{
						{
							Updated: false,
						},
						{
							Updated: true,
						},
						{
							Updated: false,
						},
					},
				},
			},
			want: true,
		},
		{
			name: "no updated",
			t: Targets{
				&Target{
					Tags: []*TargetTag{
						{
							Updated: false,
						},
						{
							Updated: false,
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.AreUpdated(); got != tt.want {
				t.Errorf("AreUpdated() = %v, want %v", got, tt.want)
			}
		})
	}
}
