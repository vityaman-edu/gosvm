package fp_test

import (
	"strconv"
	"testing"

	"github.com/vityaman-edu/gosvm/internal/fp"
)

func TestFold(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name   string
		values []int
		want   int
	}{
		{name: "initial and left-to-right order", values: []int{3, 2}, want: 15},
		{name: "empty", values: []int{}, want: 20},
		{name: "nil", values: nil, want: 20},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := 20
			calls := 0

			fp.Fold(test.values, &result, func(accumulator, value int) int {
				calls++

				return accumulator - value
			})

			if result != test.want {
				t.Errorf("Fold result = %d, want %d", result, test.want)
			}

			if calls != len(test.values) {
				t.Errorf("reducer calls = %d, want %d", calls, len(test.values))
			}
		})
	}
}

func TestFoldDifferentTypes(t *testing.T) {
	t.Parallel()

	result := "seed:"
	fp.Fold([]int{3, 1, 2}, &result, func(accumulator string, value int) string {
		return accumulator + strconv.Itoa(value)
	})

	if result != "seed:312" {
		t.Errorf("Fold result = %q, want seed:312", result)
	}
}

func TestFold1(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name   string
		values []int
		want   int
	}{
		{name: "first value seed", values: []int{7}, want: 7},
		{name: "left-to-right order", values: []int{20, 3, 2}, want: 2032},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := 99
			calls := 0

			fp.Fold1(test.values, &result, func(accumulator, value int) int {
				calls++

				return accumulator*10 + value
			})

			if result != test.want {
				t.Errorf("Fold1 result = %d, want %d", result, test.want)
			}

			if want := len(test.values) - 1; calls != want {
				t.Errorf("reducer calls = %d, want %d", calls, want)
			}
		})
	}
}

func TestFold1Empty(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name   string
		values []int
	}{
		{name: "empty", values: []int{}},
		{name: "nil", values: nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := 99

			checkPanic(t, "Fold1: empty slice", func() {
				fp.Fold1(test.values, &result, func(accumulator, value int) int {
					t.Error("reducer called for empty slice")

					return accumulator + value
				})
			})

			if result != 99 {
				t.Errorf("Fold1 changed output to %d after panic", result)
			}
		})
	}
}

func TestFoldReducerPanic(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		fold func([]int, *int, func(int, int) int)
		name string
	}{
		{name: "Fold", fold: fp.Fold[int, int]},
		{name: "Fold1", fold: fp.Fold1[int]},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := 99

			checkPanic(t, "reducer failed", func() {
				test.fold([]int{1, 2, 3}, &result, func(accumulator, value int) int {
					if value == 3 {
						panic("reducer failed")
					}

					return accumulator + value
				})
			})

			if result != 99 {
				t.Errorf("%s changed output to %d after panic", test.name, result)
			}
		})
	}
}

func checkPanic(t *testing.T, want string, action func()) {
	t.Helper()

	defer func() {
		if got := recover(); got != want {
			t.Errorf("panic = %v, want %q", got, want)
		}
	}()

	action()
}
