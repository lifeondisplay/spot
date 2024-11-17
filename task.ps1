$curVer = [regex]::Match((Get-Content ".\src\spot.go"), "version = `"([\d\.]*)`"").Captures.Groups[1].Value
Write-Host "versão atual: $curVer"

function BumpVersion {
    param (
        [Parameter(Mandatory=$true)][int16]$major,
        [Parameter(Mandatory=$true)][int16]$minor,
        [Parameter(Mandatory=$true)][int16]$patch
    )

    $ver = "$($major).$($minor).$($patch)"

    (Get-Content ".\src\spot.go") -replace "version = `"[\d\.]*`"", "version = `"$($ver)`"" |
        Set-Content ".\src\spot.go"
}

function Dist {
    param (
        [Parameter(Mandatory=$true)][int16]$major,
        [Parameter(Mandatory=$true)][int16]$minor,
        [Parameter(Mandatory=$true)][int16]$patch
    )

    BumpVersion $major $minor $patch

    $nameVersion="spot-$($major).$($minor).$($patch)"
    $env:GOARCH="amd64"

    if (Test-Path "./bin") {
        Remove-Item -Recurse "./bin"
    }

    Write-Host "construindo binários do linux:"
    $env:GOOS="linux"

    go build -o "./bin/linux/spot" "./src/spot.go"

    7z a -bb0 "./bin/linux/$(nameVersion)-linux-amd64.tar" "./bin/linux/*" "./Themes" "./jsHelper" >$null 2>&1
    7z a -bb0 -sdel -mx9 "./bin/$(nameVersion)-linux-amd64.tar.gz" "./bin/linux/$($nameVersion)-linux-amd64.tar" >$null 2>&1
    Write-Host "✔" -ForegroundColor Green

    Write-Host "construindo binários do macos:"
    $env:GOOS="darwin"

    go build -o "./bin/darwin/spot" "./src/spot.go"

    7z a -bb0 "./bin/darwin/$(nameVersion)-darwin-amd64.tar" "./bin/darwin/*" "./Themes" "./jsHelper" >$null 2>&1
    7z a -bb0 -sdel -mx9 "./bin/$(nameVersion)-darwin-amd64.tar.gz" "./bin/darwin/$($nameVersion)-darwin-amd64.tar" >$null 2>&1
    Write-Host "✔" -ForegroundColor Green

    Write-Host "construindo binários do windows:"
    $env:GOOS="windows"

    go build -o "./bin/windows/spot.exe" "./src/spot.go"

    7z a -bb0 -mx9 "./bin/$(nameVersion)-windows-amd64.tar.gz" "./bin/windows/$($nameVersion)-darwin-amd64.tar" >$null 2>&1
    Write-Host "✔" -ForegroundColor Green
}