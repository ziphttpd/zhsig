package zhsig

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/xorvercom/util/pkg/json"
)

func TestCreateSiteFile(t *testing.T) {
	err := TempSpace(func(tempdir string) error {
		host := NewHost(tempdir, "abcd.com")
		err := CreateSiteFile(host)
		if err != nil {
			return err
		}
		if pkey, err := DecodePublicKeyFile(host.MyPublicKeyFile()); err == nil {
			var bytes []byte
			if elem, err := json.LoadFromJSONFile(host.MySiteFile()); err == nil {
				if oe, ok := elem.AsObject(); ok {
					if se, ok := oe.Child("site").AsString(); ok {
						bytes, err = base64.StdEncoding.DecodeString(se.Text())
						if err != nil {
							return err
						}
					} else {
						return fmt.Errorf("invalid site value")
					}
				} else {
					return fmt.Errorf("invalid json")
				}
			} else {
				return err
			}

			if false == VerifySiteFileBytes(pkey, host, bytes) {
				return fmt.Errorf("verify miss")
			}
		} else {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
