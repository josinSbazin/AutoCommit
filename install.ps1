$ErrorActionPreference = "Stop"

$Repo = "josinSbazin/AutoCommit"
$InstallDir = "$env:LOCALAPPDATA\Programs\autocommit"

function Get-Arch {
    if ([Environment]::Is64BitOperatingSystem) {
        if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
            return "arm64"
        }
        return "amd64"
    }
    return "386"
}

function Main {
    $Arch = Get-Arch
    Write-Host "Detected: windows/$Arch"

    $LatestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
    $Version = $LatestRelease.tag_name

    if (-not $Version) {
        Write-Error "Failed to get latest version"
        exit 1
    }

    Write-Host "Latest version: $Version"

    $Filename = "autocommit_windows_$Arch.zip"
    $Url = "https://github.com/$Repo/releases/download/$Version/$Filename"

    Write-Host "Downloading $Url..."

    $TmpDir = New-TemporaryFile | ForEach-Object { Remove-Item $_; New-Item -ItemType Directory -Path $_ }
    $ZipPath = Join-Path $TmpDir $Filename

    Invoke-WebRequest -Uri $Url -OutFile $ZipPath
    Expand-Archive -Path $ZipPath -DestinationPath $TmpDir -Force

    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    Copy-Item (Join-Path $TmpDir "autocommit.exe") $InstallDir -Force
    Remove-Item $TmpDir -Recurse -Force

    Write-Host "Installed to $InstallDir\autocommit.exe"

    $CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($CurrentPath -notlike "*$InstallDir*") {
        $NewPath = "$CurrentPath;$InstallDir"
        [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
        Write-Host "Added to PATH"
        Write-Host ""
        Write-Host "Restart your terminal to use 'autocommit' command"
    }

    Write-Host "Done! Run 'autocommit --help' to get started."
}

Main
