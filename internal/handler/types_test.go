package handler

import (
	"reflect"
	"testing"
)

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
							mode: executed,
						},
						{
							mode: notExecuted,
						},
						{
							mode: notMatched,
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
							mode: notMatched,
						},
						{
							mode: notMatched,
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

func TestTarget_Done(t1 *testing.T) {
	tests := []struct {
		name string
		tags []*TargetTag
		want bool
	}{
		{
			name: "Not Done",
			tags: []*TargetTag{
				{
					mode: notMatched,
				},
				{
					mode: executed,
				},
			},
			want: false,
		},
		{
			name: "Done",
			tags: []*TargetTag{
				{
					mode: notExecuted,
				},
				{
					mode: executed,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Target{
				Tags: tt.tags,
			}
			if got := t.Done(); got != tt.want {
				t1.Errorf("Done() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTargets_Merge(t *testing.T) {
	type args struct {
		ts1 *Targets
	}
	tests := []struct {
		name string
		ts   Targets
		args args
		want *Targets
	}{
		{
			name: "test",
			ts: Targets{
				{
					Image: "library/php",
					Tags: []*TargetTag{
						{
							Pattern: "pattern1",
						},
						{
							Pattern: "pattern3",
						},
					},
				},
			},
			args: args{
				ts1: &Targets{
					{
						Image: "library/php",
						Tags: []*TargetTag{
							{
								Tag:     "lock1",
								Pattern: "pattern1",
							},
							{
								Tag:     "lock2",
								Pattern: "pattern2",
							},
						},
					},
				},
			},
			want: &Targets{
				{
					Image: "library/php",
					Tags: []*TargetTag{
						{
							Tag:     "lock1",
							Pattern: "pattern1",
						},
						{
							Pattern: "pattern3",
						},
					},
				},
			},
		},
		{
			name: "merge empty",
			ts: Targets{
				{
					Image: "library/php",
					Tags: []*TargetTag{
						{
							Tag:     "lock1",
							Pattern: "pattern1",
						},
					},
				},
			},
			args: args{
				ts1: &Targets{
				},
			},
			want: &Targets{
				{
					Image: "library/php",
					Tags: []*TargetTag{
						{
							Tag:     "lock1",
							Pattern: "pattern1",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ts.Merge(tt.args.ts1); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}
