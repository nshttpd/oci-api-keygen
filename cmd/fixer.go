package cmd

import (
	"crypto/md5"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func fixFingerprint(cfg *Config) *Config {
	fmt.Println("fixing fingerprint issue")

	for i, a := range cfg.ApiKeys {
		f, err := ioutil.ReadFile(cfg.KeyPath + "/" + a.Tenancy + ".pub.pem")
		if err != nil {
			fmt.Printf("error reading pub key for tenancy : %s\n", a.Tenancy)
			fmt.Printf("error : %v\n", err)
		}
		b, _ := pem.Decode(f)
		md5sum := md5.Sum(b.Bytes)
		cfg.ApiKeys[i].Fingerprint = rfc4716hex(md5sum[:])
	}
	cfg.Version = 1
	cfg.SaveConfig()
	return cfg
}
