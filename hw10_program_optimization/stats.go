package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result *users, err error) {
	result = new(users)
	scanner := bufio.NewScanner(r)
	var line []byte
	var user User
	var index int
	for scanner.Scan() {
		line = scanner.Bytes()
		err = user.UnmarshalJSON(line)
		if err != nil {
			return result, err
		}
		result[index] = user
		index++
	}
	if err := scanner.Err(); err != nil {
		return result, err
	}
	return result, nil
}

func countDomains(u *users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	var email, key string
	var endsWith bool
	for _, user := range u {
		email = strings.ToLower(user.Email)
		endsWith = strings.HasSuffix(email, "."+domain)
		if endsWith {
			key = strings.SplitN(email, "@", 2)[1]
			result[key]++
		}
	}
	return result, nil
}
