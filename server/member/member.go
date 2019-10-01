package member

import (
	"errors"
	"time"
)

// ClanMembers is a struct wrapping the json output from Bungo
type ClanMembers struct {
	Response ClanMembersResponse
}

// ClanMembersResponse ...
type ClanMembersResponse struct {
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

// FindByName a Member in MembersCollection
func (m ClanMembers) FindByName(n string) (Member, error) {
	for _, member := range m.Response.Results {
		if member.DestinyUserInfo.DisplayName == n {
			return member, nil
		}
	}
	return Member{}, errors.New("Unable to find member: " + n)
}

// FindByID a Member in MembersCollection
func (m ClanMembers) FindByID(id string) (Member, error) {
	for _, member := range m.Response.Results {
		if member.DestinyUserInfo.MembershipID == id {
			return member, nil
		}
	}
	return Member{}, errors.New("Unable to find member: " + id)
}

// All returns all the clan members in an array
func (m ClanMembers) All() []Member {
	return m.Response.Results
}
