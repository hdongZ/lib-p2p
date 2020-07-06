package main

import (
	"crypto/rand"
	"fmt"

	"github.com/hdongZ/blsAggregateSig/crypto"
)

type SeqRead struct {
}

func (self *SeqRead) Read(buf []byte) (int, error) {
	n := len(buf)
	for i := 0; i < n; i++ {
		buf[i] = byte(i)
	}
	return n, nil
}

func testReadRand() {
	s1 := new(SeqRead)
	p2p.SetRandFunc(s1)
	var sec p2p.SecretKey
	sec.SetByCSPRNG()
	buf := sec.GetLittleEndian()
	fmt.Printf("1. buf=%x\n", buf)
	for i := 0; i < len(buf); i++ {
		if buf[i] != byte(i) {
			fmt.Printf("err %d\n", i)
		}
	}
	p2p.SetRandFunc(rand.Reader)
	sec.SetByCSPRNG()
	buf = sec.GetLittleEndian()
	fmt.Printf("2. (cr.Read) buf=%x\n", buf)
	p2p.SetRandFunc(nil)
	sec.SetByCSPRNG()
	buf = sec.GetLittleEndian()
	fmt.Printf("3. (cr.Read) buf=%x\n", buf)
}

func main() {
	p2p.Init(p2p.BLS12_381)
	var sec p2p.SecretKey
	sec.SetByCSPRNG()
	fmt.Printf("sec:%s\n", sec.SerializeToHexStr())
	pub := sec.GetPublicKey()
	fmt.Printf("pub:%s\n", pub.SerializeToHexStr())
	msgTbl := []string{"abc", "def", "123"}
	n := len(msgTbl)
	sigVec := make([]*p2p.Sign, n)
	for i := 0; i < n; i++ {
		m := msgTbl[i]
		sigVec[i] = sec.Sign(m)
		fmt.Printf("%d. sign(%s)=%s\n", i, m, sigVec[i].SerializeToHexStr())
	}
	agg := sigVec[0]
	for i := 1; i < n; i++ {
		agg.Add(sigVec[i])
	}
	hashPt := p2p.HashAndMapToSignature([]byte(msgTbl[0]))
	for i := 1; i < n; i++ {
		hashPt.Add(p2p.HashAndMapToSignature([]byte(msgTbl[i])))
	}
	fmt.Printf("verify %t\n", p2p.VerifyPairing(agg, hashPt, pub))
	testReadRand()
}
