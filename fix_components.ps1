$files = Get-ChildItem -Path frontend\src\components\atoms -Recurse -Filter *.vue
foreach ($file in $files) {
    $name = $file.BaseName
    $content = [System.IO.File]::ReadAllText($file.FullName)
    $importPattern = "import\s+\{[^}]*\b$name\b(?!\s+as\s+)[^}]*\}\s+from\s+['""]reka-ui['""]"
    if ($content -match $importPattern) {
        Write-Host "Fixing naming conflict: $name in $($file.FullName)"
        # Alias the import
        $content = $content -replace "(?<!as\s+)\b$name\b(?=[^}]*\}\s+from\s+['""]reka-ui['""])", "$name as Reka$name"
        # Update template tags
        $content = $content -replace "<$name(\b|(?=>))", "<Reka$name"
        $content = $content -replace "</$name>", "</Reka$name>"
        [System.IO.File]::WriteAllText($file.FullName, $content)
    }
}
