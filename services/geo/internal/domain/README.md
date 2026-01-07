# 도메인 계층

이 디렉토리는 API 서비스의 도메인 계층을 포함하고 있습니다. 도메인 계층은 비즈니스 엔티티와 비즈니스 규칙을 정의합니다.

## 구조

- **entity**: 핵심 비즈니스 객체 (태스크, 프로젝트, 사용자 등)
- **repository**: 데이터 저장소 인터페이스

## 사용 방법

### 엔티티 정의

엔티티는 비즈니스 객체를 표현하며, 유효성 검사 메서드를 포함합니다:

```go
package entity

import (
    "time"
    "errors"
)

type ItemType string

const (
    ItemTypeTask    ItemType = "task"
    ItemTypeProject ItemType = "project"
)

type Item struct {
    ID          string
    ParentID    *string
    Type        ItemType
    Name        string
    Contents    string
    Objective   string
    Deliverable string
    Role        string
    CreatedBy   string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func NewItem(name, contents string, itemType ItemType, createdBy string, parentID *string) (*Item, error) {
    if name == "" {
        return nil, errors.New("이름은 필수입니다")
    }
    
    if itemType == "" {
        return nil, errors.New("아이템 유형은 필수입니다")
    }
    
    if createdBy == "" {
        return nil, errors.New("생성자는 필수입니다")
    }
    
    now := time.Now()
    
    return &Item{
        Name:      name,
        Contents:  contents,
        Type:      itemType,
        ParentID:  parentID,
        CreatedBy: createdBy,
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}
```

### 레포지토리 인터페이스 정의

레포지토리 인터페이스는 도메인 엔티티에 대한 데이터 접근 계약을 정의합니다:

```go
package repository

import (
    "context"
    
    "github.com/wekeepgrowing/semo-backend-monorepo/services/api/internal/domain/entity"
)

type ItemRepository interface {
    Create(ctx context.Context, item *entity.Item) error
    GetByID(ctx context.Context, id string) (*entity.Item, error)
    GetChildren(ctx context.Context, parentID string) ([]*entity.Item, error)
    Update(ctx context.Context, item *entity.Item) error
    Delete(ctx context.Context, id string) error
}
```

## 가이드라인

1. 도메인 계층은 외부 종속성(데이터베이스, 외부 API 등)을 가져오지 않아야 합니다.
2. 도메인 계층은 비즈니스 규칙만 포함해야 합니다.
3. 레포지토리 인터페이스는 도메인 엔티티만 사용해야 합니다.
4. 비즈니스 로직은 도메인 서비스에 구현하세요. 