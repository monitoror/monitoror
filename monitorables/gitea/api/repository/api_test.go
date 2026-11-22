package repository

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/monitoror/monitoror/monitorables/gitea/api/models"
	"github.com/monitoror/monitoror/monitorables/gitea/config"
	"github.com/monitoror/monitoror/pkg/gravatar"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func initRepository(t *testing.T, token string, handler http.HandlerFunc) *giteaRepository {
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	conf := &config.Gitea{
		URL:       server.URL + "/",
		Token:     token,
		Timeout:   config.Default.Timeout,
		SSLVerify: true,
	}
	repository := NewGiteaRepository(conf)
	assert.Equal(t, server.URL, conf.URL)

	return repository.(*giteaRepository)
}

func TestGiteaRepository_Authorization(t *testing.T) {
	repository := initRepository(t, "TOKEN", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "token TOKEN", r.Header.Get("Authorization"))
		w.Header().Set("X-Total-Count", "0")
	})
	_, err := repository.GetCount(&models.CountParams{Owner: "o", Repository: "r"})
	assert.NoError(t, err)

	repository = initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get("Authorization"))
		w.Header().Set("X-Total-Count", "0")
	})
	_, err = repository.GetCount(&models.CountParams{Owner: "o", Repository: "r"})
	assert.NoError(t, err)
}

func TestGiteaRepository_GetCount_Success(t *testing.T) {
	repository := initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/repos/owner/repo/issues", r.URL.Path)
		assert.Equal(t, "issues", r.URL.Query().Get("type"))
		assert.Equal(t, "open", r.URL.Query().Get("state"))
		assert.Equal(t, "1", r.URL.Query().Get("limit"))
		assert.False(t, r.URL.Query().Has("labels"))
		assert.False(t, r.URL.Query().Has("q"))
		w.Header().Set("X-Total-Count", "42")
		_, _ = fmt.Fprint(w, `[]`)
	})

	count, err := repository.GetCount(&models.CountParams{Owner: "owner", Repository: "repo"})
	if assert.NoError(t, err) {
		assert.Equal(t, 42, count)
	}
}

func TestGiteaRepository_GetCount_WithFilters(t *testing.T) {
	repository := initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "pulls", r.URL.Query().Get("type"))
		assert.Equal(t, "all", r.URL.Query().Get("state"))
		assert.Equal(t, "bug,ui", r.URL.Query().Get("labels"))
		assert.Equal(t, "crash", r.URL.Query().Get("q"))
		w.Header().Set("X-Total-Count", "3")
	})

	count, err := repository.GetCount(&models.CountParams{Owner: "owner", Repository: "repo", Type: "pulls", State: "all", Labels: "bug,ui", Query: "crash"})
	if assert.NoError(t, err) {
		assert.Equal(t, 3, count)
	}
}

func TestGiteaRepository_GetCount_MissingOrInvalidHeader(t *testing.T) {
	repository := initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `[]`)
	})
	_, err := repository.GetCount(&models.CountParams{Owner: "owner", Repository: "repo"})
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "X-Total-Count")
	}

	repository = initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Total-Count", "abc")
	})
	_, err = repository.GetCount(&models.CountParams{Owner: "owner", Repository: "repo"})
	assert.Error(t, err)
}

func TestGiteaRepository_ErrorStatus(t *testing.T) {
	repository := initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"message":"`+strings.Repeat("x", 500)+`"}`)
	})

	_, err := repository.GetCount(&models.CountParams{Owner: "owner", Repository: "repo"})
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "gitea returned status 404")
		assert.Less(t, len(err.Error()), 300)
	}

	_, err = repository.GetChecks("owner", "repo", "main")
	assert.Error(t, err)
	_, err = repository.GetPullRequest("owner", "repo", 1)
	assert.Error(t, err)
	_, err = repository.GetPullRequests("owner", "repo")
	assert.Error(t, err)
	_, err = repository.GetCommit("owner", "repo", "sha")
	assert.Error(t, err)
}

func TestGiteaRepository_MalformedBody(t *testing.T) {
	repository := initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{not json`)
	})

	_, err := repository.GetChecks("owner", "repo", "main")
	assert.Error(t, err)
	_, err = repository.GetPullRequest("owner", "repo", 1)
	assert.Error(t, err)
	_, err = repository.GetPullRequests("owner", "repo")
	assert.Error(t, err)
	_, err = repository.GetCommit("owner", "repo", "sha")
	assert.Error(t, err)
}

