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
// VerifyEngineRequestBody function is declared in handler/utils.go

type EngineHandler struct {
	service service.EngineServiceInterface
}

func NewEngineHandler(service service.EngineServiceInterface) *EngineHandler {
	return &EngineHandler{
		service: service,
	}
}

func (e *EngineHandler) GetEngineByID(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	ctx := r.Context()

	// Get engine id
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if id is valid uuid
	result, _ := uuid.Parse(id)
	if result.Version() != 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Invalid engine ID"})
		log.Println("Invalid Engine ID")
		return
	}

	// Get data from service layer
	resp, err := e.service.GetEngineByID(ctx, id)
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
	log.Println("Engine data populated successfully based on ID")
}

func (e *EngineHandler) GetAllEngine(w http.ResponseWriter, r *http.Request) {

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
	resp, err := e.service.GetAllEngine(ctx)
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
	log.Println("All engine data populated successfully")
}

func (e *EngineHandler) CreateEngine(w http.ResponseWriter, r *http.Request) {

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

	var engineReq models.Engine
	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&engineReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Value type is incorrect"})
		log.Println(err)
		return
	}

	// verify request field
	doesKeyExists := routes.VerifyEngineRequestBody(body, 1)
	if !doesKeyExists {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Missing required fields"})
		log.Println("Missing required fields")
		return
	}

	// validate request body
	if err := models.ValidateEngineReq(engineReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: err.Error()})
		log.Println(err)
		return
	}

	// Pass data to service layer to create engine
	createdEngine, err := e.service.CreateEngine(ctx, &engineReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	// send response
	if createdEngine > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusCreated, Message: "engine data inserted into DB successfully!"})
		log.Println("engine data inserted into DB successfully!")
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "No rows inserted - Possibly data already exists"})
		log.Println("No rows inserted - Possibly data already exists")
	}

}

func (e *EngineHandler) UpdateEngine(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		var r interface{}
		if r = recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	ctx := r.Context()

	// Get id
	params := mux.Vars(r)
	id := params["id"]

	// Check if id is valid uuid
	result, _ := uuid.Parse(id)
	if result.Version() != 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Invalid engine ID"})
		log.Println("Invalid Engine ID")
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

	var engineReq models.Engine
	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&engineReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Value type is incorrect"})
		log.Println(err)
		return
	}

	// verify request field
	doesKeyExists := routes.VerifyEngineRequestBody(body, 1)
	if !doesKeyExists {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Missing required fields"})
		return
	}

	// validate request body
	if err := models.ValidateEngineReq(engineReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: err.Error()})
		log.Println(err)
		return
	}

	// Pass data to service layer to update engine
	updatedEngine, err := e.service.UpdateEngine(ctx, id, &engineReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	if updatedEngine > 0 {
		// data is updated successfully
		log.Println("Engine data updated successfully!")
		// Get the updated result
		e.GetEngineByID(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "No data present for provided Engine ID or data already exists"})
		log.Println("value of updatedEngine is ", updatedEngine)
		return
	}
}

func (e *EngineHandler) UpdateEnginePartial(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		var r interface{}
		if r = recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	ctx := r.Context()

	// Get id
	params := mux.Vars(r)
	id := params["id"]

	// Check if id is valid uuid
	result, _ := uuid.Parse(id)
	if result.Version() != 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Invalid engine ID"})
		log.Println("Invalid Engine ID")
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

	var engineReq models.Engine
	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&engineReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Value type is incorrect"})
		log.Println(err)
		return
	}

	// verify request field
	doesKeyExists := routes.VerifyEngineRequestBody(body, 0)
	if !doesKeyExists {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Missing required fields"})
		return
	}

	// validate request body
	if err := models.ValidateEnginePatchReq(body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: err.Error()})
		log.Println(err)
		return
	}

	// Pass data to service layer to update engine
	updatedEngine, err := e.service.UpdateEngine(ctx, id, &engineReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	if updatedEngine > 0 {
		// data is updated successfully
		log.Println("Engine data updated successfully!")
		// Get the updated result
		e.GetEngineByID(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "No data present for provided Engine ID or data already exists"})
		log.Println("value of updatedEngine is ", updatedEngine)
		return
	}
}

func (e *EngineHandler) DeleteEngine(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		var r interface{}
		if r = recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	ctx := r.Context()
	params := mux.Vars(r)

	// Get id
	id := params["id"]

	// Check if id is valid uuid
	result, _ := uuid.Parse(id)
	if result.Version() != 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "Invalid engine ID"})
		log.Println("Invalid Engine ID")
		return
	}

	// Pass data to service layer to delete engine
	deletedEngine, err := e.service.DeleteEngine(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusInternalServerError, Message: "Error occured while deleting data"})
		panic(err)
	}

	if deletedEngine > 0 {
		// data is deleted successfully
		w.WriteHeader(http.StatusNoContent)
		w.Header().Set("Content-Type", "application/json")
		log.Println("value of deletedEngine is ", deletedEngine)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes.Response{Code: http.StatusBadRequest, Message: "No data present for provided Engine ID or data already deleted"})
		log.Println("value of deletedEngine is ", deletedEngine)
		return
	}

}
