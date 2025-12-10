# clean-go – Scaffold interactivo para servicios Go con Clean Architecture

Este repositorio contiene un script que ayuda a crear la estructura base de un
servicio en Go, siguiendo una arquitectura limpia y con opciones de:

- Framework HTTP (`net/http`, `chi`, `gin`, `fiber`)
- Base de datos principal (Postgres, MySQL, MongoDB, Oracle)
- Extras opcionales (Redis, Kafka)
- Logger estructurado con `zap`
- Configuración centralizada via variables de entorno

La idea es similar a los scaffolds de Clean Architecture para Java (como el de
Bancolombia), pero adaptado al ecosistema Go y 100% personalizable.

---

## Requisitos

- Go instalado (1.20 o superior recomendado)  
- Git Bash o WSL en Windows (o cualquier shell bash en Linux/macOS)

Comprueba tu versión de Go:

```bash
go version
