package models

type Error struct {
	UserMessage string  `json:"userMessage"`
	Code        float64 `json:"code"`
}

type OkResponse struct {
	Ok     bool    `json:"ok"`
	Errors []Error `json:"errors"`
}

type LoginResponse struct {
	IdToken string  `json:"idToken"`
	Errors  []Error `json:"errors"`
}

type QueryEventsResponse struct {
	Events *[]Event `json:"events"`
	Count  int64    `json:"count"`
}

type GetEventResponse struct {
	Event      *Event    `json:"event"`
	Tags       *[]string `json:"tags"`
	IsTimeSet  bool      `json:"isTimeSet"`
	IsActivist bool      `json:"isActivist"`
	IsJoined   bool      `json:"isJoined"`
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

type JoinEventVolonteurResponse struct {
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
	AsVolonteur bool  `json:"asVolonteur"`
	FormId      int64 `json:"formId"`
}
