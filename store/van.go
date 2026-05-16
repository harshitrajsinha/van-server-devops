package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/harshitrajsinha/van-server-devops/models"
)

type vanQueryResponse struct {
	VanID       uuid.UUID `json:"van-id"`
	Name        string    `json:"name"`
	Brand       string    `json:"brand"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	FuelType    string    `json:"fuel-type"`
	EngineID    string    `json:"engine-id"`
	Price       int64     `json:"price"`
	Image       string    `json:"image-url"`
	CreatedAt   string    `json:"-"`
	UpdatedAt   string    `json:"-"`
}

type VanStore struct {
	db *sql.DB
}

func NewVanStore(db *sql.DB) VanStore {
	return VanStore{db: db}
}

// Query for van price > or < or range

func (v VanStore) GetVanById(ctx context.Context, id string) (interface{}, error) {
	var queryData vanQueryResponse

	// DB transaction
	tx, err := v.db.BeginTx(ctx, nil)
	if err != nil {
		return vanQueryResponse{}, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println("Transaction rollback error: ", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				log.Println("Commit rollback error: ", cmErr)
			}
		}
	}()
	err = tx.QueryRowContext(ctx, "SELECT * FROM van WHERE van_id=$1", id).Scan(
		&queryData.VanID, &queryData.Name, &queryData.Brand, &queryData.Description, &queryData.Category, &queryData.FuelType, &queryData.EngineID, &queryData.Price, &queryData.Image, &queryData.CreatedAt, &queryData.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return vanQueryResponse{}, err // return empty model
		}
		return vanQueryResponse{}, err // return empty model
	}
	return queryData, nil
}

func (v VanStore) GetAllVan(ctx context.Context) (interface{}, error) {

	// DB transaction
	tx, err := v.db.BeginTx(ctx, nil)
	if err != nil {
		return vanQueryResponse{}, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println("Transaction rollback error: ", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				log.Println("Commit rollback error: ", cmErr)
			}
		}
	}()
	rows, err := tx.QueryContext(ctx, "SELECT * FROM van;")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return vanQueryResponse{}, err // return empty model
		}
		return vanQueryResponse{}, err // return empty model
	}
	defer rows.Close()

	// slice to store all rows
	vanData := make([]interface{}, 0)

	// Get each row data into a slice
	for rows.Next() {
		var queryData vanQueryResponse
		if err := rows.Scan(
			&queryData.VanID, &queryData.Name, &queryData.Brand, &queryData.Description, &queryData.Category, &queryData.FuelType, &queryData.EngineID, &queryData.Price, &queryData.Image, &queryData.CreatedAt, &queryData.UpdatedAt); err != nil {
			log.Fatal(err)
		}
		vanData = append(vanData, queryData)
	}

	return vanData, nil
}

func (v VanStore) CreateVan(ctx context.Context, vanReq *models.Van) (int64, error) {

	// DB transaction
	tx, err := v.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error while inserting data ", err)
		return -1, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println("Transaction rollback error: ", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				log.Println("Commit rollback error: ", cmErr)
			}
		}
	}()

	var query string = "INSERT INTO van (name, brand, description, category, fuel_type, engine_id, price, image_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	result, err := tx.ExecContext(ctx, query, vanReq.Name, vanReq.Brand, vanReq.Description, vanReq.Category, vanReq.FuelType, vanReq.EngineID, vanReq.Price, vanReq.ImageURL)

	if err != nil {
		log.Println("Error while inserting data ", err)
		return -1, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error while inserting data ", err)
		return -1, err
	}

	return rowsAffected, nil
}

func (v VanStore) UpdateVan(ctx context.Context, vanID string, vanReq *models.Van) (int64, error) {

	// DB transaction
	tx, err := v.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error while updating data ", err)
		return -1, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	var query strings.Builder
	var args []interface{}
	argCount := 1

	query.WriteString("UPDATE van SET ")

	if vanReq.Name != "" {
		query.WriteString(fmt.Sprintf("name=$%d ", argCount))
		args = append(args, vanReq.Name)
		argCount++
	}
	if vanReq.Brand != "" {
		if argCount > 1 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("brand=$%d ", argCount))
		args = append(args, vanReq.Brand)
		argCount++
	}
	if vanReq.Description != "" {
		if argCount > 1 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("description=$%d ", argCount))
		args = append(args, vanReq.Description)
		argCount++
	}
	if vanReq.Category != "" {
		if argCount > 1 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("category=$%d ", argCount))
		args = append(args, vanReq.Category)
		argCount++
	}
	if vanReq.FuelType != "" {
		if argCount > 1 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("fuel_type=$%d ", argCount))
		args = append(args, vanReq.FuelType)
		argCount++
	}
	if vanReq.EngineID.Version() == 4 {
		if argCount > 1 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("engine_id=$%d ", argCount))
		args = append(args, vanReq.EngineID)
		argCount++
	}
	if vanReq.Price > 0 {
		if argCount > 1 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("price=$%d ", argCount))
		args = append(args, vanReq.Price)
		argCount++
	}
	if vanReq.ImageURL != "" {
		if argCount > 1 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("image_url=$%d ", argCount))
		args = append(args, vanReq.ImageURL)
		argCount++
	}

	query.WriteString(fmt.Sprintf("WHERE van_id=$%d ", argCount))
	args = append(args, vanID)

	result, err := tx.ExecContext(ctx, query.String(), args...)
	if err != nil {
		log.Println("Error while updating data ", err)
		return -1, err
	}

	rowAffected, err := result.RowsAffected()
	return rowAffected, nil
}

func (v VanStore) DeleteVan(ctx context.Context, id string) (int64, error) {

	// DB transaction
	tx, err := v.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println("Transaction rollback error: ", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				log.Println("Commit rollback error: ", cmErr)
			}
		}
	}()

	var query string = "DELETE FROM van WHERE van_id=$1"
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return -1, err
	}
	rowAffected, err := result.RowsAffected()

	return rowAffected, nil

}