func TestGiteaRepository_GetChecks_Success(t *testing.T) {
	repository := initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		// Ref with slash must stay a single path segment
		assert.Equal(t, "/api/v1/repos/owner/repo/commits/feature%2Ftest/statuses", r.URL.EscapedPath())
		assert.Equal(t, "50", r.URL.Query().Get("limit"))
		_, _ = fmt.Fprint(w, `[
			{"id":1,"status":"success","context":"build","url":"http://gitea/api/v1/repos/owner/repo/statuses/abc123","created_at":"2026-10-29T21:08:46Z","updated_at":"2026-10-29T21:10:46Z"},
			{"id":2,"status":"failure","context":"test","url":"http://gitea/api/v1/repos/owner/repo/statuses/abc123","created_at":"2026-10-29T21:09:00Z","updated_at":"2026-10-29T21:12:00Z"},
			{"id":3,"status":"","context":"ignored","created_at":"2026-10-29T21:09:00Z","updated_at":"2026-10-29T21:12:00Z"}
		]`)
	})

	checks, err := repository.GetChecks("owner", "repo", "feature/test")
	if assert.NoError(t, err) {
		assert.Equal(t, pointer.ToString("abc123"), checks.HeadCommit)
		assert.Equal(t, []models.Status{
			{
				ID: 1, Title: "build", State: "success",
				CreatedAt: time.Date(2026, 10, 29, 21, 8, 46, 0, time.UTC),
				UpdatedAt: time.Date(2026, 10, 29, 21, 10, 46, 0, time.UTC),
			},
			{
				ID: 2, Title: "test", State: "failure",
				CreatedAt: time.Date(2026, 10, 29, 21, 9, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 10, 29, 21, 12, 0, 0, time.UTC),
			},
		}, checks.Statuses)
	}
}

func TestGiteaRepository_GetChecks_Empty(t *testing.T) {
	repository := initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `[]`)
	})

	checks, err := repository.GetChecks("owner", "repo", "main")
	if assert.NoError(t, err) {
		assert.Nil(t, checks.HeadCommit)
		assert.Empty(t, checks.Statuses)
	}
}

func TestGiteaRepository_GetPullRequest_Success(t *testing.T) {
	repository := initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/repos/owner/repo/pulls/12", r.URL.Path)
		_, _ = fmt.Fprint(w, `{"number":12,"title":"Add feature","user":{"login":"jdoe","full_name":"John Doe","avatar_url":"https://gitea.example.com/avatars/1"},
			"head":{"ref":"feature","sha":"abc123","repo":{"name":"repo-fork","owner":{"login":"fork-owner"}}}}`)
	})

	pullRequest, err := repository.GetPullRequest("owner", "repo", 12)
	if assert.NoError(t, err) {
		expected := &models.PullRequest{
			ID: 12, Title: "Add feature",
			SourceOwner: "fork-owner", SourceRepository: "repo-fork", SourceBranch: "feature", CommitSHA: "abc123",
		}
		expected.Author.Name = "John Doe"
		expected.Author.AvatarURL = "https://gitea.example.com/avatars/1"
		assert.Equal(t, expected, pullRequest)
	}
}

