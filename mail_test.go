package main

import "testing"

func TestSendMail(t *testing.T) {
	loadEnv()
	err := sendEmail("test", "this is just a test")
	if err != nil {
		panic(err)
	}
}
