# ============================================================
# B.DEV Clean System Script  
# Nettoyage intelligent du systeme de developpement
# ============================================================

param(
    [switch]$All,
    [switch]$Downloads,
    [switch]$NodeModules,
    [switch]$Cache,
    [switch]$Temp,
    [switch]$DryRun
)

$ErrorActionPreference = "SilentlyContinue"

# Header
Write-Host ""
Write-Host "============================================================" -ForegroundColor Magenta
Write-Host "  B.DEV SYSTEM CLEANER" -ForegroundColor Magenta
Write-Host "  $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -ForegroundColor DarkMagenta
Write-Host "============================================================" -ForegroundColor Magenta
Write-Host ""

if ($DryRun) {
    Write-Host "[DRY RUN] Aucun fichier ne sera supprime" -ForegroundColor Yellow
    Write-Host ""
}

$totalFreed = 0

function Format-Size {
    param([long]$bytes)
    if ($bytes -ge 1GB) { return "{0:N2} GB" -f ($bytes / 1GB) }
    if ($bytes -ge 1MB) { return "{0:N2} MB" -f ($bytes / 1MB) }
    if ($bytes -ge 1KB) { return "{0:N2} KB" -f ($bytes / 1KB) }
    return "$bytes bytes"
}

function Remove-FolderContents {
    param(
        [string]$Path,
        [string]$Name,
        [int]$OlderThanDays = 0
    )
    
    if (-not (Test-Path $Path)) {
        Write-Host "[SKIP] $Name : Dossier non trouve" -ForegroundColor DarkGray
        return 0
    }
    
    $items = if ($OlderThanDays -gt 0) {
        Get-ChildItem $Path -Recurse -File | Where-Object { $_.LastWriteTime -lt (Get-Date).AddDays(-$OlderThanDays) }
    } else {
        Get-ChildItem $Path -Recurse -File
    }
    
    $size = ($items | Measure-Object -Property Length -Sum).Sum
    $count = $items.Count
    
    if ($count -eq 0) {
        Write-Host "[ OK ] $Name : Deja propre" -ForegroundColor Green
        return 0
    }
    
    if (-not $DryRun) {
        $items | Remove-Item -Force -ErrorAction SilentlyContinue
        Write-Host "[ OK ] $Name : $(Format-Size $size) liberes ($count fichiers)" -ForegroundColor Cyan
    } else {
        Write-Host "[DRY ] $Name : $(Format-Size $size) a liberer ($count fichiers)" -ForegroundColor Yellow
    }
    
    return $size
}

# Si aucun switch, activer All
if (-not ($All -or $Downloads -or $NodeModules -or $Cache -or $Temp)) {
    $All = $true
}

# DOWNLOADS
if ($All -or $Downloads) {
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host " DOWNLOADS (fichiers > 30 jours)" -ForegroundColor White
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    
    $downloadsPath = "$HOME\Downloads"
    $totalFreed += Remove-FolderContents -Path $downloadsPath -Name "Downloads" -OlderThanDays 30
    Write-Host ""
}

# NODE_MODULES
if ($All -or $NodeModules) {
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host " NODE_MODULES (projets inactifs > 30 jours)" -ForegroundColor White
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    
    $projectsPath = "$HOME\Dev\Projects"
    $inactiveProjects = Get-ChildItem $projectsPath -Directory | 
        Where-Object { 
            (Test-Path "$($_.FullName)\node_modules") -and 
            $_.LastWriteTime -lt (Get-Date).AddDays(-30)
        }
    
    foreach ($proj in $inactiveProjects) {
        $nmPath = "$($proj.FullName)\node_modules"
        $size = (Get-ChildItem $nmPath -Recurse -File -ErrorAction SilentlyContinue | Measure-Object -Property Length -Sum).Sum
        
        if (-not $DryRun) {
            Remove-Item $nmPath -Recurse -Force -ErrorAction SilentlyContinue
            Write-Host "[ OK ] $($proj.Name)/node_modules : $(Format-Size $size)" -ForegroundColor Cyan
        } else {
            Write-Host "[DRY ] $($proj.Name)/node_modules : $(Format-Size $size)" -ForegroundColor Yellow
        }
        $totalFreed += $size
    }
    
    if ($inactiveProjects.Count -eq 0) {
        Write-Host "[ OK ] Aucun node_modules inactif trouve" -ForegroundColor Green
    }
    Write-Host ""
}

# CACHE
if ($All -or $Cache) {
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host " CACHE" -ForegroundColor White
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    
    $totalFreed += Remove-FolderContents -Path "$HOME\AppData\Local\npm-cache" -Name "npm cache"
    $totalFreed += Remove-FolderContents -Path "$HOME\Dev\.bdev\cache" -Name "bdev cache"
    $totalFreed += Remove-FolderContents -Path "$HOME\AppData\Local\pip\cache" -Name "pip cache"
    Write-Host ""
}

# TEMP
if ($All -or $Temp) {
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host " FICHIERS TEMPORAIRES" -ForegroundColor White
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    
    $totalFreed += Remove-FolderContents -Path "$HOME\AppData\Local\Temp" -Name "Temp utilisateur" -OlderThanDays 7
    Write-Host ""
}

# RESUME
Write-Host "============================================================" -ForegroundColor Magenta
Write-Host " RESUME" -ForegroundColor Magenta
Write-Host "============================================================" -ForegroundColor Magenta
Write-Host ""

if ($DryRun) {
    Write-Host "  Espace recuperable : $(Format-Size $totalFreed)" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "  Executez sans -DryRun pour nettoyer effectivement" -ForegroundColor DarkGray
} else {
    Write-Host "  Espace libere : $(Format-Size $totalFreed)" -ForegroundColor Green
}

Write-Host ""
