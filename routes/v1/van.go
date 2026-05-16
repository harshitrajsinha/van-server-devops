package routes

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/harshitrajsinha/van-server-devops/models"
	"github.com/harshitrajsinha/van-server-devops/routes"
	"github.com/harshitrajsinha/van-server-devops/service"
)

// Response type is declared in handler/utils.go
// VerifyVanRequestBody function is declared in handler/utils.go

type VanHandler struct {
	service service.VanServiceInterface
}

func NewVanHandler(service service.VanServiceInterface) *VanHandler {
	return &VanHandler{
		service: service,
	}
}

func (v *VanHandler) GetVanByID(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	ctx := r.Context()

	// Get van id
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if id is valid uuid
	result, _ := uuid.Parse(id)
	if result.Version() != 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Invalid engine ID"})
		log.Println("Invalid Van ID")
		return
	}

	// Get data from service layer
	resp, err := v.service.GetVanById(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	// Send response
	var respData []interface{}
	respData = append(respData, resp) // enclose data in an array
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(routes.Response{Code: http.StatusOK, Data: respData})
	log.Println("Van data populated successfully based on ID")
}

func (v *VanHandler) GetAllVan(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		var r interface{}
		if r = recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	ctx := r.Context()

	// Get data from service layer
	resp, err := v.service.GetAllVan(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(routes.Response{Code: http.StatusOK, Data: resp})
	log.Println("All van data populated successfully")
}

func (v *VanHandler) CreateVan(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		var r interface{}
		if r = recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	ctx := r.Context()

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}
	defer r.Body.Close()

	var vanReq models.Van
	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&vanReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Value type is incorrect"})
		log.Println(err)
		return
	}

	// Verify request field
	doesKeyExists := routes.VerifyVanRequestBody(body, 1)
	if !doesKeyExists {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Missing required fields"})
		log.Println("Missing required fields")
		return
	}

	// validate request body
	if err := models.ValidateVanReq(vanReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: err.Error()})
		log.Println(err)
		return
	}

	// Pass data to service layer to create van
	createdVan, err := v.service.CreateVan(ctx, &vanReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	// send response
	if createdVan > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusCreated, Message: "van data inserted into DB successfully!"})
		log.Println("van data inserted into DB successfully!")
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "No rows inserted - Possibly data already exists"})
		log.Println("No rows inserted - Possibly data already exists")
	}
}

func (v *VanHandler) UpdateVan(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		var r interface{}
		if r = recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	ctx := r.Context()

	// Get van id
	params := mux.Vars(r)
	id := params["id"]

	// Check if id is valid uuid
	result, _ := uuid.Parse(id)
	if result.Version() != 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Invalid van ID"})
		log.Println("Invalid van ID")
		return
	}
	defer r.Body.Close()

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	var vanReq models.Van
	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&vanReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Value type is incorrect"})
		log.Println(err)
		return
	}

	// Verify request field
	doesKeyExists := routes.VerifyVanRequestBody(body, 1)
	if !doesKeyExists {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Missing required fields"})
		log.Println("Missing required fields")
		return
	}

	// validate request body
	if err := models.ValidateVanReq(vanReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: err.Error()})
		log.Println(err)
		return
	}

	// Pass data to service layer to update van
	updatedVan, err := v.service.UpdateVan(ctx, id, &vanReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	if updatedVan > 0 {
		// data is updated successfully
		log.Println("Van data updated successfully!")
		// Get the updated result
		v.GetVanByID(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "No data present for provided Van ID or data already exists"})
		log.Println("value of updatedVan is ", updatedVan)
		return
	}
}

func (v *VanHandler) UpdateVanPartial(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		var r interface{}
		if r = recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	ctx := r.Context()

	// Get van id
	params := mux.Vars(r)
	id := params["id"]

	// Check if id is valid uuid
	result, _ := uuid.Parse(id)
	if result.Version() != 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Invalid van ID"})
		log.Println("Invalid van ID")
		return
	}
	defer r.Body.Close()

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	var vanReq models.Van
	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&vanReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Value type is incorrect"})
		log.Println(err)
		return
	}

	// Verify request field
	doesKeyExists := routes.VerifyVanRequestBody(body, 0)
	if !doesKeyExists {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Missing required fields"})
		log.Println("Missing required fields")
		return
	}

	// validate request body
	if err := models.ValidateVanPatchReq(body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: err.Error()})
		log.Println(err)
		return
	}

	// Pass data to service layer to update van
	updatedVan, err := v.service.UpdateVan(ctx, id, &vanReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	if updatedVan > 0 {
		// data is updated successfully
		log.Println("Van data updated successfully!")
		// Get the updated result
		v.GetVanByID(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "No data present for provided Van ID or data already exists"})
		log.Println("value of updatedVan is ", updatedVan)
		return
	}
}

func (v *VanHandler) DeleteVan(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		var r interface{}
		if r = recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	ctx := r.Context()

	// Get van id
	params := mux.Vars(r)
	id := params["id"]

	// Check if id is valid uuid
	result, _ := uuid.Parse(id)
	if result.Version() != 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Invalid van ID"})
		log.Println("Invalid van ID")
		return
	}

	// Pass data to service layer to delete van
	deletedVan, err := v.service.DeleteVan(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusInternalServerError, Message: "Error occured while deleting data"})
		panic(err)
	}

	if deletedVan > 0 {
		// data is deleted successfully
		w.WriteHeader(http.StatusNoContent)
		w.Header().Set("Content-Type", "application/json")
		log.Println("value of deletedVan is ", deletedVan)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "No data present for provided Van ID or data already deleted"})
		log.Println("value of deletedVan is ", deletedVan)
		return
	}

}