func TestGiteaRepository_GetPullRequest_DeletedForkAndLoginFallback(t *testing.T) {
	repository := initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"number":3,"title":"Fix","user":{"login":"jdoe","full_name":""},"head":{"ref":"fix","sha":"def456","repo":null}}`)
	})

	pullRequest, err := repository.GetPullRequest("owner", "repo", 3)
	if assert.NoError(t, err) {
		assert.Equal(t, "", pullRequest.SourceOwner)
		assert.Equal(t, "", pullRequest.SourceRepository)
		assert.Equal(t, "jdoe", pullRequest.Author.Name)
	}
}

func TestGiteaRepository_GetPullRequests_Pagination(t *testing.T) {
	PageSize = 2
	t.Cleanup(func() { PageSize = 50 })

	var pages []string
	repository := initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/repos/owner/repo/pulls", r.URL.Path)
		assert.Equal(t, "open", r.URL.Query().Get("state"))
		assert.Equal(t, "2", r.URL.Query().Get("limit"))

		page := r.URL.Query().Get("page")
		pages = append(pages, page)
		switch page {
		case "1":
			_, _ = fmt.Fprint(w, `[{"number":1,"head":{"ref":"a"}},{"number":2,"head":{"ref":"b"}}]`)
		default:
			_, _ = fmt.Fprint(w, `[{"number":3,"head":{"ref":"c"}}]`)
		}
	})

	pullRequests, err := repository.GetPullRequests("owner", "repo")
	if assert.NoError(t, err) {
		assert.Equal(t, []string{"1", "2"}, pages)
		if assert.Len(t, pullRequests, 3) {
			assert.Equal(t, 1, pullRequests[0].ID)
			assert.Equal(t, 3, pullRequests[2].ID)
		}
	}
}

func TestGiteaRepository_GetPullRequests_ServerCapsPageSize(t *testing.T) {
	// The instance returns 2 items per page (MAX_RESPONSE_ITEMS) while 50 are requested
	var pages []string
	repository := initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		pages = append(pages, page)
		w.Header().Set("X-Total-Count", "5")
		switch page {
		case "1":
			_, _ = fmt.Fprint(w, `[{"number":1},{"number":2}]`)
		case "2":
			_, _ = fmt.Fprint(w, `[{"number":3},{"number":4}]`)
		default:
			_, _ = fmt.Fprint(w, `[{"number":5}]`)
		}
	})

	pullRequests, err := repository.GetPullRequests("owner", "repo")
	if assert.NoError(t, err) {
		assert.Equal(t, []string{"1", "2", "3"}, pages)
		assert.Len(t, pullRequests, 5)
	}
}

func TestGiteaRepository_GetPullRequests_StopsOnEmptyPageAndMaxPages(t *testing.T) {
	calls := 0
	repository := initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("X-Total-Count", "1000")
		if r.URL.Query().Get("page") == "3" {
			_, _ = fmt.Fprint(w, `[]`)
			return
		}
		_, _ = fmt.Fprint(w, `[{"number":1}]`)
	})

	pullRequests, err := repository.GetPullRequests("owner", "repo")
	if assert.NoError(t, err) {
		assert.Equal(t, 3, calls)
		assert.Len(t, pullRequests, 2)
	}

	// A server ignoring "page" never returns an empty page: stop after maxPages
	calls = 0
	repository = initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("X-Total-Count", "1000")
		_, _ = fmt.Fprint(w, `[{"number":1}]`)
	})

	pullRequests, err = repository.GetPullRequests("owner", "repo")
	if assert.NoError(t, err) {
		assert.Equal(t, maxPages, calls)
		assert.Len(t, pullRequests, maxPages)
	}
}

func TestTruncate_KeepsValidUTF8(t *testing.T) {
	s := strings.Repeat("a", maxErrorBodyLength-1) + "é" + "tail"
	result := truncate(s)
	assert.True(t, utf8.ValidString(result))
	assert.Equal(t, strings.Repeat("a", maxErrorBodyLength-1)+"...", result)
	assert.Equal(t, "short", truncate("short"))
}

func TestGiteaRepository_GetCommit_GiteaAvatar(t *testing.T) {
	repository := initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/repos/owner/repo/git/commits/abc123", r.URL.Path)
		_, _ = fmt.Fprint(w, `{"sha":"abc123","commit":{"author":{"name":"John Doe","email":"john@example.com"}},
			"author":{"login":"jdoe","full_name":"John D.","avatar_url":"https://gitea.example.com/avatars/1"}}`)
	})

	commit, err := repository.GetCommit("owner", "repo", "abc123")
	if assert.NoError(t, err) {
		assert.Equal(t, "abc123", commit.SHA)
		assert.Equal(t, "John Doe", commit.Author.Name)
		assert.Equal(t, "https://gitea.example.com/avatars/1", commit.Author.AvatarURL)
	}
}

func TestGiteaRepository_GetCommit_GravatarFallback(t *testing.T) {
	repository := initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"sha":"abc123","commit":{"author":{"name":"John Doe","email":"john@example.com"}},"author":null}`)
	})

	commit, err := repository.GetCommit("owner", "repo", "abc123")
	if assert.NoError(t, err) {
		assert.Equal(t, "John Doe", commit.Author.Name)
		assert.Equal(t, gravatar.GetGravatarURL("john@example.com"), commit.Author.AvatarURL)
	}
}

func TestGiteaRepository_GetCommit_NameFallback(t *testing.T) {
	repository := initRepository(t, "", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"sha":"abc123","commit":{"author":{"name":"","email":""}},"author":{"login":"jdoe","full_name":""}}`)
	})

	commit, err := repository.GetCommit("owner", "repo", "abc123")
	if assert.NoError(t, err) {
		assert.Equal(t, "jdoe", commit.Author.Name)
		assert.Equal(t, "", commit.Author.AvatarURL)
	}
}
