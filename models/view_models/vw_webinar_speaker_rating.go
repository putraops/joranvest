package view_models

import "database/sql"

type WebinarSpeakerRatingViewModel struct {
	Id                        string       `json:"id"`
	CreatedAt                 sql.NullTime `json:"created_at"`
	SpeakerId                 string       `json:"speaker_id"`
	UserId                    string       `json:"user_id"`
	ObjectRatedId             string       `json:"object_rated_id"`
	ReferenceId               string       `json:"reference_id"`
	Rating                    int          `json:"rating"`
	Comment                   string       `json:"comment"`
	OrganizationName          string       `json:"organization_name"`
	OrganizationFilepath      string       `json:"organization_filepath"`
	OrganizationFilepathThumb string       `json:"organization_filepath_thumb"`
	SpeakerFullname           string       `json:"speaker_fullname"`
	SpeakerFilepath           string       `json:"speaker_filepath"`
	SpeakerFilepathThumb      string       `json:"speaker_filepath_thumb"`
}
