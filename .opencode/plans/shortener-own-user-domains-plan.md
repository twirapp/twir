I understand! I'm in plan mode and should only create a plan document, not make any changes. Let me show you the comprehensive plan as text output since I cannot write files in plan mode:

---

# Custom Domains for Short Links - Implementation Plan

## Overview

Add support for users to configure custom domains for their short links. Instead of using `twir.app/s/Qsqe1`, users can configure their own domain like `links.278.co/Qsqe1` which will point directly to their short links without the `/s/` prefix.

## Requirements Summary

1. **Custom Domain per User**: Each user can configure ONE custom domain for all their short links
2. **Domain Verification**: Use DNS CNAME verification to ensure users have correctly pointed their domain to our service
3. **Unique Slugs per Domain**: Slugs will be unique per domain (not globally unique anymore)
   - User A with `links.example.com` can use slug `tg`
   - User B with `short.example2.com` can also use slug `tg`
   - Default domain `twir.app` maintains its own namespace
4. **Routing**: Custom domain requests go directly to `/v1/short-links/{shortId}` without the `/s/` prefix
5. **Backend**: HAProxy handles custom domain routing, backend handles domain lookup and slug resolution

## Database Changes

### 1. Create New Table: `short_links_custom_domains`

**Migration file**: `libs/migrations/postgres/YYYYMMDDHHMMSS_add_short_links_custom_domains.sql`

```sql
-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE short_links_custom_domains (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    user_id TEXT NOT NULL,
    domain TEXT NOT NULL UNIQUE,
    verified BOOLEAN NOT NULL DEFAULT false,
    verification_token TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX short_links_custom_domains_user_id_idx ON short_links_custom_domains(user_id);
CREATE INDEX short_links_custom_domains_domain_idx ON short_links_custom_domains(domain);
CREATE UNIQUE INDEX short_links_custom_domains_user_id_unique_idx ON short_links_custom_domains(user_id);

COMMENT ON TABLE short_links_custom_domains IS 'Custom domains for users short links';
COMMENT ON COLUMN short_links_custom_domains.verified IS 'Whether the domain has been verified via CNAME check';
COMMENT ON COLUMN short_links_custom_domains.verification_token IS 'Expected CNAME target: short-{token}.twir.app';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS short_links_custom_domains;
-- +goose StatementEnd
```

### 2. Modify `shortened_urls` Table

**Migration file**: `libs/migrations/postgres/YYYYMMDDHHMMSS_add_domain_to_shortened_urls.sql`

```sql
-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

-- Add domain column (nullable for backward compatibility)
ALTER TABLE shortened_urls
    ADD COLUMN domain TEXT;

-- Add index for domain lookups
CREATE INDEX shortened_urls_domain_idx ON shortened_urls(domain);

-- Create composite unique index for slug uniqueness per domain
CREATE UNIQUE INDEX shortened_urls_domain_short_id_unique_idx
    ON shortened_urls(COALESCE(domain, ''), short_id);

COMMENT ON COLUMN shortened_urls.domain IS 'Custom domain for this short link (NULL = default twir.app domain)';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS shortened_urls_domain_short_id_unique_idx;
DROP INDEX IF EXISTS shortened_urls_domain_idx;
ALTER TABLE shortened_urls DROP COLUMN IF EXISTS domain;
-- +goose StatementEnd
```

**Important Note**: The `short_id` PRIMARY KEY constraint needs special handling. Since we're allowing duplicate `short_id` values across different domains, we need to change the primary key strategy:

**Option A (Recommended)**: Add a new UUID primary key column, keep `short_id` as a regular column, enforce uniqueness via composite index `(domain, short_id)`.

```sql
-- Migration to change primary key
ALTER TABLE shortened_urls DROP CONSTRAINT shortened_urls_pkey;
ALTER TABLE shortened_urls ADD COLUMN id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text;
```

## Repository Layer

### 3. Create Custom Domains Repository

**File structure**:

```
libs/repositories/short_links_custom_domains/
├── short_links_custom_domains.go          # Interface definition
├── errors.go                               # Custom errors
├── model/
│   └── model.go                            # Model with Nil pattern
└── pgx/
    └── pgx.go                              # Repository implementation
```

#### 3.1. Interface: `libs/repositories/short_links_custom_domains/short_links_custom_domains.go`

```go
package shortlinkscustomdomains

import (
    "context"
    "github.com/twirapp/twir/libs/repositories/short_links_custom_domains/model"
)

type Repository interface {
    GetByUserID(ctx context.Context, userID string) (model.CustomDomain, error)
    GetByDomain(ctx context.Context, domain string) (model.CustomDomain, error)
    Create(ctx context.Context, input CreateInput) (model.CustomDomain, error)
    Update(ctx context.Context, id string, input UpdateInput) (model.CustomDomain, error)
    Delete(ctx context.Context, id string) error
    VerifyDomain(ctx context.Context, id string) error
}

type CreateInput struct {
    UserID            string
    Domain            string
    VerificationToken string
}

type UpdateInput struct {
    Domain            *string
    Verified          *bool
    VerificationToken *string
}
```

