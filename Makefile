# --------
# Manifest
# --------
PROJECT_NAME:="PDS Goldcard Service"
PROJECT_PKG:=srv-goldcard
DOCKER_NAMESPACE:=artifactory.pegadaian.co.id:5443

# -----------------
# Project Variables
# -----------------
BINARY_NAME:=srv-goldcard
PROJECT_ROOT?=$(shell pwd)
PROJECT_WORKDIR?=${PROJECT_ROOT}
PROJECT_CONFIG:=.env
PROJECT_CONFIG_RELEASE:=.env
PROJECT_WEB_TEMPLATES=web/templates
PROJECT_WEB_STATIC=web/static
PROJECT_MAIN_PKG=cmd
PROJECT_DOCKERFILE_DIR?=${PROJECT_ROOT}/deployment
APP_PATH=${PROJECT_ROOT}
OUTPUT_DIR:=${PROJECT_ROOT}/bin
DOCTOR_CMD:=${PROJECT_ROOT}/script/doctor.sh
SCRIPTS_DIR:=${PROJECT_ROOT}/script
DEPLOYMENT_DIR:=${PROJECT_ROOT}/deployment
DOCKER_COMPOSE_FILE:=${DEPLOYMENT_DIR}/docker-compose.yml


# ---------------
# Command Aliases
# ---------------
GO_CMD:=go
GO_BUILD:=${GO_CMD} build
GO_MOD:=${GO_CMD} mod
GO_CLEAN:=${GO_CMD} clean
GO_GET:=${GO_CMD} get
DOCKER_CMD:=docker

# ----------------------
# Debug Output Variables
# ----------------------
DEBUG_DIR:=${OUTPUT_DIR}/debug
DEBUG_BIN:=${DEBUG_DIR}/${BINARY_NAME}
# DEBUG_ENV_FILES:=$(addprefix ${DEBUG_DIR}/,${PROJECT_CONFIG} ${PROJECT_RESPONSES})


# ------------------------
# Release Output Variables
# ------------------------
RELEASE_OUTPUT_DIR:=${OUTPUT_DIR}/release
RELEASE_ENV_APP_ENV?=1
RELEASE_ENV_LOG_LEVEL?=error
RELEASE_ENV_LOG_FORMAT?=console

# ----------------
# Docker Variables
# ----------------
CI_PROJECT_PATH ?= srv-goldcard
CI_COMMIT_REF_SLUG ?= local

IMAGE_APP ?= $(DOCKER_NAMESPACE)/$(CI_PROJECT_PATH)
IMAGE_APP_TAG ?= $(CI_COMMIT_REF_SLUG)

# -------------------
# Migration Variables
# -------------------
MIGRATION_DIR := ${PROJECT_ROOT}/migrations
MIGRATION_SRC_DIR := ${MIGRATION_DIR}/sql

# -----------
# API Version
# -----------
CI_COMMIT_TAG?=$$(git describe --tags $$(git rev-list --tags --max-count=1))
CI_COMMIT_SHA?=$$(git rev-parse HEAD)

# --------
# Commands
# --------

# Initialize CLI environment
-include ${PROJECT_CONFIG}
export

# Initialize DB configuration
MIGRATION_URL := "${MIGRATION_DB_DRIVER}://${MIGRATION_DB_USER}:${MIGRATION_DB_PASS}@${MIGRATION_DB_HOST}:${MIGRATION_DB_PORT}/${MIGRATION_DB_NAME}?sslmode=disable"
MIGRATION_BIN := migrate -source "file://${MIGRATION_SRC_DIR}" -database ${MIGRATION_URL}

# ------------
# Common Rules
# ------------

## help: Show command help
.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "${PROJECT_NAME}":"
	@echo
	@sed -n 's/^
	@echo

## clean: Clean everything
.PHONY: clean
clean:
	@-echo "  > Deleting output dir..."
	@-rm -rf ${OUTPUT_DIR}
	@-echo "  > Done"

## doctor: Check for prerequisites
.PHONY: doctor
doctor: $(DOCTOR_CMD)
	@-echo "  > Checking dependencies..."
	@-${DOCTOR_CMD}

# ---------
# API Rules
# ---------

