// nolint
package requestid

//

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	requestId := "123456abc"
	ctx := context.Background()
	nowtime := time.Now()

	want := ""
	got := String(ctx)
	if want != got {
		t.Fatalf("want != got, want = %s , got = %s", want, got)
	}

	want = ""
	got = String(nil)
	if want != got {
		t.Fatalf("want != got, want = %s , got = %s", want, got)
	}

	got2 := CreateTime(ctx)
	if !reflect.DeepEqual(got2, time.Time{}) {
		t.Fatalf("got2 not nil")
	}

	got3 := GetID(nil)
	want3 := ""
	if got3 != want3 {
		t.Fatalf("want3 != got3, want3 = %s , got3 = %s", want3, got3)
	}

	got3 = GetID(ctx)
	want3 = ""
	if got3 != want3 {
		t.Fatalf("want3 != got3, want3 = %s , got3 = %s", want3, got3)
	}

	got4 := Cost(ctx)
	if got4 != 0 {
		t.Fatalf("got4 not nil")
	}

	ctx = NewContext(ctx, requestId, nowtime)
	got5 := String(ctx)
	if got5 == "" {
		t.Fatalf("got5 is nil")
	}

	got6 := CreateTime(ctx)
	want6 := nowtime
	if want6 != got6 {
		t.Fatalf("want6 != got6, want6 = %s , got6 = %s", want6, got6)
	}

	got7 := GetID(ctx)
	want7 := "123456abc"
	if got7 != want7 {
		t.Fatalf("want7 != got7, want7 = %s , got7 = %s", want7, got7)
	}

	ctx = WithContext(ctx)
	got8 := String(ctx)
	if got8 == "" {
		t.Fatalf("got8 is nil")
	}

	got9 := Cost(ctx)
	if got9 == 0 {
		t.Fatalf("got9 is nil")
	}
}
