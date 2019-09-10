package main

import (
	"errors"
	"time"
)

// Member contains pertinent clan user info
type Member struct {
	DestinyUserInfo DestinyUserInfo
	IsOnline        bool
	JoinDate        time.Time
}

// Members is an array of members
type Members struct {
	Response MemberResponse
}

// MemberResponse is a struct wrapping the json output from Bungo
type MemberResponse struct {
	Results []Member
}

// DestinyUserInfo contains non-internal destiny info about a member
type DestinyUserInfo struct {
	MembershipType int
	MembershipID   string
	DisplayName    string
	IconPath       string
}

// Find a Member in MembersCollection
func (m Members) Find(n string) (Member, error) {
	for _, member := range m.Response.Results {
		if member.DestinyUserInfo.DisplayName == n {
			return member, nil
		}
	}
	return Member{}, errors.New("Unable to find member: " + n)
}