## setup: Make env from env example and grant permission.
.PHONY: setup
setup:
	@-echo "  > Creating env file..."
	@cp config/.env.sample .env
	@-echo "  > Fix script permission..."
	@chmod +x script/*.sh
	@-echo "  > Removing tmp..."
	@-rm -rf tmp
	@-echo "  > Make new directory temp..."
	@-mkdir tmp

## configure: Configure project
.PHONY: configure
configure: --permit-exec --copy-env db-configure vendor
	@-echo "  > Configure: Done"

## serve: Run server in development mode
.PHONY: serve
serve: --dev-build ${DEBUG_ENV_FILES}
	@-echo "  > Starting Server...\n"
	@LOG_LEVEL=debug;LOG_FORMAT=console; ${DEBUG_BIN} -dir=${PROJECT_ROOT} -load-env-file

## vendor: Download dependencies to vendor folder
vendor: go.mod
	@-echo "  > Vendoring..."
	@${GO_MOD} vendor
	@-echo "  > Vendoring: Done"

## release: Compile binary for deployment.
.PHONY: release
release: vendor
	@-echo "  > Compiling for release..."
	@-echo "  >   Version: ${CI_COMMIT_TAG}"
	@-echo "  >   CommitHash: ${CI_COMMIT_SHA}"
	@CGO_ENABLED=0 GOOS=linux ${GO_BUILD} -a -v -mod=vendor \
		-ldflags "-X main.AppVersion=${CI_COMMIT_TAG} -X main.BuildHash=${CI_COMMIT_SHA}" \
		-o ${RELEASE_OUTPUT_DIR}/${BINARY_NAME} ${PROJECT_ROOT}/${PROJECT_MAIN_PKG}
	@-echo "  > Copying required file for release..."
	@-echo "  > Output: ${RELEASE_OUTPUT_DIR}"

## image: Build a docker image from release
.PHONY: image
image:
	@-echo "  > Building image ${IMAGE_APP}:${IMAGE_APP_TAG}..."
	${DOCKER_CMD} build -t ${IMAGE_APP}:$(IMAGE_APP_TAG) \
		--build-arg ARG_PORT=${PORT} \
	    --progress plain -f ${PROJECT_DOCKERFILE_DIR}/Dockerfile .

## image-push: Push app image
.PHONY: image-push
image-push: image
	@-echo "  > Push image ${IMAGE_APP}:${IMAGE_APP_TAG} to Container Registry..."
	@${DOCKER_CMD} push ${IMAGE_APP}:${IMAGE_APP_TAG}

## docker-serve: Run application with docker compose
.PHONY: docker-serve
docker-serve:
	@-echo "  > Run application using docker container"
	@docker-compose -f ${DOCKER_COMPOSE_FILE} pull
	@docker-compose -f ${DOCKER_COMPOSE_FILE} build
	@docker-compose -f ${DOCKER_COMPOSE_FILE} up --force-recreate


# ---------------
# Migration Rules
# ---------------

## db-configure: Generate a configuration for database migration tool
.PHONY: db-configure
db-configure:
	@-echo "  > Installing golang-migrate..."
	@-go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.1

## db-status: Prints the details and status information about all the migrations.
.PHONY: db-status
db-status:
	@-${MIGRATION_BIN} version

## db-up: Upgrade database
.PHONY: db-up
db-up:
	@-echo "  > Running up script..."
	@${MIGRATION_BIN} up

## db-down: (Experimental) undo to previous migration version
.PHONY: db-down
db-down:
	@${MIGRATION_BIN} down 1

## db-clean: Clean database
.PHONY: db-clean
db-clean: --clean-prompt
	@-echo "  > Cleaning database..."
	@${MIGRATION_BIN} drop

# -------------
# Private Rules
# -------------
.PHONY: --copy-env
--copy-env:
	@-echo "  > Copy .env (did not overwrite existing file)..."
	@-cp -n $(PROJECT_ROOT)/config/.env.sample $(PROJECT_CONFIG)

.PHONY: --permit-exec
--permit-exec: $(shell find $(SCRIPTS_DIR) -type f -name "*.sh")
	@-echo "  > Set executable permission to script..."
	@-chmod +x $(SCRIPTS_DIR)/*.sh

.PHONY: --clean-release
--clean-release:
	@-echo "  > Cleaning ${RELEASE_OUTPUT_DIR}..."
	@rm -rf ${RELEASE_OUTPUT_DIR}

.PHONY: --dev-build
--dev-build:
	@-echo "  > Compiling..."
	@${GO_BUILD} -ldflags "-X main.AppVersion=${CI_COMMIT_TAG} -X main.BuildHash=${CI_COMMIT_SHA}" \
		-o ${DEBUG_BIN} ${PROJECT_ROOT}/${PROJECT_MAIN_PKG}
	@-echo "  > Output: ${DEBUG_BIN}"

.PHONY: --clean-prompt
--clean-prompt:
	@echo -n "Are you sure want to clean all data in database? [y/N] " && read ans && [ $${ans:-N} = y ]

${DEBUG_ENV_FILES}: $(PROJECT_ENV_FILES)
	@-echo "  > Copying environment files..."
	@-cp -R ${PROJECT_ENV_FILES} ${DEBUG_DIR}
