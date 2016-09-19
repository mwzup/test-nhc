package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type FAQ struct {
	ID       bson.ObjectID `bson:"_id" json:"id"`
	Question string        `bson:"question" json:"question"`
	Answers  string        `bson:"answer" json:"answer"`
	Category string        `bson:"category" json:"category"`
}

// GetFaqs gets and returns all frequently asked questions
func GetFaqs(w http.ResponseWriter, r *http.Request) {
	db := GetDB(w, r)
	faqs, errM := FindFaqs(db)
	if errM != nil {
		HandleModelError(w, r, errM)
		return
	}

	b, _ := json.Marshal(faqs)
	ServeJSONArray(w, r, string(b), http.StatusOK)
}

// FindFaqs finds and returns faq collection to GetFaqs()
func FindFaqs(db *mgo.Database) (faqs []FAQ, errM *Error) {
	c := db.C("faqs")
	err := c.Find(nil).All(&faqs)
	if err != nil {
		errM = &Error{Reason: errors.New(fmt.Sprintf("Error retrieving faqs from DB: %s", err)), Internal: true}
		return
	}

	return
}

/****************************************************************/
// AddFaq /admin creates a new frequently asked question
func AddFaq(w http.ResponseWriter, r *http.Request) {
	// When do we need the token data for authentication?
	tokenData := GetToken(w, r)
	if tokenData == nil {
		return
	}

	if !IsAuthorized(w, r, "global_admin") {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var faq FAQ
	err := decoder.Decode(&faq)
	if err != nil {
		BR(w, r, errors.New(PARSE_ERROR), http.StatusBadRequest)
		return
	}

	// Validate faq data
	if faq.Question == "" || faq.Answer == "" || faq.Category == "" {
		BR(w, r, errors.New(MISSING_FIELDS_ERROR), http.StatusBadRequest)
		return
	}

	// 	// DO I need to create a FAQ here and then save it?

	// 	db := GetDB(w, r)
	// 	errM := CreateFAQ(db, faq.ID, faq.Question, faq.Answer, faq.Category)
	// 	if errM != nil {
	// 		HandleModelError(w, r, errM)
	// 		return
	// 	}

	// 	ServeJSON(w, r, &Response{"status": "FAQ successfully created."}, http.StatusOK)
	// }

	// Save faq
	db := GetDB(w, r)
	faq.ID = bson.NewObjectId()
	errM := faq.Save(db)
	if errM != nil {
		HandleModelError(w, r, errM)
		return
	}

	ServeJSON(w, r, &Response{"status": "FAQ successfully added."}, http.StatusOK)
}

//----------------------------------------------
func (f *FAQ) Save(db *mgo.Database) *Error {
	c := db.C("faqs")
	_, err := c.UpsertId(f.ID, bson.M{"$set": f})
	if err != nil {
		return &Error{Internal: true, Reason: errors.New(fmt.Sprintf("Error saving faq: %s\n",
			err))}
	}

	return nil
}

//---------------------------

/****************************************************************/
// UpdateFaq /admin updates an existing frequently asked question
/// func UpdateFaq(w http.ResponseWriter, r *http.Request) {

/****************************************************************/
// DeleteFaq /admin deletes an existing frequently asked question
// func DeleteFaq(w http.ResponseWriter, r *http.Request) {

/***************************************************************/

///
// ADD TO main.go ---------------------------
/*
api.HandleFunc("/faq", GetFaqs).Methods("GET") // All users
api.HandleFunc("admin/faq", AddFaq).Methods("POST")
api.HandleFunc("admin/faq", UpdateFaq).Methods("PUT")
api.HandleFunc("admin/faq/{id}", DeleteFaq).Methods("DELETE")
*/
