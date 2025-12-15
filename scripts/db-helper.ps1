# ============================================================
# B.DEV Database Helper Script
# Gestion des bases de donnees locales
# ============================================================

param(
    [Parameter(Position=0)]
    [ValidateSet("list", "create", "drop", "backup", "restore")]
    [string]$Action = "list",
    
    [string]$Name,
    [string]$Type = "mysql",
    [string]$File
)

$ErrorActionPreference = "SilentlyContinue"

# Header
Write-Host ""
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  B.DEV DATABASE HELPER" -ForegroundColor Cyan
Write-Host "  $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -ForegroundColor DarkCyan
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""

# Configuration par defaut
$mysqlUser = "root"
$mysqlPass = ""
$mysqlHost = "localhost"
$backupDir = "$HOME\Dev\Databases\backups"

# Creer le dossier de backup si necessaire
if (-not (Test-Path $backupDir)) {
    New-Item -ItemType Directory -Path $backupDir -Force | Out-Null
}

switch ($Action) {
    "list" {
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host " BASES DE DONNEES MYSQL" -ForegroundColor White
        Write-Host "------------------------------------------------------------" -ForegroundColor DarkGray
        Write-Host ""
        
        $databases = mysql -u$mysqlUser -h$mysqlHost -e "SHOW DATABASES;" 2>$null
        
        if ($databases) {
            foreach ($db in $databases) {
                if ($db -notmatch "Database|information_schema|performance_schema|mysql|sys") {
                    Write-Host "  - $db" -ForegroundColor White
                }
            }
        } else {
            Write-Host "[WARN] MySQL n'est pas accessible" -ForegroundColor Yellow
            Write-Host "       Verifiez que MySQL est demarre" -ForegroundColor DarkGray
        }
    }
    
    "create" {
        if (-not $Name) {
            Write-Host "[FAIL] Nom de base de donnees requis" -ForegroundColor Red
            Write-Host "       Usage: .\db-helper.ps1 create -Name mabase" -ForegroundColor DarkGray
            exit 1
        }
        
        Write-Host "[....] Creation de la base '$Name'" -NoNewline
        
        mysql -u$mysqlUser -h$mysqlHost -e "CREATE DATABASE IF NOT EXISTS ``$Name`` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>$null
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "`r[ OK ] Creation de la base '$Name'" -ForegroundColor Green
        } else {
            Write-Host "`r[FAIL] Creation de la base '$Name'" -ForegroundColor Red
        }
    }
    
    "drop" {
        if (-not $Name) {
            Write-Host "[FAIL] Nom de base de donnees requis" -ForegroundColor Red
            exit 1
        }
        
        $confirm = Read-Host "Confirmer la suppression de '$Name'? (oui/non)"
        if ($confirm -ne "oui") {
            Write-Host "[INFO] Annule" -ForegroundColor Yellow
            exit
        }
        
        Write-Host "[....] Suppression de la base '$Name'" -NoNewline
        
        mysql -u$mysqlUser -h$mysqlHost -e "DROP DATABASE IF EXISTS ``$Name``;" 2>$null
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "`r[ OK ] Suppression de la base '$Name'" -ForegroundColor Green
        } else {
            Write-Host "`r[FAIL] Suppression de la base '$Name'" -ForegroundColor Red
        }
    }
    
    "backup" {
        if (-not $Name) {
            Write-Host "[FAIL] Nom de base de donnees requis" -ForegroundColor Red
            exit 1
        }
        
        $timestamp = Get-Date -Format "yyyyMMdd_HHmm"
        $backupFile = "$backupDir\${Name}_$timestamp.sql"
        
        Write-Host "[....] Backup de '$Name'" -NoNewline
        
        mysqldump -u$mysqlUser -h$mysqlHost $Name > $backupFile 2>$null
        
        if ($LASTEXITCODE -eq 0 -and (Test-Path $backupFile)) {
            $size = [math]::Round((Get-Item $backupFile).Length / 1MB, 2)
            Write-Host "`r[ OK ] Backup de '$Name' - $size MB" -ForegroundColor Green
            Write-Host "       Fichier: $backupFile" -ForegroundColor DarkGray
        } else {
            Write-Host "`r[FAIL] Backup de '$Name'" -ForegroundColor Red
        }
    }
    
    "restore" {
        if (-not $Name -or -not $File) {
            Write-Host "[FAIL] Nom de base et fichier requis" -ForegroundColor Red
            Write-Host "       Usage: .\db-helper.ps1 restore -Name mabase -File backup.sql" -ForegroundColor DarkGray
            exit 1
        }
        
        if (-not (Test-Path $File)) {
            Write-Host "[FAIL] Fichier non trouve: $File" -ForegroundColor Red
            exit 1
        }
        
        Write-Host "[....] Restauration de '$Name'" -NoNewline
        
        mysql -u$mysqlUser -h$mysqlHost $Name < $File 2>$null
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "`r[ OK ] Restauration de '$Name'" -ForegroundColor Green
        } else {
            Write-Host "`r[FAIL] Restauration de '$Name'" -ForegroundColor Red
        }
    }
}

Write-Host ""
