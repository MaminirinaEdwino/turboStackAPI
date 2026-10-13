package turbojwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// base64Encode allow to generate base64.URLEncoding
func base64Encode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// generateHmac generate a hashed data for the jwt signature
func generateHmac(secret string, unsignedToken string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(unsignedToken))
	return base64Encode(h.Sum(nil))
}

// Encode generate jwt token by using the secret and payload variable <br>
// exp(hour) : the token duration \n
// nbf(minutes) : 0 if don't want to use it \n
// nbf > if you want to use not before a duration
func Encode(secret string, payload map[string]interface{}, exp float64, nbf float64) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	expiration := time.Now().Add(time.Hour * time.Duration(exp)).Unix()
	iat := time.Now().Unix()

	payload["exp"] = expiration
	payload["iat"] = iat
	if nbf > 0 {
		notBefore := time.Now().Add(time.Minute * time.Duration(nbf)).Unix()
		payload["nbf"] = notBefore
	}
	headerJson, _ := json.Marshal(header)
	encodedHeader := base64Encode(headerJson)
	payloadJson, _ := json.Marshal(payload)
	encodedPayload := base64Encode(payloadJson)
	unsignedToken := encodedHeader + "." + encodedPayload
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(unsignedToken))
	signature := base64Encode(h.Sum(nil))
	return unsignedToken + "." + signature, nil
}

// Verify return the payload by using the secret string and the token | 
//
// secret : the secret key \n
// token : the hashed token \n
// leeway(second) : delay if > 0 == you have a delay \n
// if = 0, no delay
func Verify(secret string, token string, leeway float64) (map[string]any, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalide Token")
	}

	unsignedPart := parts[0] + "." + parts[1]
	signature := generateHmac(secret, unsignedPart)

	if !hmac.Equal([]byte(signature), []byte(parts[2])) {
		return nil, fmt.Errorf("Invalide signature.")
	}
	payloadPart, _ := base64.RawURLEncoding.DecodeString(parts[1])
	var payload map[string]any
	if err := json.Unmarshal(payloadPart, &payload); err != nil {
		return nil, err
	}
	exp, ok := payload["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("Invalide token, missing exp time")
	}
	if time.Now().Unix() > int64(exp) {
		return nil, fmt.Errorf("Expired Token")
	}
	now := time.Now().Unix()
	leewaySec := int64(time.Duration(leeway).Seconds())

	if exp, ok := payload["exp"].(float64); ok {
		if now > (int64(exp) + leewaySec) {
			return nil, fmt.Errorf("Expired token")
		}
	}

	if val, ok := payload["nbf"]; ok {
		nbf, ok := val.(float64)
		if !ok {
			return nil, fmt.Errorf("Token not yet active")
		}
		if time.Now().Unix() < (int64(nbf) - leewaySec) {
			return nil, fmt.Errorf("Inactive token")
		}
	}
	return payload, nil
}
