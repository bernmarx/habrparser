package storage

import (
	"errors"
	"testing"

	"github.com/bernmarx/habrparser/internal/scraper"
	"github.com/golang/mock/gomock"
)

func TestAddPageData(t *testing.T) {
	ctrl := gomock.NewController(t)

	m := NewMockstorage(ctrl)

	testPage := scraper.Page{ID: 123}
	testPageJSON := testPage.GetJSON()

	m.EXPECT().Exec(`SELECT addpage($1)`, testPageJSON).Return(nil, nil)
	m.EXPECT().Exec(`SELECT addpagejson($1, $2)`, int(123), testPageJSON).Return(nil, nil)

	s := Storage{m}

	err := s.AddPageData(&testPage)

	if err != nil {
		t.Error(err)
	}

	m.EXPECT().Exec(`SELECT addpage($1)`, testPageJSON).Return(nil, nil)
	m.EXPECT().Exec(`SELECT addpagejson($1, $2)`, int(123), testPageJSON).Return(nil, errors.New("err"))

	err = s.AddPageData(&testPage)

	if err == nil {
		t.Error(err)
	}

	m.EXPECT().Exec(`SELECT addpage($1)`, testPageJSON).Return(nil, errors.New("err"))

	err = s.AddPageData(&testPage)

	if err == nil {
		t.Error(err)
	}
}
