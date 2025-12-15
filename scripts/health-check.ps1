# ============================================================
# B.DEV Health Check Script
# Diagnostic complet du systeme de developpement
# ============================================================

param(
    [switch]$Detailed,
    [switch]$Json
)

$ErrorActionPreference = "SilentlyContinue"

# Couleurs et formatage
function Write-Section {
    param($title) 
    Write-Host ""
    Write-Host "----------------------------------------------------" -ForegroundColor DarkGray
    Write-Host " $title" -ForegroundColor Cyan
    Write-Host "----------------------------------------------------" -ForegroundColor DarkGray
}

function Get-StatusIcon { 
    param($good) 
    if ($good) { "[OK]" } else { "[!!]" } 
}

function Get-StatusColor { 
    param($good) 
    if ($good) { "Green" } else { "Yellow" } 
}

# Header
Write-Host ""
Write-Host "===========================================================" -ForegroundColor Cyan
Write-Host "           B.DEV HEALTH CHECK                              " -ForegroundColor Cyan
Write-Host "           $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')        " -ForegroundColor DarkCyan
Write-Host "===========================================================" -ForegroundColor Cyan

# ==================== SYSTEME ====================
Write-Section "SYSTEME"

# OS Info
$os = Get-CimInstance Win32_OperatingSystem
Write-Host "   OS          : $($os.Caption) $($os.Version)" -ForegroundColor White

# Uptime
$uptime = (Get-Date) - $os.LastBootUpTime
$uptimeStr = "{0}j {1}h {2}m" -f $uptime.Days, $uptime.Hours, $uptime.Minutes
Write-Host "   Uptime      : $uptimeStr" -ForegroundColor White

# ==================== PERFORMANCE ====================
Write-Section "PERFORMANCE"

# CPU
$cpu = Get-CimInstance Win32_Processor
$cpuLoad = (Get-Counter '\Processor(_Total)\% Processor Time' -ErrorAction SilentlyContinue).CounterSamples.CookedValue
$cpuLoadStr = if ($cpuLoad) { "{0:N1}%" -f $cpuLoad } else { "N/A" }
$cpuGood = $cpuLoad -lt 80
Write-Host "   $(Get-StatusIcon $cpuGood) CPU         : $($cpu.Name)" -ForegroundColor White
Write-Host "                 Charge: $cpuLoadStr" -ForegroundColor $(Get-StatusColor $cpuGood)

# RAM
$ram = Get-CimInstance Win32_OperatingSystem
$totalRAM = [math]::Round($ram.TotalVisibleMemorySize / 1MB, 1)
$freeRAM = [math]::Round($ram.FreePhysicalMemory / 1MB, 1)
$usedRAM = $totalRAM - $freeRAM
$ramPercent = [math]::Round(($usedRAM / $totalRAM) * 100, 1)
$ramGood = $ramPercent -lt 85
Write-Host "   $(Get-StatusIcon $ramGood) RAM         : $usedRAM GB / $totalRAM GB ($ramPercent%)" -ForegroundColor $(Get-StatusColor $ramGood)

# Disques
$drives = Get-PSDrive -PSProvider FileSystem | Where-Object { $_.Used -gt 0 }
foreach ($drive in $drives) {
    $total = [math]::Round(($drive.Used + $drive.Free) / 1GB, 1)
    $used = [math]::Round($drive.Used / 1GB, 1)
    $free = [math]::Round($drive.Free / 1GB, 1)
    $percent = [math]::Round(($drive.Used / ($drive.Used + $drive.Free)) * 100, 1)
    $diskGood = $percent -lt 90
    Write-Host "   $(Get-StatusIcon $diskGood) Disque $($drive.Name): : $used GB / $total GB ($percent%) - $free GB libres" -ForegroundColor $(Get-StatusColor $diskGood)
}

# ==================== OUTILS DEV ====================
Write-Section "OUTILS DE DEVELOPPEMENT"

# Node.js
$nodeVersion = node --version 2>$null
if ($nodeVersion) {
    Write-Host "   [OK] Node.js    : $nodeVersion" -ForegroundColor Green
    $npmVersion = npm --version 2>$null
    Write-Host "        npm        : v$npmVersion" -ForegroundColor DarkGray
} else {
    Write-Host "   [X] Node.js    : Non installe" -ForegroundColor Red
}

# Git
$gitVersion = git --version 2>$null
if ($gitVersion) {
    Write-Host "   [OK] Git        : $($gitVersion -replace 'git version ','')" -ForegroundColor Green
} else {
    Write-Host "   [X] Git        : Non installe" -ForegroundColor Red
}

# PHP
$phpVersion = php --version 2>$null | Select-Object -First 1
if ($phpVersion) {
    $phpVer = ($phpVersion -split ' ')[1]
    Write-Host "   [OK] PHP        : $phpVer" -ForegroundColor Green
} else {
    Write-Host "   [ ] PHP        : Non installe" -ForegroundColor DarkGray
}

