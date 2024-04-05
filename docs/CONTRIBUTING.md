# Contributing :: SnippetBox

## Branches :: SnippetBox

1. Branches should be named for the package they pertain to, and then the feature being worked on, e.g. `web_authentication`.
2. No long lived branches. Keep PR's tight, focused and small and merge back into main frequently.
3. No separate develop branch. We practice trunk(ish) based development. You may branch off of main to do some work and open a PR,
   but if the tests pass it should be merged back to main and deployed automatically.

## Merge Request Rules

### Overview

1. Make [atomic](https://www.aleksandrhovhannisyan.com/blog/atomic-git-commits/) commits.
2. Use [conventional commit messages](https://www.conventionalcommits.org/en/v1.0.0/).
3. Use the commit body to provide context. See Note 5 below.
   1. Use the body of a commit message where it makes sense to provide additional details and context about the changes in the commit if necessary, not what is in the commit. Reasoning for a change, other technical considerations etc that would be useful to you or someone else understanding why that change was made in the future.
   2. [See these suggestions.](https://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html)
4. Follow the same “conventional” naming and detail for MR titles and descriptions.
5. All tests must pass on every commit that makes it to the main branch.
6. Avoid mixing feature implementation, refactoring, bug fixes, etc in the same MR.
   - This is often caused by scope creep due to (often valid) input from reviewers.
   - Consider making subtasks in response - to action in a later MR.
   - Keeping MRs tight and focused helps it get reviewed and merged faster too ☺️
7. Keep a linear or [semi-linear](https://dev.to/akorda/semi-linear-git-history-1191) git history
   - Don’t merge the parent branch, creating “train tracks”.
   - Fully [rebase](https://www.atlassian.com/git/tutorials/rewriting-history/git-rebase) on parent and target branches before merging.
   - There are pros and cons, for example [see this discussion](https://stackoverflow.com/questions/20348629/what-are-the-advantages-of-keeping-linear-history-in-git), but I’ve found this approach still works very well even on teams of ~15 contributors.

### Notes and Caveats

1. It is OK to relax on these points while in _Draft_, but fix before marking it as _Ready_.
2. If you plan to use a _squash commit_, make this clear in the MR. It is OK then to have smaller commits that break some of the rules - but the final squashed commit should satisfy them!
3. It is on both the submitters and reviewers to check and uphold these standards.
4. Using automatic merge, without semi-linear protection configured in Gitlab, might result in (7) being violated. Use with caution.

## Related skills

Perhaps we want to include tutorials and / or videos around the following skills?

- [Manual squashing](https://stackoverflow.com/a/5201642)
- [Interactive rebasing](https://git-scm.com/book/en/v2/Git-Tools-Rewriting-History)
- [Branch stacking](https://graphite.dev/)
- [Rebasing - what can go wrong?](https://jvns.ca/blog/2023/11/06/rebasing-what-can-go-wrong-/)
