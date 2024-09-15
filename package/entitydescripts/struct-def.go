package entitydescripts

import "github.com/google/uuid"

type Tender struct {
	TenderID       uuid.UUID `json:"tender_id"`
	Name           string    `json:"name"`
	ServiceType    string    `json:"service_type"`
	Description    string    `json:"description"`
	Status         string    `json:"status"`
	OrganizationID uuid.UUID `json:"organization_id"`
	AuthorID       uuid.UUID `json:"author_id"` //внешний ключ к таблице employees
	Version        uint      `json:"version"`
}

type Offer struct {
	OfferID        uuid.UUID `json:"offer_id"`
	Name           string    `json:"name"`
	ServiceType    string    `json:"service_type"`
	Description    string    `json:"description"`
	Status         string    `json:"status"`
	OrganizationID uuid.UUID `json:"organization_id"`
	TenderID       uuid.UUID `json:"tender_id"`
	AuthorID       uuid.UUID `json:"author_id"`
	Version        uint      `json:"version"`
}
