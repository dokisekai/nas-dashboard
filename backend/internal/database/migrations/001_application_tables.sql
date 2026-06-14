-- 应用包管理相关表

-- 应用实例表
CREATE TABLE IF NOT EXISTS app_instances (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255),
    package_name VARCHAR(255) NOT NULL,
    version VARCHAR(50),
    description TEXT,
    category VARCHAR(50),
    author VARCHAR(255),
    website VARCHAR(500),
    status VARCHAR(50) DEFAULT 'stopped',
    container_id VARCHAR(255),
    pid INTEGER,
    exit_code INTEGER DEFAULT 0,
    last_exit_time TIMESTAMP,
    config TEXT,
    env_vars TEXT,
    ports TEXT,
    volumes TEXT,
    resources TEXT,
    permissions TEXT,
    install_path VARCHAR(500),
    data_path VARCHAR(500),
    config_path VARCHAR(500),
    backup_paths TEXT
);

-- 应用包表
CREATE TABLE IF NOT EXISTS app_packages (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255),
    version VARCHAR(50),
    description TEXT,
    author VARCHAR(255),
    website VARCHAR(500),
    category VARCHAR(50),
    license VARCHAR(100),
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT,
    file_hash VARCHAR(255),
    download_url VARCHAR(500),
    architecture VARCHAR(50),
    min_os_version VARCHAR(50),
    max_os_version VARCHAR(50),
    min_ram INTEGER DEFAULT 0,
    min_disk_space INTEGER DEFAULT 0,
    dependencies TEXT,
    install_path VARCHAR(500),
    data_path VARCHAR(500),
    config_path VARCHAR(500),
    backup_paths TEXT,
    resources TEXT,
    permissions TEXT,
    repository_id INTEGER REFERENCES app_repositories(id),
    download_count INTEGER DEFAULT 0,
    install_count INTEGER DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0.00
);

-- 应用仓库表
CREATE TABLE IF NOT EXISTS app_repositories (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(255) UNIQUE NOT NULL,
    url VARCHAR(500) NOT NULL,
    type VARCHAR(50) DEFAULT 'custom',
    enabled BOOLEAN DEFAULT true,
    priority INTEGER DEFAULT 0,
    description TEXT,
    auto_update BOOLEAN DEFAULT false
);

-- 应用安装日志表
CREATE TABLE IF NOT EXISTS app_install_logs (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    app_instance_id INTEGER NOT NULL REFERENCES app_instances(id),
    step VARCHAR(100),
    message TEXT,
    status VARCHAR(50),
    details TEXT
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_app_instances_name ON app_instances(name);
CREATE INDEX IF NOT EXISTS idx_app_instances_status ON app_instances(status);
CREATE INDEX IF NOT EXISTS idx_app_instances_package_name ON app_instances(package_name);

CREATE INDEX IF NOT EXISTS idx_app_packages_name ON app_packages(name);
CREATE INDEX IF NOT EXISTS idx_app_packages_category ON app_packages(category);
CREATE INDEX IF NOT EXISTS idx_app_packages_repository_id ON app_packages(repository_id);

CREATE INDEX IF NOT EXISTS idx_app_repositories_enabled ON app_repositories(enabled);
CREATE INDEX IF NOT EXISTS idx_app_repositories_type ON app_repositories(type);

CREATE INDEX IF NOT EXISTS idx_app_install_logs_app_instance_id ON app_install_logs(app_instance_id);
CREATE INDEX IF NOT EXISTS idx_app_install_logs_status ON app_install_logs(status);

-- 插入默认应用仓库
INSERT INTO app_repositories (name, url, type, enabled, priority, description, auto_update)
VALUES
    ('官方仓库', 'https://apps.nas-dashboard.com/official', 'official', true, 100, '官方应用仓库', true),
    ('社区仓库', 'https://apps.nas-dashboard.com/community', 'community', true, 50, '社区应用仓库', true)
ON CONFLICT (name) DO NOTHING;