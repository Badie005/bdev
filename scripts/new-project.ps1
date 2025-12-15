# ============================================================
# B.DEV New Project Script
# Creation rapide de projets depuis templates
# ============================================================

param(
    [Parameter(Mandatory=$true, Position=0)]
    [string]$Template,
    
    [Parameter(Mandatory=$true, Position=1)]
    [string]$Name,
    
    [switch]$NoGit,
    [switch]$NoOpen
)

$ErrorActionPreference = "Stop"

$BDEV_PATH = "$HOME\Dev\.bdev"
$PROJECTS_PATH = "$HOME\Dev\Projects"
$TEMPLATES_PATH = "$BDEV_PATH\templates"

# Header
Write-Host ""
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  B.DEV NEW PROJECT" -ForegroundColor Cyan
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""

# Verifier le template
$templatePath = "$TEMPLATES_PATH\$Template"
if (-not (Test-Path $templatePath)) {
    Write-Host "[FAIL] Template '$Template' non trouve" -ForegroundColor Red
    Write-Host ""
    Write-Host "Templates disponibles:" -ForegroundColor Yellow
    Get-ChildItem $TEMPLATES_PATH -Directory | ForEach-Object {
        Write-Host "  - $($_.Name)" -ForegroundColor White
    }
    Write-Host ""
    exit 1
}

# Verifier que le projet n'existe pas
$destPath = "$PROJECTS_PATH\$Name"
if (Test-Path $destPath) {
    Write-Host "[FAIL] Le projet '$Name' existe deja" -ForegroundColor Red
    Write-Host "       Chemin: $destPath" -ForegroundColor DarkGray
    exit 1
}

Write-Host "[INFO] Template  : $Template" -ForegroundColor Cyan
Write-Host "[INFO] Projet    : $Name" -ForegroundColor Cyan
Write-Host "[INFO] Chemin    : $destPath" -ForegroundColor DarkGray
Write-Host ""

# Copier le template
Write-Host "[....] Copie du template" -NoNewline
Copy-Item $templatePath -Destination $destPath -Recurse
Write-Host "`r[ OK ] Copie du template" -ForegroundColor Green

# Remplacer les placeholders dans les fichiers
$filesToProcess = @("package.json", "composer.json", "README.md", "pyproject.toml")
foreach ($fileName in $filesToProcess) {
    $filePath = "$destPath\$fileName"
    if (Test-Path $filePath) {
        $content = Get-Content $filePath -Raw
        $content = $content -replace '\{\{PROJECT_NAME\}\}', $Name
        $content = $content -replace '\{\{DATE\}\}', (Get-Date -Format "yyyy-MM-dd")
        $content = $content -replace 'nextjs-starter|laravel-api|angular-app|python-cli', $Name
        Set-Content $filePath -Value $content
        Write-Host "[ OK ] $fileName mis a jour" -ForegroundColor Green
    }
}

# Initialiser Git
if (-not $NoGit) {
    Push-Location $destPath
    git init --quiet
    git add -A
    git commit -m "Initial commit from B.DEV template: $Template" --quiet
    Pop-Location
    Write-Host "[ OK ] Git initialise" -ForegroundColor Green
}

# Ouvrir dans VSCode
if (-not $NoOpen) {
    code $destPath
    Write-Host "[ OK ] Ouvert dans VSCode" -ForegroundColor Green
}

Write-Host ""
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  PROJET CREE AVEC SUCCES" -ForegroundColor Green
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Prochaines etapes:" -ForegroundColor White

# Instructions specifiques par template
switch ($Template) {
    "nextjs-starter" {
        Write-Host "    cd $Name" -ForegroundColor DarkGray
        Write-Host "    npm install" -ForegroundColor DarkGray
        Write-Host "    npm run dev" -ForegroundColor DarkGray
    }
    "laravel-api" {
        Write-Host "    cd $Name" -ForegroundColor DarkGray
        Write-Host "    composer install" -ForegroundColor DarkGray
        Write-Host "    cp .env.example .env" -ForegroundColor DarkGray
        Write-Host "    php artisan key:generate" -ForegroundColor DarkGray
    }
    "angular-app" {
        Write-Host "    cd $Name" -ForegroundColor DarkGray
        Write-Host "    npm install" -ForegroundColor DarkGray
        Write-Host "    ng serve" -ForegroundColor DarkGray
    }
    "python-cli" {
        Write-Host "    cd $Name" -ForegroundColor DarkGray
        Write-Host "    python -m venv .venv" -ForegroundColor DarkGray
        Write-Host "    .venv\Scripts\activate" -ForegroundColor DarkGray
        Write-Host "    pip install -r requirements.txt" -ForegroundColor DarkGray
    }
    default {
        Write-Host "    cd $Name" -ForegroundColor DarkGray
    }
}

Write-Host ""
