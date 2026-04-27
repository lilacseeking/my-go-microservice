# API Test Script (PowerShell) - Test HTTP endpoints using curl-like commands
# Usage: .\scripts\test-api.ps1 [base_url]
# Default base_url: http://localhost:8080

param(
    [string]$BaseUrl = "http://localhost:8080"
)

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "API Test Script" -ForegroundColor Cyan
Write-Host "Base URL: $BaseUrl" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

# Test counters
$script:Total = 0
$script:Passed = 0
$script:Failed = 0

# Test result function
function Print-Result {
    param([int]$Success)
    $script:Total++
    if ($Success -eq 0) {
        Write-Host "[PASS] Test passed" -ForegroundColor Green
        $script:Passed++
    } else {
        Write-Host "[FAIL] Test failed" -ForegroundColor Red
        $script:Failed++
    }
    Write-Host ""
}

# ==========================================
# Test 1: Health Check
# ==========================================
$healthUrl = $BaseUrl -replace ':8080$', ':8081'
Write-Host "[Test 1] Health Check" -ForegroundColor Yellow
Write-Host "GET $healthUrl/healthz"
try {
    $response = Invoke-RestMethod -Uri "$healthUrl/healthz" -Method Get
    Write-Host "HTTP Status: 200"
    Write-Host "Response: OK"
    Print-Result 0
} catch {
    Write-Host "HTTP Status: $($_.Exception.Response.StatusCode.value__)"
    Write-Host "Error: $_"
    Print-Result 1
}

# ==========================================
# Test 2: Create User - Success
# ==========================================
Write-Host "[Test 2] Create User - Success" -ForegroundColor Yellow
Write-Host "POST $BaseUrl/v1/users"
$createBody = @{
    name = "John Doe"
    email = "john@example.com"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/v1/users" -Method Post `
        -ContentType "application/json" -Body $createBody
    Write-Host "HTTP Status: 201"
    Write-Host "Response: $($response | ConvertTo-Json -Compress)"
    Print-Result 0
    $userId = $response.data.user_id
} catch {
    Write-Host "HTTP Status: $($_.Exception.Response.StatusCode.value__)"
    Write-Host "Error: $_"
    Print-Result 1
    $userId = $null
}

# ==========================================
# Test 3: Create User - Validation Error
# ==========================================
Write-Host "[Test 3] Create User - Invalid Data" -ForegroundColor Yellow
Write-Host "POST $BaseUrl/v1/users"
$invalidBody = @{
    name = ""
    email = "invalid-email"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/v1/users" -Method Post `
        -ContentType "application/json" -Body $invalidBody
    Write-Host "HTTP Status: 200 (expected 400)"
    Write-Host "Response: $($response | ConvertTo-Json -Compress)"
    Print-Result 1
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    Write-Host "HTTP Status: $statusCode"
    if ($statusCode -eq 400) {
        Print-Result 0
    } else {
        Print-Result 1
    }
}

# ==========================================
# Test 4: Get User - Test Data
# ==========================================
Write-Host "[Test 4] Get User - test-user-001" -ForegroundColor Yellow
Write-Host "GET $BaseUrl/v1/users/test-user-001"
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/v1/users/test-user-001" -Method Get
    Write-Host "HTTP Status: 200"
    Write-Host "Response: $($response | ConvertTo-Json -Compress)"
    Print-Result 0
} catch {
    Write-Host "HTTP Status: $($_.Exception.Response.StatusCode.value__)"
    Write-Host "Error: $_"
    Print-Result 1
}

# ==========================================
# Test 5: Get User - Newly Created
# ==========================================
if ($userId) {
    Write-Host "[Test 5] Get User - Newly Created" -ForegroundColor Yellow
    Write-Host "GET $BaseUrl/v1/users/$userId"
    try {
        $response = Invoke-RestMethod -Uri "$BaseUrl/v1/users/$userId" -Method Get
        Write-Host "HTTP Status: 200"
        Write-Host "Response: $($response | ConvertTo-Json -Compress)"
        Print-Result 0
    } catch {
        Write-Host "HTTP Status: $($_.Exception.Response.StatusCode.value__)"
        Write-Host "Error: $_"
        Print-Result 1
    }
}

# ==========================================
# Test 6: Get User - Not Found
# ==========================================
Write-Host "[Test 6] Get User - Not Found" -ForegroundColor Yellow
Write-Host "GET $BaseUrl/v1/users/non-existent-id"
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/v1/users/non-existent-id" -Method Get
    Write-Host "HTTP Status: 200 (expected 404)"
    Write-Host "Response: $($response | ConvertTo-Json -Compress)"
    Print-Result 1
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    Write-Host "HTTP Status: $statusCode"
    if ($statusCode -eq 404) {
        Print-Result 0
    } else {
        Print-Result 1
    }
}

# ==========================================
# Test 7: Create Second User
# ==========================================
Write-Host "[Test 7] Create Second User" -ForegroundColor Yellow
Write-Host "POST $BaseUrl/v1/users"
$createBody2 = @{
    name = "Jane Smith"
    email = "jane@example.com"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/v1/users" -Method Post `
        -ContentType "application/json" -Body $createBody2
    Write-Host "HTTP Status: 201"
    Write-Host "Response: $($response | ConvertTo-Json -Compress)"
    Print-Result 0
} catch {
    Write-Host "HTTP Status: $($_.Exception.Response.StatusCode.value__)"
    Write-Host "Error: $_"
    Print-Result 1
}

# ==========================================
# Summary
# ==========================================
Write-Host ""
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Test Summary" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Total Tests: $script:Total"
Write-Host "Passed: $script:Passed" -ForegroundColor Green
Write-Host "Failed: $script:Failed" -ForegroundColor Red
Write-Host ""

if ($script:Failed -eq 0) {
    Write-Host "All tests passed!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "Some tests failed" -ForegroundColor Red
    exit 1
}
