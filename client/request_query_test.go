package ekaweb_client_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/inaneverb/ekaweb/client"
)

type T struct {
	Name            string `form:"not_a_name" json:"not_a_name"`
	Position        string
	Empty           string `form:"empty,omitempty"`
	NameFromGoEmpty string `form:",omitempty"`
	Age             int
	unexported      bool `form:"should_not_be_added" json:"should_not_be_added"`
}

func emptyT() *T {
	return new(T)
}

func filledT() *T {
	return &T{
		Name:            "John Doe",
		Position:        "Demigod",
		Age:             322,
		NameFromGoEmpty: "Lorem Ipsum",
	}
}

func TestRequestQuery(t *testing.T) {

	var q = filledT()
	var wq = ekaweb_client.RequestQuery(q)

	var data, err = wq.Data()
	require.NoError(t, err)

	fmt.Printf("%s\n", data)
}

func BenchmarkRequestQuery(b *testing.B) {
	b.ReportAllocs()

	var q = filledT()
	var wq = ekaweb_client.RequestQuery(q)

	for i := 0; i < b.N; i++ {
		_, _ = wq.Data()
	}
}
