package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type DockerRepo struct {
	Id           string `json:"id" gorm:"column:id"`
	RepoName     string `json:"repoName" gorm:"column:repo_name"`
	RepoAddr     string `json:"repoAddr" gorm:"column:repo_addr"`
	RepoDesc     string `json:"repoDesc" gorm:"column:repo_desc"`
	RepoUsername string `json:"repoUsername" gorm:"column:repo_username"`
	RepoPassword string `json:"repoPassword" gorm:"column:repo_password"`
	CreateAt     int64  `json:"createAt" gorm:"column:create_at"`
}

func (*DockerRepo) TableName() string {
	return "d_repo"
}

// 创建repo记录
func (r *DockerRepo) Create(db *gorm.DB) error {
	if err := db.Model(r).Create(r).Error; err != nil {
		return fmt.Errorf("创建repo失败: %v", err)
	}
	return nil
}

func (r *DockerRepo) FindByAddr(db *gorm.DB, add string) (*DockerRepo, error) {
	var res DockerRepo
	err := db.Where("repo_addr = ?", add).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 或者返回一个自定义的错误，表示未找到
		}
		return nil, err // 返回其他类型的错误
	}
	return &res, nil
}

func (r *DockerRepo) GetRepoLs(db *gorm.DB) []DockerRepo {
	var res []DockerRepo
	db.Find(&res).Order("create_time DESC")
	return res
}

// 根据ID删除repo
func (r *DockerRepo) DeleteByID(db *gorm.DB, id string) error {
	result := db.Where("id = ?", id).Delete(&DockerRepo{})
	if result.Error != nil {
		return fmt.Errorf("删除repo失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// 根据ID获取repo详情（排除敏感字段）
func (r *DockerRepo) GetByID(db *gorm.DB, id string) (*DockerRepo, error) {
	var repo DockerRepo
	err := db.Model(&DockerRepo{}).
		Select("id, repo_name, repo_addr, repo_desc, create_at,repo_password,repo_username").
		Where("id = ?", id).
		First(&repo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("repo不存在")
		}
		return nil, fmt.Errorf("查询失败: %v", err)
	}
	return &repo, nil
}

// 分页获取repo列表（带模糊查询）
func (r *DockerRepo) GetList(db *gorm.DB, repoName string, page, pageSize int) (total int64, list []DockerRepo, err error) {
	query := db.Model(&DockerRepo{}).
		Select("id, repo_name, repo_addr, repo_desc, create_at,repo_password,repo_username")

	if repoName != "" {
		query = query.Where("repo_name LIKE ?", "%"+repoName+"%")
	}

	// 获取总数
	if err = query.Count(&total).Error; err != nil {
		return 0, nil, fmt.Errorf("查询总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Order("create_at DESC").
		Find(&list).Error

	if err != nil {
		return 0, nil, fmt.Errorf("查询列表失败: %v", err)
	}
	return total, list, nil
}

// 更新repo信息
func (r *DockerRepo) Update(db *gorm.DB) error {
	updateData := map[string]interface{}{
		"repo_name":     r.RepoName,
		"repo_addr":     r.RepoAddr,
		"repo_desc":     r.RepoDesc,
		"repo_username": r.RepoUsername,
		"repo_password": r.RepoPassword,
	}

	return db.Model(r).Updates(updateData).Error
}
