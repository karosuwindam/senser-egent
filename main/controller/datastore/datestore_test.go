package datastore_test

import (
	"senseregent/controller/datastore"
	"testing"
	"time"
)

type a struct {
	str string
}

type b struct {
	v int
}

func TestDatastore(t *testing.T) {
	if v := datastore.Read(a{}); v != nil {
		t.Fatalf("bb")
	}

	if v := datastore.Read(b{}); v != nil {
		t.Fatalf("bb")
	}
	if ti := datastore.ReadTime(a{}); ti != (time.Time{}) {
		t.Fatalf("ta")
	}

	if ti := datastore.ReadTime(b{}); ti != (time.Time{}) {
		t.Fatalf("tb")
	}

	datastore.Add(a{}, a{"abc"})
	datastore.Add(b{}, b{10})

	if v := datastore.Read(a{}); v != nil {
		if tmp, ok := v.(a); !ok {
			t.Fatalf("cc")
		} else if tmp.str != "abc" {
			t.Fatalf("cc2")
		}
	}
	if ti := datastore.ReadTime(a{}); ti == (time.Time{}) {
		t.Fatalf("ta")
	}

	if v := datastore.Read(b{}); v != nil {
		if tmp, ok := v.(b); !ok {
			t.Fatalf("dd")
		} else if tmp.v != 10 {
			t.Fatalf("dd2")
		}
	}

	if ti := datastore.ReadTime(b{}); ti == (time.Time{}) {
		t.Fatalf("tb")
	}
}
