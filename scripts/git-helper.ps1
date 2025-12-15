# ============================================================
# B.DEV Git Helper Script
# Utilitaires Git avances
# ============================================================

param(
    [Parameter(Position=0)]
    [ValidateSet("status-all", "pull-all", "clean-branches", "find-dirty")]
    [string]$Action = "status-all"
)

$ErrorActionPreference = "SilentlyContinue"
$projectsPath = "$HOME\Dev\Projects"

# Header
Write-Host ""
Write-Host "============================================================" -ForegroundColor Magenta
Write-Host "  B.DEV GIT HELPER" -ForegroundColor Magenta
Write-Host "  $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -ForegroundColor DarkMagenta
Write-Host "============================================================" -ForegroundColor Magenta
Write-Host ""

$projects = Get-ChildItem $projectsPath -Directory | 
    Where-Object { Test-Path "$($_.FullName)\.git" }

switch ($Action) {
    "status-all" {
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host " STATUS DE TOUS LES REPOS" -ForegroundColor White
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host ""
        
        foreach ($proj in $projects) {
            Push-Location $proj.FullName
            
            $branch = git branch --show-current 2>$null
            $status = git status --porcelain 2>$null
            $ahead = git rev-list --count "@{u}..HEAD" 2>$null
            $behind = git rev-list --count "HEAD..@{u}" 2>$null
            
            $statusStr = if ($status) { "dirty" } else { "clean" }
            $syncStr = ""
            if ($ahead -gt 0) { $syncStr += "+$ahead " }
            if ($behind -gt 0) { $syncStr += "-$behind" }
            if (-not $syncStr) { $syncStr = "synced" }
            
            $statusColor = if ($status) { "Yellow" } else { "Green" }
            
            Write-Host ("  {0,-25} {1,-15} {2,-8} {3}" -f $proj.Name, $branch, $statusStr, $syncStr) -ForegroundColor $statusColor
            
            Pop-Location
        }
    }
    
    "pull-all" {
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host " PULL SUR TOUS LES REPOS" -ForegroundColor White
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host ""
        
        foreach ($proj in $projects) {
            Write-Host "[....] $($proj.Name)" -NoNewline
            
            Push-Location $proj.FullName
            
            $status = git status --porcelain 2>$null
            if ($status) {
                Write-Host "`r[SKIP] $($proj.Name) - dirty, pull ignore" -ForegroundColor Yellow
            } else {
                $result = git pull 2>&1
                if ($LASTEXITCODE -eq 0) {
                    if ($result -match "Already up to date") {
                        Write-Host "`r[ OK ] $($proj.Name) - deja a jour" -ForegroundColor Green
                    } else {
                        Write-Host "`r[ OK ] $($proj.Name) - mis a jour" -ForegroundColor Cyan
                    }
                } else {
                    Write-Host "`r[FAIL] $($proj.Name) - erreur pull" -ForegroundColor Red
                }
            }
            
            Pop-Location
        }
    }
    
    "clean-branches" {
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host " NETTOYAGE DES BRANCHES MERGEES" -ForegroundColor White
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host ""
        
        foreach ($proj in $projects) {
            Push-Location $proj.FullName
            
            $mergedBranches = git branch --merged main 2>$null | 
                Where-Object { $_ -notmatch '^\*' -and $_ -notmatch 'main|master|develop' }
            
            if ($mergedBranches) {
                Write-Host "  $($proj.Name):" -ForegroundColor Cyan
                foreach ($branch in $mergedBranches) {
                    $branchName = $branch.Trim()
                    git branch -d $branchName 2>$null
                    Write-Host "    Supprime: $branchName" -ForegroundColor Yellow
                }
            }
            
            Pop-Location
        }
        
        Write-Host ""
        Write-Host "[ OK ] Nettoyage termine" -ForegroundColor Green
    }
    
    "find-dirty" {
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host " PROJETS AVEC CHANGEMENTS NON COMMITES" -ForegroundColor White
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host ""
        
        $dirtyCount = 0
        
        foreach ($proj in $projects) {
            Push-Location $proj.FullName
            
            $status = git status --porcelain 2>$null
            if ($status) {
                $dirtyCount++
                $changes = ($status | Measure-Object).Count
                Write-Host "  $($proj.Name) - $changes fichiers modifies" -ForegroundColor Yellow
            }
            
            Pop-Location
        }
        
        Write-Host ""
        if ($dirtyCount -eq 0) {
            Write-Host "[ OK ] Tous les repos sont clean" -ForegroundColor Green
        } else {
            Write-Host "[WARN] $dirtyCount repos avec changements" -ForegroundColor Yellow
        }
    }
}

Write-Host ""
