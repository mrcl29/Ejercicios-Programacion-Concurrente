# L’estanquer il·legal - Taller 3 Programació Concurrent

Marc Llobera Villalonga

## Com Executar

-----------------------------------------------

(Si no conté go.mod i go.sum)

go mod init exercici3
go mod tidy
go get github.com/rabbitmq/amqp091-go

-----------------------------------------------

go run ./estanquer.go
go run ./fumadorTabac.go
go run ./fumadorMistos.go
go run ./fumadorXivato.go
