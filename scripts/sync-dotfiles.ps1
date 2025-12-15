# ============================================================
# B.DEV Sync Dotfiles Script
# Synchronisation des fichiers de configuration
# ============================================================

param(
    [switch]$Push,
    [switch]$Pull,
    [switch]$Status,
    [string]$RepoPath = "$HOME\Dev\Projects\dotfiles"
)

$ErrorActionPreference = "Stop"

# Header
Write-Host ""
Write-Host "============================================================" -ForegroundColor Green
Write-Host "  B.DEV DOTFILES SYNC" -ForegroundColor Green
Write-Host "  $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -ForegroundColor DarkGreen
Write-Host "============================================================" -ForegroundColor Green
Write-Host ""

# Fichiers a synchroniser
$dotfiles = @(
    @{ Source = "$HOME\.gitconfig"; Dest = "git\.gitconfig" },
    @{ Source = "$HOME\.gitignore_global"; Dest = "git\.gitignore_global" },
    @{ Source = "$HOME\Dev\.bdev\config.json"; Dest = "bdev\config.json" },
    @{ Source = "$HOME\Dev\.bdev\powershell\profile.ps1"; Dest = "powershell\profile.ps1" },
    @{ Source = "$HOME\Dev\.bdev\powershell\themes\bdev.omp.json"; Dest = "powershell\themes\bdev.omp.json" },
    @{ Source = "$HOME\Dev\.bdev\ai\prompts\assistant.md"; Dest = "ai\prompts\assistant.md" },
    @{ Source = "$HOME\AppData\Roaming\Code\User\settings.json"; Dest = "vscode\settings.json" },
    @{ Source = "$HOME\AppData\Roaming\Code\User\keybindings.json"; Dest = "vscode\keybindings.json" }
)

# Verifier que le repo existe
if (-not (Test-Path $RepoPath)) {
    Write-Host "[WARN] Repo dotfiles non trouve: $RepoPath" -ForegroundColor Yellow
    Write-Host "[INFO] Creation du dossier..." -ForegroundColor Cyan
    New-Item -ItemType Directory -Path $RepoPath -Force | Out-Null
    
    Push-Location $RepoPath
    git init
    Pop-Location
    
    Write-Host "[ OK ] Repo initialise" -ForegroundColor Green
}

if ($Status -or (-not $Push -and -not $Pull)) {
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host " STATUS DES DOTFILES" -ForegroundColor White
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host ""
    
    foreach ($file in $dotfiles) {
        $sourceExists = Test-Path $file.Source
        $destPath = Join-Path $RepoPath $file.Dest
        $destExists = Test-Path $destPath
        
        $status = if ($sourceExists -and $destExists) {
            $sourceTime = (Get-Item $file.Source).LastWriteTime
            $destTime = (Get-Item $destPath).LastWriteTime
            if ($sourceTime -gt $destTime) { "Local plus recent" }
            elseif ($destTime -gt $sourceTime) { "Repo plus recent" }
            else { "Synchronise" }
        } elseif ($sourceExists) {
            "Pas dans repo"
        } elseif ($destExists) {
            "Pas en local"
        } else {
            "Manquant"
        }
        
        $color = switch ($status) {
            "Synchronise" { "Green" }
            "Local plus recent" { "Yellow" }
            "Repo plus recent" { "Cyan" }
            default { "DarkGray" }
        }
        
        $fileName = Split-Path $file.Source -Leaf
        Write-Host ("  {0,-30} {1}" -f $fileName, $status) -ForegroundColor $color
    }
}

if ($Push) {
    Write-Host ""
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host " PUSH: Local -> Repo" -ForegroundColor White
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host ""
    
    foreach ($file in $dotfiles) {
        if (Test-Path $file.Source) {
            $destPath = Join-Path $RepoPath $file.Dest
            $destDir = Split-Path $destPath -Parent
            
            if (-not (Test-Path $destDir)) {
                New-Item -ItemType Directory -Path $destDir -Force | Out-Null
            }
            
            Copy-Item $file.Source -Destination $destPath -Force
            $fileName = Split-Path $file.Source -Leaf
            Write-Host "[ OK ] $fileName" -ForegroundColor Green
        }
    }
    
    Write-Host ""
    Write-Host "[INFO] N'oubliez pas de commit et push le repo" -ForegroundColor Cyan
}

if ($Pull) {
    Write-Host ""
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host " PULL: Repo -> Local" -ForegroundColor White
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host ""
    
    foreach ($file in $dotfiles) {
        $sourcePath = Join-Path $RepoPath $file.Dest
        
        if (Test-Path $sourcePath) {
            $destDir = Split-Path $file.Source -Parent
            
            if (-not (Test-Path $destDir)) {
                New-Item -ItemType Directory -Path $destDir -Force | Out-Null
            }
            
            Copy-Item $sourcePath -Destination $file.Source -Force
            $fileName = Split-Path $file.Source -Leaf
            Write-Host "[ OK ] $fileName" -ForegroundColor Green
        }
    }
}

Write-Host ""