#### 3.2. Model: `libs/repositories/short_links_custom_domains/model/model.go`

```go
package model

import "time"

type CustomDomain struct {
    ID                string
    UserID            string
    Domain            string
    Verified          bool
    VerificationToken string
    CreatedAt         time.Time
    UpdatedAt         time.Time

    isNil bool
}

func (c CustomDomain) IsNil() bool {
    return c.isNil
}

var Nil = CustomDomain{
    isNil: true,
}
```

#### 3.3. Errors: `libs/repositories/short_links_custom_domains/errors.go`

```go
package shortlinkscustomdomains

import "errors"

var ErrNotFound = errors.New("custom domain not found")
var ErrDomainAlreadyExists = errors.New("domain already exists")
var ErrUserAlreadyHasDomain = errors.New("user already has a custom domain")
```

### 4. Update Shortened URLs Repository

**File**: `libs/repositories/shortened_urls/shortened_urls.go`

Update interface to handle domains:

```go
type Repository interface {
    // Updated: add domain parameter
    GetByShortID(ctx context.Context, domain *string, id string) (model.ShortenedUrl, error)
    GetManyByShortIDs(ctx context.Context, domain *string, ids []string) ([]model.ShortenedUrl, error)
    GetByUrl(ctx context.Context, domain *string, url string) (model.ShortenedUrl, error)

    Create(ctx context.Context, input CreateInput) (model.ShortenedUrl, error)
    Update(ctx context.Context, id string, input UpdateInput) (model.ShortenedUrl, error)
    GetList(ctx context.Context, input GetListInput) (GetListOutput, error)
    Delete(ctx context.Context, id string) error
    Count(ctx context.Context, input CountInput) (int, error)
}

type CreateInput struct {
    ShortID         string
    URL             string
    CreatedByUserID *string
    UserIp          *string
    UserAgent       *string
    Domain          *string  // NEW
}

type UpdateInput struct {
    Views   *int
    ShortID *string
    URL     *string
    Domain  *string  // NEW
}
```

**File**: `libs/repositories/shortened_urls/model/model.go`

```go
type ShortenedUrl struct {
    ShortID         string
    CreatedAt       time.Time
    UpdatedAt       time.Time
    URL             string
    CreatedByUserId *string
    Views           int
    UserAgent       *string
    UserIp          *netip.Addr
    Domain          *string  // NEW

    isNil bool  // ADD: Following project pattern
}

func (s ShortenedUrl) IsNil() bool {
    return s.isNil
}

var Nil = ShortenedUrl{
    isNil: true,
}
```

**File**: `libs/repositories/shortened_urls/datasource/postgres/pgx.go`

Key changes needed:

1. Add `domain` to SELECT queries
2. Update `GetByShortID` to filter by `(domain, short_id)` pair
3. Update `GetByUrl` to filter by `(domain, url)` pair
4. Update `Create` to insert domain
5. Handle NULL domain for default domain queries

```go
func (r *Repository) GetByShortID(ctx context.Context, domain *string, id string) (model.ShortenedUrl, error) {
    query := r.queryBuilder.
        Select("short_id", "created_at", "updated_at", "url", "created_by_user_id", "views", "user_agent", "user_ip", "domain").
        From("shortened_urls").
        Where(squirrel.Eq{"short_id": id})

    // Handle domain filtering
    if domain == nil {
        query = query.Where("domain IS NULL")
    } else {
        query = query.Where(squirrel.Eq{"domain": *domain})
    }

    // ... rest of implementation
}
```

## Entity Layer

### 5. Create Custom Domain Entity

**File**: `libs/entities/short_links_custom_domain/entity.go`

```go
package shortlinkscustomdomain

import (
    "errors"
    "time"
)

type Entity struct {
    ID                string
    UserID            string
    Domain            string
    Verified          bool
    VerificationToken string
    CreatedAt         time.Time
    UpdatedAt         time.Time

    isNil bool
}

func (c Entity) IsNil() bool {
    return c.isNil
}

var Nil = Entity{
    isNil: true,
}

// GetVerificationTarget returns the expected CNAME target
// Format: short-{token}.twir.app
func (c Entity) GetVerificationTarget() string {
    return "short-" + c.VerificationToken + ".twir.app"
}

// Validate performs domain validation
func (c Entity) Validate() error {
    if c.Domain == "" {
        return errors.New("domain cannot be empty")
    }
    if c.UserID == "" {
        return errors.New("user ID cannot be empty")
    }
    // Add regex validation for domain format
    return nil
}
```

## Service Layer

### 6. Create Custom Domains Service

**File**: `apps/api-gql/internal/services/shortlinkscustomdomains/service.go`

