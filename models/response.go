package models

type Error struct {
	UserMessage string  `json:"userMessage"`
	Code        float64 `json:"code"`
}

type DefaultResponse struct {
	Ok     bool     `json:"ok"`
	Errors *[]Error `json:"errors,omitempty"`
	Error  *Error   `json:"error,omitempty"`
}

type IndexPageResponse struct {
	SoonerEvents *[]Event       `json:"soonerEvents,omitempty"`
	EventsByTags *[]EventsByTag `json:"eventsByTags,omitempty"`
}

type EventsByTag struct {
	Tag    string   `json:"tag,omitempty"`
	Events *[]Event `json:"events,omitempty"`
}

type LoginResponse struct {
	Ok      bool    `json:"ok"`
	IdToken string  `json:"idToken,omitempty"`
	Errors  []Error `json:"errors,omitempty"`
}

type QueryEventsResponse struct {
	Events *[]Event `json:"events"`
	Count  int64    `json:"count"`
}

type GetEventResponse struct {
	Event       *Event    `json:"event"`
	Organizer   *User     `json:"organizer"`
	Tags        *[]string `json:"tags,omitempty"`
	AsVolunteer bool      `json:"asVolunteer,omitempty"`
	IsJoined    bool      `json:"isJoined,omitempty"`
	HasForm     bool      `json:"hasForm,omitempty"`
}

type GetUserInfoResponse struct {
	User   *User   `json:"user"`
	Errors []Error `json:"errors"`
}

type UserInfo struct {
	Email      *string `json:"email"`
	Group      *int64  `json:"group"`
	FirstName  *string `json:"firstName"`
	SecondName *string `json:"secondName"`
	LastName   *string `json:"lastName"`
	Gender     *int64  `json:"gender"`
}

type AddEventResponse struct {
	Ok      bool    `json:"ok"`
	Errors  []Error `json:"errors"`
	EventId int64   `json:"eventId"`
}

type JoinEventVolunteerResponse struct {
	Ok          bool  `json:"ok"`
	HasForm     bool  `json:"hasForm"`
	OrganizerId int64 `json:"organizerId"`
}

type GetTagStatusResponse struct {
	Ok        bool `json:"ok"`
	HasStatus bool `json:"hasStatus"`
	Status    bool `json:"status"`
}

type GetJoinedUsersResponse struct {
	Ok    bool          `json:"ok"`
	Users *[]JoinedUser `json:"users"`
}

type JoinedUser struct {
	User        User  `json:"user"`
	AsVolunteer bool  `json:"asVolunteer"`
	FormId      int64 `json:"formId"`
}

type GenerateTemplateTokenResponse struct {
	Ok    bool   `json:"ok"`
	Token string `json:"token"`
}

func AppendError(errors *[]Error, message string, code float64) {
	*errors = append(*errors, Error{
		UserMessage: message,
		Code:        code,
	})
}
