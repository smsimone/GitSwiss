# GitSwiss

A collection of git commands to help managing mutliple git repositories at once.

> [!IMPORTANT]
> Each of the following commands will iterate over all the repositories in the given directory and execute the given command.

## Available commands

- [**align-branch**](#align-branch)
- [**find-branch**](#find-branch)
- [**create-branch**](#create-branch)
- [**update-branches**](#create-branch)

## Align Branch

For each repository in the given directory, it will align the target branch to the given source branch.

```bash
swit align-branch -source <source-branch> -target <target-branch> [-directory <directory>] [-strategy <pull|merge>] [-remote <remote_name>]
```
*Default strategy: merge*
*Default remote:   origin*

## Find Branch

Iterates all the given repositories in search of a given branch.

```bash
swit find-branch -branch <branch> [-directory <directory>]
```

## Create branch

Createe a branch in the current repository and push it to the remote.

```bash
swit create-branch -target <branch> [-source <source-branch>] [-directory <directory>]
```

## Update branches

Update the branch in the given directory. If no branch is given, it will update the current branch, otherwise will update all the matching branches (it will search with a like)

```bash
swit update-branch -directory <directory> [-branch <branch>]
```
