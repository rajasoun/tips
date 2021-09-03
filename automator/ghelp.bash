#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/automator/src/lib/os.sh"

case "$OSTYPE" in
darwin*)
	echo "${GREEN}Welcome $(git config user.name) | OS: OSX${NC}"
	;;
linux*)
	echo "${GREEN}${BOLD}Welcome $(git config user.name) | OS: Linux${NC}"
	;;
msys*)
	echo "${GREEN}${BOLD}Welcome $(git config user.name) | OS: Windows${NC}"
	;;
*)
	echo "unknown: $OSTYPE"
	;;
esac

function ghelp() {
	echo "
- - - - - - - - - - - - - -
Git Convenience Shortcuts:
- - - - - - - - - - - - - -
ghelp 		- List all Git Convenience commands and prompt symbols
gsetup		- Install Git Flow & pre-commit hooks
glogin		- Web Login to GitHub
gstatus		- GitHub Login status
code_churn	- Frequency of change to code base
alias		- List all Alias
pretty		- Code prettier
- - - - - - - - - - - - - -
"
	gextras
}

function gextras() {
	echo "
- - - - - - - - - - - - - -
Extra Command For Quick Fix:
- - - - - - - - - - - - - -
git-ssh-check		- Check Git SSH Works
git-ssh-fix		- Fix Git SSH Permission denied (publickey) Issue
init-debug		- Initialize Debug for Open Source Sentry
release			- Release through Automation
- - - - - - - - - - - - - -
"
}

function _git_tag() {
	CUR_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
	if [ "$CUR_BRANCH" != "main" ]; then
		echo "${RED} Need to be in main branch ${NC}"
		return 1
	fi

	GIT_CLEAN="$(git status --porcelain)"
	if [ -z "$GIT_CLEAN" ]; then
		git fetch --prune --tags
		VERSION=$(git describe --tags --abbrev=0 | awk -F. '{OFS="."; $NF+=1; print $0}')
		git tag -a "$VERSION" -m "tip : $VERSION | For Release"
		git push origin "$VERSION" --no-verify
		git fetch --prune --tags
	else
		echo "${RED} Git Not Clean... ${NC}"
	fi
}

function _teardown_git_flow() {
	git config --remove-section "gitflow.path"
	git config --remove-section "gitflow.prefix"
	git config --remove-section "gitflow.branch"
}

function _init_git_flow() {
	if (git flow config >/dev/null 2>&1); then
		prompt "Git Flow Already Initialized..."
	else
		prompt "Git Flow Not Initialized. Initializing..."
		git flow init -fd
	fi
}

function _install_git_hooks() {
	prompt "Git Repository..."
	prompt "Installing Git Hooks"
	git config --unset-all core.hooksPath
	pre-commit uninstall --hook-type commit-msg
	pre-commit install --hook-type commit-msg
	pre-commit install-hooks
	npm --prefix ./shift-left install
}

function _check_gg_api() {
	prompt "Checking Git Guardian API Validity"
	# shellcheck source=/dev/null
	source ".env"
	curl -H "Authorization: Token ${GITGUARDIAN_API_KEY}" "${GITGUARDIAN_API_URL}/v1/health"
	prompt ""
}

function _populate_dot_env() {
	prompt "Populating .env File"
	if [ -f "$(git rev-parse --show-toplevel)/.env" ]; then
		mv .env .env.bak
	fi
	cp .env.sample .env

	prompt "${BLUE}To Get the GitHub Key  ${NC}"
	prompt "${YELLOW} Visit https://www.$(dotenv -f .env.sample get GITHUB_URL)/settings/tokens ${NC}"
	prompt "${BOLD}Enter Git Token: ${NC}"
	read -r GITTOKEN
	_file_replace_text "1__________FILL_ME__________1" "$GITTOKEN" "$(git rev-parse --show-toplevel)/.env"

	prompt "${BLUE}To Get the GG Key - Register to Git Guardian ${NC}"
	prompt "${YELLOW} Visit $(dotenv -f .env.sample get GITGUARDIAN_URL) ${NC}"
	prompt "${BOLD}Enter Git Guardian API Key: ${NC}"
	read -r GG_KEY
	_file_replace_text "2__________FILL_ME__________2" "$GG_KEY" "$(git rev-parse --show-toplevel)/.env"
	_check_gg_api

	prompt "${BLUE}To Get the Sentry DSN  ${NC}"
	prompt "${YELLOW} Visit $(dotenv -f .env.sample get SENTRY_URL) ${NC}"
	prompt "${BOLD}Enter Sentry DSN: ${NC}"
	read -r SENTRYDSN
	_file_replace_text "3__________FILL_ME__________3" "$SENTRYDSN" "$(git rev-parse --show-toplevel)/.env"
}

