#!/bin/bash

for i in {1..10000}; do
	./test_upload
    s3cmd get s3://test/test-upload --force && s3cmd rm s3://test/test-upload
    #s3cmd get s3://test/test-upload /tmp/ --force && s3cmd rm s3://test/test-upload
	if [[ $? -ne 0 ]]; then
		echo "failed"
		exit 0
	fi
    ./test_file
done

