package mock

import (
	"fmt"
	"io/ioutil"
	"os"
)

var grpcCert = map[string]string{
	"key": `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDVm9llcSG8opFu
lP64J13TWkV840kTEYCqKYwM3EiA+b3Yf0FmZJwB7+ENguiz40GZxtb0us2Tn1gK
/NWgKJZpTyAfMVNbEGDVOThPw9S+sbbCol/i/4ATK4cT7Az4n3YJ1au7VnZFxsJd
k6HA7kq0snPQDgdBgO+FPCyKnAgQMX8ds1mEBhTCkoWZwPRioEiqlzJRJVlajzJ6
HhRDYrbNOMm32I9A1+Fc9aQpJwRVG+TlimSZIbTO3h/T9oR6AufmKUwVs/yer2TP
ALTEO+2Ek0K/l1qxq+mAIYx2s1PFWbj9ctFjaxOSKhXdyVcEnzONiiuMun3Ji3Pz
mrVqbpbNAgMBAAECggEAKd6omDe6syzycMiHvTUKMRlqsXYjprcxjykhqIutMorB
aaRX/2gNJFwOswVL86bB/xB4pfPPw/Xh3hV9Ei2iffXih/q1Kg5lzKWwogHyu4Y+
GpyVSvtl8VwA2CjWyg0HqBnX1Zq3CCpTguAjPpv2zMCF7uwxA+fwCx+mY2z+L54I
RQE93bqQdy9NaL3AC2dky6qfYc2Gx+hx+gZBD8s1Q8h9+A1zfQyU9mpcHkDenIIg
L0mBKjXW59wukNhoY2SBb1Q9y8TUpHf4VP+xBoH96S/QT3sBofQ8Q1offep3WKKD
Y2vrG34o3sQjGrnnufJzuQHQvwT43kzEUBX6RI90iQKBgQD4FLuDsld4nQIFpZbT
Wd1Ld+qXNUuB7oS4rWIGOQPDDWRP297XEEXrEF4TbGUa+8I4LSzhNPwB/948XM0V
9lfjR3u3w2j5IaJFhBN9KVRkCRxUibDhuHA7LMJQMDglj8szwsCrrSZWfdpnKLDL
WyXSxFIWzmqxUx08cnwIly66IwKBgQDcbWq5lMW5Fnae9iJhlEWEF4sa9fgKkcz6
pbXQx5r+xRiOah3CV94WmmftCIub/9tvO56riFtnNM/aRyenRjsfpWFkA8105Rrl
j488BnUjyMCYiFZyqDT7DuPxEJh9bCnNSqcESK1QxSK6AnAP3q3qdklN121VtDc9
++molvOiTwKBgHEXYnQi9OUzDhzs49jteohQ3kyYKxfMWAyoXatgimp8zGHrZaa8
8GK8T2ajX2PxqRRa4762nLt8nR7/Xy7H4kDl8WxQVKZdws/V6dyA3svLq3KOYmhD
4EXZnatYj//vkT7DZXndsUB0lv+3+QB7SL7QaGulJdY4gXdw6UIxSUfpAoGBANT5
XwnKBbRMUPZLyHJRiU0UVlIJX8wOjWeLnn0HrukD1DMdsn0o2qsqKsmp3QIwFnuF
tkvz5qR0MXOsFlMXl15/MvcoeWW9StyMdY9AigO2HugBqs0DWpVMEM7FAyED1evF
elO4SMTmhCQG4PFkbNNB0JfGUpxhEJLyCBPdLa8fAoGBAKkqZaweygSP/wlTLJQ4
pC1NJ3LLUoU1jG+tQDlmifI8KXK2wOzC7B+UtK8QnYANcBdUC2I/8ILBuhElty1J
2eX9ZC6pIgbsG7rT2VdkV1JgRLUHXJHvcHf7xjzJLKI7GUZ8xG/u+SfRbZSxvwUB
t/Qx6keY3Kd45t3Ho0nJr3Sw
-----END PRIVATE KEY-----`,
	"cert": `-----BEGIN CERTIFICATE-----
MIIERTCCAq2gAwIBAgIQOHMkg1DAYOcN8fS0sE3qvjANBgkqhkiG9w0BAQsFADBr
MR4wHAYDVQQKExVta2NlcnQgZGV2ZWxvcG1lbnQgQ0ExIDAeBgNVBAsMF2pvaGFu
bmVzQE1pc3Npb25Db250cm9sMScwJQYDVQQDDB5ta2NlcnQgam9oYW5uZXNATWlz
c2lvbkNvbnRyb2wwHhcNMTkwNjAxMDAwMDAwWhcNMjkxMjE2MTAyMDU4WjBLMScw
JQYDVQQKEx5ta2NlcnQgZGV2ZWxvcG1lbnQgY2VydGlmaWNhdGUxIDAeBgNVBAsM
F2pvaGFubmVzQE1pc3Npb25Db250cm9sMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEA1ZvZZXEhvKKRbpT+uCdd01pFfONJExGAqimMDNxIgPm92H9BZmSc
Ae/hDYLos+NBmcbW9LrNk59YCvzVoCiWaU8gHzFTWxBg1Tk4T8PUvrG2wqJf4v+A
EyuHE+wM+J92CdWru1Z2RcbCXZOhwO5KtLJz0A4HQYDvhTwsipwIEDF/HbNZhAYU
wpKFmcD0YqBIqpcyUSVZWo8yeh4UQ2K2zTjJt9iPQNfhXPWkKScEVRvk5YpkmSG0
zt4f0/aEegLn5ilMFbP8nq9kzwC0xDvthJNCv5dasavpgCGMdrNTxVm4/XLRY2sT
kioV3clXBJ8zjYorjLp9yYtz85q1am6WzQIDAQABo4GEMIGBMA4GA1UdDwEB/wQE
AwIFoDATBgNVHSUEDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMB8GA1UdIwQY
MBaAFJups8ufkPjmCBgODGdDt8U348v0MCsGA1UdEQQkMCKCFW5mYy1jYXNoLXN5
c3RlbS5sb2NhbIIJbG9jYWxob3N0MA0GCSqGSIb3DQEBCwUAA4IBgQCBXyjblcIX
iH3/vniEfK4ppJw80AQMwZVgzaciZITzE4XHvioZTCILP0gYGlGaJCQyWDucEY62
EHTzm2lMcoFMd+wzpSbsrbucy4197qOAWYzZDNuW+Jj9snAMdO5xtYYAUauf1Eku
Ox/If3uVGSFcqmB45j6lFyqFqNWNMbe91Hu9uW+B/IWoTMCOLzxJxRfabTY3goj2
Nyf5FinruLSZXhp27f7mkDzHzD+abXgWZIoVlDznlupZvtOroEpwhFvJYv0zAKPr
cLlCnhCkxHIS06gBIpirTHiSXkP/TlMEkgIlm+QETQrwgfYFV4JbRZJCmiUAROoe
/0shj2/Y40CLQ15GIh8keMiYEgjSydpUOBb703lhzqsB8O8EnFjMO/sxtMnh2/Gj
98zo760nqncgN/YzfuKxh/8QLAGSS3ZIf4AzFoHptdsqvQlu0ZK0oVv1/dlwfZS2
MKeTut5WZH44Yioy/spG2v5Bg4kGF+XGS+lKqKjC3vmo30yt9IVoHJw=
-----END CERTIFICATE-----`,
}

func CertFiles() (files []string, teardown func(), err error) {
	certFile, err := ioutil.TempFile("", "*.pem")
	if err != nil {
		return nil, nil, fmt.Errorf("could not create temp cert file: %w", err)
	}
	defer certFile.Close()

	err = ioutil.WriteFile(certFile.Name(), []byte(grpcCert["cert"]), 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("could not write temp cert file: %w", err)
	}

	keyFile, err := ioutil.TempFile("", "*.pem")
	if err != nil {
		return nil, nil, fmt.Errorf("could not create temp key file: %w", err)
	}
	defer certFile.Close()

	err = ioutil.WriteFile(keyFile.Name(), []byte(grpcCert["key"]), 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("could not write temp key file: %w", err)
	}

	return []string{certFile.Name(), keyFile.Name()}, func() {
		os.Remove(certFile.Name())
		os.Remove(keyFile.Name())
	}, nil
}
