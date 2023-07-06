# BLockchain-DNS
      
      
      -c '{"function": "CreateDomain","args": ["{\"IP\":\"192.168.0.1\",\"URL\":\"example.com\",\"IdentityProofs\":{\"idno\":\"865148558452\",\"comno\":\"65298549854568\"}}"]}'


    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "GetIPAddressByURL","Args":["example.com"]}'
