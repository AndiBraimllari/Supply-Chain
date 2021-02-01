This README file contains vital instructions for interacting with the project. Below are mentioned some transactions, interfaces and commands which fetch current products, request a new one and finally consume it. This document has been limited to technical aspects only given that the report will cover further concepts. Note that the chaincode is contained in the supply-chain folder. supply-chain/chaincode-go/chaincode/smartcontract.go has the main business logic of the demo.

transactions (additional transactions not mentioned in this simple demo):
	-request product (retailer asks for some products, pushes request to state)
	-manufacture product (producer checks to see whats in the todo list for creating products, then builds them and pushes them to the state)
	-consume product (retailer checks state for available products, consumes them)
	
custom interfaces:
	-Product {
			ID,
			Description,
			Status # REQUESTED, IN_PROGRESS, MANUFACTURED 
		}
=================================================
# Basic binaries for interacting with the chaincode
. ./addBinaries.sh

# Invoke chaincode as a retailer
. ./switchToOrg1.sh

# Manage the simple network
./network.sh down && ./network.sh up createChannel && ./network.sh deployCC

# Invoke the initialization of certain objects and state
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n supply-chain --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'

# Get all the products along with their statuses from the state
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n supply-chain --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"GetAllProducts","Args":[]}'

# Create the specified product with a REQUESTED status and push it in the state
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n supply-chain --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"RequestProduct","Args":["test_id", "shoes"]}'

# Invoke chaincode as a producer
. ./switchToOrg2.sh

# Fetch all REQUESTED-state products and push them back in the state with a MANUFACTURED state
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n supply-chain --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"ManufactureRequestedProducts","Args":[]}'

# Invoke chaincode as a retailer
. ./switchToOrg1.sh

# The specified product with the given ID is retrieved from the state and consumed
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n supply-chain --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"ConsumeProduct","Args":["test_id"]}'

