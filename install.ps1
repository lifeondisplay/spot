param (
    [string] $version
)

$ErrorActionPreference = "Stop"

# habilita o tls 1.2, pois ele é necessário para conexões com o github
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

# funções auxiliares para um output bonito do terminal
function Write-Part ([string] $Text) {
    Write-Host $Text -NoNewline
}

function Write-Emphasized ([string] $Text) {
    Write-Host $Text -NoNewline -ForegroundColor "Cyan"
}

function Write-Done {
    Write-Host " > " -NoNewline
    Write-Host "OK" -ForegroundColor "Green"
}

if (-not $version) {
    # determina a release mais recente do spot por meio da api do github
    $latest_release_uri = "https://api.github.com/repos/lifeondisplay/spot/releases/latest"

    Write-Part "BAIXANDO    "; Write-Emphasized $latest_release_uri
    $latest_release_json = Invoke-WebRequest -Uri $latest_release_uri

    Write-Done

    $version = ($latest_release_json | ConvertFrom-Json).tag_name -replace "v", ""
}

# cria o diretório ~\spot caso ainda não exista
$sp_dir = "${HOME}\spot"

if (-not (Test-Path $sp_dir)) {
    Write-Part "CRIANDO PASTA  "; Write-Emphasized $sp_dir
    New-Item -Path $sp_dir -ItemType Directory | Out-Null
    Write-Done
}

# baixa a release
$zip_file = "${sp_dir}\spot-${version}-windows-x64.zip"
$download_uri = "https://github.com/lifeondisplay/spot/releases/download/" + "v${version}/spot-${version}-windows-x64.zip"
Write-Part "BAIXANDO    "; Write-Emphasized $download_uri
Invoke-WebRequest -Uri $download_uri -OutFile $zip_file
Write-Done

# extrai o spot.exe e os assets do arquivo .zip
Write-Part "EXTRAINDO     "; Write-Emphasized $zip_file
Write-Part " em "; Write-Emphasized ${sp_dir};

# utilizando -force para reescrever o spot.exe e seus assets caso já exista
Expand-Archive -Path $zip_file -DestinationPath $sp_dir -Force
Write-Done

# remove o arquivo .zip
Write-Part "REMOVENDO       "; Write-Emphasized $zip_file
Remove-Item -Path $zip_file
Write-Done

# obtém a variável de ambiente path para o usuário atual
$user = [EnvironmentVariableTarget]::User
$path = [Environment]::GetEnvironmentVariable("PATH", $user)

# verifica se o diretório spot está no path
$paths = $path -split ";"
$is_in_path = $paths -contains $sp_dir -or $paths -contains "${sp_dir}\"

# adiciona o diretório spot ao path se ainda não tiver sido adicionado
if (-not $is_in_path) {
    Write-Part "ADICIONANDO         "; Write-Emphasized $sp_dir; Write-Part " a "
    Write-Emphasized "PATH"; Write-Part " variável de ambiente..."
    [Environment]::SetEnvironmentVariable("PATH", "${path};${sp_dir}", $user)

    # Add Deno to the PATH variable of the current terminal session
    # so `deno` can be used immediately without restarting the terminal.
    $env:PATH += ";${sp_dir}"
    Write-Done
}

Write-Host ""
Write-Done "spot foi instalado com sucesso."
Write-Part "rode "; Write-Emphasized "spot --help"; Write-Host " para iniciar."
Write-Host ""