# ============================================================
# B.DEV Port Manager Script
# Gestion des ports et processus
# ============================================================

param(
    [Parameter(Position=0)]
    [ValidateSet("list", "find", "kill")]
    [string]$Action = "list",
    
    [int]$Port
)

$ErrorActionPreference = "SilentlyContinue"

# Header
Write-Host ""
Write-Host "============================================================" -ForegroundColor Yellow
Write-Host "  B.DEV PORT MANAGER" -ForegroundColor Yellow
Write-Host "  $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -ForegroundColor DarkYellow
Write-Host "============================================================" -ForegroundColor Yellow
Write-Host ""

switch ($Action) {
    "list" {
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host " PORTS EN ECOUTE" -ForegroundColor White
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host ""
        
        $connections = Get-NetTCPConnection -State Listen -ErrorAction SilentlyContinue | 
            Select-Object LocalPort, OwningProcess, @{n='Process';e={(Get-Process -Id $_.OwningProcess -ErrorAction SilentlyContinue).Name}} |
            Sort-Object LocalPort |
            Where-Object { $_.LocalPort -lt 65000 }
        
        # Ports communs de dev
        $devPorts = @(80, 443, 3000, 3001, 4200, 5000, 5173, 5174, 8000, 8080, 8888, 9000)
        
        Write-Host ("  {0,-8} {1,-25} {2}" -f "PORT", "PROCESS", "PID") -ForegroundColor DarkGray
        Write-Host ("  " + "-" * 50) -ForegroundColor DarkGray
        
        foreach ($conn in $connections) {
            if ($devPorts -contains $conn.LocalPort) {
                $color = "Yellow"
            } elseif ($conn.Process -match "node|php|python|docker|nginx|apache") {
                $color = "Cyan"
            } else {
                $color = "White"
            }
            
            Write-Host ("  {0,-8} {1,-25} {2}" -f $conn.LocalPort, $conn.Process, $conn.OwningProcess) -ForegroundColor $color
        }
        
        Write-Host ""
        Write-Host "[INFO] Ports de developpement en jaune" -ForegroundColor DarkGray
    }
    
    "find" {
        if (-not $Port) {
            Write-Host "[FAIL] Port requis" -ForegroundColor Red
            Write-Host "       Usage: .\port-manager.ps1 find -Port 3000" -ForegroundColor DarkGray
            exit 1
        }
        
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host " RECHERCHE DU PORT $Port" -ForegroundColor White
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host ""
        
        $connection = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue | Select-Object -First 1
        
        if ($connection) {
            $process = Get-Process -Id $connection.OwningProcess -ErrorAction SilentlyContinue
            
            Write-Host "  Port      : $Port" -ForegroundColor Yellow
            Write-Host "  Status    : $($connection.State)" -ForegroundColor White
            Write-Host "  Process   : $($process.Name)" -ForegroundColor White
            Write-Host "  PID       : $($connection.OwningProcess)" -ForegroundColor White
            Write-Host "  Chemin    : $($process.Path)" -ForegroundColor DarkGray
            
            if ($process.StartTime) {
                Write-Host "  Demarre   : $($process.StartTime)" -ForegroundColor DarkGray
            }
        } else {
            Write-Host "[ OK ] Port $Port est libre" -ForegroundColor Green
        }
    }
    
    "kill" {
        if (-not $Port) {
            Write-Host "[FAIL] Port requis" -ForegroundColor Red
            Write-Host "       Usage: .\port-manager.ps1 kill -Port 3000" -ForegroundColor DarkGray
            exit 1
        }
        
        $connection = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue | Select-Object -First 1
        
        if ($connection) {
            $process = Get-Process -Id $connection.OwningProcess -ErrorAction SilentlyContinue
            
            Write-Host "[INFO] Process trouve: $($process.Name) (PID: $($connection.OwningProcess))" -ForegroundColor Cyan
            
            $confirm = Read-Host "Terminer ce processus? (o/n)"
            if ($confirm -eq "o") {
                Stop-Process -Id $connection.OwningProcess -Force -ErrorAction SilentlyContinue
                
                Start-Sleep -Milliseconds 500
                
                $check = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue
                if (-not $check) {
                    Write-Host "[ OK ] Port $Port libere" -ForegroundColor Green
                } else {
                    Write-Host "[WARN] Le processus peut encore etre en cours d'arret" -ForegroundColor Yellow
                }
            } else {
                Write-Host "[INFO] Annule" -ForegroundColor Yellow
            }
        } else {
            Write-Host "[ OK ] Port $Port est deja libre" -ForegroundColor Green
        }
    }
}

Write-Host ""
