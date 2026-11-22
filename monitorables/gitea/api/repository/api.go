package repository

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/monitoror/monitoror/monitorables/gitea/api"
	"github.com/monitoror/monitoror/monitorables/gitea/api/models"
	"github.com/monitoror/monitoror/monitorables/gitea/config"
	"github.com/monitoror/monitoror/pkg/gravatar"
)

type (
	giteaRepository struct {
		client *http.Client
		config *config.Gitea
	}

	user struct {
		Login     string `json:"login"`
		FullName  string `json:"full_name"`
		AvatarURL string `json:"avatar_url"`
	}

	commitStatus struct {
		ID        int64     `json:"id"`
		Status    string    `json:"status"`
		Context   string    `json:"context"`
		URL       string    `json:"url"` // .../repos/{owner}/{repo}/statuses/{sha}
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	repository struct {
		Name  string `json:"name"`
		Owner *user  `json:"owner"`
	}

	pullRequestBranch struct {
		Ref  string      `json:"ref"`
		SHA  string      `json:"sha"`
		Repo *repository `json:"repo"` // nil when the fork has been deleted
	}

	pullRequest struct {
		Number int               `json:"number"`
		Title  string            `json:"title"`
		User   *user             `json:"user"`
		Head   pullRequestBranch `json:"head"`
	}

	commit struct {
		SHA    string `json:"sha"`
		Commit struct {
			Author struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
		} `json:"commit"`
		Author *user `json:"author"` // nil when the commit email is not linked to a Gitea user
	}
)

// PageSize is the number of items requested per page. 50 is Gitea's default MAX_RESPONSE_ITEMS.
var PageSize = 50

const (
	apiPath            = "/api/v1"
	maxErrorBodyLength = 200
	maxBodySize        = 10 << 20 // 10 MiB
	maxPages           = 100
	totalCountHeader   = "X-Total-Count"
)

func NewGiteaRepository(config *config.Gitea) api.Repository {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.SSLVerify},
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(config.Timeout) * time.Millisecond}

	// Remove last /
	config.URL = strings.TrimRight(config.URL, "/")

	return &giteaRepository{client, config}
}

func (r *giteaRepository) GetCount(params *models.CountParams) (int, error) {
	query := url.Values{}
	query.Set("type", defaultString(params.Type, "issues"))
	query.Set("state", defaultString(params.State, "open"))
	if params.Labels != "" {
		query.Set("labels", params.Labels)
	}
	if params.Query != "" {
		query.Set("q", params.Query)
	}
	// Only the total count header is used
	query.Set("limit", "1")

	header, err := r.get(fmt.Sprintf("/repos/%s/%s/issues", url.PathEscape(params.Owner), url.PathEscape(params.Repository)), query, nil)
	if err != nil {
		return 0, err
	}

	rawCount := header.Get(totalCountHeader)
	if rawCount == "" {
		return 0, fmt.Errorf("gitea response is missing the %s header", totalCountHeader)
	}

	count, err := strconv.Atoi(rawCount)
	if err != nil {
		return 0, fmt.Errorf("invalid %s header: %q", totalCountHeader, rawCount)
	}

	return count, nil
}

// GetChecks loads the commit status history of ref (one row per state change, most recent first).
// The history is used instead of the combined status: the combined status only keeps the last row
// of each context, which hides when a build started and makes every duration zero.
// Only the most recent page is loaded, which is enough to cover the current build.
func (r *giteaRepository) GetChecks(owner, repository, ref string) (*models.Checks, error) {
	query := url.Values{}
	query.Set("limit", strconv.Itoa(PageSize))

	var result []commitStatus
	path := fmt.Sprintf("/repos/%s/%s/commits/%s/statuses", url.PathEscape(owner), url.PathEscape(repository), url.PathEscape(ref))
	if _, err := r.get(path, query, &result); err != nil {
		return nil, err
	}

	checks := &models.Checks{Statuses: []models.Status{}}
	for _, s := range result {
		if s.Context == "" || s.Status == "" {
			continue
		}

		if checks.HeadCommit == nil && s.URL != "" {
			sha := s.URL[strings.LastIndex(s.URL, "/")+1:]
			checks.HeadCommit = &sha
		}

		checks.Statuses = append(checks.Statuses, models.Status{
			ID:        s.ID,
			Title:     s.Context,
			State:     s.Status,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		})
	}

	return checks, nil
}

