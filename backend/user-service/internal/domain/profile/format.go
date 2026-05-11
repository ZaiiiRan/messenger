package profile

import (
	"fmt"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

func FormatPhone(phone, defaultRegion string) (string, error) {
	if phone == "" {
		return "", nil
	}

	num, err := parseAny(phone, defaultRegion)
	if err != nil {
		return "", err
	}

	cc := num.GetCountryCode()
	national := phonenumbers.GetNationalSignificantNumber(num)

	areaLen := phonenumbers.GetLengthOfGeographicalAreaCode(num)
	if areaLen == 0 {
		areaLen = phonenumbers.GetLengthOfNationalDestinationCode(num)
	}

	if areaLen == 0 || areaLen >= len(national) {
		return phonenumbers.Format(num, phonenumbers.INTERNATIONAL), nil
	}

	area := national[:areaLen]
	body := national[areaLen:]

	return fmt.Sprintf("+%d(%s)-%s", cc, area, splitGroups(body)), nil
}

func parseAny(phone, defaultRegion string) (*phonenumbers.PhoneNumber, error) {
	region := defaultRegion
	if strings.HasPrefix(strings.TrimSpace(phone), "+") {
		region = ""
	}

	num, err := phonenumbers.Parse(phone, region)
	if err != nil {
		return nil, fmt.Errorf("parse phone %q: %w", phone, err)
	}
	if !phonenumbers.IsValidNumber(num) {
		return nil, fmt.Errorf("phone %q is not a valid number", phone)
	}
	return num, nil
}

func splitGroups(s string) string {
	switch len(s) {
	case 7:
		return s[:3] + "-" + s[3:5] + "-" + s[5:7]
	case 6:
		return s[:3] + "-" + s[3:6]
	case 8:
		return s[:4] + "-" + s[4:8]
	}

	var groups []string
	for len(s) > 4 {
		groups = append(groups, s[:3])
		s = s[3:]
	}
	groups = append(groups, s)
	return strings.Join(groups, "-")
}
