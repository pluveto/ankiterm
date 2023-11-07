package xmisc_test

import (
	"testing"

	"github.com/pluveto/ankiterm/x/xmisc"
)

func TestPurgeStyle(t *testing.T){
	cases := []struct{
		in string
		want string
	}{
		{
			in: `<style type="text/css">xx</style>`,
			want: ``,
		},
	}
	for _, c := range cases {
		got := xmisc.PurgeStyle(c.in)
		if got != c.want {
			t.Errorf("PurgeStyle(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
