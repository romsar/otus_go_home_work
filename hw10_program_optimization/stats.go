package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	json "github.com/json-iterator/go"
)

var ErrEmptyDomain = errors.New("domain is not passed")

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
	result := make(DomainStat)

	if domain == "" {
		return result, ErrEmptyDomain
	}

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	user := new(User)

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), user); err != nil {
			return result, err
		}

		if matched := strings.HasSuffix(user.Email, "."+domain); matched {
			d := getDomainFromEmail(user.Email)
			result[d]++
		}

		user = &User{}
	}

	return result, nil
}

func getDomainFromEmail(email string) string {
	return strings.ToLower(strings.SplitN(email, "@", 2)[1])
}
