package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	secure "gopkg.in/unrolled/secure.v1"
)

func main() {
	var tlsPort string = os.Getenv("TLS_PORT")
	var useTls bool = len(tlsPort) > 0
	engine := GetEngine(useTls)
	if useTls {
		_ = engine.RunTLS(":"+tlsPort, os.Getenv("TLS_CERT_PATH"), os.Getenv("TLS_KEY_PATH"))
	} else {
		_ = engine.Run()
	}
}

func GetEngine(useTls bool) *gin.Engine {
	engine := gin.Default()
	engine.Use(CORSMiddleware())

	if useTls {
		engine.Use(Https())
	}

	v0 := engine.Group("/v0")
	{
		v0.GET("/integrations/figma/token", FigmaTokenExchangeHandler)
		v0.GET("/ping", func(ctx *gin.Context) {
			ctx.String(StatusOK, "pong")
		})
	}

	return engine // listen on PORT specified in env
}

func FigmaTokenExchangeHandler(ctx *gin.Context) {
	queryString := url.Values{}
	queryString.Set("client_id", os.Getenv("FIGMA_CLIENT_ID"))
	queryString.Set("client_secret", os.Getenv("FIGMA_CLIENT_SECRET"))
	queryString.Set("redirect_uri", os.Getenv("FIGMA_REDIRECT_URI"))
	queryString.Set("code", ctx.Query("Code"))
	queryString.Set("grant_type", "authorization_code")
	tokenExchangeUrl := fmt.Sprintf("%s?%s", os.Getenv("FIGMA_TOKEN_EXCHANGE_ENDPOINT"), queryString.Encode())
	response, err := http.Post(tokenExchangeUrl, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil || response.StatusCode != StatusOK {
		ctx.AbortWithStatus(StatusBadRequest)
		return
	}
	responseBody, errBody := ioutil.ReadAll(response.Body)
	if errBody != nil {
		ctx.AbortWithStatus(StatusBadRequest)
		return
	}
	accessTokenResponseRaw := new(OAuthAccessTokenResponseRaw)
	err = json.Unmarshal([]byte(responseBody), accessTokenResponseRaw)
	if err != nil {
		ctx.AbortWithStatus(StatusBadRequest)
		return
	}
	ctx.JSON(StatusOK, OAuthAccessTokenResponse(*accessTokenResponseRaw))
}

type OAuthAccessTokenResponseRaw struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    uint   `json:"expires_in"`
}

type OAuthAccessTokenResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    uint
}

type JSONWebTokenResponse struct {
	AccessToken string
}

//Support for TLS/HTTPS (must also specify relevant env)
func Https() gin.HandlerFunc {
	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:            []string{"*"},
		SSLRedirect:             true,
		SSLTemporaryRedirect:    false,
		SSLHost:                 "",
		SSLProxyHeaders:         map[string]string{},
		STSSeconds:              0,
		STSIncludeSubdomains:    false,
		STSPreload:              false,
		ForceSTSHeader:          false,
		FrameDeny:               false,
		CustomFrameOptionsValue: "",
		ContentTypeNosniff:      false,
		BrowserXssFilter:        false,
		ContentSecurityPolicy:   "",
		PublicKey:               "",
		IsDevelopment:           gin.Mode() == "debug",
	})

	return func(ctx *gin.Context) {
		err := secureMiddleware.Process(ctx.Writer, ctx.Request)

		// If there was an error, do not continue.
		if err != nil {
			ctx.String(StatusInternalServerError, "error securing session")
			ctx.Abort()
			return
		}

		// Avoid header rewrite if response is a redirection.
		if status := ctx.Writer.Status(); status > 300 && status < 399 {
			ctx.Abort()
		}
	}

}

//simple, permissive CORS config
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(StatusNoContent)
			return
		}

		c.Next()
	}
}

//Error codes
//copied from iris webserver project (Apache license)
const (
	// StatusContinue http status '100'
	StatusContinue = 100
	// StatusSwitchingProtocols http status '101'
	StatusSwitchingProtocols = 101
	// StatusOK http status '200'
	StatusOK = 200
	// StatusCreated http status '201'
	StatusCreated = 201
	// StatusAccepted http status '202'
	StatusAccepted = 202
	// StatusNonAuthoritativeInfo http status '203'
	StatusNonAuthoritativeInfo = 203
	// StatusNoContent http status '204'
	StatusNoContent = 204
	// StatusResetContent http status '205'
	StatusResetContent = 205
	// StatusPartialContent http status '206'
	StatusPartialContent = 206
	// StatusMultipleChoices http status '300'
	StatusMultipleChoices = 300
	// StatusMovedPermanently http status '301'
	StatusMovedPermanently = 301
	// StatusFound http status '302'
	StatusFound = 302
	// StatusSeeOther http status '303'
	StatusSeeOther = 303
	// StatusNotModified http status '304'
	StatusNotModified = 304
	// StatusUseProxy http status '305'
	StatusUseProxy = 305
	// StatusTemporaryRedirect http status '307'
	StatusTemporaryRedirect = 307
	// StatusBadRequest http status '400'
	StatusBadRequest = 400
	// StatusUnauthorized http status '401'
	StatusUnauthorized = 401
	// StatusPaymentRequired http status '402'
	StatusPaymentRequired = 402
	// StatusForbidden http status '403'
	StatusForbidden = 403
	// StatusNotFound http status '404'
	StatusNotFound = 404
	// StatusMethodNotAllowed http status '405'
	StatusMethodNotAllowed = 405
	// StatusNotAcceptable http status '406'
	StatusNotAcceptable = 406
	// StatusProxyAuthRequired http status '407'
	StatusProxyAuthRequired = 407
	// StatusRequestTimeout http status '408'
	StatusRequestTimeout = 408
	// StatusConflict http status '409'
	StatusConflict = 409
	// StatusGone http status '410'
	StatusGone = 410
	// StatusLengthRequired http status '411'
	StatusLengthRequired = 411
	// StatusPreconditionFailed http status '412'
	StatusPreconditionFailed = 412
	// StatusRequestEntityTooLarge http status '413'
	StatusRequestEntityTooLarge = 413
	// StatusRequestURITooLong http status '414'
	StatusRequestURITooLong = 414
	// StatusUnsupportedMediaType http status '415'
	StatusUnsupportedMediaType = 415
	// StatusRequestedRangeNotSatisfiable http status '416'
	StatusRequestedRangeNotSatisfiable = 416
	// StatusExpectationFailed http status '417'
	StatusExpectationFailed = 417
	// StatusTeapot http status '418'
	StatusTeapot = 418
	// StatusPreconditionRequired http status '428'
	StatusPreconditionRequired = 428
	// StatusTooManyRequests http status '429'
	StatusTooManyRequests = 429
	// StatusRequestHeaderFieldsTooLarge http status '431'
	StatusRequestHeaderFieldsTooLarge = 431
	// StatusUnavailableForLegalReasons http status '451'
	StatusUnavailableForLegalReasons = 451
	// StatusInternalServerError http status '500'
	StatusInternalServerError = 500
	// StatusNotImplemented http status '501'
	StatusNotImplemented = 501
	// StatusBadGateway http status '502'
	StatusBadGateway = 502
	// StatusServiceUnavailable http status '503'
	StatusServiceUnavailable = 503
	// StatusGatewayTimeout http status '504'
	StatusGatewayTimeout = 504
	// StatusHTTPVersionNotSupported http status '505'
	StatusHTTPVersionNotSupported = 505
	// StatusNetworkAuthenticationRequired http status '511'
	StatusNetworkAuthenticationRequired = 511
)