```go
package shortlinkscustomdomains

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "net"

    shortlinkscustomdomain "github.com/twirapp/twir/libs/entities/short_links_custom_domain"
    shortlinkscustomdomainsrepo "github.com/twirapp/twir/libs/repositories/short_links_custom_domains"
    "github.com/twirapp/twir/libs/repositories/short_links_custom_domains/model"
    "go.uber.org/fx"
)

type Opts struct {
    fx.In
    Repository shortlinkscustomdomainsrepo.Repository
}

func New(opts Opts) *Service {
    return &Service{
        repository: opts.Repository,
    }
}

type Service struct {
    repository shortlinkscustomdomainsrepo.Repository
}

type CreateInput struct {
    UserID string
    Domain string
}

func (s *Service) Create(ctx context.Context, input CreateInput) (shortlinkscustomdomain.Entity, error) {
    // Generate verification token
    token := generateVerificationToken()

    // Check if user already has a custom domain
    existing, err := s.repository.GetByUserID(ctx, input.UserID)
    if err != nil && !errors.Is(err, shortlinkscustomdomainsrepo.ErrNotFound) {
        return shortlinkscustomdomain.Nil, err
    }

    if !existing.IsNil() {
        return shortlinkscustomdomain.Nil, shortlinkscustomdomainsrepo.ErrUserAlreadyHasDomain
    }

    // Create custom domain
    created, err := s.repository.Create(ctx, shortlinkscustomdomainsrepo.CreateInput{
        UserID:            input.UserID,
        Domain:            input.Domain,
        VerificationToken: token,
    })
    if err != nil {
        return shortlinkscustomdomain.Nil, err
    }

    return mapToEntity(created), nil
}

// VerifyDomain checks DNS records and verifies the domain
func (s *Service) VerifyDomain(ctx context.Context, userID string) error {
    customDomain, err := s.repository.GetByUserID(ctx, userID)
    if err != nil {
        return err
    }

    if customDomain.IsNil() {
        return errors.New("custom domain not found")
    }

    if customDomain.Verified {
        return nil // Already verified
    }

    // Perform DNS CNAME check
    expectedTarget := "short-" + customDomain.VerificationToken + ".twir.app."

    cname, err := net.LookupCNAME(customDomain.Domain)
    if err != nil {
        return errors.New("DNS lookup failed: " + err.Error())
    }

    if cname != expectedTarget {
        return errors.New("CNAME record does not match expected target")
    }

    // Mark as verified
    return s.repository.VerifyDomain(ctx, customDomain.ID)
}

func (s *Service) GetByUserID(ctx context.Context, userID string) (shortlinkscustomdomain.Entity, error) {
    customDomain, err := s.repository.GetByUserID(ctx, userID)
    if err != nil {
        if errors.Is(err, shortlinkscustomdomainsrepo.ErrNotFound) {
            return shortlinkscustomdomain.Nil, nil
        }
        return shortlinkscustomdomain.Nil, err
    }

    return mapToEntity(customDomain), nil
}

func (s *Service) GetByDomain(ctx context.Context, domain string) (shortlinkscustomdomain.Entity, error) {
    customDomain, err := s.repository.GetByDomain(ctx, domain)
    if err != nil {
        if errors.Is(err, shortlinkscustomdomainsrepo.ErrNotFound) {
            return shortlinkscustomdomain.Nil, nil
        }
        return shortlinkscustomdomain.Nil, err
    }

    return mapToEntity(customDomain), nil
}

func (s *Service) Delete(ctx context.Context, userID string) error {
    customDomain, err := s.repository.GetByUserID(ctx, userID)
    if err != nil {
        return err
    }

    if customDomain.IsNil() {
        return errors.New("custom domain not found")
    }

    return s.repository.Delete(ctx, customDomain.ID)
}

func generateVerificationToken() string {
    b := make([]byte, 16)
    rand.Read(b)
    return hex.EncodeToString(b)
}

func mapToEntity(m model.CustomDomain) shortlinkscustomdomain.Entity {
    return shortlinkscustomdomain.Entity{
        ID:                m.ID,
        UserID:            m.UserID,
        Domain:            m.Domain,
        Verified:          m.Verified,
        VerificationToken: m.VerificationToken,
        CreatedAt:         m.CreatedAt,
        UpdatedAt:         m.UpdatedAt,
    }
}
```

### 7. Update Shortened URLs Service

**File**: `apps/api-gql/internal/services/shortenedurls/shortenedurls.go`

