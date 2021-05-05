package errors

import "testing"

func TestFileNotFound_Error(t *testing.T) {
	tests := []struct {
		name   string
		fields FileNotFound
		want   string
	}{
		{
			name: "Success",
			fields: FileNotFound{
				FileName: "openapi.json",
				Path:     "go/src/github.com/zopsmart/gofr",
			},
			want: "File openapi.json not found at location go/src/github.com/zopsmart/gofr",
		},
	}
	for _, tt := range tests {
		f := FileNotFound{
			FileName: tt.fields.FileName,
			Path:     tt.fields.Path,
		}
		if got := f.Error(); got != tt.want {
			t.Errorf("Error() = %v, want %v", got, tt.want)
		}
	}
}
