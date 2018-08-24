package gocd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApiVersion_Int(t *testing.T) {
	for _, tt := range []struct {
		name string
		av   APIVersion
		want int
	}{
		{
			name: "v1",
			av:   APIVersion0,
			want: 0,
		},
		{
			name: "v5",
			av:   APIVersion5,
			want: 5,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.av.Int()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAPIVersion_LessThan(t *testing.T) {
	type args struct {
		version APIVersion
	}
	for _, tt := range []struct {
		name string
		av   APIVersion
		args args
		want bool
	}{
		{
			name: "success",
			av:   APIVersion0,
			args: args{
				version: APIVersion1,
			},
			want: true,
		},
		{
			name: "fail",
			av:   APIVersion1,
			args: args{
				version: APIVersion0,
			},
			want: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.av.LessThan(tt.args.version)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAPIVersion_Equal(t *testing.T) {
	type args struct {
		version APIVersion
	}
	for _, tt := range []struct {
		name string
		av   APIVersion
		args args
		want bool
	}{
		{
			name: "equal",
			av:   APIVersion1,
			args: args{
				version: APIVersion1,
			},
			want: true,
		},
		{
			name: "not-equal",
			av:   APIVersion1,
			args: args{
				version: APIVersion0,
			},
			want: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.av.Equal(tt.args.version)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAPIVersion_String(t *testing.T) {
	for _, tt := range []struct {
		name string
		av   APIVersion
		want string
	}{
		{
			name: "v1",
			av:   APIVersion1,
			want: "application/vnd.go.cd.v1+json",
		},
		{
			name: "v4",
			av:   APIVersion4,
			want: "application/vnd.go.cd.v4+json",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.av.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewAPIVersionRange(t *testing.T) {
	type args struct {
		min APIVersion
		max APIVersion
	}
	for _, tt := range []struct {
		name string
		args args
		want *APIVersionRange
	}{
		{
			name: "basic",
			args: args{
				min: APIVersion0,
				max: APIVersion1,
			},
			want: &APIVersionRange{
				min: APIVersion0,
				max: APIVersion1,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAPIVersionRange(tt.args.min, tt.args.max)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAPIVersionRange_Contains(t *testing.T) {
	type fields struct {
		min APIVersion
		max APIVersion
	}
	type args struct {
		version APIVersion
	}
	for _, tt := range []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "success",
			fields: fields{
				min: APIVersion0,
				max: APIVersion3,
			},
			args: args{
				version: APIVersion2,
			},
			want: true,
		},
		{
			name: "fail",
			fields: fields{
				min: APIVersion0,
				max: APIVersion3,
			},
			args: args{
				version: APIVersion5,
			},
			want: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			avr := &APIVersionRange{
				min: tt.fields.min,
				max: tt.fields.max,
			}
			got := avr.Contains(tt.args.version)
			assert.Equal(t, tt.want, got)
		})
	}
}
