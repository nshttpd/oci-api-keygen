// Copyright © 2018 Steve Brunton <sbrunton@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package cmd

import (
	"crypto/rand"
	"crypto/rsa"

	"os"

	"crypto/x509"
	"encoding/pem"

	"crypto/md5"
	"fmt"

	"crypto"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	bitsize = 2048
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Args:  cobra.ExactArgs(1),
	Short: "Create API key for a tenancy",
	Long: `Create an API public and private key for accessing the OCI API for
a specific tenancy.

	oci-api-keygen create [tenancy]

the generated keys will be stored in the same directory as the config file
named as the tenancy supplied.`,
	Run: func(cmd *cobra.Command, args []string) {
		t := args[0]
		reader := rand.Reader
		if key, err := rsa.GenerateKey(reader, bitsize); err != nil {
			log.WithFields(log.Fields{"tenancy": t, "error": err}).Error("error generating key")
		} else {
			f := savePrivate(t+".pem", key)
			savePublic(t+".pub.pem", key.Public())
			a := &ApiKey{Tenancy: t, Fingerprint: f}
			cfg.ApiKeys = append(cfg.ApiKeys, *a)
		}

	},
}

func savePrivate(filename string, key *rsa.PrivateKey) string {
	f, err := os.Create(cfg.KeyPath + "/" + filename)
	if err != nil {
		log.WithFields(log.Fields{"filename": filename, "error": err}).Fatal("error creating private file")
	}
	defer f.Close()
	pk := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(f, pk)
	if err != nil {
		log.WithFields(log.Fields{"filename": filename, "error": err}).Fatal("error writing private file")
	}

	md5sum := md5.Sum(pk.Bytes)
	return rfc4716hex(md5sum[:])

}

func rfc4716hex(data []byte) string {
	var fingerprint string
	for i := 0; i < len(data); i++ {
		fingerprint = fmt.Sprintf("%s%0.2x", fingerprint, data[i])
		if i != len(data)-1 {
			fingerprint = fingerprint + ":"
		}
	}
	return fingerprint
}

func savePublic(filename string, key crypto.PublicKey) {

	//	asn1Bytes, err := asn1.Marshal(key)
	asn1Bytes, err := x509.MarshalPKIXPublicKey(key)

	if err != nil {
		log.WithFields(log.Fields{"filename": filename, "error": err}).Fatal("error marshalling public key")
	}

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	f, err := os.Create(cfg.KeyPath + "/" + filename)
	if err != nil {
		log.WithFields(log.Fields{"filename": filename, "error": err}).Fatal("error creating public file")
	}
	defer f.Close()

	err = pem.Encode(f, pemkey)
	if err != nil {
		log.WithFields(log.Fields{"filename": filename, "error": err}).Fatal("error writing public file")
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
}
