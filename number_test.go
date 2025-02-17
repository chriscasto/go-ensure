package valid

import "testing"

func TestNumberValidator_InRange(t *testing.T) {
	rangeMin := 1
	rangeMax := 10

	// Check to make sure the method panics if max < min
	t.Run("panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		bad := Number[int]().InRange(rangeMax, rangeMin)
		if err := bad.Validate(rangeMin); err != nil {
			t.Errorf("validation occured and generated an error: %s", err.Error())
		}

	})

	testCases := []struct {
		name      string
		i         int
		expectErr bool
	}{
		{"less than", rangeMin - 1, true},
		{"bottom of range", rangeMin, false},
		{"top of range", rangeMax, true},
		{"greater than", rangeMax + 1, true},
	}

	iv := Number[int]().InRange(rangeMin, rangeMax)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := iv.Validate(tc.i)
			if err != nil && !tc.expectErr {
				t.Errorf("InRange(%d, %d).Validate(%d); expected no error, got %s", rangeMin, rangeMax, tc.i, err)
			} else if err == nil && tc.expectErr {
				t.Errorf("InRange(%d, %d).Validate(%d); expected error but got none", rangeMin, rangeMax, tc.i)
			}
		})
	}
}

func TestNumberValidator_LessThan(t *testing.T) {
	testCases := []struct {
		name      string
		i         int
		expectErr bool
	}{
		{"less than", 2, false},
		{"equal to", 10, true},
		{"greater than", 20, true},
	}

	valMax := 10
	iv := Number[int]().LessThan(valMax)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := iv.Validate(tc.i)
			if err != nil && !tc.expectErr {
				t.Errorf("LessThan(%d).Validate(%d); expected no error, got %s", valMax, tc.i, err)
			} else if err == nil && tc.expectErr {
				t.Errorf("LessThan(%d).Validate(%d); expected error but got none", valMax, tc.i)
			}
		})
	}
}

func TestNumberValidator_GreaterThan(t *testing.T) {
	testCases := []struct {
		name      string
		i         int
		expectErr bool
	}{
		{"less than", 2, true},
		{"equal to", 10, true},
		{"greater than", 20, false},
	}

	valMin := 10
	iv := Number[int]().GreaterThan(valMin)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := iv.Validate(tc.i)
			if err != nil && !tc.expectErr {
				t.Errorf("GreaterThan(%d).Validate(%d); expected no error, got %s", valMin, tc.i, err)
			} else if err == nil && tc.expectErr {
				t.Errorf("GreaterThan(%d).Validate(%d); expected error but got none", valMin, tc.i)
			}
		})
	}
}
