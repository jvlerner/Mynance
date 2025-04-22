CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_password_change TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'plus', 'pro')),
    active BOOLEAN DEFAULT TRUE
);

-- Já cobre autenticação e recuperação de dados
CREATE UNIQUE INDEX idx_users_email ON users(email);

-- Acesso frequente por ID
CREATE INDEX idx_users_id ON users(id);

-- Se você quiser listar usuários por status (ativos/inativos)
CREATE INDEX idx_users_active ON users(active);

-- Se você listar usuários por tipo de conta
CREATE INDEX idx_users_role ON users(role);

-- Se gerar relatórios ou ordenar por criação
CREATE INDEX idx_users_created_at ON users(created_at);