```go
type CreateInput struct {
    CreatedByUserID *string
    ShortID         string
    URL             string
    UserIp          *string
    UserAgent       *string
    Domain          *string  // NEW
}

func (c *Service) Create(ctx context.Context, input CreateInput) (model.ShortenedUrl, error) {
    shortId := input.ShortID
    if input.ShortID == "" {
        shortId = genId()
    }

    return c.repository.Create(
        ctx,
        shortenedurlsrepository.CreateInput{
            ShortID:         shortId,
            URL:             input.URL,
            CreatedByUserID: input.CreatedByUserID,
            UserIp:          input.UserIp,
            UserAgent:       input.UserAgent,
            Domain:          input.Domain,  // NEW
        },
    )
}

// Update GetByShortID to accept domain parameter
func (c *Service) GetByShortID(ctx context.Context, domain *string, id string) (model.ShortenedUrl, error) {
    link, err := c.repository.GetByShortID(ctx, domain, id)
    if err != nil {
        if errors.Is(err, shortenedurlsrepository.ErrNotFound) {
            return model.Nil, nil
        }
        return model.Nil, err
    }

    return link, nil
}
```

## HTTP Routes Layer

### 8. Custom Domain Management Routes

**File**: `apps/api-gql/internal/delivery/http/routes/shortlinks/custom_domain_get.go`

```go
package shortlinks

import (
    "context"
    "net/http"
    "time"

    "github.com/danielgtaylor/huma/v2"
    httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
)

type getCustomDomain struct {
    customDomainsService *shortlinkscustomdomains.Service
    sessions             *auth.Auth
}

type customDomainOutputDto struct {
    ID                 string    `json:"id"`
    Domain             string    `json:"domain"`
    Verified           bool      `json:"verified"`
    VerificationToken  string    `json:"verification_token"`
    VerificationTarget string    `json:"verification_target"`
    CreatedAt          time.Time `json:"created_at"`
}

func (c *getCustomDomain) GetMeta() huma.Operation {
    return huma.Operation{
        OperationID: "short-links-get-custom-domain",
        Method:      http.MethodGet,
        Path:        "/v1/short-links/custom-domain",
        Tags:        []string{"Short links"},
        Summary:     "Get custom domain configuration",
    }
}

func (c *getCustomDomain) Handler(ctx context.Context, input *struct{}) (*httpbase.BaseOutputJson[customDomainOutputDto], error) {
    user, err := c.sessions.GetAuthenticatedUserModel(ctx)
    if err != nil {
        return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
    }

    customDomain, err := c.customDomainsService.GetByUserID(ctx, user.ID)
    if err != nil {
        return nil, huma.NewError(http.StatusInternalServerError, "Cannot get custom domain", err)
    }

    if customDomain.IsNil() {
        return nil, huma.NewError(http.StatusNotFound, "Custom domain not configured")
    }

    return httpbase.CreateBaseOutputJson(customDomainOutputDto{
        ID:                 customDomain.ID,
        Domain:             customDomain.Domain,
        Verified:           customDomain.Verified,
        VerificationToken:  customDomain.VerificationToken,
        VerificationTarget: customDomain.GetVerificationTarget(),
        CreatedAt:          customDomain.CreatedAt,
    }), nil
}
```

**File**: `apps/api-gql/internal/delivery/http/routes/shortlinks/custom_domain_create.go`

```go
package shortlinks

import (
    "context"
    "net/http"

    "github.com/danielgtaylor/huma/v2"
    httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
)

type createCustomDomain struct {
    customDomainsService *shortlinkscustomdomains.Service
    sessions             *auth.Auth
}

type createCustomDomainInput struct {
    Body struct {
        Domain string `json:"domain" required:"true" minLength:"3" maxLength:"255" example:"links.example.com"`
    }
}

func (c *createCustomDomain) GetMeta() huma.Operation {
    return huma.Operation{
        OperationID: "short-links-create-custom-domain",
        Method:      http.MethodPost,
        Path:        "/v1/short-links/custom-domain",
        Tags:        []string{"Short links"},
        Summary:     "Configure custom domain",
    }
}

func (c *createCustomDomain) Handler(ctx context.Context, input *createCustomDomainInput) (*httpbase.BaseOutputJson[customDomainOutputDto], error) {
    user, err := c.sessions.GetAuthenticatedUserModel(ctx)
    if err != nil {
        return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
    }

    customDomain, err := c.customDomainsService.Create(ctx, shortlinkscustomdomains.CreateInput{
        UserID: user.ID,
        Domain: input.Body.Domain,
    })
    if err != nil {
        return nil, huma.NewError(http.StatusBadRequest, "Cannot create custom domain", err)
    }

    return httpbase.CreateBaseOutputJson(customDomainOutputDto{
        ID:                 customDomain.ID,
        Domain:             customDomain.Domain,
        Verified:           customDomain.Verified,
        VerificationToken:  customDomain.VerificationToken,
        VerificationTarget: customDomain.GetVerificationTarget(),
        CreatedAt:          customDomain.CreatedAt,
    }), nil
}
```

**File**: `apps/api-gql/internal/delivery/http/routes/shortlinks/custom_domain_verify.go`

