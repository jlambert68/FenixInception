#!/usr/bin/env bash

cd main
go build -o IntegrationTest2
mv IntegrationTest2 ../integrationTests/integrationTest2/IntegrationTest2
cd ..
#chmod +x gateway2
./integrationTests/integrationTest2/IntegrationTest2 -gatewayUsedInIntegrationTest=true -datbasePath=integrationTests/integrationTest2/databases/IntegrationTest2_parentForToBeTested.db -logdPath=integrationTests/integrationTest2/logFiles/IntegrationTest2_parentForToBeTested.log -configPath=integrationTests/integrationTest2/configFiles/integrationTest2_gatewayConfig_parentForToBeTested.toml