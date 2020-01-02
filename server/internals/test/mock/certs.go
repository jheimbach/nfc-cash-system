package mock

import (
	"io/ioutil"
	"os"
	"testing"
)

var grpcCert = map[string]string{
	"key": `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA0gPH9kM1kBoE78f6uF9UPmeRh2nUI4NY4mEY4KwRy/tnhVNO
XJJQCIrRE7Ddz6tbG3b4owll/L89IJwajZenpwothNFeKJmwPOvq63q1Noj3VSEV
d7xD2G8vG3BRK/h8JJ7WzoYzC5fZtU6XhZBpvQWqs3SBBygjjO6PzSsmvEJZESy1
mzcOTO/Aga/j5dfVVX8iSco5cZfPZjP05B0DzSg9ht3S2HNWkPdFyAqgWm0+/HiW
wCn16HOmUrKZxEcvW6iEV42omW94AJDEjugCWswEV35B3TkMtgQcwdZHyxgF3H6Q
a/nyoKnx5Ho1C17dHhqM59NwHoqnpj7HG/6H8wIDAQABAoIBADqKS0bd3SRZ3F5q
Q/z4tabff7VbToLHrhMMNb8Kt7tATM7hNcqgDicTaswuVOX6Qd8Z/pyTlhYpyIQ6
fzxQta7eK2oGYlS24mVY7ZbOEY0uCKN2IWHK3K6L73maiUXUceZTUFUpGzl0Dn93
dM0KQC8sUTTMmNrB2YcDSJpMGPQImRDSjC8WoJyw9xag+gx2/CHNyZSPzMeV3MS0
fKAbGqgj7o73wvoufYRrh/Ach9ZDkWZFn3EGGx8AKXLF9vYA3/0WR+4X+eXjx7M5
Sxvxu4HjSH4uZKi/HtKQ98eVQHSxHMTppgm8nY+abVFDKFX5dWjK0Z++fpJdRcWl
/K2mzCECgYEA9HDHbxd0QZWUa8SbvkRBd+t0t8APB/NdmeHDxQXPV4o/etmNnoXP
eWCYatobDTNy6IqDuy5wMiRrYG6hKMS+JJ8FWzxic8kaQV2gCXeDQ2WwkdN4v7Bs
Bbcjv5layv3gbYVQgCuPri3xzjO2JMyUSk25gIqmgxS0XbZXZyKWyQcCgYEA2/I9
hNL7eYqMGfF9EjdZl/rbrO3v3ZahSbeKy0x0MxUh9D6S02CwbKWtcDUBkafa5Vua
wSFtXFR+ipKEcMZACATNsQhV+dGK5CN1faN2IWTlfV6mHs2C/FJYHntQ8RZKPZ3G
Q/+UQwsv6iCIEB4Hyx30fklUos98mZr8I3DO6rUCgYEA7NbVNVtJRj5y91Qg/uJN
eK7HgT5ykeaIO5AWyjBN7GjEBvkBkaXfF6CzLVy0Nz8xSATljBh7lunYrC+ksMan
4P2/B95jGgKxEMJxTJrisQu3YCPA7CI5F/SRi5Q/90yzBgrUq8sJRN+5WWybP96E
k9XpNZWhroICHIaO+xv1c0MCgYATPYSeKuquviz1VAex4T+oKNywqvvRsYyYa0Lr
99suYMngmNy1Ov8T8gZTC4AAouNmLvZBsM/lRMrxClIln2IYkXsA4o7K1MbGoEd1
3yfFUhM1PWNgzG+J7RYiTH1PKbkC3NtsNV5d4wuk+oWMi2P5I7ywk2+g8m+e7Ezk
OMkRBQKBgQDg5Yr9FdDUdL3SuMlfd1AxZGXB8FBx88a9tslGPDcHLGQ0ruTzoSxo
5nrRJ0jhtFMtf5aO6qiCsZUQckRCrnz77G7wjVB6m6C0UTC3ry7Yq+ZUIzuWNfd3
fc+MHSZ8htqAzwloKQd7nbwgQtPemEkumwQQADjIbNNrC8hFLYVnpg==
-----END RSA PRIVATE KEY-----`,
	"cert": `-----BEGIN CERTIFICATE-----
MIIC+zCCAeOgAwIBAgIJAJ4Oj7+Y+TFgMA0GCSqGSIb3DQEBBQUAMBQxEjAQBgNV
BAMMCWxvY2FsaG9zdDAeFw0yMDAxMDIxMTQwMzhaFw0yOTEyMzAxMTQwMzhaMBQx
EjAQBgNVBAMMCWxvY2FsaG9zdDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBANIDx/ZDNZAaBO/H+rhfVD5nkYdp1CODWOJhGOCsEcv7Z4VTTlySUAiK0ROw
3c+rWxt2+KMJZfy/PSCcGo2Xp6cKLYTRXiiZsDzr6ut6tTaI91UhFXe8Q9hvLxtw
USv4fCSe1s6GMwuX2bVOl4WQab0FqrN0gQcoI4zuj80rJrxCWREstZs3DkzvwIGv
4+XX1VV/IknKOXGXz2Yz9OQdA80oPYbd0thzVpD3RcgKoFptPvx4lsAp9ehzplKy
mcRHL1uohFeNqJlveACQxI7oAlrMBFd+Qd05DLYEHMHWR8sYBdx+kGv58qCp8eR6
NQte3R4ajOfTcB6Kp6Y+xxv+h/MCAwEAAaNQME4wHQYDVR0OBBYEFH7tHoS/QDTj
jPnIXKRjTPLQow2MMB8GA1UdIwQYMBaAFH7tHoS/QDTjjPnIXKRjTPLQow2MMAwG
A1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAHOOoSsAuvCbtVUdIVX3sAe3
aycyylhiFbYaaY7+XC5q0pRJwcuPdZXNqdUSVQZfSs9hrdMiwG3IUIL/Qs65oAfi
oJxaC3FlAyinYDoOKHoyUmeyFsLGK0ethpgMm699LPvK1r5lBQWHrTa8bCiOyMEF
FFU+gsBnSqZE87lpGsqfADNKWMObK77EAa3vvEbhUoMkkJuJ+Ao3HeCsDiTizkme
DHaz+hm0qiTIeTg5/jMJCrYQnIDDM7irWdZqYwCvECLQlewbKfV76tl6nAkqN3GT
8QVL3aTe6BJ2Icarre7hXwrT2teJS43LCYV0WcBkP8Rui9MsaTiy4RkpSltIKL8=
-----END CERTIFICATE-----`,
}

func CertFiles(t *testing.T) (files []string, teardown func()) {
	certFile, err := ioutil.TempFile("", "*.pem")
	if err != nil {
		t.Skipf("could not create temp file: %v", err)
	}
	defer certFile.Close()

	err = ioutil.WriteFile(certFile.Name(), []byte(grpcCert["cert"]), 0666)
	if err != nil {
		t.Skipf("could not create temp file: %v", err)
	}

	keyFile, err := ioutil.TempFile("", "*.pem")
	if err != nil {
		t.Skipf("could not create temp file: %v", err)
	}
	defer certFile.Close()

	err = ioutil.WriteFile(keyFile.Name(), []byte(grpcCert["key"]), 0666)
	if err != nil {
		t.Skipf("could not create temp file: %v", err)
	}

	return []string{certFile.Name(), keyFile.Name()}, func() {
		os.Remove(certFile.Name())
		os.Remove(keyFile.Name())
	}
}
