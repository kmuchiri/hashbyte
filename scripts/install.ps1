$ErrorActionPreference = "Stop"

$Repo = "kmuchiri/hashbyte"
Write-Host "Fetching latest release for $Repo..."

# Get latest release tag
$ReleaseInfo = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
$LatestTag = $ReleaseInfo.tag_name

if ([string]::IsNullOrEmpty($LatestTag)) {
    Write-Error "Could not determine latest release tag. Have you published a release yet?"
}

# Determine architecture
$Arch = $env:PROCESSOR_ARCHITECTURE
if ($Arch -eq "AMD64") {
    $ArchName = "x86_64"
} elseif ($Arch -eq "ARM64") {
    $ArchName = "arm64"
} else {
    $ArchName = "i386"
}

$FileName = "hashbyte_Windows_${ArchName}.tar.gz"
$DownloadUrl = "https://github.com/$Repo/releases/download/${LatestTag}/${FileName}"
$TempFile = Join-Path $env:TEMP $FileName
$ExtractDir = Join-Path $env:TEMP "hashbyte_extract"

Write-Host "Downloading $DownloadUrl..."
Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempFile

Write-Host "Extracting..."
if (Test-Path $ExtractDir) { Remove-Item -Recurse -Force $ExtractDir }
New-Item -ItemType Directory -Path $ExtractDir | Out-Null

# Use native tar.exe built into Windows 10+
tar -xzf $TempFile -C $ExtractDir

$InstallDir = Join-Path $env:LOCALAPPDATA "hashbyte"
if (-Not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

Write-Host "Installing to $InstallDir..."
Move-Item -Path (Join-Path $ExtractDir "hashbyte.exe") -Destination (Join-Path $InstallDir "hashbyte.exe") -Force

# Cleanup
Remove-Item $TempFile
Remove-Item -Recurse -Force $ExtractDir

# Add to PATH if not exists
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notlike "*$InstallDir*") {
    Write-Host "Adding $InstallDir to User PATH..."
    $NewPath = if ($UserPath.EndsWith(";")) { "$UserPath$InstallDir" } else { "$UserPath;$InstallDir" }
    [Environment]::SetEnvironmentVariable("PATH", $NewPath, "User")
}

Write-Host ""
Write-Host "hashbyte successfully installed!"
Write-Host "Please close and reopen your PowerShell/Terminal to apply the PATH changes, then run 'hashbyte'."
