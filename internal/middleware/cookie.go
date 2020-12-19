package middleware

import (
	"fmt"
	"strings"
)

func parseCookie(cookies string) map[string]string {
	cookiesMap := make(map[string]string)
	for _, c := range strings.Split(cookies, ";") {
		c = strings.Trim(c, " ")
		splited := strings.SplitN(c, "=", 2)
		if len(splited) < 2 {
			return nil
		}
		k := splited[0]
		v := splited[1]
		cookiesMap[k] = v
	}

	return cookiesMap
}

// SetCookie appends 'Set-Cookie' to headers in Response
func SetCookie(name string, value string, resp *Response) {
	cookie := fmt.Sprintf("%s=%s", name, value)
	cookie = cookie + "; SameSite=Lax"
	cookie = cookie + "; HttpOnly"
	cookie = cookie + "; Secure"
	cookie = fmt.Sprintf("%s; Domain=%s", cookie, domain)
	cookie = cookie + "; Path=/"

	var cookies []string
	headers := resp.MultiValueHeaders
	if headers == nil {
		resp.MultiValueHeaders = make(map[string][]string)
	}

	cookies, ok := resp.MultiValueHeaders["Set-Cookie"]
	if !ok {
		cookies = []string{}
	}
	resp.MultiValueHeaders["Set-Cookie"] = append(cookies, cookie)
}
