package main

type MemberProfile struct {
	Response ProfileResponse
}

type ProfileResponse struct {
	Profile Profile
}

type Profile struct {
	Data ProfileData
}

type ProfileData struct {
	CharacterIDs []string
}

// Characters returns a Profile's characterIds
func (p MemberProfile) Characters() []string {
	return p.Response.Profile.Data.CharacterIDs
}
