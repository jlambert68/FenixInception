#!/usr/bin/env bash

go build main/main.go -o main/IntegrationTest2
#chmod +x gateway2
./main/IntegrationTest2 -datbasePath=go_code/TestExecutionGateway/integrationTests/integrationTest2/databases