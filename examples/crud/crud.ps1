$url = 'http://localhost:8080/create'

$users = @(
    @{ name = 'John'; id = (New-Guid).Guid }
    @{ name = 'Jane'; id = (New-Guid).Guid }
    @{ name = 'Jack'; id = (New-Guid).Guid }
    @{ name = 'Jill'; id = (New-Guid).Guid }
    @{ name = 'Joe'; id = (New-Guid).Guid }
    @{ name = 'Jenny'; id = (New-Guid).Guid }
)

foreach ($user in $users) {
    Invoke-RestMethod -Uri $url -Method Post -Body $user -TimeoutSec 10
}

$john = $users[0]
$update = "http://localhost:8080/users/$($john.id)/update"


Invoke-RestMethod -Uri $update -Method Post -Body @{ name = 'Doe John' } -TimeoutSec 10