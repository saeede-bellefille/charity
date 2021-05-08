package handler

import (
	"SABKAD/pkg/model"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Router is exported
func Router(db *gorm.DB) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	cbt := New(func() model.Model { return &model.CommonBaseType{} }, db, "CommonBaseType")
	router.HandleFunc("/CBT", cbt.Create).Methods("POST")
	router.HandleFunc("/CBT/{id:[0-9]+}", cbt.Update).Methods("PUT")
	router.HandleFunc("/CBT/{id:[0-9]+}", cbt.Delete).Methods("DELETE")
	router.HandleFunc("/CBT", cbt.Search).Methods("GET")

	cbd := New(func() model.Model { return &model.CommonBaseData{} }, db, "CommonBaseData")
	router.HandleFunc("/CBD", cbd.Create).Methods("POST")
	router.HandleFunc("/CBD/{id:[0-9]+}", cbd.Update).Methods("PUT")
	router.HandleFunc("/CBD/{id:[0-9]+}", cbd.Delete).Methods("DELETE")
	router.HandleFunc("/CBD", cbd.Search).Methods("GET")

	ca := New(func() model.Model { return &model.CharityAccount{} }, db, "CharityAccount")
	router.HandleFunc("/CA", ca.Create).Methods("POST")
	router.HandleFunc("/CA/{id:[0-9]+}", ca.Update).Methods("PUT")
	router.HandleFunc("/CA/{id:[0-9]+}", ca.Delete).Methods("DELETE")
	router.HandleFunc("/CA", ca.Search).Methods("GET")

	p := New(func() model.Model { return &model.Personal{} }, db, "Personal")
	router.HandleFunc("/P", p.Create).Methods("POST")
	router.HandleFunc("/P/{id:[0-9]+}", p.Update).Methods("PUT")
	router.HandleFunc("/P/{id:[0-9]+}", p.Delete).Methods("DELETE")
	router.HandleFunc("/P", p.Search).Methods("GET")

	na := New(func() model.Model { return &model.NeedyAccount{} }, db, "NeedyAccount")
	router.HandleFunc("/NA", na.Create).Methods("POST")
	router.HandleFunc("/NA/{id:[0-9]+}", na.Update).Methods("PUT")
	router.HandleFunc("/NA/{id:[0-9]+}", na.Delete).Methods("DELETE")
	router.HandleFunc("/NA", na.Search).Methods("GET")

	plan := New(func() model.Model { return &model.Plan{} }, db, "Plan")
	router.HandleFunc("/PLAN", plan.Create).Methods("POST")
	router.HandleFunc("/PLAN/{id:[0-9]+}", plan.Update).Methods("PUT")
	router.HandleFunc("/PLAN/{id:[0-9]+}", plan.Delete).Methods("DELETE")
	router.HandleFunc("/PLAN", plan.Search).Methods("GET")

	antp := New(func() model.Model { return &model.AssignNeedyToPlan{} }, db, "AssignNeedyToPlan")
	router.HandleFunc("/ANTP", antp.Create).Methods("POST")
	router.HandleFunc("/ANTP/{id:[0-9]+}", antp.Update).Methods("PUT")
	router.HandleFunc("/ANTP/{id:[0-9]+}", antp.Delete).Methods("DELETE")
	router.HandleFunc("/ANTP", antp.Search).Methods("GET")

	cad := New(func() model.Model { return &model.CashAssistanceDetail{} }, db, "CashAssistanceDetail")
	router.HandleFunc("/CAD", cad.Create).Methods("POST")
	router.HandleFunc("/CAD/{id:[0-9]+}", cad.Update).Methods("PUT")
	router.HandleFunc("/CAD/{id:[0-9]+}", cad.Delete).Methods("DELETE")
	router.HandleFunc("/CAD", cad.Search).Methods("GET")

	pay := New(func() model.Model { return &model.Payment{} }, db, "Payment")
	router.HandleFunc("/PAY", pay.Create).Methods("POST")
	router.HandleFunc("/PAY/{id:[0-9]+}", pay.Update).Methods("PUT")
	router.HandleFunc("/PAY/{id:[0-9]+}", pay.Delete).Methods("DELETE")
	router.HandleFunc("/PAY", pay.Search).Methods("GET")

	return router
}
