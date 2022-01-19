package veil

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func GetEnvToken() string {
	token := os.Getenv("VEIL_API_TOKEN")
	if token == "" {
		token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJ1c2VyX2lkIjo1NiwidXNlcm5hbWUiOiJidXIiLCJleHAiOjE5NTcwNjg2MDYsInNzbyI6ZmFsc2UsIm9yaWdfaWF0IjoxNjQyNTcyNjA2fQ.S249CYNV6yg25hx5tMDNYH5AV00ZNS3fay-dkb4MISa-AOGfRB9tGtico5wdFOq8DpzpqCoRRBZPoPB9yTxz0c-hgMOWK3mrR6QAroOi1RGw8Cuaqhg7jdlgkc6S1fhcBPZB0r9rPQ38N3MNkrcUnxNrj2JhOCXG8av6ZBuK-BxzO4Gb6PeuVCtqP9UT-dIhyR54_vf_hD_fRsBxmfzWib0htL_TjtQYe99fs0q6uOHQGtw6lIdgmMQVdgdedFgziON4PnwVY__RTDWX9SS70u6dRUIqLX9Hhh3317n3uVezxKb-PizCk3WKJOWdqCA5ISIRV3v3YL32UEoOhhDotnhhcw-2n_mHcZjw0vjYihRsv3cOKrYNjBRmLo0ax9qONi5hqWxqH6Jr9HAMEUu_WFaKZG0Sy2RKxlrYpmffBeB0Rc7NyoR-P8gSjo1tE_8CO5Gc-9_BfyNbRj0JpdRLRRs7CZdH2qeiWx2NHdmRk1vuLyumulojX0Aixo9gf5-OCKiJg1TcKWmm-moVma_zrcx4C8k1lB0pJtJ8op8Kn9fFryX9kwrXNoN2BecwlAiZ6KO6ciX3VVG0TpRZV3sKFAxlelD_Sm3HGrPOTPdsyNy7Wku4sQplrgPq1TFhTGK_pL4BvuIJi6Soj5kdaMPRw_MW46rhz8o4keaouC5-bdg"
		//return token, errors.New("Token is empty")
		return token
	}

	return token
}

func GetEnvUrl() string {
	url := os.Getenv("VEIL_API_URL")
	if url == "" {
		url := "http://192.168.11.105"
		//return url, errors.New("Url is empty")
		return url
	}

	return url
}

func IsSuccess(code int) bool {

	successCodes := []int{200, 202}

	for _, i := range successCodes {
		if i == code {
			return true
		}
	}
	return false
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func NameGenerator(class string) string {
	return fmt.Sprint(class, "_", StringWithCharset(3, charset), "__TEST")
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
