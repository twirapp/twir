package moderation

import "testing"

func TestIsTooMuchCaps(t *testing.T) {
	res := IsTooMuchCaps("QWERTyuiop", 60)
	if res {
		t.Errorf(`IsTooMuchCaps("QWERTyuiop", 60) = true; want false`)
	}

	res = IsTooMuchCaps("QWERTyuiop", 50)
	if !res {
		t.Errorf(`IsTooMuchCaps("QWERTyuiop", 50) = false; want true`)
	}
}
