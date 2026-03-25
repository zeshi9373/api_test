package fn

import "net/url"

func UrlEncode(s string) string {
	return url.QueryEscape(s)
}