```go
package shortlinks

import (
    "context"
    "net/http"

    "github.com/danielgtaylor/huma/v2"
    httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
)

type verifyCustomDomain struct {
    customDomainsService *shortlinkscustomdomains.Service
    sessions             *auth.Auth
}

func (c *verifyCustomDomain) GetMeta() huma.Operation {
    return huma.Operation{
        OperationID: "short-links-verify-custom-domain",
        Method:      http.MethodPost,
        Path:        "/v1/short-links/custom-domain/verify",
        Tags:        []string{"Short links"},
        Summary:     "Verify custom domain DNS configuration",
    }
}

func (c *verifyCustomDomain) Handler(ctx context.Context, input *struct{}) (*httpbase.BaseOutputJson[customDomainOutputDto], error) {
    user, err := c.sessions.GetAuthenticatedUserModel(ctx)
    if err != nil {
        return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
    }

    err = c.customDomainsService.VerifyDomain(ctx, user.ID)
    if err != nil {
        return nil, huma.NewError(http.StatusBadRequest, "Domain verification failed", err)
    }

    customDomain, err := c.customDomainsService.GetByUserID(ctx, user.ID)
    if err != nil {
        return nil, huma.NewError(http.StatusInternalServerError, "Cannot get custom domain", err)
    }

    return httpbase.CreateBaseOutputJson(customDomainOutputDto{
        ID:                 customDomain.ID,
        Domain:             customDomain.Domain,
        Verified:           customDomain.Verified,
        VerificationToken:  customDomain.VerificationToken,
        VerificationTarget: customDomain.GetVerificationTarget(),
        CreatedAt:          customDomain.CreatedAt,
    }), nil
}
```

**File**: `apps/api-gql/internal/delivery/http/routes/shortlinks/custom_domain_delete.go`

```go
package shortlinks

import (
    "context"
    "net/http"

    "github.com/danielgtaylor/huma/v2"
    httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
)

type deleteCustomDomain struct {
    customDomainsService *shortlinkscustomdomains.Service
    sessions             *auth.Auth
}

func (c *deleteCustomDomain) GetMeta() huma.Operation {
    return huma.Operation{
        OperationID: "short-links-delete-custom-domain",
        Method:      http.MethodDelete,
        Path:        "/v1/short-links/custom-domain",
        Tags:        []string{"Short links"},
        Summary:     "Delete custom domain configuration",
    }
}

func (c *deleteCustomDomain) Handler(ctx context.Context, input *struct{}) (*struct{ Body struct{} }, error) {
    user, err := c.sessions.GetAuthenticatedUserModel(ctx)
    if err != nil {
        return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
    }

    err = c.customDomainsService.Delete(ctx, user.ID)
    if err != nil {
        return nil, huma.NewError(http.StatusBadRequest, "Cannot delete custom domain", err)
    }

    return &struct{ Body struct{} }{}, nil
}
```

### 9. Update Redirect Route

**File**: `apps/api-gql/internal/delivery/http/routes/shortlinks/redirect.go`

Update the `redirect.go:71` Handler function:

```go
func (r *redirect) Handler(ctx context.Context, input *redirectRequestDto) (
    *redirectResponseDto,
    error,
) {
    // Extract domain from request Host header
    var domain *string
    if host, err := humahelpers.GetHostFromCtx(ctx); err == nil {
        // Check if this is NOT a default domain
        if !isDefaultDomain(host) {
            domain = &host
        }
    }

    // Look up link by short ID and domain
    link, err := r.service.GetByShortID(ctx, domain, input.ShortId)
    if err != nil {
        return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
    }

    if link.IsNil() {
        return nil, huma.NewError(http.StatusNotFound, "Link not found")
    }

    // ... rest of the handler remains the same
}

func isDefaultDomain(host string) bool {
    return host == "twir.app" || host == "cf.twir.app"
}
```

**Also need to add helper**: `apps/api-gql/internal/server/huma_helpers/host.go`

```go
package humahelpers

import (
    "context"
    "errors"

    "github.com/danielgtaylor/huma/v2"
)

func GetHostFromCtx(ctx context.Context) (string, error) {
    ginCtx, err := huma.Context(ctx)
    if err != nil {
        return "", err
    }

    host := ginCtx.Header("Host")
    if host == "" {
        return "", errors.New("host header not found")
    }

    return host, nil
}
```

### 10. Update Create Route

**File**: `apps/api-gql/internal/delivery/http/routes/shortlinks/create.go`

Update the `create.go:80` Handler function:

