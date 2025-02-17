package valid

import "testing"

func TestStringValidator_HasLength(t *testing.T) {
	testCases := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"fewer letters", "a", true},
		{"same letters", "abc", false},
		{"more letters", "wxyz", true},
	}

	strLen := 3
	sv := String().HasLength(strLen)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := sv.Validate(tc.value)
			if err != nil && !tc.expectErr {
				t.Errorf("HasLength(%d).Validate(%s); expected no error, got %s", strLen, tc.value, err)
			} else if err == nil && tc.expectErr {
				t.Errorf("HasLength(%d).Validate(%s); expected error but got none", strLen, tc.value)
			}
		})
	}
}

func TestStringValidator_ShorterThan(t *testing.T) {
	testCases := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"fewer letters", "a", false},
		{"same letters", "abc", true},
		{"more letters", "wxyz", true},
	}

	strLen := 3
	sv := String().ShorterThan(strLen)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := sv.Validate(tc.value)
			if err != nil && !tc.expectErr {
				t.Errorf("ShorterThan(%d).Validate(%s); expected no error, got %s", strLen, tc.value, err)
			} else if err == nil && tc.expectErr {
				t.Errorf("ShorterThan(%d).Validate(%s); expected error but got none", strLen, tc.value)
			}
		})
	}
}

func TestStringValidator_LongerThan(t *testing.T) {
	testCases := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"fewer letters", "a", true},
		{"same letters", "abc", true},
		{"more letters", "wxyz", false},
	}

	strLen := 3
	sv := String().LongerThan(strLen)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := sv.Validate(tc.value)
			if err != nil && !tc.expectErr {
				t.Errorf("LongerThan(%d).Validate(%s); expected no error, got %s", strLen, tc.value, err)
			} else if err == nil && tc.expectErr {
				t.Errorf("LongerThan(%d).Validate(%s); expected error but got none", strLen, tc.value)
			}
		})
	}
}
