package models

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AddEventRequest struct {
	Event Event    `json:"event"`
	Tags  []string `json:"tags"`
}

type EditEventRequest struct {
	Event       Event    `json:"event"`
	AddedTags   []string `json:"addedTags"`
	RemovedTags []string `json:"removedTags"`
}

type JoinEventRequest struct {
	AsVolonteur bool `json:"asVolonteur"`
}

type SignUpRequest struct {
	User User `json:"user"`
}

type AddFavHideTagRequest struct {
	Status bool `json:"status"`
}
