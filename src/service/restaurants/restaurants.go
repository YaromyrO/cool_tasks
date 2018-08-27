//Package restaurants implements restaurant handlers
package restaurants

import (
	"github.com/Nastya-Kruglikova/cool_tasks/src/model"
	"github.com/Nastya-Kruglikova/cool_tasks/src/service/common"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"net/http"
	"encoding/json"
)

type successAdd struct {
	Status string `json:"message"`
}

type successDelete struct {
	Status string `json:"message"`
}

type tripRestaurant struct {
	RestaurantID string
	TripID       string
}

//GetRestaurantHandler used for getting restaurants
func GetRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if val, ok := query["id"]; ok {
		_, err := uuid.FromString(val[0])
		if err != nil {
			common.SendNotFound(w, r, "ERROR: Invalid ID", err)
			return
		}
		restaurants, err := model.GetRestaurants(query)

		if err != nil {
			common.SendNotFound(w, r, "ERROR: Can't get restaurant", err)
			return
		}
		restaurant := restaurants[0]
		common.RenderJSON(w, r, restaurant)
	}

	restaurants, err := model.GetRestaurants(query)

	if err != nil {
		common.SendNotFound(w, r, "ERROR: Can't get restaurants", err)
		return
	}

	common.RenderJSON(w, r, restaurants)
}

//AddRestaurantToTrip saves Restaurant to Trip
func AddRestaurantToTripHandler(w http.ResponseWriter, r *http.Request) {
	var newRestaurant tripRestaurant

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newRestaurant)

	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Can't decode JSON POST Body", err)
		return
	}

	restaurantID, err := uuid.FromString(newRestaurant.RestaurantID)
	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Wrong restaurant ID (can't convert string to uuid)", err)
		return
	}

	tripID, err := uuid.FromString(newRestaurant.TripID)
	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Wrong trip ID (can't convert string to uuid)", err)
		return
	}

	err = model.AddRestaurantToTrip(tripID, restaurantID)
	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Can't add new restaurant to trip", err)
		return
	}

	common.RenderJSON(w, r, successAdd{Status: "201 Created"})
}

//DeleteRestaurantHandler deletes Restaurant from DB
func DeleteRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	itemID, err := uuid.FromString(params["id"])

	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Wrong item ID (can't convert string to uuid)", err)
		return
	}

	err = model.DeleteRestaurant(itemID)

	if err != nil {
		common.SendNotFound(w, r, "ERROR: Can't delete this item", err)
		return
	}

	common.RenderJSON(w, r, successDelete{Status: "204 No Content"})
}

//GetRestaurantFromTrip gets Restaurant from Trip by tripID
func GetRestaurantFromTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	tripID, err := uuid.FromString(params["id"])
	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Wrong trip ID (can't convert string to uuid)", err)
		return
	}

	items, err := model.GetRestaurantsFromTrip(tripID)
	if err != nil {
		common.SendNotFound(w, r, "ERROR: Can't get restaurants by trip ID", err)
		return
	}

	common.RenderJSON(w, r, items)
}