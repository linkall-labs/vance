// Copyright 2022 Linkall Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"context"

	"github.com/google/go-github/v58/github"
	lodash "github.com/samber/lo"
)

func (s *GitHubAPISource) startPR(ctx context.Context) {
	for i := range s.config.PRConfigs {
		config := s.config.PRConfigs[i]
		s.listPullRequests(ctx, config)
		s.logger.Info().Str("organization", config.Organization).
			Str("repo", config.Repo).
			Int("prs", s.numPRs).
			Msg("listPullRequests")
	}
}

func (s *GitHubAPISource) listPullRequests(ctx context.Context, config PRConfig) {
	listOption := &github.PullRequestListOptions{
		State: "all",
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 250,
		},
	}

	for {
		s.Limiter.Take()
		prs, resp, err := s.client.PullRequests.List(ctx, config.Organization, config.Repo, listOption)
		if err != nil {
			s.logger.Warn().Err(err).Msg("pull request list error")
		}
		if len(prs) == 0 {
			break
		}

		for _, pr := range prs {
			s.Limiter.Take()
			prDetail, _, _ := s.client.PullRequests.Get(ctx, config.Organization, config.Repo, *pr.Number)
			s.prInfo(ctx, config, prDetail)
		}

		if resp.NextPage <= listOption.ListOptions.Page {
			break
		}
		listOption.ListOptions.Page = resp.NextPage
	}
}

func (s *GitHubAPISource) prInfo(ctx context.Context, config PRConfig, pr *github.PullRequest) {
	if !lodash.Contains(config.UserList, pr.GetUser().GetLogin()) {
		return
	}

	s.numPRs += 1
	data := make(map[string]interface{})
	data["repo"] = config.Repo
	data["org"] = config.Organization
	data["uid"] = pr.GetUser().GetLogin()
	data["username"] = pr.GetUser().GetName()
	data["email"] = pr.GetUser().GetEmail()
	data["company"] = pr.GetUser().GetCompany()
	data["type"] = "PR"
	data["state"] = pr.GetState()
	data["link"] = pr.GetHTMLURL()
	data["lines"] = pr.GetAdditions() + pr.GetDeletions()
	data["updateAt"] = pr.GetUpdatedAt()
	data["branch"] = pr.GetBase().GetRef()

	s.sendEvent("prs", config.Organization, data)
}
