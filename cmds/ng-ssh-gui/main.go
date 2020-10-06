package main

import (
	"fmt"
	"os"
	"path/filepath"
	"net"
	"net/http"
	"time"
	"bytes"
	"strconv"
	
	"github.com/spf13/afero"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/mitchellh/go-homedir"
	
	"github.com/wanglu119/me-deps/viper"
	"github.com/wanglu119/me-deps/mux"

	"github.com/wanglu119/ng-ssh/lib/config"
	"github.com/wanglu119/ng-ssh/lib/ssh"
	ngSshHttp "github.com/wanglu119/ng-ssh/lib/http"
)

func init() {
	var home string
	var err error
	vper := viper.GetViper()
	
	home, err = homedir.Dir()
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		os.Exit(-1)
	}
	home = filepath.Join(home, ".ng-ssh")
	vper.AddConfigPath(home)
	vper.SetConfigName("ng-ssh")
	
	if err := vper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			panic(err)
		}
		
		osFs := afero.NewOsFs()
		if exist, _ := afero.DirExists(osFs, home); !exist {
			if err := osFs.Mkdir(home, 0777); err != nil {
				log.Error(fmt.Sprintf("%v", err))
				os.Exit(-1)
			}
		}
		cfgFile := filepath.Join(home, "ng-ssh.yaml")
		vper.SetConfigFile(cfgFile)
	}
	
	// set log file location
	logFile := filepath.Join(home, "logs/ng-ssh.log")
	log.SetFileLocation(logFile)
	
	server := config.GetServer()
	if server == nil {
		server = &config.Server{}
		vper := viper.GetViper()
		vper.Set("server", server)
	}
	server.Mode = config.SERVER_MODE_LOCAL
}

func startServer() {
	server := config.GetServer()
	
	vper := viper.GetViper()
	err := vper.WriteConfig()
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		os.Exit(-1)
	}
	
	r := mux.NewRouter()
	r.Use(CorsMiddleware)
	
	ngSshHttp.Setup(r)
	ssh.RunSshServer()
	
	var listener net.Listener
	listener, err = net.Listen("tcp", server.GetFullAddress())
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		os.Exit(-1)
	}
	defer listener.Close()
	log.Info(fmt.Sprintf("Listening on: %s", listener.Addr().String()))
	err = http.Serve(listener, r)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		os.Exit(-1)
	}
}

func checkServerStarted() bool {
	server := config.GetServer()
	addr := fmt.Sprintf("http://localhost:%d/ngssh_api/auth/login", server.Port)
	
	client := &http.Client{}
	
	req := fmt.Sprintf(`{"username":"%s", "password":"%s"}`, server.Username, server.Password)
    reqBuf := bytes.NewBuffer([]byte(req))
    request, err := http.NewRequest("POST", addr, reqBuf)
    if err != nil {
    	fmt.Println(err)
    	return false
    }
    request.Header.Set("Content-Type", "application/json")
    response, err := client.Do(request)
    if err != nil {
    	fmt.Println(err)
    	return false
    }
    
    if response.StatusCode == 200 {
    	return true
    } else {
    	return false
    }
}

func main() {
	app := app.New()

	w := app.NewWindow("ng-ssh")
	
	server := config.GetServer()
	
	usernameInput :=  widget.NewMultiLineEntry()
	usernameInput.MultiLine = false
	passwordInput := widget.NewPasswordEntry()
	serverPort :=  widget.NewMultiLineEntry()
	serverPort.MultiLine = false
	serverPort.Text = "20019"
	statusLable := widget.NewLabel("")
	
	form := widget.NewForm(
		widget.NewFormItem("username", usernameInput),
		widget.NewFormItem("password", passwordInput),
		widget.NewFormItem("ServerPort", serverPort),
	)
	
	isStart := false
	if server.Username != "" {
		isStart = checkServerStarted()
	}
	
	var start = func() {
		server.Username = usernameInput.Text
		server.Password = passwordInput.Text
		server.Port,_ = strconv.Atoi(serverPort.Text)
		if isStart { 
			statusLable.Text = "server has start"
			time.Sleep(1*time.Second)
			app.Quit()
		} else {
			statusLable.Text = "server starting..."
			time.Sleep(1*time.Second)
			w.Close()
		}
	}
	
	if isStart {
		statusLable.Text = fmt.Sprintf("server has start, port: %d, username: %s", server.Port, server.Username)
		w.SetContent(
			widget.NewVBox(
				statusLable,
			),
		)
		go func() {
			time.Sleep(3*time.Second)
			app.Quit()
		}()
	} else {
		usernameInput.Text = server.Username
		passwordInput.Text = server.Password
		serverPort.Text = fmt.Sprintf("%d", server.Port)
		
		w.SetContent(
			widget.NewVBox(
				form,
				widget.NewButton("Start", start),
				statusLable,
			),
		)
	}
	
	
	w.Resize(fyne.Size{450,200})
	
	w.CenterOnScreen()
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	
	w.ShowAndRun()
	if !isStart {
		startServer()
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Connection", "close")
    	w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE,PUT,PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Origin,X-Requested-With,Content-Type,Accept,Authorization,token,X-Auth,x-auth")
		
		if r.Method == http.MethodOptions {
			return
		}
        next.ServeHTTP(w, r)
    })
}
