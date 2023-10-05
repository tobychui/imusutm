package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"imuslab.com/utm/pkg/utils"
)

func totpInit() error {
	if !utils.FileExists("./totp.json") {
		var exampleTotpEntry = TotpEntry{
			Name:   "Example",
			Secret: "JBSWY3DPEHPK3PXP",
			Link:   "example.com",
		}
		thisConfig := TotpConfig{
			Entries: []*TotpEntry{&exampleTotpEntry},
		}

		js, _ := json.MarshalIndent(thisConfig, "", " ")
		os.WriteFile("./totp.json", js, 0775)
	}

	//Load the config file
	configContent, err := os.ReadFile("./totp.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(configContent, &totpConfig)
	if err != nil {
		return err
	}

	//Config loaded
	return nil
}

func HandleTOTPUpdate(w http.ResponseWriter, r *http.Request) {
	//Generate a list of TOTP code from config
	codes := []*TotpCode{}
	for _, entry := range totpConfig.Entries {
		thisCode, validFor, err := GenerateTOTP(entry.Secret)
		if err != nil {
			codes = append(codes, &TotpCode{
				Name:     entry.Name,
				Code:     "000000",
				Link:     entry.Link,
				ValidFor: 0,
				Succ:     false,
			})
		} else {
			codes = append(codes, &TotpCode{
				Name:     entry.Name,
				Code:     thisCode,
				Link:     entry.Link,
				ValidFor: validFor,
				Succ:     true,
			})
		}
	}

	js, _ := json.Marshal(codes)
	utils.SendJSONResponse(w, string(js))
}

// Generate totp base on secret, return the code, remaining valid time in sec and err if any
func GenerateTOTP(secret string) (string, int, error) {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", 0, err
	}

	counter := make([]byte, 8)
	now := time.Now()
	binary.BigEndian.PutUint64(counter, uint64(now.Unix()/30))

	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(counter)
	hash := hmacSha1.Sum(nil)

	offset := hash[len(hash)-1] & 0xf
	truncatedHash := hash[offset : offset+4]

	code := binary.BigEndian.Uint32(truncatedHash)
	code &= 0x7FFFFFFF
	code %= 1000000

	remainingTime := 30 - (uint64(now.Unix()) % 30)

	return fmt.Sprintf("%06d", code), int(remainingTime), nil
}

/*
func main() {
	secret := "JBSWY3DPEHPK3PXP" // Replace this with your 6-digit TOTP secret
	totpCode, err := generateTOTP(secret)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("TOTP Code:", totpCode)
}
*/
