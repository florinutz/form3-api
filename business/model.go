package business

import "gopkg.in/mgo.v2/bson"

type BeneficiaryParty struct {
	AccountName       string `json:"account_name,omitempty"`
	AccountNumber     string `json:"account_number,omitempty"`
	AccountNumberCode string `json:"account_number_code,omitempty"`
	Address           string `json:"address,omitempty"`
	// ...
}

type DebtorParty struct {
	AccountName       string `json:"account_name,omitempty"`
	AccountNumber     string `json:"account_number,omitempty"`
	AccountNumberCode string `json:"account_number_code,omitempty"`
	Address           string `json:"address,omitempty"`
	// ...
}

type SponsorParty struct {
	AccountNumber string `json:"account_number,omitempty"`
	BankId        int    `json:"bank_id,omitempty"`
	BankIdCode    string `json:"bank_id_code,omitempty"`
}

type SenderCharge struct {
	Amount   float64 `json:"amount,omitempty"`
	Currency string  `json:"currency,omitempty"`
}

type ChargesInformation struct {
	BearerCode    string         `json:"bearer_code,omitempty"`
	SenderCharges []SenderCharge `json:"sender_charges,omitempty"`
	// ...
}

type Attributes struct {
	Amount             float32             `json:"amount,omitempty"`
	Currency           string              `json:"currency,omitempty"`
	BeneficiaryParty   *BeneficiaryParty   `json:"beneficiary_party,omitempty"`
	DebtorParty        *DebtorParty        `json:"debtor_party,omitempty"`
	SponsorParty       *SponsorParty       `json:"sponsor_party,omitempty"`
	ChargesInformation *ChargesInformation `json:"charges_information,omitempty"`
	EndToEndReference  string              `json:"end_to_end_reference,omitempty"`
	ProcessingDate     int64               `json:"processing_date,omitempty"`
	// ...
}

type Payment struct {
	Id             bson.ObjectId `bson:"_id" json:"id,omitempty"`
	Version        int           `json:"version,omitempty"`
	OrganisationId bson.ObjectId `json:"organisation_id,omitempty"`
	Attributes     *Attributes   `json:"attributes,omitempty"`
}
