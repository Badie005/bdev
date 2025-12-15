# ============================================================
# B.DEV Backup Script
# Sauvegarde intelligente des configurations et projets
# ============================================================

param(
    [string]$Destination = "D:\Backups\B.LAPTOP",
    [switch]$Full,
    [switch]$DryRun
)

$ErrorActionPreference = "Stop"

# Configuration
$timestamp = Get-Date -Format "yyyy-MM-dd_HHmm"
$backupDir = "$Destination\$timestamp"
$logFile = "$HOME\Dev\.bdev\data\logs\backup_$timestamp.log"

# Header
Write-Host ""
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  B.DEV BACKUP SYSTEM" -ForegroundColor Cyan
Write-Host "  $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -ForegroundColor DarkCyan
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""

if ($DryRun) {
    Write-Host "[DRY RUN] Aucune modification ne sera effectuee" -ForegroundColor Yellow
    Write-Host ""
}

# Ce qu'on sauvegarde
$targets = @(
    @{ Path = "$HOME\.config"; Name = "config" },
    @{ Path = "$HOME\Dev\.bdev"; Name = "bdev" },
    @{ Path = "$HOME\Dev\Projects"; Name = "projects" }
)

if ($Full) {
    $targets += @{ Path = "$HOME\Documents\Notes"; Name = "notes" }
    $targets += @{ Path = "$HOME\Dev\Design"; Name = "design" }
    Write-Host "[INFO] Mode FULL active - backup complet" -ForegroundColor Cyan
    Write-Host ""
}

# Creer le dossier de backup
if (-not $DryRun) {
    if (-not (Test-Path $Destination)) {
        Write-Host "[WARN] Creation du dossier de destination: $Destination" -ForegroundColor Yellow
        New-Item -ItemType Directory -Path $Destination -Force | Out-Null
    }
    New-Item -ItemType Directory -Path $backupDir -Force | Out-Null
}

# Statistiques
$totalSize = 0
$backedUp = 0
$failed = 0

Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host " SAUVEGARDE EN COURS" -ForegroundColor White
Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host ""

foreach ($target in $targets) {
    $sourcePath = $target.Path
    $destPath = "$backupDir\$($target.Name)"
    
    if (-not (Test-Path $sourcePath)) {
        Write-Host "[SKIP] $($target.Name) : Source non trouvee" -ForegroundColor DarkGray
        continue
    }
    
    Write-Host "[....] $($target.Name)" -NoNewline
    
    if (-not $DryRun) {
        try {
            $robocopyArgs = @(
                $sourcePath,
                $destPath,
                "/MIR",
                "/XD", "node_modules", "vendor", ".git", "dist", "build", ".next", "__pycache__", ".venv",
                "/XF", "*.log",
                "/NFL", "/NDL", "/NJH", "/NJS", "/NC", "/NS"
            )
            
            $null = robocopy @robocopyArgs
            
            if ($LASTEXITCODE -le 3) {
                $size = (Get-ChildItem $destPath -Recurse -File -ErrorAction SilentlyContinue | Measure-Object -Property Length -Sum).Sum
                $sizeMB = [math]::Round($size / 1MB, 2)
                $totalSize += $size
                $backedUp++
                Write-Host "`r[ OK ] $($target.Name) - $sizeMB MB" -ForegroundColor Green
            } else {
                $failed++
                Write-Host "`r[FAIL] $($target.Name)" -ForegroundColor Red
            }
        } catch {
            $failed++
            Write-Host "`r[FAIL] $($target.Name) - $_" -ForegroundColor Red
        }
    } else {
        Write-Host "`r[DRY ] $($target.Name) - serait copie vers $destPath" -ForegroundColor Yellow
        $backedUp++
    }
}

# Resume
Write-Host ""
Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host " RESUME" -ForegroundColor White
Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host ""

if (-not $DryRun) {
    $totalSizeMB = [math]::Round($totalSize / 1MB, 2)
    $totalSizeGB = [math]::Round($totalSize / 1GB, 2)
    
    Write-Host "  Destination : $backupDir" -ForegroundColor White
    Write-Host "  Elements    : $backedUp sauvegardes" -ForegroundColor Green
    Write-Host "  Taille      : $totalSizeMB MB ($totalSizeGB GB)" -ForegroundColor Yellow
    
    if ($failed -gt 0) {
        Write-Host "  Erreurs     : $failed" -ForegroundColor Red
    }
    
    # Log
    $logsDir = Split-Path $logFile -Parent
    if (-not (Test-Path $logsDir)) {
        New-Item -ItemType Directory -Path $logsDir -Force | Out-Null
    }
    
    @"
B.DEV Backup Log
================
Date: $timestamp
Destination: $backupDir
Mode: $(if ($Full) { 'FULL' } else { 'Standard' })
Elements: $backedUp
Taille: $totalSizeMB MB
Erreurs: $failed
"@ | Out-File $logFile -Encoding UTF8

    Write-Host "  Log         : $logFile" -ForegroundColor DarkGray
} else {
    Write-Host "  [DRY RUN] Aucune modification effectuee" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "[ OK ] Backup termine" -ForegroundColor Green
Write-Host ""
