package middleware

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"strings"

	"ai-matching/src/domain/interface/repository"
	"ai-matching/src/infrastructure/external/cognito"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	jwtValidator   *cognito.CognitoJWTValidator
	tenantUserRepo repository.TenantUserRepository
}

func NewAuthMiddleware(userRepo repository.UserRepository, tenantRepo repository.TenantRepository, tenantUserRepo repository.TenantUserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtValidator:   cognito.NewCognitoJWTValidator(userRepo, tenantRepo),
		tenantUserRepo: tenantUserRepo,
	}
}

func (m *AuthMiddleware) FiberMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		token, err := m.jwtValidator.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token: " + err.Error(),
			})
		}

		userInfo, err := m.jwtValidator.GetUserInfoFromToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Failed to extract user info: " + err.Error(),
			})
		}

		// 型変換を行ってからLocalsに保存
		if userIDStr, ok := userInfo["user_id"].(string); ok {
			if userID, err := uuid.Parse(userIDStr); err == nil {
				c.Locals("user_id", userID)
			}
		}
		c.Locals("token", tokenString)

		// Set organization_id and tenant_id if available with proper type conversion
		if orgIDStr, ok := userInfo["organization_id"].(string); ok {
			if orgID, err := uuid.Parse(orgIDStr); err == nil {
				c.Locals("organization_id", orgID)
			}
		}
		if tenantIDStr, ok := userInfo["tenant_id"].(string); ok {
			if tenantID, err := uuid.Parse(tenantIDStr); err == nil {
				c.Locals("tenant_id", tenantID)
			}
		}
		if tenantName, ok := userInfo["tenant_name"]; ok {
			c.Locals("tenant_name", tenantName)
		}
		if tenantSubdomain, ok := userInfo["tenant_subdomain"]; ok {
			c.Locals("tenant_subdomain", tenantSubdomain)
		}
		if tenantIsActive, ok := userInfo["tenant_is_active"]; ok {
			c.Locals("tenant_is_active", tenantIsActive)
		}

		// URLパスから organizationId を抽出
		path := c.Path()
		organizationId := extractOrganizationIdFromPath(path)

		if organizationId != "" {
			if orgID, ok := c.Locals("organization_id").(uuid.UUID); ok {
				if organizationId != orgID.String() {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"error": "Invalid organization ID",
					})
				}
			}
		}

		// URLパスから tenantId を抽出
		tenantId := extractTenantIdFromPath(path)
		if tenantId != "" {
			if tenantID, ok := c.Locals("tenant_id").(uuid.UUID); ok {
				if userID, ok := c.Locals("user_id").(uuid.UUID); ok {
					_, err := m.tenantUserRepo.GetTenantUser(context.Background(), tenantID, userID)
					if err != nil {
						return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
							"error": "User not authorized for this tenant",
						})
					}
				}
			}
		}

		return c.Next()
	}
}

// URLパスから organizationId を抽出する関数
func extractOrganizationIdFromPath(path string) string {
	// /api/v1/auth/organizations/{organizationId} のパターンにマッチ
	re := regexp.MustCompile(`/organizations/([a-fA-F0-9-]+)`)
	matches := re.FindStringSubmatch(path)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// URLパスから tenantId を抽出する関数
func extractTenantIdFromPath(path string) string {
	// 必要に応じてテナントIDのパターンも追加
	re := regexp.MustCompile(`/tenants/([a-fA-F0-9-]+)`)
	matches := re.FindStringSubmatch(path)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

type Tenant struct {
	ID        uuid.UUID
	Name      string
	Subdomain string
	IsActive  bool
}

type UserContext struct {
	UserID         uuid.UUID
	Email          string
	Token          string
	OrganizationID uuid.UUID
	Tenant         *Tenant
}

func GetUserFromContext(ctx context.Context) (*UserContext, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("user_id not found in context")
	}

	token, ok := ctx.Value("token").(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	userContext := &UserContext{
		UserID: userID,
		Token:  token,
		Tenant: &Tenant{},
	}

	// Organization ID and Tenant ID are optional
	if orgID, ok := ctx.Value("organization_id").(uuid.UUID); ok {
		userContext.OrganizationID = orgID
	}
	if tenantID, ok := ctx.Value("tenant_id").(uuid.UUID); ok {
		userContext.Tenant.ID = tenantID
	}
	if tenantName, ok := ctx.Value("tenant_name").(string); ok {
		userContext.Tenant.Name = tenantName
	}
	if tenantSubdomain, ok := ctx.Value("tenant_subdomain").(string); ok {
		userContext.Tenant.Subdomain = tenantSubdomain
	}
	if tenantIsActive, ok := ctx.Value("tenant_is_active").(bool); ok {
		userContext.Tenant.IsActive = tenantIsActive
	}

	return userContext, nil
}