# Python
$pythonVersion = python --version 2>$null
if ($pythonVersion) {
    Write-Host "   [OK] Python     : $($pythonVersion -replace 'Python ','')" -ForegroundColor Green
} else {
    Write-Host "   [ ] Python     : Non installe" -ForegroundColor DarkGray
}

# Docker
$dockerVersion = docker --version 2>$null
$dockerRunning = Get-Process -Name "Docker Desktop" -ErrorAction SilentlyContinue
if ($dockerVersion) {
    $dockerVer = ($dockerVersion -split ',')[0] -replace 'Docker version ',''
    $dockerStatus = if ($dockerRunning) { "Running" } else { "Stopped" }
    $dockerColor = if ($dockerRunning) { "Green" } else { "Yellow" }
    $dockerIcon = if ($dockerRunning) { "[OK]" } else { "[!!]" }
    Write-Host "   $dockerIcon Docker      : v$dockerVer ($dockerStatus)" -ForegroundColor $dockerColor
} else {
    Write-Host "   [ ] Docker      : Non installe" -ForegroundColor DarkGray
}

# Ollama (AI)
$ollamaRunning = Get-Process -Name "ollama" -ErrorAction SilentlyContinue
if ($ollamaRunning) {
    Write-Host "   [OK] Ollama     : Running (IA locale active)" -ForegroundColor Green
} else {
    $ollamaInstalled = Get-Command ollama -ErrorAction SilentlyContinue
    if ($ollamaInstalled) {
        Write-Host "   [ ] Ollama     : Installe (non demarre)" -ForegroundColor DarkGray
    } else {
        Write-Host "   [ ] Ollama     : Non installe" -ForegroundColor DarkGray
    }
}

# ==================== PROJETS ====================
Write-Section "PROJETS"

$projectsPath = "$HOME\Dev\Projects"
$dirtyProjects = @()

if (Test-Path $projectsPath) {
    $projects = Get-ChildItem $projectsPath -Directory | Where-Object { -not $_.Name.StartsWith('.') }
    $projectCount = $projects.Count
    Write-Host "   Total projets : $projectCount" -ForegroundColor White
    
    # Projets avec changements non commites
    foreach ($proj in $projects) {
        if (Test-Path "$($proj.FullName)\.git") {
            Push-Location $proj.FullName
            $status = git status --porcelain 2>$null
            if ($status) {
                $dirtyProjects += $proj.Name
            }
            Pop-Location
        }
    }
    
    if ($dirtyProjects.Count -gt 0) {
        Write-Host ""
        Write-Host "   [!!] Projets avec changements non commites:" -ForegroundColor Yellow
        foreach ($dirty in $dirtyProjects) {
            Write-Host "        - $dirty" -ForegroundColor Yellow
        }
    } else {
        Write-Host "   [OK] Tous les projets sont clean" -ForegroundColor Green
    }
    
    # Projets recents (modifies dans les 7 derniers jours)
    if ($Detailed) {
        $recentProjects = $projects | 
            Where-Object { $_.LastWriteTime -gt (Get-Date).AddDays(-7) } |
            Sort-Object LastWriteTime -Descending |
            Select-Object -First 5
        
        if ($recentProjects.Count -gt 0) {
            Write-Host ""
            Write-Host "   Projets recents (7 derniers jours):" -ForegroundColor White
            foreach ($recent in $recentProjects) {
                $daysAgo = [math]::Round(((Get-Date) - $recent.LastWriteTime).TotalDays, 0)
                Write-Host "        - $($recent.Name) (il y a ${daysAgo}j)" -ForegroundColor DarkGray
            }
        }
    }
} else {
    Write-Host "   [!!] Dossier Projects non trouve" -ForegroundColor Yellow
}

# ==================== RESUME ====================
Write-Host ""
Write-Host "===========================================================" -ForegroundColor Cyan
Write-Host " HEALTH CHECK TERMINE - $(Get-Date -Format 'HH:mm:ss')" -ForegroundColor Cyan
Write-Host "===========================================================" -ForegroundColor Cyan
Write-Host ""

# Score global
$issues = 0
if ($ramPercent -ge 85) { $issues++ }
if ($drives | Where-Object { [math]::Round(($_.Used / ($_.Used + $_.Free)) * 100, 1) -ge 90 }) { $issues++ }
if ($dirtyProjects.Count -gt 3) { $issues++ }

if ($issues -eq 0) {
    Write-Host "   [*] Systeme en excellent etat !" -ForegroundColor Green
} elseif ($issues -le 2) {
    Write-Host "   [!] Systeme fonctionnel avec quelques points d'attention" -ForegroundColor Yellow
} else {
    Write-Host "   [X] Maintenance recommandee" -ForegroundColor Red
}

Write-Host ""
