package security

import (
	"blockchain_go/structure/member"
	"errors"
)

type Permission bool

func (p Permission) Level(address string) error {
	members, _ := member.CreateMembers()
	member := members.GetAllAdresses()

	for _, addr := range member {
		if addr == address {
			return nil
		}
	}

	return errors.New("permission denied: address not found in members")
}
