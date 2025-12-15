# ============================================================
# B.DEV Watch Script
# Surveillance continue du systeme
# ============================================================

param(
    [int]$Interval = 5,
    [switch]$Compact
)

$ErrorActionPreference = "SilentlyContinue"

function Get-SystemStats {
    $cpu = (Get-Counter '\Processor(_Total)\% Processor Time' -ErrorAction SilentlyContinue).CounterSamples.CookedValue
    
    $ram = Get-CimInstance Win32_OperatingSystem
    $totalRAM = [math]::Round($ram.TotalVisibleMemorySize / 1MB, 1)
    $freeRAM = [math]::Round($ram.FreePhysicalMemory / 1MB, 1)
    $usedRAM = $totalRAM - $freeRAM
    $ramPercent = [math]::Round(($usedRAM / $totalRAM) * 100, 1)
    
    $disk = Get-PSDrive C
    $diskFree = [math]::Round($disk.Free / 1GB, 1)
    $diskTotal = [math]::Round(($disk.Used + $disk.Free) / 1GB, 1)
    $diskPercent = [math]::Round(($disk.Used / ($disk.Used + $disk.Free)) * 100, 1)
    
    return @{
        CPU = if ($cpu) { [math]::Round($cpu, 1) } else { 0 }
        RAMUsed = $usedRAM
        RAMTotal = $totalRAM
        RAMPercent = $ramPercent
        DiskFree = $diskFree
        DiskTotal = $diskTotal
        DiskPercent = $diskPercent
    }
}

function Get-Bar {
    param([int]$percent, [int]$width = 20)
    $filled = [math]::Round($percent / 100 * $width)
    $empty = $width - $filled
    return ("[" + ("#" * $filled) + ("-" * $empty) + "]")
}

function Get-Color {
    param([int]$percent)
    if ($percent -ge 90) { return "Red" }
    if ($percent -ge 75) { return "Yellow" }
    return "Green"
}

# Header
Clear-Host
Write-Host ""
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  B.DEV SYSTEM WATCH" -ForegroundColor Cyan
Write-Host "  Intervalle: ${Interval}s | Ctrl+C pour quitter" -ForegroundColor DarkCyan
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""

try {
    while ($true) {
        $stats = Get-SystemStats
        $time = Get-Date -Format "HH:mm:ss"
        
        # Effacer les lignes precedentes
        $cursorTop = [Console]::CursorTop
        [Console]::SetCursorPosition(0, $cursorTop - 4)
        
        if ($Compact) {
            Write-Host ("  [{0}] CPU: {1,5}% | RAM: {2,5}% | Disk: {3,5}%" -f $time, $stats.CPU, $stats.RAMPercent, $stats.DiskPercent) -ForegroundColor White
            Write-Host ""
        } else {
            Write-Host ("  {0}" -f $time) -ForegroundColor DarkGray
            Write-Host ""
            
            $cpuBar = Get-Bar $stats.CPU
            Write-Host ("  CPU  {0} {1,5}%" -f $cpuBar, $stats.CPU) -ForegroundColor (Get-Color $stats.CPU)
            
            $ramBar = Get-Bar $stats.RAMPercent
            Write-Host ("  RAM  {0} {1,5}% ({2}/{3} GB)" -f $ramBar, $stats.RAMPercent, $stats.RAMUsed, $stats.RAMTotal) -ForegroundColor (Get-Color $stats.RAMPercent)
            
            $diskBar = Get-Bar $stats.DiskPercent
            Write-Host ("  Disk {0} {1,5}% ({2} GB libre)" -f $diskBar, $stats.DiskPercent, $stats.DiskFree) -ForegroundColor (Get-Color $stats.DiskPercent)
            
            Write-Host ""
        }
        
        Start-Sleep -Seconds $Interval
    }
} finally {
    Write-Host ""
    Write-Host "[INFO] Surveillance arretee" -ForegroundColor Cyan
}
