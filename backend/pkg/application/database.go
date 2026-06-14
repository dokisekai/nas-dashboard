package application

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// PostgreSQLDatabase PostgreSQL数据库实现
type PostgreSQLDatabase struct {
	db *sql.DB
}

// NewPostgreSQLDatabase 创建PostgreSQL数据库连接
func NewPostgreSQLDatabase(host, port, user, password, dbname string) (*PostgreSQLDatabase, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("打开数据库连接失败: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	return &PostgreSQLDatabase{db: db}, nil
}

// Close 关闭数据库连接
func (p *PostgreSQLDatabase) Close() error {
	return p.db.Close()
}

// CreateAppInstance 创建应用实例
func (p *PostgreSQLDatabase) CreateAppInstance(instance *AppInstance) error {
	instance.CreatedAt = time.Now()
	instance.UpdatedAt = time.Now()

	query := `
		INSERT INTO app_instances (
			name, display_name, package_name, version, description, category, author, website,
			status, container_id, pid, exit_code, last_exit_time,
			config, env_vars, ports, volumes, resources, permissions,
			install_path, data_path, config_path, backup_paths,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12, $13,
			$14, $15, $16, $17, $18, $19,
			$20, $21, $22, $23,
			$24, $25
		) RETURNING id
	`

	err := p.db.QueryRow(
		query,
		instance.Name, instance.DisplayName, instance.PackageName, instance.Version,
		instance.Description, instance.Category, instance.Author, instance.Website,
		instance.Status, instance.ContainerID, instance.PID, instance.ExitCode, instance.LastExitTime,
		instance.Config, instance.EnvVars, instance.Ports, instance.Volumes, instance.Resources, instance.Permissions,
		instance.InstallPath, instance.DataPath, instance.ConfigPath, instance.BackupPaths,
		instance.CreatedAt, instance.UpdatedAt,
	).Scan(&instance.ID)

	return err
}

// GetAppInstance 获取应用实例
func (p *PostgreSQLDatabase) GetAppInstance(id uint) (*AppInstance, error) {
	var instance AppInstance

	query := `
		SELECT id, created_at, updated_at, name, display_name, package_name, version, description,
		       category, author, website, status, container_id, pid, exit_code, last_exit_time,
		       config, env_vars, ports, volumes, resources, permissions,
		       install_path, data_path, config_path, backup_paths
		FROM app_instances WHERE id = $1
	`

	err := p.db.QueryRow(query, id).Scan(
		&instance.ID, &instance.CreatedAt, &instance.UpdatedAt,
		&instance.Name, &instance.DisplayName, &instance.PackageName, &instance.Version,
		&instance.Description, &instance.Category, &instance.Author, &instance.Website,
		&instance.Status, &instance.ContainerID, &instance.PID, &instance.ExitCode, &instance.LastExitTime,
		&instance.Config, &instance.EnvVars, &instance.Ports, &instance.Volumes, &instance.Resources, &instance.Permissions,
		&instance.InstallPath, &instance.DataPath, &instance.ConfigPath, &instance.BackupPaths,
	)

	if err != nil {
		return nil, err
	}

	return &instance, nil
}

// GetAppInstanceByName 根据名称获取应用实例
func (p *PostgreSQLDatabase) GetAppInstanceByName(name string) (*AppInstance, error) {
	var instance AppInstance

	query := `
		SELECT id, created_at, updated_at, name, display_name, package_name, version, description,
		       category, author, website, status, container_id, pid, exit_code, last_exit_time,
		       config, env_vars, ports, volumes, resources, permissions,
		       install_path, data_path, config_path, backup_paths
		FROM app_instances WHERE name = $1
	`

	err := p.db.QueryRow(query, name).Scan(
		&instance.ID, &instance.CreatedAt, &instance.UpdatedAt,
		&instance.Name, &instance.DisplayName, &instance.PackageName, &instance.Version,
		&instance.Description, &instance.Category, &instance.Author, &instance.Website,
		&instance.Status, &instance.ContainerID, &instance.PID, &instance.ExitCode, &instance.LastExitTime,
		&instance.Config, &instance.EnvVars, &instance.Ports, &instance.Volumes, &instance.Resources, &instance.Permissions,
		&instance.InstallPath, &instance.DataPath, &instance.ConfigPath, &instance.BackupPaths,
	)

	if err != nil {
		return nil, err
	}

	return &instance, nil
}

// ListAppInstances 列出所有应用实例
func (p *PostgreSQLDatabase) ListAppInstances() ([]AppInstance, error) {
	var instances []AppInstance

	query := `
		SELECT id, created_at, updated_at, name, display_name, package_name, version, description,
		       category, author, website, status, container_id, pid, exit_code, last_exit_time,
		       config, env_vars, ports, volumes, resources, permissions,
		       install_path, data_path, config_path, backup_paths
		FROM app_instances ORDER BY created_at DESC
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var instance AppInstance
		err := rows.Scan(
			&instance.ID, &instance.CreatedAt, &instance.UpdatedAt,
			&instance.Name, &instance.DisplayName, &instance.PackageName, &instance.Version,
			&instance.Description, &instance.Category, &instance.Author, &instance.Website,
			&instance.Status, &instance.ContainerID, &instance.PID, &instance.ExitCode, &instance.LastExitTime,
			&instance.Config, &instance.EnvVars, &instance.Ports, &instance.Volumes, &instance.Resources, &instance.Permissions,
			&instance.InstallPath, &instance.DataPath, &instance.ConfigPath, &instance.BackupPaths,
		)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// UpdateAppInstance 更新应用实例
func (p *PostgreSQLDatabase) UpdateAppInstance(instance *AppInstance) error {
	instance.UpdatedAt = time.Now()

	query := `
		UPDATE app_instances SET
			display_name = $2, package_name = $3, version = $4, description = $5, category = $6,
			author = $7, website = $8, status = $9, container_id = $10, pid = $11, exit_code = $12,
			last_exit_time = $13, config = $14, env_vars = $15, ports = $16, volumes = $17,
			resources = $18, permissions = $19, install_path = $20, data_path = $21,
			config_path = $22, backup_paths = $23, updated_at = $24
		WHERE id = $1
	`

	_, err := p.db.Exec(
		query,
		instance.ID, instance.DisplayName, instance.PackageName, instance.Version,
		instance.Description, instance.Category, instance.Author, instance.Website,
		instance.Status, instance.ContainerID, instance.PID, instance.ExitCode, instance.LastExitTime,
		instance.Config, instance.EnvVars, instance.Ports, instance.Volumes, instance.Resources, instance.Permissions,
		instance.InstallPath, instance.DataPath, instance.ConfigPath, instance.BackupPaths,
		instance.UpdatedAt,
	)

	return err
}

// DeleteAppInstance 删除应用实例
func (p *PostgreSQLDatabase) DeleteAppInstance(id uint) error {
	query := `DELETE FROM app_instances WHERE id = $1`
	_, err := p.db.Exec(query, id)
	return err
}

// CreateAppPackage 创建应用包记录
func (p *PostgreSQLDatabase) CreateAppPackage(pkg *AppPackage) error {
	pkg.CreatedAt = time.Now()
	pkg.UpdatedAt = time.Now()

	query := `
		INSERT INTO app_packages (
			name, display_name, version, description, author, website, category, license,
			file_path, file_size, file_hash, download_url, architecture, min_os_version,
			max_os_version, min_ram, min_disk_space, dependencies, install_path, data_path,
			config_path, backup_paths, resources, permissions, repository_id,
			download_count, install_count, rating, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25,
			$26, $27, $28, $29, $30
		) RETURNING id
	`

	err := p.db.QueryRow(
		query,
		pkg.Name, pkg.DisplayName, pkg.Version, pkg.Description, pkg.Author, pkg.Website,
		pkg.Category, pkg.License, pkg.FilePath, pkg.FileSize, pkg.FileHash, pkg.DownloadURL,
		pkg.Architecture, pkg.MinOSVersion, pkg.MaxOSVersion, pkg.MinRAM, pkg.MinDiskSpace,
		pkg.Dependencies, pkg.InstallPath, pkg.DataPath, pkg.ConfigPath, pkg.BackupPaths,
		pkg.Resources, pkg.Permissions, pkg.RepositoryID, pkg.DownloadCount, pkg.InstallCount,
		pkg.Rating, pkg.CreatedAt, pkg.UpdatedAt,
	).Scan(&pkg.ID)

	return err
}

// GetAppPackage 获取应用包
func (p *PostgreSQLDatabase) GetAppPackage(name string) (*AppPackage, error) {
	var pkg AppPackage

	query := `
		SELECT id, created_at, updated_at, name, display_name, version, description, author,
		       website, category, license, file_path, file_size, file_hash, download_url,
		       architecture, min_os_version, max_os_version, min_ram, min_disk_space,
		       dependencies, install_path, data_path, config_path, backup_paths,
		       resources, permissions, repository_id, download_count, install_count, rating
		FROM app_packages WHERE name = $1
	`

	err := p.db.QueryRow(query, name).Scan(
		&pkg.ID, &pkg.CreatedAt, &pkg.UpdatedAt, &pkg.Name, &pkg.DisplayName, &pkg.Version,
		&pkg.Description, &pkg.Author, &pkg.Website, &pkg.Category, &pkg.License,
		&pkg.FilePath, &pkg.FileSize, &pkg.FileHash, &pkg.DownloadURL, &pkg.Architecture,
		&pkg.MinOSVersion, &pkg.MaxOSVersion, &pkg.MinRAM, &pkg.MinDiskSpace, &pkg.Dependencies,
		&pkg.InstallPath, &pkg.DataPath, &pkg.ConfigPath, &pkg.BackupPaths, &pkg.Resources,
		&pkg.Permissions, &pkg.RepositoryID, &pkg.DownloadCount, &pkg.InstallCount, &pkg.Rating,
	)

	if err != nil {
		return nil, err
	}

	return &pkg, nil
}

// ListAppPackages 列出所有应用包
func (p *PostgreSQLDatabase) ListAppPackages() ([]AppPackage, error) {
	var packages []AppPackage

	query := `
		SELECT id, created_at, updated_at, name, display_name, version, description, author,
		       website, category, license, file_path, file_size, file_hash, download_url,
		       architecture, min_os_version, max_os_version, min_ram, min_disk_space,
		       dependencies, install_path, data_path, config_path, backup_paths,
		       resources, permissions, repository_id, download_count, install_count, rating
		FROM app_packages ORDER BY name
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pkg AppPackage
		err := rows.Scan(
			&pkg.ID, &pkg.CreatedAt, &pkg.UpdatedAt, &pkg.Name, &pkg.DisplayName, &pkg.Version,
			&pkg.Description, &pkg.Author, &pkg.Website, &pkg.Category, &pkg.License,
			&pkg.FilePath, &pkg.FileSize, &pkg.FileHash, &pkg.DownloadURL, &pkg.Architecture,
			&pkg.MinOSVersion, &pkg.MaxOSVersion, &pkg.MinRAM, &pkg.MinDiskSpace, &pkg.Dependencies,
			&pkg.InstallPath, &pkg.DataPath, &pkg.ConfigPath, &pkg.BackupPaths, &pkg.Resources,
			&pkg.Permissions, &pkg.RepositoryID, &pkg.DownloadCount, &pkg.InstallCount, &pkg.Rating,
		)
		if err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}

// UpdateAppPackage 更新应用包
func (p *PostgreSQLDatabase) UpdateAppPackage(pkg *AppPackage) error {
	pkg.UpdatedAt = time.Now()

	query := `
		UPDATE app_packages SET
			display_name = $2, version = $3, description = $4, author = $5, website = $6,
			category = $7, license = $8, file_path = $9, file_size = $10, file_hash = $11,
			download_url = $12, architecture = $13, min_os_version = $14, max_os_version = $15,
			min_ram = $16, min_disk_space = $17, dependencies = $18, install_path = $19,
			data_path = $20, config_path = $21, backup_paths = $22, resources = $23,
			permissions = $24, repository_id = $25, download_count = $26, install_count = $27,
			rating = $28, updated_at = $29
		WHERE id = $1
	`

	_, err := p.db.Exec(
		query,
		pkg.ID, pkg.DisplayName, pkg.Version, pkg.Description, pkg.Author, pkg.Website,
		pkg.Category, pkg.License, pkg.FilePath, pkg.FileSize, pkg.FileHash, pkg.DownloadURL,
		pkg.Architecture, pkg.MinOSVersion, pkg.MaxOSVersion, pkg.MinRAM, pkg.MinDiskSpace,
		pkg.Dependencies, pkg.InstallPath, pkg.DataPath, pkg.ConfigPath, pkg.BackupPaths,
		pkg.Resources, pkg.Permissions, pkg.RepositoryID, pkg.DownloadCount, pkg.InstallCount,
		pkg.Rating, pkg.UpdatedAt,
	)

	return err
}

// DeleteAppPackage 删除应用包
func (p *PostgreSQLDatabase) DeleteAppPackage(name string) error {
	query := `DELETE FROM app_packages WHERE name = $1`
	_, err := p.db.Exec(query, name)
	return err
}

// CreateAppRepository 创建应用仓库
func (p *PostgreSQLDatabase) CreateAppRepository(repo *AppRepository) error {
	repo.CreatedAt = time.Now()
	repo.UpdatedAt = time.Now()

	query := `
		INSERT INTO app_repositories (name, url, type, enabled, priority, description, auto_update, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id
	`

	err := p.db.QueryRow(
		query,
		repo.Name, repo.URL, repo.Type, repo.Enabled, repo.Priority,
		repo.Description, repo.AutoUpdate, repo.CreatedAt, repo.UpdatedAt,
	).Scan(&repo.ID)

	return err
}

// ListAppRepositories 列出应用仓库
func (p *PostgreSQLDatabase) ListAppRepositories() ([]AppRepository, error) {
	var repos []AppRepository

	query := `
		SELECT id, created_at, updated_at, name, url, type, enabled, priority, description, auto_update
		FROM app_repositories ORDER BY priority DESC
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var repo AppRepository
		err := rows.Scan(
			&repo.ID, &repo.CreatedAt, &repo.UpdatedAt, &repo.Name, &repo.URL,
			&repo.Type, &repo.Enabled, &repo.Priority, &repo.Description, &repo.AutoUpdate,
		)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}

	return repos, nil
}

// UpdateAppRepository 更新应用仓库
func (p *PostgreSQLDatabase) UpdateAppRepository(repo *AppRepository) error {
	repo.UpdatedAt = time.Now()

	query := `
		UPDATE app_repositories SET
			name = $2, url = $3, type = $4, enabled = $5, priority = $6,
			description = $7, auto_update = $8, updated_at = $9
		WHERE id = $1
	`

	_, err := p.db.Exec(
		query,
		repo.ID, repo.Name, repo.URL, repo.Type, repo.Enabled, repo.Priority,
		repo.Description, repo.AutoUpdate, repo.UpdatedAt,
	)

	return err
}

// DeleteAppRepository 删除应用仓库
func (p *PostgreSQLDatabase) DeleteAppRepository(id uint) error {
	query := `DELETE FROM app_repositories WHERE id = $1`
	_, err := p.db.Exec(query, id)
	return err
}

// CreateInstallLog 创建安装日志
func (p *PostgreSQLDatabase) CreateInstallLog(log *AppInstallLog) error {
	log.CreatedAt = time.Now()
	log.UpdatedAt = time.Now()

	query := `
		INSERT INTO app_install_logs (app_instance_id, step, message, status, details, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`

	err := p.db.QueryRow(
		query,
		log.AppInstanceID, log.Step, log.Message, log.Status, log.Details,
		log.CreatedAt, log.UpdatedAt,
	).Scan(&log.ID)

	return err
}

// GetInstallLogs 获取安装日志
func (p *PostgreSQLDatabase) GetInstallLogs(appID uint) ([]AppInstallLog, error) {
	var logs []AppInstallLog

	query := `
		SELECT id, created_at, updated_at, app_instance_id, step, message, status, details
		FROM app_install_logs WHERE app_instance_id = $1 ORDER BY created_at
	`

	rows, err := p.db.Query(query, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var log AppInstallLog
		err := rows.Scan(
			&log.ID, &log.CreatedAt, &log.UpdatedAt, &log.AppInstanceID,
			&log.Step, &log.Message, &log.Status, &log.Details,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// 辅助函数：序列化字符串数组到JSON字符串
func serializeStringArray(arr []string) string {
	if arr == nil {
		return ""
	}
	data, _ := json.Marshal(arr)
	return string(data)
}

// 辅助函数：从JSON字符串解析字符串数组
func parseStringArray(str string) []string {
	if str == "" {
		return []string{}
	}
	var arr []string
	json.Unmarshal([]byte(str), &arr)
	return arr
}