function gsetup() {
	if [ "$(git rev-parse --is-inside-work-tree)" = true ]; then
		if [[ $(git diff --stat) != '' ]]; then
			prompt "${RED} Git Working Tree Not Clean. Aborting setup !!! ${NC}"
			EXIT_CODE=1
			log_sentry "$EXIT_CODE" "gsetup | Git Working Tree Not Clean. Aborting setup"
		else
			start=$(date +%s)
			prompt "Git Working Tree Clean - Shell will Abort on Error"
			_install_git_hooks || prompt "_install_git_hooks ❌"
			_populate_dot_env || prompt "_populate_dot_env ❌"
			_init_git_flow || prompt "_init_git_flow ❌ [Proceeding...]"
			end=$(date +%s)
			runtime=$((end - start))
			prompt "gsetup DONE in $(_display_time $runtime)"
			.devcontainer/tests/system/e2e_tests.sh
			EXIT_CODE="$?"
			log_sentry "$EXIT_CODE" "gsetup "
		fi
	fi
}

# Gits Churn -  "frequency of change to code base"
#
# $ ./git-churn.bash
# 30 src/multipass/actions.bash
# 38 test/test_integration.bats
# 97 .github/workflows/pipeline.yml
#
# This means that
# actions.bash has changed 30 times.
# pipeline.yml has changed 97 times.
#
# Show churn for specific directories:
#   $ $ ./git-churn.bash src
#
# Show churn for a time range:
#   $ $ ./git-churn.bash --since='1 month ago'
#
# All standard arguments to git log are applicable
function code_churn() {
	git log --all -M -C --name-only --format='format:' "$@" | sort | grep -v '^$' | uniq -c | sort -n
}

function git_push() {
	git push
	EXIT_CODE="$?"
	log_sentry "$EXIT_CODE" "git push | Branch: $(git rev-parse --abbrev-ref HEAD)"
}

function check_git_config() {
	user_name=$(git config user.name)
	user_email=$(git config user.email)
	if [ -z "$user_name" ] || [ -z "$user_email" ]; then
		log_sentry "1" "git config | user_name and user_email not set"
		_git_config
	else
		log_sentry "0" "git config "
	fi
}

function git_hub_login() {
	AUTH_TYPE=$1
	case "$AUTH_TYPE" in
	web)
		gh auth login --hostname "$(dotenv get GITHUB_URL)" --web
		EXIT_CODE="$?"
		log_sentry "$EXIT_CODE" "Github Login via Web"
		;;
	token)
		GT="$(dotenv get GITHUB_TOKEN)"
		[ "$GT" ] || echo "GITHUB_TOKEN Not set in .env"
		if [ ! -f "token.txt" ]; then
			echo "$GT" >automator/token.txt
		fi
		gh auth login --hostname "$(dotenv get GITHUB_URL)" --with-token <automator/token.txt
		EXIT_CODE="$?"
		log_sentry "$EXIT_CODE" "Github Login via Token"
		;;
	*) gh auth login --hostname "$(dotenv get GITHUB_URL)" --web ;;
	esac
}

#-------------------------------------------------------------
# Git Alias Commands
#-------------------------------------------------------------
alias gss="git status -s"
alias gaa="git add --all"
alias gc="git commit"
alias gp="git_push"
alias gclean="git fetch --prune origin && git gc"
#shellcheck disable=SC2139
#shellcheck disable=SC2145
alias glogin="git_hub_login $@"
alias gstatus="gh auth status --hostname dotenv get GITHUB_URL "
#alias release="npm --prefix shift-left run release"
# alias gtag="_git_tag"
alias release="_git_tag"

#-------------------------------------------------------------
# Generic Alias Commands
#-------------------------------------------------------------
alias pretty="npx prettier --config shift-left/.prettierrc.yml --write ."
alias git-ssh-check='ssh -T git@$(dotenv get GITHUB_URL)'
alias init-debug='init_debug'

if ! [ -f "$(git rev-parse --show-toplevel)/.env" ]; then
	prompt "${YELLOW} Starting gsetup ${NC}"
	gsetup
fi
glogin token #git-ssh-fix
init-debug
EXIT_CODE="$?"
log_sentry "$EXIT_CODE" "DevContainer Initialization"
check_git_config
