package response

type MessageResponse struct {
	Message string `json:"message" doc:"Response message"`
}

type TenantUserInfo struct {
	UserID    int64  `json:"userId" doc:"User ID"`
	Email     string `json:"email" doc:"User email"`
	FirstName string `json:"firstName" doc:"User first name"`
	LastName  string `json:"lastName" doc:"User last name"`
	Role      string `json:"role" doc:"User role in tenant"`
}

type TenantUsersResponse struct {
	Users []TenantUserInfo `json:"users" doc:"List of users in the tenant"`
}

type TenantDetails struct {
	ID        int64  `json:"id" doc:"Tenant ID"`
	Name      string `json:"name" doc:"Tenant name"`
	Subdomain string `json:"subdomain" doc:"Tenant subdomain"`
}

type UserTenantsResponse struct {
	Tenants []TenantDetails `json:"tenants" doc:"List of tenants user belongs to"`
}

type UserDetails struct {
	ID        int64  `json:"id" doc:"User ID"`
	Email     string `json:"email" doc:"User email"`
	FirstName string `json:"firstName" doc:"User first name"`
	LastName  string `json:"lastName" doc:"User last name"`
}

type UsersResponse struct {
	Users []UserDetails `json:"users" doc:"List of users"`
}