func (r *giteaRepository) GetPullRequest(owner, repository string, id int) (*models.PullRequest, error) {
	result := &pullRequest{}
	path := fmt.Sprintf("/repos/%s/%s/pulls/%d", url.PathEscape(owner), url.PathEscape(repository), id)
	if _, err := r.get(path, nil, result); err != nil {
		return nil, err
	}

	return toPullRequest(result), nil
}

func (r *giteaRepository) GetPullRequests(owner, repository string) ([]models.PullRequest, error) {
	var results []models.PullRequest
	path := fmt.Sprintf("/repos/%s/%s/pulls", url.PathEscape(owner), url.PathEscape(repository))

	for page := 1; page <= maxPages; page++ {
		query := url.Values{}
		query.Set("state", "open")
		query.Set("page", strconv.Itoa(page))
		query.Set("limit", strconv.Itoa(PageSize))

		var pullRequests []pullRequest
		header, err := r.get(path, query, &pullRequests)
		if err != nil {
			return nil, err
		}

		if len(pullRequests) == 0 {
			return results, nil
		}

		for i := range pullRequests {
			results = append(results, *toPullRequest(&pullRequests[i]))
		}

		// Gitea may return fewer items than requested (MAX_RESPONSE_ITEMS), so rely on the total when available
		if total, err := strconv.Atoi(header.Get(totalCountHeader)); err == nil {
			if len(results) >= total {
				return results, nil
			}
		} else if len(pullRequests) < PageSize {
			return results, nil
		}
	}

	return results, nil
}

func (r *giteaRepository) GetCommit(owner, repository, sha string) (*models.Commit, error) {
	result := &commit{}
	path := fmt.Sprintf("/repos/%s/%s/git/commits/%s", url.PathEscape(owner), url.PathEscape(repository), url.PathEscape(sha))
	if _, err := r.get(path, nil, result); err != nil {
		return nil, err
	}

	c := &models.Commit{SHA: sha}
	c.Author.Name = result.Commit.Author.Name

	// Prefer the Gitea avatar: self-hosted instances often don't use Gravatar
	if result.Author != nil && result.Author.AvatarURL != "" {
		c.Author.AvatarURL = result.Author.AvatarURL
	} else if result.Commit.Author.Email != "" {
		c.Author.AvatarURL = gravatar.GetGravatarURL(result.Commit.Author.Email)
	}

	if c.Author.Name == "" && result.Author != nil {
		c.Author.Name = defaultString(result.Author.FullName, result.Author.Login)
	}

	return c, nil
}

func (r *giteaRepository) get(path string, query url.Values, result interface{}) (http.Header, error) {
	rawURL := r.config.URL + apiPath + path
	if len(query) > 0 {
		rawURL += "?" + query.Encode()
	}

	request, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	if r.config.Token != "" {
		request.Header.Set("Authorization", "token "+r.config.Token)
	}

	response, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(io.LimitReader(response.Body, maxBodySize))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gitea returned status %d: %s", response.StatusCode, truncate(strings.TrimSpace(string(data))))
	}

	if result != nil {
		if err = json.Unmarshal(data, result); err != nil {
			return nil, err
		}
	}

	return response.Header, nil
}

func toPullRequest(pr *pullRequest) *models.PullRequest {
	result := &models.PullRequest{
		ID:           pr.Number,
		Title:        pr.Title,
		SourceBranch: pr.Head.Ref,
		CommitSHA:    pr.Head.SHA,
	}

	if pr.Head.Repo != nil {
		result.SourceRepository = pr.Head.Repo.Name
		if pr.Head.Repo.Owner != nil {
			result.SourceOwner = pr.Head.Repo.Owner.Login
		}
	}

	if pr.User != nil {
		result.Author.Name = defaultString(pr.User.FullName, pr.User.Login)
		result.Author.AvatarURL = pr.User.AvatarURL
	}

	return result
}

func defaultString(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func truncate(s string) string {
	if len(s) <= maxErrorBodyLength {
		return s
	}
	return strings.ToValidUTF8(s[:maxErrorBodyLength], "") + "..."
}
