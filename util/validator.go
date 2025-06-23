package util

import (
	"net"
	"net/mail"
	"net/url"
	"time"
)

func ValidateEmail(email string) bool {

	e, err := mail.ParseAddress(email)

	if err != nil {
		return false
	}

	return e.Address == email

}

func ValidatePassword(password string) bool {

	if len(password) < 8 || len(password) > 100 {
		return false
	}

	return true

}

func ValidateUrl(domain string) bool {

	u, err := url.Parse(domain)

	if err != nil {
		return false
	}

	if u.Scheme == "" || u.Host == "" {
		return false
	}

	return true

}

func IsValidIP(ip string) bool {

	trial := net.ParseIP(ip)

	if trial == nil {
		return false
	}

	return true

}

func IsValidCIDR(cidr string) bool {

	_, _, err := net.ParseCIDR(cidr)

	if err != nil {
		return false
	}

	return true

}

func ValidateDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}
