# ============================================================
# B.DEV NPM Helper Script
# Utilitaires NPM avances
# ============================================================

param(
    [Parameter(Position=0)]
    [ValidateSet("outdated-all", "update-all", "clean-all", "audit-all")]
    [string]$Action = "outdated-all"
)

$ErrorActionPreference = "SilentlyContinue"
$projectsPath = "$HOME\Dev\Projects"

# Header
Write-Host ""
Write-Host "============================================================" -ForegroundColor Green
Write-Host "  B.DEV NPM HELPER" -ForegroundColor Green
Write-Host "  $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -ForegroundColor DarkGreen
Write-Host "============================================================" -ForegroundColor Green
Write-Host ""

# Trouver les projets Node.js
$nodeProjects = Get-ChildItem $projectsPath -Directory | 
    Where-Object { Test-Path "$($_.FullName)\package.json" }

Write-Host "[INFO] $($nodeProjects.Count) projets Node.js trouves" -ForegroundColor Cyan
Write-Host ""

switch ($Action) {
    "outdated-all" {
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host " PACKAGES OBSOLETES PAR PROJET" -ForegroundColor White
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host ""
        
        foreach ($proj in $nodeProjects) {
            if (-not (Test-Path "$($proj.FullName)\node_modules")) {
                continue
            }
            
            Push-Location $proj.FullName
            
            $outdated = npm outdated --json 2>$null | ConvertFrom-Json
            $count = ($outdated.PSObject.Properties | Measure-Object).Count
            
            if ($count -gt 0) {
                Write-Host "  $($proj.Name) - $count packages obsoletes" -ForegroundColor Yellow
            }
            
            Pop-Location
        }
        
        Write-Host ""
        Write-Host "[INFO] Utilisez 'npm outdated' dans chaque projet pour les details" -ForegroundColor DarkGray
    }
    
    "update-all" {
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host " MISE A JOUR DES PACKAGES" -ForegroundColor White
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host ""
        
        foreach ($proj in $nodeProjects) {
            if (-not (Test-Path "$($proj.FullName)\node_modules")) {
                continue
            }
            
            Write-Host "[....] $($proj.Name)" -NoNewline
            
            Push-Location $proj.FullName
            
            npm update --silent 2>$null
            
            if ($LASTEXITCODE -eq 0) {
                Write-Host "`r[ OK ] $($proj.Name)" -ForegroundColor Green
            } else {
                Write-Host "`r[WARN] $($proj.Name) - verifier manuellement" -ForegroundColor Yellow
            }
            
            Pop-Location
        }
    }
    
    "clean-all" {
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host " NETTOYAGE NPM CACHE" -ForegroundColor White
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host ""
        
        Write-Host "[....] Nettoyage du cache npm global" -NoNewline
        npm cache clean --force 2>$null
        Write-Host "`r[ OK ] Nettoyage du cache npm global" -ForegroundColor Green
        
        Write-Host ""
        Write-Host "[INFO] Pour supprimer les node_modules inactifs:" -ForegroundColor DarkGray
        Write-Host "       .\clean-system.ps1 -NodeModules" -ForegroundColor DarkGray
    }
    
    "audit-all" {
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host " AUDIT DE SECURITE" -ForegroundColor White
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host ""
        
        $projectsWithIssues = @()
        
        foreach ($proj in $nodeProjects) {
            if (-not (Test-Path "$($proj.FullName)\node_modules")) {
                continue
            }
            
            Push-Location $proj.FullName
            
            $audit = npm audit --json 2>$null | ConvertFrom-Json
            $vulns = $audit.metadata.vulnerabilities
            
            $total = $vulns.low + $vulns.moderate + $vulns.high + $vulns.critical
            
            if ($total -gt 0) {
                $color = if ($vulns.critical -gt 0 -or $vulns.high -gt 0) { "Red" } else { "Yellow" }
                Write-Host ("  {0,-25} {1} vulnerabilites" -f $proj.Name, $total) -ForegroundColor $color
                
                if ($vulns.critical -gt 0) {
                    Write-Host "    - $($vulns.critical) critiques" -ForegroundColor Red
                }
                if ($vulns.high -gt 0) {
                    Write-Host "    - $($vulns.high) hautes" -ForegroundColor Red
                }
                
                $projectsWithIssues += $proj.Name
            }
            
            Pop-Location
        }
        
        Write-Host ""
        
        if ($projectsWithIssues.Count -eq 0) {
            Write-Host "[ OK ] Aucune vulnerabilite trouvee" -ForegroundColor Green
        } else {
            Write-Host "[WARN] $($projectsWithIssues.Count) projets avec vulnerabilites" -ForegroundColor Yellow
            Write-Host "[INFO] Executez 'npm audit fix' dans chaque projet" -ForegroundColor DarkGray
        }
    }
}

Write-Host ""
