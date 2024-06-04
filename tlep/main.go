package tlep

import (
	"fmt"
)

func main() {
	tlep1, err := InitTLEP("user2")
	if err != nil {
		panic(err)
	}
	err = tlep1.CBESInitRandom()
	if err != nil {
		panic(err)
	}
	keyA, err := tlep1.ECDHGetPublicKey()
	if err != nil {
		panic(err)
	}
	tlep2, err := InitTLEP("user1")
	if err != nil {
		panic(err)
	}
	keyB, err := tlep2.ECDHGetPublicKey()
	if err != nil {
		panic(err)
	}
	// Acception
	err = tlep1.ECDHApplyOtherKeyBytes(keyB)
	if err != nil {
		panic(err)
	}
	err = tlep2.ECDHApplyOtherKeyBytes(keyA)
	if err != nil {
		panic(err)
	}
	// First EA encryption
	fmt.Printf("TLEP1 CBES specs: %#v\n", tlep1.CBES)
	cbesSpecs, err := tlep1.CBESGetBytes()
	if err != nil {
		panic(err)
	}
	encryptedCBESSpecs, err := tlep1.EncryptMessageEA(cbesSpecs)
	if err != nil {
		panic(err)
	}
	realSpecs, err := tlep2.DecryptMessageEA(encryptedCBESSpecs)
	if err != nil {
		panic(err)
	}
	tlep2.CBESSetFromBytes(realSpecs)
	fmt.Printf("TLEP2 CBES specs: %#v\n", tlep2.CBES)
	// Second layer CAFEA encryption
	someMsg := "lol real msg i don't wanna joke xdd"
	fmt.Printf("Sending %s to TLEP2\n", someMsg)
	encryptedMsg, err := tlep1.EncryptMessageCAFEA([]byte(someMsg))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Encrypted Message = %v\n", encryptedMsg)
	decryptedMsg, err := tlep2.DecryptMessageCAFEA(encryptedMsg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Getting %s on TLEP2 from TLEP1\n", decryptedMsg)
}
