package dota

import "testing"

func TestGetDataResponseReportsWinProbabilityAvailability(t *testing.T) {
	response := GetDataResponse{WinProbabilityAvailable: true}
	if !response.WinProbabilityAvailable {
		t.Fatal("GetDataResponse did not retain win probability availability")
	}
}
