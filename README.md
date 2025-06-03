# Documentação — **user-service**

> Micro-serviço em **Go + Gin** responsável por CRUD de usuários.  
> Porta padrão: **http://localhost:8001**

---

## Sumário

| Método | Rota                     | Descrição                     |
| ------ | ------------------------ | ----------------------------- |
| GET    | `/health`               | Verifica se o serviço está OK |
| POST   | `/users`                | Cria um novo usuário          |
| GET    | `/users/{id}`           | Busca usuário por ID          |
| PUT    | `/users/{id}`           | Atualiza usuário existente    |
| DELETE | `/users/{id}`           | Remove usuário                |

---

