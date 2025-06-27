package cognito

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWK struct {
	Keys []JWKKey `json:"keys"`
}

type JWKKey struct {
	Alg string `json:"alg"`
	E   string `json:"e"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
	Use string `json:"use"`
}

type CognitoJWTValidator struct {
	userPoolID   string
	clientID     string
	region       string
	jwksURL      string
	jwkSet       *JWK
	jwkSetMutex  sync.RWMutex
	lastFetched  time.Time
	cacheDuration time.Duration
}

func NewCognitoJWTValidator() *CognitoJWTValidator {
	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
	region := os.Getenv("AWS_REGION")
	clientID := os.Getenv("COGNITO_CLIENT_ID")
	cognitoEndpoint := os.Getenv("COGNITO_ENDPOINT")

	var jwksURL string
	if cognitoEndpoint != "" {
		jwksURL = fmt.Sprintf("%s/%s/.well-known/jwks.json", cognitoEndpoint, userPoolID)
	} else {
		jwksURL = fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPoolID)
	}

	return &CognitoJWTValidator{
		userPoolID:    userPoolID,
		clientID:      clientID,
		region:        region,
		jwksURL:       jwksURL,
		cacheDuration: 1 * time.Hour,
	}
}

func (v *CognitoJWTValidator) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found")
		}

		keys, err := v.getJWKSet()
		if err != nil {
			return nil, err
		}

		for _, key := range keys.Keys {
			if key.Kid == kid {
				return v.convertJWKToRSAPublicKey(&key)
			}
		}

		return nil, fmt.Errorf("unable to find appropriate key")
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if err := v.validateClaims(claims); err != nil {
		return nil, err
	}

	return token, nil
}

func (v *CognitoJWTValidator) validateClaims(claims jwt.MapClaims) error {
	iss, ok := claims["iss"].(string)
	if !ok {
		return errors.New("iss claim not found")
	}

	expectedIss := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", v.region, v.userPoolID)
	if v.jwksURL != "" && v.jwksURL[:4] == "http" && v.jwksURL[:8] != "https://" {
		expectedIss = fmt.Sprintf("http://cognito-local:9229/%s", v.userPoolID)
	}

	if iss != expectedIss {
		return fmt.Errorf("invalid issuer: expected %s, got %s", expectedIss, iss)
	}

	tokenUse, ok := claims["token_use"].(string)
	if !ok {
		return errors.New("token_use claim not found")
	}

	if tokenUse != "access" && tokenUse != "id" {
		return fmt.Errorf("invalid token_use: %s", tokenUse)
	}

	if tokenUse == "access" {
		clientID, ok := claims["client_id"].(string)
		if !ok {
			return errors.New("client_id claim not found")
		}
		if clientID != v.clientID {
			return fmt.Errorf("invalid client_id: expected %s, got %s", v.clientID, clientID)
		}
	}

	if tokenUse == "id" {
		aud, ok := claims["aud"].(string)
		if !ok {
			return errors.New("aud claim not found")
		}
		if aud != v.clientID {
			return fmt.Errorf("invalid audience: expected %s, got %s", v.clientID, aud)
		}
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("exp claim not found")
	}

	if time.Now().Unix() > int64(exp) {
		return errors.New("token has expired")
	}

	return nil
}

func (v *CognitoJWTValidator) getJWKSet() (*JWK, error) {
	v.jwkSetMutex.RLock()
	if v.jwkSet != nil && time.Since(v.lastFetched) < v.cacheDuration {
		defer v.jwkSetMutex.RUnlock()
		return v.jwkSet, nil
	}
	v.jwkSetMutex.RUnlock()

	v.jwkSetMutex.Lock()
	defer v.jwkSetMutex.Unlock()

	if v.jwkSet != nil && time.Since(v.lastFetched) < v.cacheDuration {
		return v.jwkSet, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", v.jwksURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JWKS: %s", resp.Status)
	}

	var jwkSet JWK
	if err := json.NewDecoder(resp.Body).Decode(&jwkSet); err != nil {
		return nil, err
	}

	v.jwkSet = &jwkSet
	v.lastFetched = time.Now()

	return v.jwkSet, nil
}

func (v *CognitoJWTValidator) convertJWKToRSAPublicKey(jwk *JWKKey) (*rsa.PublicKey, error) {
	decodedE, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}

	decodedN, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}

	var e int
	if len(decodedE) == 3 {
		e = int(decodedE[0])<<16 | int(decodedE[1])<<8 | int(decodedE[2])
	} else {
		return nil, errors.New("invalid exponent")
	}

	n := new(big.Int).SetBytes(decodedN)

	return &rsa.PublicKey{
		E: e,
		N: n,
	}, nil
}

func (v *CognitoJWTValidator) GetUserInfoFromToken(token *jwt.Token) (map[string]interface{}, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userInfo := make(map[string]interface{})
	
	if sub, ok := claims["sub"].(string); ok {
		userInfo["user_id"] = sub
	}
	
	if email, ok := claims["email"].(string); ok {
		userInfo["email"] = email
	}
	
	if username, ok := claims["cognito:username"].(string); ok {
		userInfo["username"] = username
	}
	
	if groups, ok := claims["cognito:groups"].([]interface{}); ok {
		userInfo["groups"] = groups
	}

	return userInfo, nil
}