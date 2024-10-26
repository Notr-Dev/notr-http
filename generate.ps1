# Find all subdirectories containing a Go file
$directories = Get-ChildItem -Recurse -Filter *.go | ForEach-Object { $_.DirectoryName } | Sort-Object -Unique

# Loop through each directory and run 'go generate'
foreach ($dir in $directories) {
    Write-Output "Running 'go generate' in $dir"
    Set-Location $dir
    go generate
    Set-Location -  # Return to the previous directory
}
