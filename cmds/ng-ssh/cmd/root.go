package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"net"
	"net/http"
	
	"github.com/spf13/cobra"
	"github.com/spf13/afero"
	"github.com/mitchellh/go-homedir"
	
	"github.com/wanglu119/me-deps/viper"
	"github.com/wanglu119/me-deps/mux"
	
	"github.com/wanglu119/ng-ssh/lib/config"
	"github.com/wanglu119/ng-ssh/lib/ssh"
	ngSshHttp "github.com/wanglu119/ng-ssh/lib/http"
)

var rootCmd = &cobra.Command {
	Use: "ng-ssh",
	Short: "ng-ssh",
	Run: func(cmd *cobra.Command, args []string) {
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
		fmt.Printf("server running: %s \n", listener.Addr().String())
		err = http.Serve(listener, r)
		if err != nil {
			log.Error(fmt.Sprintf("%v", err))
			os.Exit(-1)
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	
	flags := rootCmd.Flags()
	server := &config.Server{}
	server.Mode = config.SERVER_MODE_LOCAL
	flags.StringVarP(&server.Address, "address", "a", "0.0.0.0", "address to listen on")
	flags.IntVarP(&server.Port, "port", "p", 20019, "port to listen on")
	flags.StringVarP(&server.Username, "username", "", "admin", "username for verify whether you have permission to use ng-ssh, default(admin)")
	flags.StringVarP(&server.Password, "password", "", "admin", "password for verify whether you have permission to use ng-ssh, default(admin)")
	
	vper := viper.GetViper()
	vper.Set("server", server)
}

func initConfig() {
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
}

// Execute executes the command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		os.Exit(-1)
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

