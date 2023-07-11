package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Domain struct {
	IP             string            `json:"IP"`
	URL            string            `json:"URL"`
	IdentityProofs map[string]string `json:"IdentityProofs"`
}

type CreateDomainResponse struct {
	TxID   string `json:"txId"`
	Domain Domain `json:"domain"`
}

func (s *SmartContract) CreateDomain(ctx contractapi.TransactionContextInterface, domainJSON string) (*CreateDomainResponse, error) {
	var response CreateDomainResponse

	var domain Domain
	err := json.Unmarshal([]byte(domainJSON), &domain)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal domain JSON: %v", err)
	}

	exists, err := s.DomainExists(ctx, domain.IP)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("the domain already exists: %v", domain.IP)
	}

	domainBytes, err := json.Marshal(domain)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal domain: %v", err)
	}

	err = ctx.GetStub().PutState(domain.IP, domainBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to put domain state: %v", err)
	}

	response.TxID = ctx.GetStub().GetTxID()
	response.Domain = domain

	return &response, nil
}

func (s *SmartContract) DomainExists(ctx contractapi.TransactionContextInterface, ip string) (bool, error) {
	domainBytes, err := ctx.GetStub().GetState(ip)
	if err != nil {
		return false, fmt.Errorf("failed to read domain state: %v", err)
	}
	return domainBytes != nil, nil
}

func (s *SmartContract) GetIPAddressByURL(ctx contractapi.TransactionContextInterface, url string) (string, error) {
	// Query the domain data using the URL
	queryString := fmt.Sprintf(`{
		"selector": {
			"URL": "%s"
		},
		"use_index": ["_design/indexURLDoc", "indexURL"]
	}`, url)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return "", fmt.Errorf("failed to query domain data: %v", err)
	}
	defer resultsIterator.Close()

	// Iterate over the query results
	if resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return "", fmt.Errorf("failed to retrieve query result: %v", err)
		}

		// Extract the IP address from the query result
		var domain Domain
		err = json.Unmarshal(queryResult.Value, &domain)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal domain data: %v", err)
		}

		return domain.IP, nil
	}

	return "", fmt.Errorf("domain not found for URL: %s", url)
}

func (s *SmartContract) ReadDomain(ctx contractapi.TransactionContextInterface, ip string) (*Domain, error) {
	domainBytes, err := ctx.GetStub().GetState(ip)
	if err != nil {
		return nil, fmt.Errorf("failed to read domain state: %v", err)
	}
	if domainBytes == nil {
		return nil, fmt.Errorf("domain does not exist: %v", ip)
	}

	var domain Domain
	err = json.Unmarshal(domainBytes, &domain)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal domain: %v", err)
	}

	return &domain, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err.Error())
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err.Error())
	}
}