```go
func (c *create) Handler(
    ctx context.Context,
    input *createLinkInput,
) (*httpbase.BaseOutputJson[linkOutputDto], error) {
    baseUrl, err := gincontext.GetBaseUrlFromContext(ctx, c.config.SiteBaseUrl)
    if err != nil {
        return nil, huma.NewError(http.StatusInternalServerError, "Cannot get base URL", err)
    }

    var customDomain *string
    var createdByUserID *string

    user, _ := c.sessions.GetAuthenticatedUserModel(ctx)
    if user != nil {
        createdByUserID = &user.ID

        // Check if user has a verified custom domain
        userDomain, err := c.customDomainsService.GetByUserID(ctx, user.ID)
        if err == nil && !userDomain.IsNil() && userDomain.Verified {
            customDomain = &userDomain.Domain
        }
    }

    // Check for existing link by URL (within same domain namespace)
    if input.Body.Alias == "" {
        existedLink, err := c.service.GetByUrl(ctx, customDomain, input.Body.Url)
        if err != nil {
            return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
        }

        if !existedLink.IsNil() {
            var shortUrl string
            if customDomain != nil {
                shortUrl = "https://" + *customDomain + "/" + existedLink.ShortID
            } else {
                parsedBaseUrl, _ := url.Parse(baseUrl)
                parsedBaseUrl.Path = "/s/" + existedLink.ShortID
                shortUrl = parsedBaseUrl.String()
            }

            return httpbase.CreateBaseOutputJson(
                linkOutputDto{
                    Id:        existedLink.ShortID,
                    Url:       existedLink.URL,
                    ShortUrl:  shortUrl,
                    Views:     existedLink.Views,
                    CreatedAt: existedLink.CreatedAt,
                },
            ), nil
        }
    }

    // Check for alias conflicts within the same domain
    if input.Body.Alias != "" {
        existedLink, err := c.service.GetByShortID(ctx, customDomain, input.Body.Alias)
        if err != nil {
            return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
        }

        if !existedLink.IsNil() {
            return nil, huma.NewError(http.StatusConflict, "Alias already in use")
        }
    }

    clientIp, err := humahelpers.GetClientIpFromCtx(ctx)
    if err != nil {
        return nil, huma.NewError(
            http.StatusInternalServerError,
            "Internal error on getting your information",
            err,
        )
    }

    clientAgent, err := humahelpers.GetClientUserAgentFromCtx(ctx)
    if err != nil {
        return nil, huma.NewError(
            http.StatusInternalServerError,
            "Internal error on getting your information",
            err,
        )
    }

    link, err := c.service.Create(
        ctx, shortenedurls.CreateInput{
            CreatedByUserID: createdByUserID,
            ShortID:         input.Body.Alias,
            URL:             input.Body.Url,
            UserIp:          &clientIp,
            UserAgent:       &clientAgent,
            Domain:          customDomain,  // NEW
        },
    )
    if err != nil {
        return nil, huma.NewError(http.StatusNotFound, "Cannot generate short id", err)
    }

    // Build short URL with custom domain if applicable
    var shortUrl string
    if customDomain != nil {
        shortUrl = "https://" + *customDomain + "/" + link.ShortID
    } else {
        parsedBaseUrl, _ := url.Parse(baseUrl)
        parsedBaseUrl.Path = "/s/" + link.ShortID
        shortUrl = parsedBaseUrl.String()
    }

    if err := c.sessions.AddLatestShortenerUrlsId(ctx, link.ShortID); err != nil {
        c.logger.Warn("Cannot save latest short links ids to session: " + err.Error())
    }

    return httpbase.CreateBaseOutputJson(
        linkOutputDto{
            Id:        link.ShortID,
            Url:       link.URL,
            ShortUrl:  shortUrl,
            Views:     link.Views,
            CreatedAt: link.CreatedAt,
        },
    ), nil
}
```

### 11. Register New Routes

**File**: `apps/api-gql/internal/delivery/http/routes/shortlinks/shortlinks.go`

Update the FxModule to include new routes:

```go
var FxModule = fx.Provide(
    httpbase.AsFxRoute(newCreate),
    httpbase.AsFxRoute(newInfo),
    httpbase.AsFxRoute(newRedirect),
    httpbase.AsFxRoute(newProfile),
    httpbase.AsFxRoute(newStatistics),
    httpbase.AsFxRoute(newTopCountries),
    httpbase.AsFxRoute(newUpdate),
    httpbase.AsFxRoute(newDelete),
    // NEW: Custom domain routes
    httpbase.AsFxRoute(newGetCustomDomain),
    httpbase.AsFxRoute(newCreateCustomDomain),
    httpbase.AsFxRoute(newVerifyCustomDomain),
    httpbase.AsFxRoute(newDeleteCustomDomain),
)
```

## HAProxy Configuration

### 12. Update HAProxy Configuration

**File**: `configs/haproxy/haproxy.cfg`

Update the configuration around line 38-70:

