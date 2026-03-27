# get all test functions in the scope (file/folder), 
# and run them sequentially, to avoid hitting Azure API rate limits

function main {
    echo "==> Running acceptance tests sequentially to avoid hitting Azure API rate limits..."
    files=$(find $TEST -type f -name "*_test.go")
    for file in $files; do
        echo "==> Running tests in $file..." 
        getFuncFromFile $file
        for func in "${functions[@]}"; do
            echo "==> Running $func in $file, command: go test -v -timeout 120m -run '^$func$' $TEST..."
            go test -v -timeout 120m -run ^$func$ $TEST
            echo "==> wait for 30 seconds before running the next test to avoid hitting Azure API rate limits..."
            sleep 30s
            #go test -v -timeout=$TESTTIMEOUT -run=^$func$ $file
        done
        # go test -v -timeout=$TESTTIMEOUT -run=TestAccBatchPool $file
    done
    exit 0
}

function getFuncFromFile {
  local file="$1"
  mapfile -t functions < <(
    awk -v run="$RUN" '
      match($0, /^[[:space:]]*func[[:space:]]+(\([^)]*\)[[:space:]]+)?([[:alnum:]_]+)[[:space:]]*\(/, m) {
        name = m[2]
        if (run == "" || index(name, run) > 0) {
          print name
        }
      }
    ' "$file"
  )
}
main