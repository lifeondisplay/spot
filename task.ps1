function Build {
    # windows
    go build -o "./bin/spot.exe" "./src/spot.go"
}

function BuildAll {
    go build -o "./bin/spot.exe" "./src/spot.go"
}