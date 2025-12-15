# ============================================================
# B.DEV Docker Helper Script
# Gestion des conteneurs Docker
# ============================================================

param(
    [Parameter(Position=0)]
    [ValidateSet("status", "start", "stop", "clean", "logs")]
    [string]$Action = "status",
    
    [string]$Container
)

$ErrorActionPreference = "SilentlyContinue"

# Header
Write-Host ""
Write-Host "============================================================" -ForegroundColor Blue
Write-Host "  B.DEV DOCKER HELPER" -ForegroundColor Blue
Write-Host "  $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -ForegroundColor DarkBlue
Write-Host "============================================================" -ForegroundColor Blue
Write-Host ""

# Verifier Docker
$dockerRunning = Get-Process -Name "Docker Desktop" -ErrorAction SilentlyContinue
if (-not $dockerRunning) {
    Write-Host "[WARN] Docker Desktop n'est pas demarre" -ForegroundColor Yellow
    Write-Host ""
    
    $start = Read-Host "Demarrer Docker Desktop? (o/n)"
    if ($start -eq "o") {
        Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe"
        Write-Host "[INFO] Demarrage de Docker Desktop..." -ForegroundColor Cyan
        Write-Host "[INFO] Attendez quelques secondes et reessayez" -ForegroundColor Cyan
    }
    exit
}

switch ($Action) {
    "status" {
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host " CONTENEURS EN COURS" -ForegroundColor White
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host ""
        
        $containers = docker ps --format "{{.ID}}\t{{.Names}}\t{{.Status}}\t{{.Ports}}" 2>$null
        
        if ($containers) {
            Write-Host ("  {0,-15} {1,-25} {2,-20} {3}" -f "ID", "NOM", "STATUS", "PORTS") -ForegroundColor DarkGray
            Write-Host ("  " + "-" * 80) -ForegroundColor DarkGray
            
            foreach ($line in $containers) {
                $parts = $line -split "`t"
                $id = $parts[0].Substring(0, 12)
                $name = $parts[1]
                $status = $parts[2]
                $ports = if ($parts[3]) { $parts[3] } else { "-" }
                
                $color = if ($status -match "Up") { "Green" } else { "Yellow" }
                Write-Host ("  {0,-15} {1,-25} {2,-20} {3}" -f $id, $name, $status, $ports) -ForegroundColor $color
            }
        } else {
            Write-Host "  Aucun conteneur en cours d'execution" -ForegroundColor DarkGray
        }
        
        Write-Host ""
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host " IMAGES DOCKER" -ForegroundColor White
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host ""
        
        $images = docker images --format "{{.Repository}}\t{{.Tag}}\t{{.Size}}" 2>$null | Select-Object -First 10
        
        if ($images) {
            Write-Host ("  {0,-30} {1,-15} {2}" -f "IMAGE", "TAG", "TAILLE") -ForegroundColor DarkGray
            Write-Host ("  " + "-" * 60) -ForegroundColor DarkGray
            
            foreach ($line in $images) {
                $parts = $line -split "`t"
                Write-Host ("  {0,-30} {1,-15} {2}" -f $parts[0], $parts[1], $parts[2]) -ForegroundColor White
            }
        }
    }
    
    "start" {
        if ($Container) {
            Write-Host "[....] Demarrage de $Container" -NoNewline
            docker start $Container 2>$null
            if ($LASTEXITCODE -eq 0) {
                Write-Host "`r[ OK ] Demarrage de $Container" -ForegroundColor Green
            } else {
                Write-Host "`r[FAIL] Demarrage de $Container" -ForegroundColor Red
            }
        } else {
            Write-Host "[INFO] Conteneurs arretes disponibles:" -ForegroundColor Cyan
            docker ps -a --filter "status=exited" --format "  - {{.Names}}" 2>$null
        }
    }
    
    "stop" {
        if ($Container) {
            Write-Host "[....] Arret de $Container" -NoNewline
            docker stop $Container 2>$null
            if ($LASTEXITCODE -eq 0) {
                Write-Host "`r[ OK ] Arret de $Container" -ForegroundColor Green
            } else {
                Write-Host "`r[FAIL] Arret de $Container" -ForegroundColor Red
            }
        } else {
            Write-Host "[INFO] Arret de tous les conteneurs..." -ForegroundColor Cyan
            docker stop $(docker ps -q) 2>$null
            Write-Host "[ OK ] Tous les conteneurs arretes" -ForegroundColor Green
        }
    }
    
    "clean" {
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host " NETTOYAGE DOCKER" -ForegroundColor White
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host ""
        
        Write-Host "[....] Suppression des conteneurs arretes" -NoNewline
        docker container prune -f 2>$null | Out-Null
        Write-Host "`r[ OK ] Suppression des conteneurs arretes" -ForegroundColor Green
        
        Write-Host "[....] Suppression des images non utilisees" -NoNewline
        docker image prune -f 2>$null | Out-Null
        Write-Host "`r[ OK ] Suppression des images non utilisees" -ForegroundColor Green
        
        Write-Host "[....] Suppression des volumes non utilises" -NoNewline
        docker volume prune -f 2>$null | Out-Null
        Write-Host "`r[ OK ] Suppression des volumes non utilises" -ForegroundColor Green
        
        Write-Host "[....] Suppression des networks non utilises" -NoNewline
        docker network prune -f 2>$null | Out-Null
        Write-Host "`r[ OK ] Suppression des networks non utilises" -ForegroundColor Green
        
        Write-Host ""
        
        $spaceFreed = docker system df 2>$null
        Write-Host "[INFO] Espace Docker:" -ForegroundColor Cyan
        $spaceFreed | ForEach-Object { Write-Host "  $_" }
    }
    
    "logs" {
        if ($Container) {
            Write-Host "[INFO] Logs de $Container (dernières 50 lignes):" -ForegroundColor Cyan
            Write-Host ""
            docker logs --tail 50 $Container 2>&1
        } else {
            Write-Host "[INFO] Conteneurs disponibles:" -ForegroundColor Cyan
            docker ps --format "  - {{.Names}}" 2>$null
            Write-Host ""
            Write-Host "Usage: .\docker-helper.ps1 logs -Container <nom>" -ForegroundColor DarkGray
        }
    }
}

Write-Host ""
