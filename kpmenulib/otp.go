package kpmenulib

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tobischo/gokeepasslib/v3"
)

// OTP constants
const (
	OTP          = "otp"
	TOTPSEED     = "TOTP Seed"
	TOTPSETTINGS = "TOTP Settings"
	TOTP         = "totp"
	OTPAUTH      = "otpauth"
)

// OTPAuth supports only TOTP (at the moment)
type OTPAuth struct {
	secret  []byte
	Type    string
	Account string
	Issuer  string
	Period  int
	Digits  int
	err     error
}

// OTPError is a structure that handle an error of otp
type OTPError struct {
	err error
}

func (o OTPError) Error() string {
	return o.err.Error()
}

// CreateOTP generates a time-sensitive TOTP code for a database entry.
//
// Modern versions of KeepassXC and Keepass2Android store this URL in the `otp` key. A historic version
// stored data ds:
//
//     TOTP Seed = SECRET
//     TOTP Settings = PERIOD;DIGITS
//
// If the `otp` key exists, it should be used and the TOTP values ignored; otherwise, the legacy values can
// be used.
//
// entry is the DB entry for which to generate a code; time is the Unix time to generate for the code --
// generally time.Now().Unix()
//
func CreateOTP(a gokeepasslib.Entry, time int64) (otp string, err error) {
	otpa, err := CreateOTPAuth(a)
	if err != nil {
		return "", err
	}
	return otpa.Create(time)
}

func (o OTPAuth) Create(time int64) (otp string, err error) {
	// Default T0 as per RFC 6238 is 0
	const T0 = 0
	m := getMessage(time, T0, o.Period)

	key, err := base32.StdEncoding.DecodeString(string(o.secret))
	if err != nil {
		return otp, fmt.Errorf("invalid key: %v", err)
	}
	hasher := hmac.New(sha1.New, key)
	_, err = hasher.Write(m)
	if err != nil {
		return otp, fmt.Errorf("failed create hash: %v", err)
	}

	h := hasher.Sum(nil)
	ofs := getOffset(h)
	r := int32(h[ofs]&0x7f)<<24 |
		int32(h[ofs+1])<<16 |
		int32(h[ofs+2])<<8 |
		int32(h[ofs+3])

	otp = fmt.Sprint(r % int32(pow(10, o.Digits)))
	if len(otp) != o.Digits {
		rpt := strings.Repeat("0", o.Digits-len(otp))
		otp = rpt + otp
	}

	return otp, nil
}

// getMessage constructs the message for HMAC with given params
func getMessage(t1 int64, t0, stepTime int) (message []byte) {
	if stepTime == 0 {
		return []byte("ERROR_ZERO_STEPTIME")
	}
	step := (t1 - int64(t0)) / int64(stepTime)
	message = make([]byte, 8)
	binary.BigEndian.PutUint64(message, uint64(step))
	return message
}

// pow returns x^y
func pow(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

// getOffset returns the offset from hash bytes as per https://tools.ietf.org/html/rfc4226#section-5.4
func getOffset(hash []byte) int {
	lastByte := hash[len(hash)-1]
	return int(lastByte & 0xf)
}

func CreateOTPAuth(a gokeepasslib.Entry) (otp OTPAuth, err error) {
	for _, vd := range a.Values {
		switch vd.Key {
		default:
			// Nothing
		case OTP:
			otp, err = parseOTPAuth(strings.TrimSpace(vd.Value.Content))
			if err != nil {
				return OTPAuth{}, OTPError{
					err: fmt.Errorf("invalid key: %v", err),
				}
			}
			return otp, nil
		case TOTPSEED:
			otp.Type = TOTP
			otp.secret = []byte(strings.TrimSpace(vd.Value.Content))
		case TOTPSETTINGS:
			parts := strings.Split(strings.TrimSpace(vd.Value.Content), ";")
			if len(parts) != 2 {
				return otp, OTPError{
					err: fmt.Errorf("wrong TOTP Settings format; expected `SECS;DIGS`, was %s", vd.Value.Content),
				}
			}
			refresh, err := strconv.Atoi(parts[0])
			if err != nil {
				return otp, OTPError{
					err: fmt.Errorf("wrong TOTP Settings format; expected `SECS;DIGS`, was %s", vd.Value.Content),
				}
			}
			otp.Period = refresh
			digits, err := strconv.Atoi(parts[1])
			if err != nil {
				return otp, OTPError{
					err: fmt.Errorf("wrong TOTP Settings format; expected `SECS;DIGS`, was %s", vd.Value.Content),
				}
			}
			otp.Digits = digits
		}
	}
	return otp, nil
}

// parseOTPAuth parses a Google Authenticator otpauth URL, which is used by
// both KeepassXC and Keepass2Android.
//
//     otpauth://TYPE/LABEL?PARAMETERS
//
// e.g., the KeepassXC format is
//
//     otpauth://totp/ISSUER:USERNAME?secret=SECRET&period=SECONDS&digits=D&issuer=ISSUER
//
// where TITLE is the record entry title, e.g. `github`; USERNAME is the entry
// user name, e.g. `xxxserxxx`; SECRET is the TOTP seed secret; SECONDS is the
// number of seconds between key refreshes, e.g. 30; D is the number of digits
// in the generated TOTP code, commonly `6`; and ISSUER is the TOTP issuer, e.g.
// `github`.
//
// The spec is at https://github.com/google/google-authenticator/wiki/Key-Uri-Format
func parseOTPAuth(s string) (OTPAuth, error) {
	otp := OTPAuth{}
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return otp, err
	}
	if u.Scheme != OTPAUTH {
		return otp, errors.New("invalid format; must start with otpauth://")
	}
	if u.Host != TOTP {
		return otp, errors.New("only totp is supported")
	}
	_, label := filepath.Split(u.Path)
	otp.Type = label
	parts := strings.Split(label, ":")
	if len(parts) > 2 {
		return otp, fmt.Errorf("invalid label %s", label)
	}
	if len(parts) == 1 {
		otp.Account = parts[0]
	} else {
		otp.Issuer = parts[0]
		otp.Account = parts[1]
	}
	ur, err := url.Parse(s)
	if err != nil {
		return otp, err
	}
	for k, vs := range ur.Query() {
		if len(vs) != 1 {
			return OTPAuth{}, OTPError{
				err: fmt.Errorf("invalid key, too many parameter values for %s", k),
			}
		}
		switch k {
		case "secret":
			otp.secret = []byte(vs[0])
		case "digits":
			otp.Digits, err = strconv.Atoi(vs[0])
			if err != nil {
				return otp, OTPError{err: err}
			}
		case "period":
			otp.Period, err = strconv.Atoi(vs[0])
			if err != nil {
				return otp, OTPError{err: err}
			}
		}
	}
	return otp, nil
}
