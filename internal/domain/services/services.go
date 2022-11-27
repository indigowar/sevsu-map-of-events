package services

type Services struct {
	Event           EventService
	Organizer       OrganizerService
	FoundingRange   RangeService
	CoFoundingRange RangeService
	Competitor      CompetitorService
	Subject         SubjectService
	Image           ImageService
}
