package main

import (
	"errors"
	"time"
)

// ClanMembers is a struct wrapping the json output from Bungo
type ClanMembers struct {
	Results []Member
}

// Member contains pertinent clan user info
type Member struct {
	DestinyUserInfo DestinyUserInfo
	IsOnline        bool
	JoinDate        time.Time
}

// DestinyUserInfo contains non-internal destiny info about a member
type DestinyUserInfo struct {
	MembershipType int
	MembershipID   string
	DisplayName    string
	IconPath       string
}

// Find a Member in MembersCollection
func (m ClanMembers) Find(n string) (Member, error) {
	for _, member := range m.Results {
		if member.DestinyUserInfo.DisplayName == n {
			return member, nil
		}
	}
	return Member{}, errors.New("Unable to find member: " + n)
}

// All returns all the clan members in an array
func (m ClanMembers) All() []Member {
	return m.Results
}
