// Copyright © 2022-2023 Obol Labs Inc. Licensed under the terms of a Business Source License 1.1

// Command verifypr provides a tool to verify charon PRs against the template defined in docs/contibuting.md.
//
//nolint:revive,cyclop
package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var titlePrefix = regexp.MustCompile(`^[*\w]+(/[*\w]+)?$`)

type PR struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	ID    string `json:"node_id"`
}

// PRFromEnv returns the PR by parsing it from "GITHUB_PR" env var or an error.
func PRFromEnv() (PR, error) {
	const prEnv = "GITHUB_PR"
	prJSON, ok := os.LookupEnv(prEnv)
	if !ok || strings.TrimSpace(prJSON) == "" {
		return PR{}, errors.New("env variable not set")
	}

	var pr PR
	if err := json.Unmarshal([]byte(prJSON), &pr); err != nil {
		return PR{}, err
	}

	if pr.Title == "" || pr.Body == "" || pr.ID == "" {
		return PR{}, errors.New("pr field not set")
	}

	return pr, nil
}

// verify returns an error if the PR doesn't correspond to the template defined in docs/contibuting.md.
func verify() error {
	pr, err := PRFromEnv()
	if err != nil {
		return err
	}

	// Skip dependabot PRs.
	if strings.Contains(pr.Title, "build(deps)") && strings.Contains(pr.Body, "dependabot") {
		return nil
	}

	log.Printf("Verifying charon PR against template\n")
	log.Printf("PR Title: %s\n", pr.Title)
	log.Printf("## PR Body:\n%s\n####\n", pr.Body)

	if err := verifyTitle(pr.Title); err != nil {
		return err
	}

	if err := verifyBody(pr.Body); err != nil {
		return err
	}

	return nil
}

func verifyTitle(title string) error {
	const maxTitleLen = 60
	if len(title) > maxTitleLen {
		return errors.New("title too long")
	}

	split := strings.SplitN(title, ":", 2)
	if len(split) < 2 {
		return errors.New("title isn't prefixed with 'package[/subpackage]:'")
	}

	if !titlePrefix.Match([]byte(split[0])) {
		return errors.New("title prefix doesn't match regex")
	}

	suffix := split[1]

	if len(suffix) < 5 {
		return errors.New("title suffix too short")
	}

	if len(suffix) < 5 {
		return errors.New("title suffix too short")
	}

	if suffix[0] != ' ' {
		return errors.New("title prefix not followed by space")
	}

	suffix = suffix[1:]

	if unicode.IsUpper(rune(suffix[0])) {
		return errors.New("title suffix shouldn't start with a capital")
	}

	if unicode.IsPunct(rune(suffix[len(suffix)-1])) {
		return errors.New("title suffix shouldn't end with punctuation")
	}

	return nil
}

//nolint:gocyclo
func verifyBody(body string) error {
	if strings.TrimSpace(body) == "" {
		return errors.New("body empty")
	}
	if strings.Contains(body, "<!--") {
		return errors.New("instructions not deleted (markdown comments present)")
	}

	var (
		prevLineEmpty bool
		foundCategory bool
	)
	for i, line := range strings.Split(body, "\n") {
		if i == 0 && strings.TrimSpace(line) == "" {
			return errors.New("first line empty")
		}

		const catTag = "category:"
		if strings.HasPrefix(line, catTag) {
			if foundCategory {
				return errors.New("multiple category tag lines")
			}
			if !prevLineEmpty {
				return errors.New("category tag not preceded by empty line")
			}

			cat := strings.TrimSpace(strings.TrimPrefix(line, catTag))

			if cat == "" {
				return errors.New("category tag empty")
			}

			var (
				ok     bool
				allows = []string{"feature", "bug", "refactor", "docs", "test", "fixbuild", "misc"}
			)
			for _, allow := range allows {
				if allow == cat {
					ok = true
					break
				}
			}

			if !ok {
				return errors.New("invalid category")
			}

			foundCategory = true
		}

		prevLineEmpty = strings.TrimSpace(line) == ""
	}

	if !foundCategory {
		return errors.New("missing category tag")
	}

	return nil
}