```haproxy
frontend http
    bind *:80

    # Rate limiting: 1000 requests per 10 seconds per IP
    stick-table type string len 45 size 100k expire 30s store http_req_rate(10s)

    # Track by X-Ru-Detected-IP if present, otherwise by source IP
    http-request track-sc0 req.hdr(X-Ru-Detected-IP) if { req.hdr(X-Ru-Detected-IP) -m found }
    http-request track-sc0 src if !{ req.hdr(X-Ru-Detected-IP) -m found }

    # Deny if rate exceeds limit
    http-request deny deny_status 429 if { sc_http_req_rate(0) gt 1000 }

    # ACLs for different services
    acl host_music_recognizer hdr(host) -i music-recognizer.twir.app
    acl host_bots hdr(host) -i services-bots.twir.app
    acl host_twir hdr(host) -i twir.app
    acl host_twir_cf hdr(host) -i cf.twir.app

    # NEW: ACL for custom domains
    # Match any domain that is NOT one of our known service domains
    acl host_custom_domain hdr(host) -m reg -v ^(?!twir\.app|cf\.twir\.app|music-recognizer\.twir\.app|services-bots\.twir\.app).*

    # Path-based ACLs
    acl path_api path_beg /api
    acl path_dashboard path_beg /dashboard
    acl path_overlays path_beg /overlays
    acl path_socket path_beg /socket
    acl path_shortener path_beg /s/

    # Combined routing rules
    use_backend music_recognizer if host_music_recognizer
    use_backend bots if host_bots

    # NEW: Custom domain routing - send all requests to API for short link handling
    use_backend api_custom_domain if host_custom_domain

    # Combined host and path rules for twir.app
    use_backend api if host_twir path_api
    use_backend api if host_twir path_shortener
    use_backend dashboard if host_twir path_dashboard
    use_backend overlays if host_twir path_overlays
    use_backend socket if host_twir path_socket
    use_backend web if host_twir !path_api !path_dashboard !path_overlays !path_socket

    # Combined host and path rules for cf.twir.app
    use_backend api if host_twir_cf path_api
    use_backend api if host_twir_cf path_shortener
    use_backend dashboard if host_twir_cf path_dashboard
    use_backend overlays if host_twir_cf path_overlays
    use_backend socket if host_twir_cf path_socket
    use_backend web if host_twir_cf !path_api !path_dashboard !path_overlays !path_socket

# NEW: Backend for custom domain short links
backend api_custom_domain
    mode http
    balance roundrobin
    option tcp-check
    # Rewrite path to add /v1/short-links/ prefix
    # Example: links.example.com/Qsqe1 -> /v1/short-links/Qsqe1
    http-request set-path /v1/short-links%[path]
    server-template api-gql- 1-30 twir_api-gql:3009 check resolvers docker init-addr none

backend api
    mode http
    balance roundrobin
    option tcp-check
    http-request set-path %[path,regsub(^/s/,/v1/short-links/)] if { path_beg /s/ }
    http-request set-path %[path,regsub(^/api/,/)] if { path_beg /api/ }
    server-template api-gql- 1-30 twir_api-gql:3009 check resolvers docker init-addr none

# ... rest remains the same
```

## Implementation Order

### 13. Step-by-Step Implementation

**Phase 1: Database Layer (Day 1)**

1. Create migration for `short_links_custom_domains` table
2. Create migration to add `domain` column to `shortened_urls`
3. Create migration to change primary key strategy (add UUID `id` column)
4. Run migrations: `bun cli m up --db postgres`

**Phase 2: Repository Layer (Day 2)**

1. Create custom domains repository:
   - Interface in `libs/repositories/short_links_custom_domains/short_links_custom_domains.go`
   - Model in `libs/repositories/short_links_custom_domains/model/model.go`
   - Errors in `libs/repositories/short_links_custom_domains/errors.go`
   - Implementation in `libs/repositories/short_links_custom_domains/pgx/pgx.go`
2. Update shortened URLs repository:
   - Update interface to accept `domain` parameter
   - Update model to include `domain` field and `isNil` pattern
   - Update pgx implementation to handle domain filtering

**Phase 3: Entity Layer (Day 2)**

1. Create `libs/entities/short_links_custom_domain/entity.go`
2. Update shortened URL entity if needed

**Phase 4: Service Layer (Day 3)**

1. Create custom domains service in `apps/api-gql/internal/services/shortlinkscustomdomains/`
2. Update shortened URLs service to handle domains
3. Wire up services in FX modules

**Phase 5: HTTP Routes (Day 4)**

1. Create custom domain management routes:
   - GET /v1/short-links/custom-domain
   - POST /v1/short-links/custom-domain
   - POST /v1/short-links/custom-domain/verify
   - DELETE /v1/short-links/custom-domain
2. Update redirect route to extract and use domain
3. Update create route to use user's custom domain
4. Register new routes in FX module
5. Add helper function for extracting Host header

**Phase 6: HAProxy Configuration (Day 4)**

1. Update HAProxy config to route custom domains
2. Test routing locally

**Phase 7: Testing (Day 5)**

1. Unit tests for all new components
2. Integration tests for DNS verification
3. End-to-end manual testing
4. Performance testing

