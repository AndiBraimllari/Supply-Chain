package chaincode

import (
	"encoding/json"
	"fmt"
	"time"
	"math/rand"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Product describes basic details of what makes up a simple product
type Product struct {
	ID             string `json:"ID"`
	Description    string `json:"description"`
	Status	       string `json:"status"` // REQUESTED, IN_PROGRESS, MANUFACTURED
}

// InitLedger adds a base set of products to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	products := []Product{
		{ID: "p1", Description: "apples", Status: "MANUFACTURED"},
		{ID: "p2", Description: "coffee", Status: "MANUFACTURED"},
		{ID: "p3", Description: "bag", Status: "MANUFACTURED"},
	}

	for _, product := range products {
		productJSON, err := json.Marshal(product)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(product.ID, productJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

func (s *SmartContract) RequestProduct(ctx contractapi.TransactionContextInterface, pId string, pDescription string) error {
	p := Product {
		ID: pId,
		Description: pDescription,
		Status: "REQUESTED",
	}

	pJSON, err := json.Marshal(p)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(pId, pJSON)
}

// UpdateProduct updates an existing product in the world state with provided parameters.
func UpdateProduct(ctx contractapi.TransactionContextInterface, id string, description string, status string) error {
	exists, err := ProductExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the product %s does not exist", id)
	}

	// overwriting original product with new product
	product := Product{
		ID:             id,
		Description:    description,
		Status:		status,
	}
	productJSON, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, productJSON)
}

// ReadProduct returns the product stored in the world state with given id.
func (s *SmartContract) ReadProduct(ctx contractapi.TransactionContextInterface, id string) (*Product, error) {
	productJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if productJSON == nil {
		return nil, fmt.Errorf("the product %s does not exist", id)
	}

	var product Product
	err = json.Unmarshal(productJSON, &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// DeleteProduct deletes an given product from the world state.
// ProductExists returns true when product with given ID exists in world state
func ProductExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	productJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return productJSON != nil, nil
}

// private
// _GetAllProducts returns all products found in world state
func _GetAllProducts(ctx contractapi.TransactionContextInterface) ([]*Product, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all products in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var products []*Product
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var product Product
		err = json.Unmarshal(queryResponse.Value, &product)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

// GetAllProducts returns all products found in world state
func (s *SmartContract) GetAllProducts(ctx contractapi.TransactionContextInterface) ([]*Product, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all products in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var products []*Product
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var product Product
		err = json.Unmarshal(queryResponse.Value, &product)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

// CreateProduct issues a new product to the world state with given details.
func (s *SmartContract) ManufactureRequestedProducts(ctx contractapi.TransactionContextInterface) error {
	products, err := _GetAllProducts(ctx)
	if err != nil {
		return err
	}

	for _, product := range products {
		if product.Status == "REQUESTED" {
			UpdateProduct(ctx, product.ID, product.Description, "IN_PROGRESS")
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second) // simulate the time it would take to manufacture a certain product
			UpdateProduct(ctx, product.ID, product.Description, "MANUFACTURED")
		}
	}
	return nil
}

func (s *SmartContract) ConsumeProduct(ctx contractapi.TransactionContextInterface, id string) error {
        exists, err := ProductExists(ctx, id)
        if err != nil {
                return err
        }
        if !exists {
                return fmt.Errorf("the product %s does not exist", id)
        }

	productJSON, err := ctx.GetStub().GetState(id)
	if productJSON == nil {
		return nil
	}

	// retailer got the product here, can utilize it how they see fit

	return ctx.GetStub().DelState(id)
}
