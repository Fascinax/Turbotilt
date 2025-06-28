# Script PowerShell pour démarrer tous les services manuellement
# Utiliser ce script pour lancer tous les microservices sans Docker

Write-Host "Démarrage des services Microservices Java-Angular Demo..." -ForegroundColor Cyan

# Vérifier si RabbitMQ est en cours d'exécution, sinon le démarrer
$rabbitMQRunning = Get-NetTCPConnection -LocalPort 5672 -ErrorAction SilentlyContinue

if (-not $rabbitMQRunning) {
    Write-Host "RabbitMQ n'est pas en cours d'exécution. Veuillez démarrer RabbitMQ avant de continuer." -ForegroundColor Yellow
    Write-Host "Vous pouvez utiliser Docker : docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management" -ForegroundColor Yellow
    $startRabbitMQ = Read-Host "Voulez-vous démarrer RabbitMQ avec Docker? (O/N)"
    
    if ($startRabbitMQ -eq "O" -or $startRabbitMQ -eq "o") {
        Write-Host "Démarrage de RabbitMQ avec Docker..." -ForegroundColor Green
        docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
        # Attendre que RabbitMQ démarre
        Write-Host "Attente du démarrage de RabbitMQ..." -ForegroundColor Yellow
        Start-Sleep -Seconds 15
    } else {
        Write-Host "Veuillez démarrer RabbitMQ manuellement et relancer ce script." -ForegroundColor Red
        exit
    }
}

# Démarrer le service utilisateur (Spring Boot)
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd $PSScriptRoot\user-service; Write-Host 'Démarrage du service utilisateur (Spring Boot)...' -ForegroundColor Green; mvn spring-boot:run"

# Attendre un peu avant de démarrer le service suivant
Start-Sleep -Seconds 5

# Démarrer le service produit (Quarkus)
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd $PSScriptRoot\product-service; Write-Host 'Démarrage du service produit (Quarkus)...' -ForegroundColor Green; mvn quarkus:dev"

# Attendre un peu avant de démarrer le service suivant
Start-Sleep -Seconds 5

# Démarrer le service commande (Micronaut)
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd $PSScriptRoot\order-service; Write-Host 'Démarrage du service commande (Micronaut)...' -ForegroundColor Green; mvn mn:run"

# Attendre un peu avant de démarrer le frontend
Start-Sleep -Seconds 5

# Démarrer le frontend Angular
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd $PSScriptRoot\frontend; Write-Host 'Démarrage du frontend Angular...' -ForegroundColor Green; npm install; npm start"

Write-Host "Tous les services ont été démarrés!" -ForegroundColor Cyan
Write-Host "URLs des services:" -ForegroundColor Green
Write-Host "- Frontend Angular: http://localhost:4200" -ForegroundColor White
Write-Host "- Service utilisateur (Spring Boot): http://localhost:8081/api/users" -ForegroundColor White
Write-Host "- Service produit (Quarkus): http://localhost:8082/api/products" -ForegroundColor White
Write-Host "- Service commande (Micronaut): http://localhost:8083/api/orders" -ForegroundColor White
Write-Host "- RabbitMQ Management: http://localhost:15672 (guest/guest)" -ForegroundColor White
