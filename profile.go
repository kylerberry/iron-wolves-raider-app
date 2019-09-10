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
