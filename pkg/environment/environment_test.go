package environment

import (
	"errors"
	"testing"
)

func TestParseRepo(t *testing.T) {
	t.Parallel()

	t.Run("gets nothing", func(t *testing.T) {
		t.Parallel()

		_, err := ParseRepo("")
		if !errors.Is(err, ErrNoRepo) {
			t.Fatalf("want %v but got error: %v", ErrNoRepo, err)
		}
	})

	t.Run("gets nothing", func(t *testing.T) {
		t.Parallel()

		_, err := ParseRepo("owner+repo+is+invalid")
		if !errors.Is(err, ErrInvalidRepo) {
			t.Fatalf("want %v but got error: %v", ErrInvalidRepo, err)
		}
	})

	t.Run("gets the parsed Repo", func(t *testing.T) {
		t.Parallel()

		got, err := ParseRepo("owner/repo")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Owner.Login != "owner" || got.Repo != "repo" {
			t.Fatalf("want %v but got %+v", "owner/repo", got)
		}
	})
}

func TestMissingOrRedacted(t *testing.T) {
	t.Parallel()

	if missingOrRedacted("token") != Redacted {
		t.Fatalf("want redacted but got something else")
	}

	if missingOrRedacted("") != Missing {
		t.Fatalf("want missingd but got something else")
	}
}
