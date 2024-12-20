param(
	[string]$Action,
	[string]$SassCommand
)

$pidFile = "setup/.pid"

function Get-SassWatcher {
	if (Test-Path $pidFile) {
		$watcherPid = Get-Content $pidFile
		$process = Get-Process -Id $watcherPid -ErrorAction SilentlyContinue
		return $process
	}
	return $null
}

function Start-SassWatcher {
	param(
		[string]$SassCommand
	)

	$process = Get-SassWatcher
	if ($null -eq $process) {
		$process = Start-Process -PassThru -NoNewWindow `
		-FilePath "cmd.exe" `
		-ArgumentList "/c $SassCommand"
		$process.Id | Out-File -FilePath $pidFile
		Write-Output "Sass watcher started with PID: $($process.Id)"
	} else {
		Write-Output "Sass watcher still running with PID: $($process.Id)"
	}
}

function Stop-SassWatcher {
	$process = Get-SassWatcher
	if ($null -eq $process) {
		Write-Output "No sass watcher process found with PID $($process.Id)."
	}
	else {
		Stop-Process -Id $process.Id -Force
		Write-Output "Sass watcher process stopped with PID: $($process.Id)"
	}
}

if ($Action -eq "start") {
	Start-SassWatcher -SassCommand $SassCommand
}
elseif ($Action -eq "stop") {
	Stop-SassWatcher
}
