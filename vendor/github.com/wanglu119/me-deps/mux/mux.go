package mux

import (
	"net/http"
	
	"github.com/gorilla/mux"
)

type Router = mux.Router

func NewRouter() *Router{
	return mux.NewRouter()
}

func Vars(r *http.Request) map[string]string{
	return mux.Vars(r)
} 
