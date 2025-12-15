# ============================================================
# B.DEV Project Statistics Script
# Analyse et statistiques des projets
# ============================================================

param(
    [string]$ProjectPath = "$HOME\Dev\Projects",
    [switch]$Detailed
)

$ErrorActionPreference = "SilentlyContinue"

# Header
Write-Host ""
Write-Host "============================================================" -ForegroundColor Blue
Write-Host "  B.DEV PROJECT STATISTICS" -ForegroundColor Blue
Write-Host "  $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -ForegroundColor DarkBlue
Write-Host "============================================================" -ForegroundColor Blue
Write-Host ""

function Detect-ProjectType {
    param([string]$Path)
    
    if (Test-Path "$Path\package.json") {
        $pkg = Get-Content "$Path\package.json" | ConvertFrom-Json
        if ($pkg.dependencies.next) { return "Next.js" }
        if ($pkg.dependencies.'@angular/core') { return "Angular" }
        if ($pkg.dependencies.react) { return "React" }
        if ($pkg.dependencies.vue) { return "Vue.js" }
        if ($pkg.devDependencies.typescript) { return "TypeScript" }
        return "Node.js"
    }
    if (Test-Path "$Path\composer.json") {
        $composer = Get-Content "$Path\composer.json" | ConvertFrom-Json
        if ($composer.require.'laravel/framework') { return "Laravel" }
        return "PHP"
    }
    if (Test-Path "$Path\requirements.txt") { return "Python" }
    if (Test-Path "$Path\Cargo.toml") { return "Rust" }
    if (Test-Path "$Path\go.mod") { return "Go" }
    return "Other"
}

function Get-ProjectSize {
    param([string]$Path)
    
    $size = Get-ChildItem $Path -Recurse -File -ErrorAction SilentlyContinue |
        Where-Object { $_.FullName -notmatch '(node_modules|vendor|\.git|dist|build)' } |
        Measure-Object -Property Length -Sum
    
    return $size.Sum
}

function Format-Size {
    param([long]$bytes)
    if ($bytes -ge 1GB) { return "{0:N1} GB" -f ($bytes / 1GB) }
    if ($bytes -ge 1MB) { return "{0:N1} MB" -f ($bytes / 1MB) }
    if ($bytes -ge 1KB) { return "{0:N1} KB" -f ($bytes / 1KB) }
    return "$bytes B"
}

# Recuperer tous les projets
$projects = Get-ChildItem $ProjectPath -Directory | 
    Where-Object { -not $_.Name.StartsWith('.') }

# Statistiques
$stats = @{
    TotalProjects = 0
    ByType = @{}
    TotalSize = 0
    WithGit = 0
    DirtyGit = 0
}

$projectData = @()

foreach ($proj in $projects) {
    $stats.TotalProjects++
    
    $type = Detect-ProjectType -Path $proj.FullName
    if (-not $stats.ByType[$type]) { $stats.ByType[$type] = 0 }
    $stats.ByType[$type]++
    
    $size = Get-ProjectSize -Path $proj.FullName
    $stats.TotalSize += $size
    
    $hasGit = Test-Path "$($proj.FullName)\.git"
    $isDirty = $false
    if ($hasGit) {
        $stats.WithGit++
        Push-Location $proj.FullName
        $gitStatus = git status --porcelain 2>$null
        if ($gitStatus) { 
            $isDirty = $true 
            $stats.DirtyGit++
        }
        Pop-Location
    }
    
    $projectData += [PSCustomObject]@{
        Name = $proj.Name
        Type = $type
        Size = $size
        SizeStr = Format-Size $size
        LastModified = $proj.LastWriteTime
        HasGit = $hasGit
        IsDirty = $isDirty
    }
}

# Trier
$recentProjects = $projectData | Sort-Object LastModified -Descending | Select-Object -First 5
$largestProjects = $projectData | Sort-Object Size -Descending | Select-Object -First 5

# AFFICHAGE

Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host " VUE D'ENSEMBLE" -ForegroundColor Cyan
Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host ""
Write-Host "  Total projets    : $($stats.TotalProjects)" -ForegroundColor White
Write-Host "  Taille totale    : $(Format-Size $stats.TotalSize)" -ForegroundColor White
Write-Host "  Avec Git         : $($stats.WithGit) ($([math]::Round($stats.WithGit / [math]::Max($stats.TotalProjects,1) * 100))%)" -ForegroundColor White
$dirtyColor = if ($stats.DirtyGit -gt 0) { 'Yellow' } else { 'Green' }
Write-Host "  Non commites     : $($stats.DirtyGit)" -ForegroundColor $dirtyColor

Write-Host ""
Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host " REPARTITION PAR TYPE" -ForegroundColor Cyan
Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host ""

$sortedTypes = $stats.ByType.GetEnumerator() | Sort-Object Value -Descending
foreach ($type in $sortedTypes) {
    $bar = "#" * [math]::Min($type.Value * 2, 20)
    Write-Host ("  {0,-12} {1,3} {2}" -f $type.Key, $type.Value, $bar) -ForegroundColor Cyan
}

Write-Host ""
Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host " PROJETS RECENTS" -ForegroundColor Cyan
Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host ""

foreach ($proj in $recentProjects) {
    $daysAgo = [math]::Round(((Get-Date) - $proj.LastModified).TotalDays)
    $dirtyMark = if ($proj.IsDirty) { " [!]" } else { "" }
    Write-Host ("  {0,-25} {1,-10} il y a {2,3}j{3}" -f $proj.Name, $proj.Type, $daysAgo, $dirtyMark) -ForegroundColor White
}

Write-Host ""
Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host " PLUS GROS PROJETS" -ForegroundColor Cyan
Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host ""

foreach ($proj in $largestProjects) {
    Write-Host ("  {0,-25} {1,-10} {2,10}" -f $proj.Name, $proj.Type, $proj.SizeStr) -ForegroundColor White
}

# Affichage detaille
if ($Detailed) {
    Write-Host ""
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host " LISTE COMPLETE" -ForegroundColor Cyan
    Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
    Write-Host ""
    Write-Host ("  {0,-25} {1,-12} {2,10} {3,-10} {4}" -f "NOM", "TYPE", "TAILLE", "GIT", "STATUS") -ForegroundColor DarkGray
    Write-Host ("  {0}" -f ("-" * 70)) -ForegroundColor DarkGray
    
    foreach ($proj in ($projectData | Sort-Object Name)) {
        $gitStr = if ($proj.HasGit) { "oui" } else { "-" }
        $statusStr = if ($proj.IsDirty) { "dirty" } elseif ($proj.HasGit) { "clean" } else { "" }
        $statusColor = if ($proj.IsDirty) { "Yellow" } else { "White" }
        Write-Host ("  {0,-25} {1,-12} {2,10} {3,-10} " -f $proj.Name, $proj.Type, $proj.SizeStr, $gitStr) -NoNewline -ForegroundColor White
        Write-Host $statusStr -ForegroundColor $statusColor
    }
}

Write-Host ""