## Testing Checklist

### Manual Testing Steps

1. **Setup Custom Domain**
   - [ ] Create custom domain via API
   - [ ] Receive verification token and target
   - [ ] Add CNAME record in DNS: `links.example.com CNAME short-{token}.twir.app`
   - [ ] Wait for DNS propagation

2. **Verify Domain**
   - [ ] Trigger verification via API
   - [ ] Verify domain is marked as verified
   - [ ] Try to verify again (should succeed/no-op)

3. **Create Short Links**
   - [ ] Create link without alias on custom domain
   - [ ] Create link with alias on custom domain
   - [ ] Verify short URL uses custom domain format

4. **Slug Uniqueness**
   - [ ] User A with custom domain creates link with alias "tg"
   - [ ] User B with different custom domain creates link with alias "tg"
   - [ ] Both should succeed
   - [ ] User A tries to create another link with alias "tg" (should fail)

5. **Redirects**
   - [ ] Visit `links.example.com/Qsqe1` (should redirect correctly)
   - [ ] Visit `twir.app/s/abc123` (default domain still works)
   - [ ] Visit custom domain link (view count increments)
   - [ ] Check analytics are recorded correctly

6. **Domain Management**
   - [ ] Delete custom domain
   - [ ] Verify existing links on that domain are inaccessible
   - [ ] Create new custom domain (should succeed)

7. **Edge Cases**
   - [ ] Try to create second custom domain (should fail)
   - [ ] Try to use alias that exists on default domain (should succeed on custom domain)
   - [ ] Try to verify domain with incorrect CNAME (should fail)
   - [ ] Try to create custom domain with invalid format (should fail)

## Security Considerations

1. **Input Validation**
   - Validate domain format (regex: `^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
   - Strip protocol if provided
   - Strip trailing slashes
   - Reject domains matching service domains

2. **DNS Verification**
   - Use timeout for DNS lookups (prevent hanging)
   - Rate limit verification attempts per user (max 10/hour)
   - Cache verification results for 5 minutes

3. **Authorization**
   - Users can only manage their own custom domains
   - Only verified domains can be used for short links
   - Unauthenticated users cannot configure custom domains

4. **Domain Blacklist**
   - Block common phishing/spam domains
   - Block domains that look like twir.app variations
   - Log all custom domain creation attempts

## Open Questions & Decisions

1. **Q: What happens to existing short links when user adds custom domain?**
   - **A**: They remain on default domain (domain = NULL). Only new links use custom domain.

2. **Q: Can user change their custom domain?**
   - **A**: Yes, by deleting and creating a new one. Existing links will break.

3. **Q: Should we support multiple custom domains per user?**
   - **A**: No, limit to one for simplicity. Can be extended later.

4. **Q: SSL/HTTPS for custom domains?**
   - **A**: Out of scope. Users must handle SSL at their CDN/proxy level.

5. **Q: Should unverified domains expire?**
   - **A**: Future enhancement - add cleanup job to delete unverified domains after 7 days.

6. **Q: Primary key strategy for shortened_urls?**
   - **A**: Add UUID `id` as primary key, keep `short_id` as regular column, enforce uniqueness via `(domain, short_id)` composite index.

## Success Criteria

- [ ] Users can configure one custom domain
- [ ] DNS CNAME verification works correctly
- [ ] Custom domain links redirect properly
- [ ] Slug uniqueness is enforced per domain
- [ ] HAProxy correctly routes custom domain requests
- [ ] No performance degradation
- [ ] Existing short links continue to work
- [ ] All tests pass
- [ ] No security vulnerabilities

## Estimated Effort

- **Database**: 4 hours
- **Repository Layer**: 8 hours
- **Entity Layer**: 2 hours
- **Service Layer**: 8 hours
- **HTTP Routes**: 8 hours
- **HAProxy Config**: 2 hours
- **Testing**: 8 hours
- **Documentation**: 2 hours

**Total**: ~42 hours (~5-6 development days)

## Dependencies & Risks

**Dependencies:**

- `net` package for DNS lookups (standard library)
- Existing repository patterns
- HAProxy routing capabilities

**Risks:**

1. **DNS Propagation Delays**: Users may need to wait for DNS changes
   - Mitigation: Clear instructions, retry mechanism
2. **HAProxy Regex Performance**: Complex domain matching may impact performance
   - Mitigation: Use efficient ACL patterns, monitor performance
3. **Database Migration**: Changing primary key strategy requires careful testing
   - Mitigation: Thorough testing in staging, backup before migration
4. **SSL Certificate Management**: Users responsible for SSL
   - Mitigation: Clear documentation, recommend Cloudflare

---

**End of Implementation Plan**

This plan provides a comprehensive guide to implement custom domain support for short links. Would you like me to proceed with the implementation, or do you have any questions or modifications to the plan?
