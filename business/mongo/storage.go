package mongo // first time I work with mongo btw

import (
	"errors"
	"time"

	"form3/business"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// DB represents the mongo db name
	DB = "payments"
	// PaymentsCollection is the mongo payments collection
	PaymentsCollection = "payments"
)

// Storage implements the Storage interface
type Storage struct {
	session *mgo.Session
}

// New is the Storage constructor
func New(connection string, dialTimeout time.Duration) (*Storage, error) {
	s, err := mgo.DialWithTimeout(connection, dialTimeout)
	if err != nil {
		return nil, err
	}
	return &Storage{s}, err
}

func (p *Storage) getFreshSession() *mgo.Session {
	return p.session.Copy()
}

// Create implements the Storage interface's Create
func (p *Storage) Create(payment business.Payment) (*business.Payment, error) {
	session := p.getFreshSession()
	defer session.Close()

	if !payment.Id.Valid() {
		payment.Id = bson.NewObjectId()
	}

	if !payment.OrganisationId.Valid() {
		payment.OrganisationId = bson.NewObjectId()
	}

	err := session.DB(DB).C(PaymentsCollection).Insert(&payment)

	return &payment, err
}

// Retrieve implements the Storage interface's Retrieve
func (p *Storage) Retrieve(id []byte) (*business.Payment, error) {
	s := p.getFreshSession()
	defer s.Close()

	if !bson.IsObjectIdHex(string(id)) {
		return nil, errors.New("object id is invalid")
	}

	payment := new(business.Payment)
	err := s.DB(DB).C(PaymentsCollection).FindId(bson.ObjectIdHex(string(id))).One(payment)

	return payment, err
}

// Update implements the Storage interface's Update
func (p *Storage) Update(payment business.Payment) error {
	panic("implement me")
}

// Delete implements the Storage interface's Delete
func (p *Storage) Delete(id []byte) error {
	panic("implement me")
}

// List implements the Storage interface's List
func (p *Storage) List() (payments []*business.Payment, err error) {
	s := p.getFreshSession()
	defer s.Close()

	err = s.DB(DB).C(PaymentsCollection).Find(nil).All(&payments)

	return
}
