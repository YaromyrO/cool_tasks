package models

import (
	"github.com/Nastya-Kruglikova/cool_tasks/src/database"
	"github.com/satori/go.uuid"
)

const (
	createTrip        = "INSERT INTO trips (user_id) VALUES ($1) RETURNING trip_id;"
	getTripIDByUserID = "SELECT trips.trip_id FROM trips WHERE trips.user_id = $1;"
	getTripsByTripID  = "SELECT trips.user_id FROM trips WHERE trip_id = $1;"
)

//Trip is a representation of Event Trip in DB
type Trip struct {
	TripID      uuid.UUID
	UserID      uuid.UUID
	Events      []Event
	Flights     []Flight
	Museums     []Museum
	Restaurants []Restaurant
	Hotels      []Hotel
	Trains      []Train
}

//CreateTrip creates Trip and saves it to DB
var CreateTrip = func(trip Trip) (uuid.UUID, error) {
	var id uuid.UUID
	err := database.DB.QueryRow(createTrip, trip.UserID).Scan(&id)

	return id, err
}

//GetTripsByTripID gets Trips from DB by tripID
var GetTripsByTripID = func(id uuid.UUID) (Trip, error) {

	var (
		trip        Trip
		err         error
		events      []Event
		flights     []Flight
		museums     []Museum
		hotels      []Hotel
		trains      []Train
		restaurants []Restaurant
	)

	trip.TripID = id

	events, err = GetEventsByTrip(id)
	if err != nil {
		return trip, err
	}
	trip.Events = events

	flights, err = GetFlightsByTrip(id)
	if err != nil {
		return trip, err
	}
	trip.Flights = flights

	museums, err = GetMuseumsByTrip(id)
	if err != nil {
		return trip, err
	}
	trip.Museums = museums

	hotels, err = GetHotelsByTrip(id)
	if err != nil {
		return trip, err
	}
	trip.Hotels = hotels

	trains, err = GetTrainFromTrip(id)
	if err != nil {
		return trip, err
	}
	trip.Trains = trains

	restaurants, err = GetRestFromTrip(id)
	if err != nil {
		return trip, err
	}
	trip.Restaurants = restaurants

	errDB := database.DB.QueryRow(getTripsByTripID, id).Scan(&trip.UserID)
	if err != nil {
		return trip, errDB
	}

	return trip, nil
}

//GetTripIDByUserID gets Trips from DB by userID
var GetTripIDByUserID = func(id uuid.UUID) ([]uuid.UUID, error) {

	var (
		tripIDs []uuid.UUID
		tripID  uuid.UUID
	)

	rows, err := database.DB.Query(getTripIDByUserID, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&tripID); err != nil {
			return nil, err
		}
		tripIDs = append(tripIDs, tripID)
	}

	return tripIDs, nil
}
