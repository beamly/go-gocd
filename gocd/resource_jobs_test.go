package gocd

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResourceJobsJSONMarshal(t *testing.T) {
	for _, test := range []struct {
		want    string
		timeout int
	}{
		{want: "5", timeout: 5},
	} {
		tf := TimeoutField(test.timeout)

		b, err := tf.MarshalJSON()

		assert.Equal(t, test.want, string(b))
		assert.NoError(t, err)
	}
}

func TestResourceJobsJSONUnmarshal(t *testing.T) {
	for _, test := range []struct {
		want    TimeoutField
		tf      TimeoutField
		timeout []byte
	}{
		// The `tf` attribute is present so that we have a value which does not match
		// the expected value and we can ensure the test is working and not just putting
		// the default in (which is `0`).
		{want: TimeoutField(0), tf: TimeoutField(-1), timeout: []byte(`"never"`)},
		{want: TimeoutField(0), tf: TimeoutField(-1), timeout: []byte(`"null"`)},
		{want: TimeoutField(1), tf: TimeoutField(-1), timeout: []byte("1")},
		{want: TimeoutField(3), tf: TimeoutField(-1), timeout: []byte("3")},
		{want: TimeoutField(10), tf: TimeoutField(-1), timeout: []byte("10")},
		{want: TimeoutField(23870), tf: TimeoutField(-1), timeout: []byte("23870")},
	} {

		err := json.Unmarshal(test.timeout, &test.tf)

		if assert.NoError(t, err) {
			assert.Equal(t, test.want, test.tf)
		}
	}
}
