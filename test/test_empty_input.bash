output=`aphid < /dev/null 2>&1`
if [[ "$output" != "" ]]; then
    echo "Expected empty output for empty input. Got '$output'"
    exit 1
fi

