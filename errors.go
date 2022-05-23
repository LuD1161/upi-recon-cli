package main

import "errors"

var (
	errCarReg       = errors.New("Car registration number invalid")
	errMobileNumber = errors.New("Mobile number invalid. Enter the 10 digit mobile number.")
	errGmail        = errors.New("Gmail id invalid.")
)
