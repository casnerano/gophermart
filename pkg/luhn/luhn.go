package luhn

func Checksum(nums string) bool {
	sum, length := 0, len(nums)
	if length < 2 {
		return false
	}
	for index, num := range nums {
		dig := int(num - '0')
		if length%2 == index%2 {
			dig *= 2
			if dig > 9 {
				dig = dig%10 + dig/10
			}
		}
		sum += dig
	}
	return sum%10 == 0
}
