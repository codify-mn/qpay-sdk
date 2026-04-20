package qpay

import (
	"context"
	"net/http"
	"testing"
)

func TestLocations(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/auth/token":
			tokenHandler(w)
		case "/v2/aimaghot":
			_, _ = w.Write([]byte(`[{"code":"UB","name":"Ulaanbaatar"}]`))
		case "/v2/sumduureg/UB":
			_, _ = w.Write([]byte(`[{"code":"BZD","name":"Bayanzurkh"}]`))
		case "/v2/bank/list":
			_, _ = w.Write([]byte(`[{"code":"KHAN","name":"Khan Bank"}]`))
		}
	})
	defer srv.Close()

	cities, err := c.GetCities(context.Background())
	if err != nil || len(cities) != 1 {
		t.Fatalf("cities: err=%v cities=%+v", err, cities)
	}
	districts, err := c.GetDistricts(context.Background(), "UB")
	if err != nil || len(districts) != 1 {
		t.Fatalf("districts: err=%v districts=%+v", err, districts)
	}
	banks, err := c.GetBanks(context.Background())
	if err != nil || len(banks) != 1 || banks[0].Code != "KHAN" {
		t.Fatalf("banks: err=%v banks=%+v", err, banks)
	}
}
