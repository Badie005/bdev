# ============================================================
# B.DEV PowerShell Profile
# Profil avance pour le poste de developpement
# ============================================================

# === CONFIGURATION ===
$BDEV_PATH = "$HOME\Dev\.bdev"
$PROJECTS_PATH = "$HOME\Dev\Projects"

# === OH MY POSH ===
$themePath = "$BDEV_PATH\powershell\themes\bdev.omp.json"
if ((Test-Path $themePath) -and (Get-Command oh-my-posh -ErrorAction SilentlyContinue)) {
    oh-my-posh init pwsh --config $themePath | Invoke-Expression
}

# === TERMINAL ICONS ===
if (Get-Module -ListAvailable -Name Terminal-Icons) {
    Import-Module Terminal-Icons
}

# === PSReadLine ===
if (Get-Module -ListAvailable -Name PSReadLine) {
    Set-PSReadLineOption -PredictionSource History
    Set-PSReadLineOption -PredictionViewStyle ListView
    Set-PSReadLineOption -EditMode Windows
    Set-PSReadLineKeyHandler -Key Tab -Function MenuComplete
    Set-PSReadLineKeyHandler -Key UpArrow -Function HistorySearchBackward
    Set-PSReadLineKeyHandler -Key DownArrow -Function HistorySearchForward
}

# ============================================================
# ALIASES
# ============================================================

Set-Alias -Name g -Value git
Set-Alias -Name c -Value code
Set-Alias -Name e -Value explorer
Set-Alias -Name d -Value docker
Set-Alias -Name dc -Value docker-compose

# ============================================================
# NAVIGATION
# ============================================================

function dev { Set-Location $PROJECTS_PATH }
function home { Set-Location $HOME }
function desk { Set-Location "$HOME\Desktop" }
function docs { Set-Location "$HOME\Documents" }

function proj {
    param([string]$name)
    if ($name) {
        $matches = Get-ChildItem $PROJECTS_PATH -Directory | Where-Object { $_.Name -like "*$name*" }
        if ($matches.Count -eq 1) {
            Set-Location $matches[0].FullName
        } elseif ($matches.Count -gt 1) {
            Write-Host "Plusieurs correspondances:" -ForegroundColor Yellow
            $matches | ForEach-Object { Write-Host "  - $($_.Name)" }
        } else {
            Write-Host "Aucun projet trouve: $name" -ForegroundColor Red
        }
    } else {
        Get-ChildItem $PROJECTS_PATH -Directory | ForEach-Object { Write-Host $_.Name }
    }
}

function mkproj {
    param([string]$name)
    $path = "$PROJECTS_PATH\$name"
    New-Item -ItemType Directory -Path $path -Force | Out-Null
    Set-Location $path
    git init
    Write-Host "Projet cree: $path" -ForegroundColor Green
}

# ============================================================
# GIT
# ============================================================

function gs { git status }
function ga { git add . }
function gaa { git add -A }
function gc { param($msg) git commit -m $msg }
function gp { git push }
function gpl { git pull }
function gco { param($branch) git checkout $branch }
function gcob { param($branch) git checkout -b $branch }
function gb { git branch }
function glog { git log --oneline -n 15 }
function gd { git diff }
function gds { git diff --staged }

function gst {
    Write-Host ""
    Write-Host "Projet: $(Split-Path -Leaf (Get-Location))" -ForegroundColor Cyan
    Write-Host ""
    git status -sb
}

function gcommit {
    param([string]$msg)
    if (-not $msg) {
        Write-Host "Usage: gcommit 'message'" -ForegroundColor Yellow
        return
    }
    git add -A
    git commit -m $msg
    Write-Host "Commit cree: $msg" -ForegroundColor Green
}

function gpush {
    param([string]$msg)
    if ($msg) {
        git add -A
        git commit -m $msg
    }
    $branch = git branch --show-current
    git push origin $branch
    Write-Host "Pousse vers: $branch" -ForegroundColor Green
}

# ============================================================
# NPM
# ============================================================

function nr { npm run $args }
function ni { npm install $args }
function nid { npm install -D $args }
function nci { npm ci }
function nrd { npm run dev }
function nrb { npm run build }
function nrt { npm run test }

# ============================================================
# LARAVEL
# ============================================================

function pa { php artisan $args }
function sail { ./vendor/bin/sail $args }
function tinker { php artisan tinker }
function migrate { php artisan migrate }
function fresh { php artisan migrate:fresh --seed }

# ============================================================
# B.DEV
# ============================================================

function bdev {
    & "$BDEV_PATH\bdev.exe" @args
}

function ai {
    param([string]$prompt)
    if (-not $prompt) {
        Write-Host "Usage: ai 'votre question'" -ForegroundColor Yellow
        return
    }
    ollama run codellama:7b $prompt
}

function health {
    & "$BDEV_PATH\scripts\health-check.ps1"
}

function backup {
    param([switch]$Full, [switch]$DryRun)
    $params = @()
    if ($Full) { $params += "-Full" }
    if ($DryRun) { $params += "-DryRun" }
    & "$BDEV_PATH\scripts\backup.ps1" @params
}

function pstats {
    param([switch]$Detailed)
    $params = @()
    if ($Detailed) { $params += "-Detailed" }
    & "$BDEV_PATH\scripts\project-stats.ps1" @params
}

function dotfiles {
    param([switch]$Push, [switch]$Pull)
    $params = @()
    if ($Push) { $params += "-Push" }
    if ($Pull) { $params += "-Pull" }
    if (-not $Push -and -not $Pull) { $params += "-Status" }
    & "$BDEV_PATH\scripts\sync-dotfiles.ps1" @params
}

# ============================================================
# UTILITAIRES
# ============================================================

function mkcd {
    param([string]$dir)
    New-Item -ItemType Directory -Path $dir -Force | Out-Null
    Set-Location $dir
}

function which { 
    param($cmd) 
    Get-Command $cmd -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Source 
}

function reload { . $PROFILE }

function ports {
    Get-NetTCPConnection -State Listen | 
        Select-Object LocalPort, OwningProcess, @{n='Process';e={(Get-Process -Id $_.OwningProcess).Name}} |
        Sort-Object LocalPort
}

function killport {
    param([int]$port)
    $process = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($process) {
        Stop-Process -Id $process.OwningProcess -Force
        Write-Host "Process sur le port $port termine" -ForegroundColor Green
    } else {
        Write-Host "Aucun process sur le port $port" -ForegroundColor Yellow
    }
}

# ============================================================
# BIENVENUE
# ============================================================

function Show-Welcome {
    $hour = (Get-Date).Hour
    $greeting = switch ($hour) {
        {$_ -lt 12} { "Bonjour" }
        {$_ -lt 18} { "Bon apres-midi" }
        default { "Bonsoir" }
    }
    
    Write-Host ""
    Write-Host "[B.DEV] $greeting ! Workstation pret." -ForegroundColor Cyan
    Write-Host "   Commandes: bdev, health, dev, gs" -ForegroundColor DarkGray
    Write-Host ""
}

if (-not $env:BDEV_WELCOMED) {
    Show-Welcome
    $env:BDEV_WELCOMED = "1"
}
