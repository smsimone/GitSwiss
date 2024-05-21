# GitSwiss

A collection of git commands to help managing mutliple git repositories at once.

> [!IMPORTANT]
> Each of the following commands will iterate over all the repositories in the given directory and execute the given command.

## Available commands

- **align-branch**
- **find-branch**

### To be added

- [ ] **create-branch**
- [ ] **update-branches**
- [ ] ...

# Align Branch

For each repository in the given directory, it will align the target branch to the given source branch.

```bash
git-swiss align-branch -source <source-branch> -target <target-branch> -directory <directory>
```

# Find Branch

Iterates all the given repositories in search of a given branch.

```bash
git-swiss find-branch -branch <branch> [-directory <directory>]
```

# Create branch

Createe a branch in the current repository and push it to the remote.

```bash
git-swiss create-branch -branch <branch> [-source <source-branch>] [-directory <directory>]
```

# Update branches

Update the branch in the given directory. If no branch is given, it will update the current branch, otherwise will update all the matching branches (it will search with a like)

```bash
git-swiss update-branch -directory <directory> [-branch <branch>]
```
