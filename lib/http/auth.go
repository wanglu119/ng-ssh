package http

import (
	"net/http"
	"time"
	"fmt"
	"strings"
	"encoding/json"
	"bytes"
	"compress/gzip"
	"path"
	
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	
	"github.com/wanglu119/ng-ssh/lib/config"
	"github.com/wanglu119/ng-ssh/lib/assets"
)

// copy from ng-common/client-common/lib/http/auth.go
type Extractor []string

func (e Extractor) ExtractToken(r *http.Request) (string, error) {
	token, _ := request.HeaderExtractor{"X-Auth"}.ExtractToken(r)
	
	// Checks if the token isn't empty and if it contains two dots.
	// The former prevents incompatibility with URLs that previously
	// used basic auth.
	if token != "" && strings.Count(token, ".") == 2 {
		return token, nil
	}

	auth := r.URL.Query().Get("auth")
	if auth == "" {
		return "", request.ErrNoTokenInRequest
	}

	return auth, nil
}

func withUser(fn handleFunc) handleFunc {
	return func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		if d.server.Mode == config.SERVER_MODE_NG {
			var err error
			
			keyFunc := func(token *jwt.Token) (interface{}, error) {
				return []byte(d.server.Token), nil
			}
			
			var tk jwt.StandardClaims
			token, err := request.ParseFromRequest(r, &Extractor{}, keyFunc, request.WithClaims(&tk))
			
			if err != nil || !token.Valid {
				log.Error(fmt.Sprintf("%v", err))
				return http.StatusForbidden, nil
			}
			
			expired := !tk.VerifyExpiresAt(time.Now().Add(time.Hour).Unix(), true)
			
			if expired {
				w.Header().Add("X-Renew-Token", "true")
			}
		} else {
			var err error
			
			keyFunc := func(token *jwt.Token) (interface{}, error) {
				return []byte(d.server.Password), nil
			}
			
			var tk jwt.StandardClaims
			token, err := request.ParseFromRequest(r, &Extractor{}, keyFunc, request.WithClaims(&tk))
			
			if err != nil || !token.Valid {
				log.Error(fmt.Sprintf("%v", err))
				return http.StatusForbidden, nil
			}
			
			expired := !tk.VerifyExpiresAt(time.Now().Add(time.Hour).Unix(), true)
			
			if expired {
				w.Header().Add("X-Renew-Token", "true")
			}
		}
		
		return fn(w, r, d)
	}
}

func GenSigned(server *config.Server) (string, error) {
	if server.Mode == config.SERVER_MODE_NG {
		claims := jwt.StandardClaims {
			IssuedAt: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			Issuer: "ng-ssh",
		}
		
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(server.Token))
		
		return signed, err
	} else {
		claims := jwt.StandardClaims {
			IssuedAt: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			Issuer: "ng-ssh",
		}
		
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(server.Password))
		
		return signed, err
	}
	
}


var renew = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	signed, err := GenSigned(d.server)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, err
	}
	
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(signed))

	return 0, nil
})

// =================================================================================
type loginParam struct {
	Username string 
	Password string
}

var login = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	var param loginParam
	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		return http.StatusBadRequest, err
	}
	
	if param.Username == d.server.Username && param.Password == d.server.Password {
		signed, err := GenSigned(d.server)
		if err != nil {
			log.Error(fmt.Sprintf("%v", err))
			return http.StatusInternalServerError, err
		}
		
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(signed))
		
		return 0,nil
	} else {
		return http.StatusNotAcceptable, nil
	}
}

var indexHtml = `
<html>
<script>
window.location="%s"
</script>
</html>
`

