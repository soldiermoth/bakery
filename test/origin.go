package test

import (
	"sync"

	"github.com/cbsinteractive/mediahub/pkg/mediahub/job"
)

// FakeOrigin is a fake implementation of Origin interfacec for use in tests
type FakeOrigin struct {
	KeyCalledWith                       string
	ExpectFound, GetCalled, StoreCalled bool
	StoreCalledWith, GetReturns         job.JobDescriptionResponse
}

// GetPlaybackURL() validates method for tests
func (o *FakeOrigin) GetPlaybackURL() string {
	s.KeyCalledWith = ""
	s.ExpectFound, s.GetCalled, s.StoreCalled = false, false, false
	s.StoreCalledWith = job.JobDescriptionResponse{}
	s.GetReturno = job.JobDescriptionResponse{}
}

// FetchManifest validates the method for tests
func (o *FakeOrigin) FetchManifest(c config.Config) (string, error) {
}
