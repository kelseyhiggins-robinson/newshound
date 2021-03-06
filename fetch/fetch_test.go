package fetch

import (
	"net/mail"
	"reflect"
	"testing"
)

func TestFindHREFs(t *testing.T) {
	tests := []struct {
		given string
		want  []string
	}{
		{
			`<html>
				<body>
					<div>
					random text
					</div>
					<div>
						<a href="http://newshound.test.link.com">a link!</a>
					</div>
					<ol>
						<li>
						<a href="http://newshound.test.link.com/123">another link!</a>
						</li>
					</ol>
				</body>
			</html>`,
			[]string{
				"http://newshound.test.link.com",
				"http://newshound.test.link.com/123",
			},
		},
		{
			`A plain text email with a simple link <a href="http://newshound.test.link.com">a link!</a>`,
			[]string{
				"http://newshound.test.link.com",
			},
		},
		{
			`A plain text email with a simple link http://newshound.test.link.com/123?abc `,
			[]string{
				"http://newshound.test.link.com/123?abc",
			},
		},
	}

	for _, test := range tests {
		got := findHREFs([]byte(test.given))
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("findHREFs() got:%s want:%s", got, test.want)
		}
	}
}

func TestScrubBody(t *testing.T) {
	tests := []struct {
		given string
		email string
		want  string
	}{
		{
			`unsubscribe Unsubscribe blah@gmail.com blah`,
			"blah@gmail.com",
			"   ",
		},
	}

	for _, test := range tests {
		got := scrubBody([]byte(test.given), test.email)
		if got != test.want {
			t.Errorf("scrubBody(%s, %s) got:%s want:%s", test.given, test.email, got, test.want)
		}
	}
}

func TestFindSender(t *testing.T) {
	tests := []struct {
		given *mail.Address
		want  string
	}{
		{
			&mail.Address{Name: "USATODAY.com"},
			"USATODAY.com",
		},
		{
			&mail.Address{Name: "FT Exclusive"},
			"FT",
		},
		{
			&mail.Address{Name: "Los Angeles Times"},
			"Los Angeles Times",
		},
		{
			&mail.Address{Name: "LA Times"},
			"Los Angeles Times",
		},
		{
			&mail.Address{Name: "NYTimes.com"},
			"NYTimes.com",
		},
		{

			&mail.Address{Name: "The Washington Post"},
			"The Washington Post",
		},
		{
			&mail.Address{Name: "POLITICO Breaking News"},
			"POLITICO",
		},
	}

	for _, test := range tests {
		got := findSender(test.given)
		if got != test.want {
			t.Errorf("findSender(%s) got:%s want:%s", test.given, got, test.want)
		}
	}
}
