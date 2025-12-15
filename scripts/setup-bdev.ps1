# ============================================================
# B.DEV Setup Script
# Installation et configuration de l'environnement complet
# ============================================================

param(
    [switch]$SkipPackages,
    [switch]$SkipProfile,
    [switch]$SkipAI
)

$ErrorActionPreference = "Stop"

# Header
Write-Host ""
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  B.DEV WORKSTATION SETUP" -ForegroundColor Cyan
Write-Host "  Installation automatique" -ForegroundColor DarkCyan
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""

$BDEV_PATH = "$HOME\Dev\.bdev"

# ==================== VERIFICATION ====================
Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host " VERIFICATION DES PREREQUIS" -ForegroundColor White
Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host ""

# Verifier winget
if (Get-Command winget -ErrorAction SilentlyContinue) {
    Write-Host "[ OK ] winget disponible" -ForegroundColor Green
} else {
    Write-Host "[FAIL] winget non trouve" -ForegroundColor Red
    exit 1
}

# Verifier Python
$python = Get-Command python -ErrorAction SilentlyContinue
if ($python) {
    $pyVersion = python --version 2>&1
    Write-Host "[ OK ] Python: $pyVersion" -ForegroundColor Green
} else {
    Write-Host "[WARN] Python non trouve - Installation..." -ForegroundColor Yellow
    winget install Python.Python.3.12 --accept-source-agreements --accept-package-agreements
}

# ==================== PACKAGES ====================
if (-not $SkipPackages) {
    Write-Host ""
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host " INSTALLATION DES PACKAGES" -ForegroundColor White
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host ""

    $packages = @(
        @{ Name = "Oh My Posh"; Id = "JanDeDobbeleer.OhMyPosh" },
        @{ Name = "Git"; Id = "Git.Git" },
        @{ Name = "Node.js LTS"; Id = "OpenJS.NodeJS.LTS" },
        @{ Name = "Visual Studio Code"; Id = "Microsoft.VisualStudioCode" }
    )

    foreach ($pkg in $packages) {
        Write-Host "[....] $($pkg.Name)" -NoNewline
        $installed = winget list --id $pkg.Id 2>$null | Select-String $pkg.Id
        if ($installed) {
            Write-Host "`r[ OK ] $($pkg.Name) - Deja installe" -ForegroundColor Green
        } else {
            winget install $pkg.Id --accept-source-agreements --accept-package-agreements --silent 2>$null
            Write-Host "`r[ OK ] $($pkg.Name) - Installe" -ForegroundColor Green
        }
    }
}

# ==================== PYTHON DEPENDENCIES ====================
Write-Host ""
Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host " DEPENDANCES PYTHON (CLI B.DEV)" -ForegroundColor White
Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host ""

$requirementsPath = "$BDEV_PATH\cli\requirements.txt"
if (Test-Path $requirementsPath) {
    Write-Host "[....] Installation de typer, rich, colorama..." -NoNewline
    python -m pip install -r $requirementsPath --quiet 2>$null
    Write-Host "`r[ OK ] Dependances Python installees" -ForegroundColor Green
}

# ==================== POWERSHELL MODULES ====================
Write-Host ""
Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host " MODULES POWERSHELL" -ForegroundColor White
Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host ""

if (-not (Get-Module -ListAvailable -Name Terminal-Icons)) {
    Write-Host "[....] Installation Terminal-Icons..." -NoNewline
    Install-Module -Name Terminal-Icons -Scope CurrentUser -Force
    Write-Host "`r[ OK ] Terminal-Icons" -ForegroundColor Green
} else {
    Write-Host "[ OK ] Terminal-Icons" -ForegroundColor Green
}

if (Get-Module -ListAvailable -Name PSReadLine) {
    Write-Host "[ OK ] PSReadLine" -ForegroundColor Green
}

# ==================== PROFIL POWERSHELL ====================
if (-not $SkipProfile) {
    Write-Host ""
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host " CONFIGURATION DU PROFIL POWERSHELL" -ForegroundColor White
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host ""

    $profileDir = Split-Path $PROFILE -Parent
    if (-not (Test-Path $profileDir)) {
        New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
    }

    $bdevProfilePath = "$BDEV_PATH\powershell\profile.ps1"
    $importLine = ". `"$bdevProfilePath`""
    
    if (Test-Path $PROFILE) {
        $currentProfile = Get-Content $PROFILE -Raw
        if ($currentProfile -notmatch [regex]::Escape($bdevProfilePath)) {
            Add-Content -Path $PROFILE -Value "`n# B.DEV Profile`n$importLine"
            Write-Host "[ OK ] Profil B.DEV ajoute a `$PROFILE" -ForegroundColor Green
        } else {
            Write-Host "[ OK ] Profil B.DEV deja configure" -ForegroundColor Green
        }
    } else {
        "# B.DEV Profile`n$importLine" | Out-File $PROFILE -Encoding UTF8
        Write-Host "[ OK ] Profil PowerShell cree avec B.DEV" -ForegroundColor Green
    }
}

# ==================== OLLAMA ====================
if (-not $SkipAI) {
    Write-Host ""
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host " IA LOCALE (OLLAMA)" -ForegroundColor White
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host ""

    $ollama = Get-Command ollama -ErrorAction SilentlyContinue
    if (-not $ollama) {
        Write-Host "[....] Installation d'Ollama..." -NoNewline
        winget install Ollama.Ollama --accept-source-agreements --accept-package-agreements 2>$null
        Write-Host "`r[ OK ] Ollama installe" -ForegroundColor Green
        Write-Host ""
        Write-Host "[INFO] Pour telecharger un modele IA :" -ForegroundColor Yellow
        Write-Host "       ollama pull phi3:mini    (leger, ~2GB)" -ForegroundColor DarkGray
        Write-Host "       ollama pull codellama:7b (code, ~4GB)" -ForegroundColor DarkGray
    } else {
        Write-Host "[ OK ] Ollama deja installe" -ForegroundColor Green
    }
}

# ==================== RESUME ====================
Write-Host ""
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  INSTALLATION TERMINEE" -ForegroundColor Cyan
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Prochaines etapes :" -ForegroundColor White
Write-Host ""
Write-Host "  1. Redemarrez votre terminal" -ForegroundColor DarkGray
Write-Host "  2. Tapez 'bdev' pour voir le dashboard" -ForegroundColor DarkGray
Write-Host "  3. Tapez 'health' pour verifier le systeme" -ForegroundColor DarkGray
Write-Host ""
Write-Host "  Commandes utiles :" -ForegroundColor White
Write-Host "     bdev list      - Lister les projets" -ForegroundColor DarkGray
Write-Host "     bdev new       - Creer un projet" -ForegroundColor DarkGray
Write-Host "     pstats         - Statistiques projets" -ForegroundColor DarkGray
Write-Host "     backup         - Sauvegarder" -ForegroundColor DarkGray
Write-Host ""
Write-Host "  Happy coding!" -ForegroundColor Cyan
Write-Host ""
