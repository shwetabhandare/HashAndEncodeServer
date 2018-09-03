#!/bin/bash

#Test 1: Verify the return code is greater than 0.
STATUS=$(curl -d "password=happyMonkey" http://localhost:8080/hash)
if [[ STATUS -ne 0 ]]; then
	echo “Test 1: Success”
else
	echo “Test 1: Failed”
fi

#Test 2: Re-run the same command. Verify the return code is not 1.
STATUS=$(curl -d "password=happyMonkey2" http://localhost:8080/hash)
if [[ STATUS -ne 1 ]]; then
	echo Test 2: “Success”
else
	echo Test 2: “Failed”
fi

echo "Sleeping for 5 seconds to complete hashing of password 1."
sleep 5

#Test 3: Validate the hash value
HASH=$(curl http://localhost:8080/hash/1)
EXPECTED="ZD5ZJT9f3qy02B3TwV8XTPX3I1qmODxU6eGTQHdrIzWcaMoeWyA68Uu+I0xrIGdxlJEkaUNM+Nl8pdPM4au6mg=="
echo $STATUS
if [ "$HASH" = "$EXPECTED" ]; then
	echo “Test 3: Success”
else
	echo Test 3: “Failed”
fi

#Test 4: Validate the hash value for a non-existent hash returns an error.
HASH=$(curl http://localhost:8080/hash/100)
EXPECTED="ERROR: Password Hash for request id: 100 does not exist."
echo $HASH
if [ "$HASH" = "$EXPECTED" ]; then
	echo Test 4: “Success”
else
	echo Test 4: “Failed”
fi

#Test 5: Hash a new value. Verify that the hash isn't computed for 5 seconds.
STATUS=$(curl -d "password=angryMonkey" http://localhost:8080/hash)
if [[ STATUS -ne 1 ]]; then
        echo “Test 5: Success”
else
        echo “Test 5: Failed”
fi

echo "Sleeping for 2 seconds..."
sleep 2
HASH=$(curl http://localhost:8080/hash/$STATUS)
EXPECTED="ERROR: Password Hash for request id: $STATUS does not exist."
echo $HASH
if [[ "$HASH" = *"ERROR: Password Hash for request id"* ]]; then
        echo “Test 5.1: Success”
else
        echo “Test 5.1: Failed”
fi
echo "Sleeping for 3 seconds..."
sleep 3
echo "Checking hash for request number: $STATUS..."
HASH=$(curl http://localhost:8080/hash/$STATUS)
echo $HASH
if [ "$HASH" != *"ERROR"* ]; then
	echo “Test 5.2: Success”
else
	echo “Test 5.2: Failed”
fi
