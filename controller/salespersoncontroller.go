package controller

import (
	"encoding/base64"
	"fmt"
	"go/build"
	"log"
	"net/http"
	"project/data/model"
	"project/data/service"
	"text/template"

	"github.com/gorilla/sessions"
)

var fm = template.FuncMap{
	"convstring": convertimagetostring,
}

func convertimagetostring(image []byte) string {
	imagestring := base64.StdEncoding.EncodeToString(image)
	return imagestring
}

var chekerror bool
var store = sessions.NewCookieStore([]byte("t0p-s3cr3t"))
var vehiclesave bool
var deletevehicle bool
var updatevehicle bool

//SalespersonIndexpageProcess is..
func SalespersonIndexpageProcess(w http.ResponseWriter, r *http.Request) {
	var message string
	var hasmessge bool
	if vehiclesave {
		vehiclesave = false
		hasmessge = true
		message = "Vehicle data stored successfully"
	}

	if deletevehicle {
		deletevehicle = false
		hasmessge = true
		message = "Vehicle deleted successfully"
	}

	if updatevehicle {
		updatevehicle = false
		hasmessge = true
		message = "Vehicle updated successfully"
	}
	vehicles := service.GetAllVehicle()
	path := build.Default.GOPATH + "/src/project/template/salesperson/*"
	tpl := template.Must(template.New("").Funcs(fm).ParseGlob(path))
	tpl.ExecuteTemplate(w, "index.html", struct {
		HasMessage bool
		Message    string
		Vehicles   []model.Vehicle
	}{hasmessge, message, vehicles})
}

//NotFound is...
func NotFound(w http.ResponseWriter, r *http.Request) {
	path := build.Default.GOPATH + "/src/project/template/salesperson/*"
	tpl := template.Must(template.New("").Funcs(fm).ParseGlob(path))
	tpl.ExecuteTemplate(w, "404.html", nil)
}

//Login is...
func Login(w http.ResponseWriter, r *http.Request) {
	var message string
	var hasmessge bool
	if chekerror {
		hasmessge = true
		message = "Username or Password is wrong"
		chekerror = false
	}
	path := build.Default.GOPATH + "/src/project/template/salesperson/*"
	tpl := template.Must(template.New("").Funcs(fm).ParseGlob(path))
	tpl.ExecuteTemplate(w, "login.html", struct {
		HasMessage bool
		Message    string
	}{hasmessge, message})
}

//LoginPost is...
func LoginPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := r.PostForm.Get("username")
	pass := r.PostForm.Get("password")
	if user == "jaimin@gmail.com" && pass == "1312" {
		session, _ := store.Get(r, "username")
		session.Values["username"] = user
		session.Save(r, w)
		//fmt.Fprintf(w, "username is save")
		http.Redirect(w, r, "/admin/vehicle", http.StatusSeeOther)
		return
	}
	chekerror = true
	//fmt.Fprintf(w, "username or password is wrong")
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

//Logout is....
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "username")
	session.Options.MaxAge = -1
	delete(session.Values, "username")
	session.Save(r, w)
	//fmt.Fprintf(w, "username is cleared")
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
	//return
}

//CreateVehicleform is....
func CreateVehicleform(w http.ResponseWriter, r *http.Request) {
	path := build.Default.GOPATH + "/src/project/template/salesperson/*"
	tpl := template.Must(template.New("").Funcs(fm).ParseGlob(path))
	tpl.ExecuteTemplate(w, "createvehicle.html", nil)
}

//Authentication is..
func Authentication(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "username")
		_, ok := session.Values["username"]
		if !ok {
			chekerror = true
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

//ServerError is...
func ServerError(w http.ResponseWriter, r *http.Request) {
	path := build.Default.GOPATH + "/src/project/template/salesperson/*"
	tpl := template.Must(template.New("").Funcs(fm).ParseGlob(path))
	tpl.ExecuteTemplate(w, "error.html", nil)
}

//SaveVehicle is....
func SaveVehicle(w http.ResponseWriter, r *http.Request) {
	err := service.SaveVehicle(r)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	vehiclesave = true
	http.Redirect(w, r, "/admin/vehicle", http.StatusSeeOther)
}

//GetoneVehicleforview is...
func GetoneVehicleforview(w http.ResponseWriter, r *http.Request) {
	vehicle := service.GetOneVehicle(r)
	path := build.Default.GOPATH + "/src/project/template/salesperson/*"
	tpl := template.Must(template.New("").Funcs(fm).ParseGlob(path))
	tpl.ExecuteTemplate(w, "viewvehicle.html", vehicle)
}

//ViewVehiclepage is...

//GetoneVehicleforedit is..
func GetoneVehicleforedit(w http.ResponseWriter, r *http.Request) {
	vehicle := service.GetOneVehicle(r)
	path := build.Default.GOPATH + "/src/project/template/salesperson/*"
	tpl := template.Must(template.New("").Funcs(fm).ParseGlob(path))
	tpl.ExecuteTemplate(w, "editvehicle.html", vehicle)
}

//DeleteVehicle is....
func DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	service.DeleteOneVehicle(r)
	deletevehicle = true
}

//UpdateVehicle is....
func UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call donw")
	data, err := service.UpdateVehicle(r)
	if err != nil {
		log.Fatalln(err)
		return
	}
	updatevehicle = true
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	//http.Redirect(w, r, "/admin/vehicle", http.StatusSeeOther)
}