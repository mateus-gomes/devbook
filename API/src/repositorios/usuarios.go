package repositorios

import (
	"api/src/models"
	"database/sql"
	"fmt"
)

//Representa um repositório de usuarios
type Usuarios struct {
	db *sql.DB
}

//Cria um repositório de usuários
func NovoRepositorioDeUsuarios(db *sql.DB) *Usuarios {
	return &Usuarios{db}
}

//Insere um usuário no banco de dados
func (repositorio Usuarios) Criar(usuario models.Usuario) (uint64, error) {
	statement, erro := repositorio.db.Prepare(
		"insert into usuarios (nome, nick, email, senha) values (?, ?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

//Busca todos usuarios no banco de dados
func (repositorio Usuarios) Buscar(nomeOuNick string) ([]models.Usuario, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick) //%nomeOuNick%

	linhas, erro := repositorio.db.Query(
		"SELECT id, nome, nick, email, criadoEm FROM usuarios WHERE nome LIKE ? OR nick LIKE ?",
		nomeOuNick, nomeOuNick,
	)

	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []models.Usuario

	for linhas.Next() {
		var usuario models.Usuario

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

//Busca um usuario no banco de dados
func (repositorio Usuarios) BuscarPorID(ID uint64) (models.Usuario, error) {
	linhas, erro := repositorio.db.Query(
		"SELECT id, nome, nick, email, criadoEm FROM usuarios WHERE id = ?",
		ID,
	)

	if erro != nil {
		return models.Usuario{}, erro
	}
	defer linhas.Close()

	var usuario models.Usuario

	if linhas.Next() {
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return models.Usuario{}, erro
		}
	}

	return usuario, nil
}

//Atualiza os dados de um usuario no banco de dados
func (repositorio Usuarios) Atualizar(ID uint64, usuario models.Usuario) error {
	statement, erro := repositorio.db.Prepare(
		"UPDATE usuarios SET nome = ?, nick = ?, email = ? where id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, ID); erro != nil {
		return erro
	}

	return nil
}

//Deleta os dados de um usuário no banco de dados
func (repositorio Usuarios) Deletar(ID uint64) error {
	statement, erro := repositorio.db.Prepare("DELETE FROM usuarios WHERE id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(ID); erro != nil {
		return erro
	}

	return nil
}

//Busca usuario pelo email e retorna seu id e senha com hash
func (repositorio Usuarios) BuscarPorEmail(email string) (models.Usuario, error) {
	linha, erro := repositorio.db.Query("SELECT id, senha FROM usuarios WHERE email = ?", email)
	if erro != nil {
		return models.Usuario{}, erro
	}
	defer linha.Close()

	var usuario models.Usuario

	if linha.Next() {
		if erro = linha.Scan(&usuario.ID, &usuario.Senha); erro != nil {
			return models.Usuario{}, erro
		}
	}

	return usuario, nil
}

//Permite que você siga outro usuário
func (repositorio Usuarios) Seguir(usuarioID, seguidorID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"INSERT IGNORE INTO seguidores (usuario_id, seguidor_id) VALUES (?, ?)",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(usuarioID, seguidorID); erro != nil {
		return erro
	}

	return nil
}

//Permite que você pare de seguir outro usuário
func (repositorio Usuarios) PararDeSeguir(usuarioID, seguidorID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"DELETE FROM seguidores WHERE usuario_id = ? AND seguidor_id = ?",
	)
	if erro != nil {
		return erro
	}

	if _, erro = statement.Exec(usuarioID, seguidorID); erro != nil {
		return erro
	}

	return nil
}

//Traz todos seguidores de um usuário
func (repositorio Usuarios) BuscarSeguidores(usuarioID uint64) ([]models.Usuario, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT u.id, u.nome, u.nick, u.email, u.criadoEm FROM usuarios u
		INNER JOIN seguidores s ON u.id = s.seguidor_id where s.usuario_id = ?
	`, usuarioID)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []models.Usuario
	for linhas.Next() {
		var usuario models.Usuario

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

//Traz todos os usuários que um usuário está seguindo
func (repositorio Usuarios) BuscarSeguindo(usuarioID uint64) ([]models.Usuario, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT u.id, u.nome, u.nick, u.email, u.criadoEm FROM usuarios u
		INNER JOIN seguidores s on u.id = s.usuario_id WHERE s.seguidor_id = ?`,
		usuarioID,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []models.Usuario
	for linhas.Next() {
		var usuario models.Usuario

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

//Traz a senha de um usuário pelo ID
func (repositorio Usuarios) BuscarSenha(usuarioID uint64) (string, error) {
	linha, erro := repositorio.db.Query("SELECT senha FROM usuarios WHERE id = ?", usuarioID)
	if erro != nil {
		return "", erro
	}
	defer linha.Close()

	var usuario models.Usuario

	if linha.Next() {
		if erro = linha.Scan(&usuario.Senha); erro != nil {
			return "", erro
		}
	}

	return usuario.Senha, nil
}

//Atualiza a senha do usuário
func (repositorio Usuarios) AtualizarSenha(usuarioID uint64, senha string) error {
	statemente, erro := repositorio.db.Prepare("UPDATE usuarios SET senha = ? WHERE id = ?")
	if erro != nil {
		return erro
	}
	defer statemente.Close()

	if _, erro = statemente.Exec(senha, usuarioID); erro != nil {
		return erro
	}

	return nil
}
