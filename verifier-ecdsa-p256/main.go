package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

const (
	challenge           = "authchallenge000"
	pk64                = "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEPTpHDMruXAF1oQiFgZofG5PFcpDzbaQXqrfV2WCpZ425U3H37QvSt9VFv6UHiKRBqv9n7gYnKH1hqIfrc7wNYg=="
	authenticatorData64 = "SZYN5YgOjGh0NBcPZHZgW4/krrmihjLHmVzzuoMdl2MFAAAAAA=="
	clientDataJSON64    = "eyJ0eXBlIjoid2ViYXV0aG4uZ2V0IiwiY2hhbGxlbmdlIjoiWVhWMGFHTm9ZV3hzWlc1blpUQXdNQSIsIm9yaWdpbiI6Imh0dHA6Ly9sb2NhbGhvc3Q6NDUwNyIsImNyb3NzT3JpZ2luIjpmYWxzZX0="
	signature64         = "MEUCIHg6wL+3EQhPK4fy+FwqQtLUrq/Wa5Ntpf9WrquyDRJgAiEAogTRbuYfqknO8AkdRK8X0iHjDbRC4urBEwMPMn3Oj6A="
)

func main() {
	fmt.Println("Expected Challenge:", challenge)

	msg := []byte{}
	authenticatorData, err := base64.StdEncoding.DecodeString(authenticatorData64)
	if err != nil {
		panic(err)
	}
	msg = append(msg, authenticatorData...)
	clientDataJSON, err := base64.StdEncoding.DecodeString(clientDataJSON64)
	if err != nil {
		panic(err)
	}
	clientData := struct {
		Challenge string `json:"challenge"`
	}{}
	err = json.Unmarshal(clientDataJSON, &clientData)
	if err != nil {
		panic(err)
	}
	challengeDecoded, err := base64.RawURLEncoding.DecodeString(clientData.Challenge)
	if err != nil {
		panic(err)
	}
	fmt.Println("Challenge:", string(challengeDecoded))
	clientDataHash := sha256.Sum256(clientDataJSON)
	msg = append(msg, clientDataHash[:]...)

	pkBytes, err := base64.StdEncoding.DecodeString(pk64)
	if err != nil {
		panic(err)
	}
	pkUntyped, err := x509.ParsePKIXPublicKey(pkBytes)
	if err != nil {
		panic(err)
	}
	pk := pkUntyped.(*ecdsa.PublicKey)

	hash := sha256.Sum256(msg)

	sig, err := base64.StdEncoding.DecodeString(signature64)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sig len:", len(sig))
	verified := ecdsa.VerifyASN1(pk, hash[:], sig)
	fmt.Println(verified)
}