var indexSsh = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	if d.server.Mode == config.SERVER_MODE_LOCAL {
		useGzip := false
		filePath := ""
		
		lastName := r.URL.Path[strings.LastIndex(r.URL.Path, "/"):]
		if(!strings.Contains(lastName, ".")) {
			filePath = "index.html"
			w.Header().Set("Content-Type", "text/html")
		} else {
			filePath = r.URL.Path[1:]
			
			if strings.HasSuffix(r.URL.Path[1:], ".css"){
				w.Header().Set("Content-Type", "text/css")
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Set("Accept-Encoding", "gzip")
				useGzip = true
			} else if strings.HasSuffix(r.URL.Path[1:], ".js"){
				w.Header().Set("Content-Type", "application/x-javascript")
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Set("Accept-Encoding", "gzip")
				useGzip = true
			} 
		}
		
		data,err := assets.Asset(path.Join("ng-ssh/dist",filePath))
		if err != nil {
			emsg := fmt.Sprintf("url %s, error: %v", r.URL.String(), err)
			log.Error(emsg)
			return http.StatusInternalServerError, err
		}
		if useGzip {
			var b bytes.Buffer
			gz := gzip.NewWriter(&b)
			
			if _, err := gz.Write(data); err != nil {
				emsg := fmt.Sprintf("index gz error: %v", err)
				log.Error(emsg)
				return http.StatusInternalServerError, err
			}
			
			if err := gz.Flush(); err != nil {
				emsg := fmt.Sprintf("index gz error: %v", err)
				log.Error(emsg)
				return http.StatusInternalServerError, err
			}
			
			if err := gz.Close(); err != nil {
				emsg := fmt.Sprintf("index gz error: %v", err)
				log.Error(emsg)
				return http.StatusInternalServerError, err
			}
			
			data = b.Bytes()
		}		
		w.Write(data)
	} else {
		login := fmt.Sprintf(indexHtml, "http://www.wl119.club:8090/ng-ssh")
		w.Write([]byte(login))
	}
	
	
	return 0,nil
}


var indexSftp = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	
	if d.server.Mode == config.SERVER_MODE_LOCAL {
		useGzip := false
		filePath := ""
		idx := len("/ng-sftp/")
		
		lastName := r.URL.Path[strings.LastIndex(r.URL.Path, "/"):]
		if(!strings.Contains(lastName, ".")) {
			filePath = "index.html"
			w.Header().Set("Content-Type", "text/html")
		} else {
			filePath = r.URL.Path[idx:]
			
			if strings.HasSuffix(r.URL.Path[idx:], ".css"){
				w.Header().Set("Content-Type", "text/css")
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Set("Accept-Encoding", "gzip")
				useGzip = true
			} else if strings.HasSuffix(r.URL.Path[idx:], ".js"){
				w.Header().Set("Content-Type", "application/x-javascript")
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Set("Accept-Encoding", "gzip")
				useGzip = true
			} 
		}
		
		filePath = path.Join("ng-sftp/dist",filePath)
		data,err := assets.Asset(filePath)
		if err != nil {
			emsg := fmt.Sprintf("url %s, error: %v", r.URL.String(), err)
			log.Error(emsg)
			return http.StatusInternalServerError, err
		}
		if useGzip {
			var b bytes.Buffer
			gz := gzip.NewWriter(&b)
			
			if _, err := gz.Write(data); err != nil {
				emsg := fmt.Sprintf("index gz error: %v", err)
				log.Error(emsg)
				return http.StatusInternalServerError, err
			}
			
			if err := gz.Flush(); err != nil {
				emsg := fmt.Sprintf("index gz error: %v", err)
				log.Error(emsg)
				return http.StatusInternalServerError, err
			}
			
			if err := gz.Close(); err != nil {
				emsg := fmt.Sprintf("index gz error: %v", err)
				log.Error(emsg)
				return http.StatusInternalServerError, err
			}
			
			data = b.Bytes()
		}		
		w.Write(data)
	} else {
		login := fmt.Sprintf(indexHtml, "http://www.wl119.club:8090/ng-sftp")
		w.Write([]byte(login))
	}
	
	return 0,nil
}






