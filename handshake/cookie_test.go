package handshake

import(
	"testing"
	"fmt"
	"net"
)

func TestCookie1(t *testing.T) {
	cookieProtector, _ := CookieProtectorNew([]byte("quic cookie"))
	fmt.Println(cookieProtector.NewToken(&net.UDPAddr { IP: net.IPv4(127, 0, 0, 1), Port: 8080 }))
}

func TestCookie2(t *testing.T) {
	cookieProtector, _ := CookieProtectorNew([]byte("quic cookie"))
	b, _ := cookieProtector.NewToken(&net.UDPAddr { IP: net.IPv4(127, 0, 0, 1), Port: 8080 })

	fmt.Println(cookieProtector.DecodeToken(b))
}