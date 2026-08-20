package order

import "strconv"

type Number struct {
	Value string
}

func NewNumber(value string) Number {
	return Number{Value: value}
}

func (o Number) CheckWithLuna() bool {
	sum := 0
	alternate := false

	for i := len(o.Value) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(o.Value[i]))
		if err != nil {
			return false
		}

		if alternate {
			digit *= 2

			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}
