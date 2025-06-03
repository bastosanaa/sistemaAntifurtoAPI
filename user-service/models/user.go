package models

type User struct {
    ID    int64  `json:"id"`    // ID (chave primária no banco)
    Name  string `json:"name"`  // Nome completo do usuário
    Phone string `json:"phone"` // Telefone do usuário (deverá ser único)
}