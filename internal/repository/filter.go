package repository

import (
	"fmt"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/pagination"
	"strings"
)

type Filter struct {
	TenantID string
	Status   string
	Search   string
}

func BuildResourceQuery(f Filter, p pagination.Query) (string, []any) {
	p = pagination.Normalize(p)
	parts := []string{"SELECT id,tenant_id,owner_id,code,name,description,status,version,created_at,updated_at FROM data_resources WHERE tenant_id=?"}
	args := []any{f.TenantID}
	if f.Status != "" {
		parts = append(parts, "AND status=?")
		args = append(args, f.Status)
	}
	if s := strings.TrimSpace(f.Search); s != "" {
		parts = append(parts, "AND (name LIKE ? OR code LIKE ?)")
		args = append(args, "%"+s+"%", "%"+s+"%")
	}
	sort := "updated_at DESC"
	if p.Sort == "name" {
		sort = "name ASC"
	}
	parts = append(parts, fmt.Sprintf("ORDER BY %s LIMIT ? OFFSET ?", sort))
	args = append(args, p.Limit, p.Offset)
	return strings.Join(parts, " "), args
}
func SafeSort(value string) string {
	switch value {
	case "name":
		return "name ASC"
	case "created":
		return "created_at DESC"
	case "status":
		return "status ASC"
	default:
		return "updated_at DESC"
	}
}
