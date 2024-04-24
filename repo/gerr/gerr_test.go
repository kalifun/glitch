package gerr

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		key    string
		code   string
		format string
	}
	tests := []struct {
		name string
		args args
		want *GError
	}{
		{
			name: "demo1",
			args: args{
				key:    "MissingParameterF",
				code:   "MissingParameter",
				format: "Missing Parameter: %s",
			},
			want: &GError{
				key:    "MissingParameterF",
				code:   "MissingParameter",
				format: "Missing Parameter: %s",
				args:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.key, tt.args.code, tt.args.format); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGError_Error(t *testing.T) {
	type fields struct {
		key    string
		code   string
		format string
		args   []interface{}
	}
	d := &GError{
		key:    "MissingParameterF",
		code:   "MissingParameter",
		format: "Missing Parameter: %s",
		args:   nil,
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "don't args",
			fields: fields{
				key:    d.key,
				code:   d.code,
				format: d.format,
				args:   nil,
			},
			want: "Code:MissingParameter, Message: Missing Parameter: %!s(MISSING)",
		},
		{
			name: "add args",
			fields: fields{
				key:    d.key,
				code:   d.code,
				format: d.format,
				args:   []interface{}{"User"},
			},
			want: "Code:MissingParameter, Message: Missing Parameter: User",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GError{
				key:    tt.fields.key,
				code:   tt.fields.code,
				format: tt.fields.format,
				args:   tt.fields.args,
			}
			if got := g.Error(); got != tt.want {
				t.Errorf("GError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGError_Args(t *testing.T) {
	type fields struct {
		key    string
		code   string
		format string
		args   []interface{}
	}
	type args struct {
		args []interface{}
	}
	d := &GError{
		key:    "MissingParameterF",
		code:   "MissingParameter",
		format: "Missing Parameter: %s",
		args:   []interface{}{"User"},
	}

	d2 := &GError{
		key:    "MissingParameterF",
		code:   "MissingParameter",
		format: "Missing Parameter: %s",
		args:   []interface{}{"User", "Name"},
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *GError
	}{
		{
			name: "test args",
			fields: fields{
				key:    d.key,
				code:   d.code,
				format: d.format,
				args:   nil,
			},
			args: args{
				args: []interface{}{"User"},
			},
			want: d,
		},
		{
			name: "mult args",
			fields: fields{
				key:    d2.key,
				code:   d2.code,
				format: d2.format,
				args:   nil,
			},
			args: args{
				args: []interface{}{"User", "Name"},
			},
			want: d2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GError{
				key:    tt.fields.key,
				code:   tt.fields.code,
				format: tt.fields.format,
				args:   tt.fields.args,
			}
			if got := g.Args(tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GError.Args() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDescribeErrorMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string]ErrorCode
	}{
		{
			name: "nil",
			want: make(map[string]ErrorCode),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DescribeErrorMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DescribeErrorMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
