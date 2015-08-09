package main

import (
	"fmt"
	"testing"
)

type tcase struct {
	apkFile,
	certFile,
	vendor,
	certType string
	result bool
}

func TestCheckCert(t *testing.T) {
	cases := []struct {
		in tcase
	}{
		{tcase{"test/app-debug.apk", "/dev/null", "android", "release", false}},
		{tcase{"/dev/null", "test/android.cert", "", "", false}},
		{tcase{"/dev/random", "test/android.cert", "", "", false}},
		{tcase{"test/app-debug.apk", "/dev/null", "", "", true}},
		{tcase{"test/app-debug.apk", "test/android.cert", "android", "release", false}},
		{tcase{"test/app-mediakey.apk", "", "android", "media", true}},
		{tcase{"test/app-sharedkey.apk", "", "android", "shared", true}},
		{tcase{"test/app-platformkey.apk", "", "android", "platform", true}},
		{tcase{"test/app-releasekey.apk", "", "android", "release", true}},
		{tcase{"test/app-unsigned.apk", "", "android", "release", false}},
	}
	for i, c := range cases {
		fmt.Printf("Testing #%d case\n", i)
		ret := Checkcert(c.in.apkFile, c.in.certFile, c.in.vendor, c.in.certType)
		if ret != c.in.result {
			t.Errorf("result incorrect in #%d test cases", i)
		}
	}
}
