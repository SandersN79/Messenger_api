package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"strconv"
)

func KeyGen() string {
	src := &cryptoSource{}
	rnd := rand.New(src)
	count := 0
	key1 := 0
	key2 := 0
	key3 := 0

	for count < 3 {
		if count == 2 {
			key3 = rnd.Intn(1000)
		} else if count == 1 {
			key2 = rnd.Intn(5000)
		} else if count == 0{
			key1 = rnd.Intn(5000)
		}
		count = count + 1
	}
	UserKey := strconv.Itoa(key1) + " " + strconv.Itoa(key2) + " " + strconv.Itoa(key3)
	fmt.Println(UserKey)
	return UserKey
}

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}