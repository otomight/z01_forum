param(
	[string]$EnvExampleFile,
	[string]$EnvFile
)

function Create-Env {
	param (
		[string]$EnvExampleFile,
		[string]$EnvFile
	)

	if (-not (Test-Path -Path $EnvFile)) {
		Copy-Item -Path $EnvExampleFile -Destination $EnvFile
		Write-Host "WARNING: $EnvFile file created. Please fill it with the right values!" -ForegroundColor Yellow
	}
}

Create-Env $EnvExampleFile $EnvFile
