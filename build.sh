#!/bin/bash

# Set version
VERSION_MAJOR=0
VERSION_MINOR=3
VERSION_PATCH=1
VERSION_SPECIAL=
VERSION=""

# Setup defaults
export GO_ARCH=amd64
export CGO_ENABLED=0

BUILD_WINDOWS=0
BUILD_LINUX=0
BUILD_MAC=0
BIN_NAME="ps-reader"
OUT_DIR=./artifacts
SCRIPT_PATH=$0
STARTING_DIR=$(pwd)

function msg() {
    echo "travis build: $*"
}

function err() {
    msg "$*" 1>&2
}

# Triggers
SHOW_VERSION_LONG=0
SHOW_VERSION_SHORT=0
UPDATE_VERSION=0
TAR_FILES=0
PUSH_VERSION=0
TEST_RUN=0

while test $# -gt 0; do
  case $1 in
    -h|--help)
      # Show help message
      echo "-w: Builds the windows binary."
      echo "-l: Builds the Linux binary."
      echo "-m: Builds the Mac binary."
      echo "-t: Tar and gzip files that are compiled."
      echo "-x86: Sets the builds to be 32bit."
      echo "--test: runs all tests."
      echo "--output-name=<bin name>: Sets the output binary to be what is supplied. Windows binarys will have a .exe suffix add to it."
      echo "--output-dir=</path/to/dir>: Sets the output directory for built binaries."
      echo "--version-major=*: Update the Major part of the version number."
      echo "--version-minor=*: Update the Minor part of the version number."
      echo "--version-patch=*: Update the Patch part of the version number."
      echo "--version-special=*: Update the Special part of the version number."
      echo "-n|--next-minor: Increments the version numer to the next patch."
      echo "-u|--update-version: Updates the buidl script with the new version number. Commits it to git."
      exit 0
      shift
      ;;
    -w)
      BUILD_WINDOWS=1
      shift
      ;;
    -l)
      BUILD_LINUX=1
      shift
      ;;
    -m)
      BUILD_MAC=1
      shift
      ;;
    -t)
      TAR_FILES=1
      shift
      ;;
    -v)
      SHOW_VERSION_SHORT=1
      shift
      ;;
    --version)
      SHOW_VERSION_LONG=1
      shift
      ;;
    -x86)
      export GO_ARCH=386
      shift
      ;;
    --test=*)
      TEST_RUN=1
      TEST_OS=`echo $1 | sed -e 's/^[^=]*=//g'`
      shift
      ;;
    --output-name=*)
      BIN_NAME=`echo $1 | sed -e 's/^[^=]*=//g'`
      shift
      ;;
    --output-dir=*)
      OUT_DIR=`echo $1 | sed -e 's/^[^=]*=//g'`
      shift
      ;;
    --version-major=*)
      VERSION_MAJOR=`echo $1 | sed -e 's/^[^=]*=//g'`
      shift
      ;;
    --version-minor=*)
      VERSION_MINOR=`echo $1 | sed -e 's/^[^=]*=//g'`
      shift
      ;;
    --version-patch=*)
      VERSION_PATCH=`echo $1 | sed -e 's/^[^=]*=//g'`
      shift
      ;;
    --version-special=*)
      VERSION_SPECIAL=`echo $1 | sed -e 's/^[^=]*=//g'`
      shift
      ;;
    -n|--next-minor)
      let VERSION_PATCH+=1
      msg "Setting Patch number to: ${VERSION_PATCH}"
      shift
      ;;
    -u|--update-version)
      UPDATE_VERSION=1
      shift
      ;;
    *)
      break
      ;;
  esac
done

# We need to set the version after all the flag are read.
VERSION=$VERSION_MAJOR.$VERSION_MINOR.$VERSION_PATCH
if [ "$SPECIAL" != "" ]; then
  VERSION=$VERSION-$VERSION_SPECIAL;
fi

# Setup functions
ensure_artifact_dir(){
  if [ ! -d $1 ]; then
    mkdir -p $1
    if [ $? -ne 0 ]; then
      err "Failed to create the output directory: ${OUT_DIR}"
    fi
  fi
}

build_bin() {
  # Set GOOS
  goos=$1
  
  # Set the binary name
  bin_name=$BIN_NAME
  if [ $1 == "windows" ]; then
    bin_name="$BIN_NAME.exe"
  fi

  # Setup where it should go
  outdir=$OUT_DIR/$BIN_NAME-$goos-$GO_ARCH-v$VERSION
  output=$outdir/$bin_name

  # Ensure the artifact directory is there
  ensure_artifact_dir $outdir
  
  # Do any fixes that need doing.
  if [ -f ./customBuildFixes.sh ]; then
    source ./customBuildFixes.sh
  fi

  # Start the build
  GOOS=$goos \
  go build \
  -ldflags "-X main.VERSION=$VERSION" \
  -a \
  -installsuffix cgo \
  -o $output

  # Check if it worked
  if [ $? -eq 0 ]; then
    msg "Binary built and store as: $output"
  else
    msg "Binary for $goos failed to build!"
    exit 1
  fi
}

function travis-branch-commit() {
  if [ -z $TRAVIS_BRANCH ]; then
    msg "Skipping branch commit becuase we are not in a travis build environment"
    return
  fi

  local head_ref branch_ref commit_message files_to_commit
  
  commit_message=$1
  files_to_commit=""
  for i in ${@:2}; do
    files_to_commit="$files_to_commit $i"
  done
  
  head_ref=$(git rev-parse HEAD)
  if [[ $? -ne 0 || ! $head_ref ]]; then
    err "failed to get HEAD reference"
    return 1
  fi

  branch_ref=$(git rev-parse "$TRAVIS_BRANCH")
  if [[ $? -ne 0 || ! $branch_ref ]]; then
    err "failed to get $TRAVIS_BRANCH reference"
    return 1
  fi
  
  if [[ $head_ref != $branch_ref ]]; then
    msg "HEAD ref ($head_ref) does not match $TRAVIS_BRANCH ref ($branch_ref)"
    msg "someone may have pushed new commits before this build cloned the repo"
    return 0
  fi
  
  if ! git checkout "$TRAVIS_BRANCH"; then
    err "failed to checkout $TRAVIS_BRANCH"
    return 1
  fi

  msg "running: git add $files_to_commit"
  if ! git add $files_to_commit; then
    err "failed to add modified files to git index"
    return 1
  fi
  # make Travis CI skip this build
  msg "Committing changes:"
  git status
  if ! git commit -m "[ci skip] $commit_message"; then
    err "failed to commit updates"
    return 1
  fi

  local remote=origin
  if [[ $GH_TOKEN ]]; then
    msg "Adding in a token for Auth"
    remote=https://$GH_TOKEN@github.com/$TRAVIS_REPO_SLUG
  fi
  if [[ $TRAVIS_BRANCH != master ]]; then
    msg "not pushing updates to branch $TRAVIS_BRANCH"
    return 0
  fi

  msg "Trying to push commit"
  if ! git push "$remote" "$TRAVIS_BRANCH"; then
    err "failed to push git changes"
    return 1
  fi
}



# Show version
if [ $SHOW_VERSION_LONG -eq 1 ]; then
  msg "Current version number in build script: ${VERSION}"
  exit 0
fi

if [ $SHOW_VERSION_SHORT -eq 1 ]; then
  echo "v$VERSION"
  exit 0
fi

if [ $TEST_RUN -eq 1 ]; then
  failed=0
  if GOOS=$TEST_OS go test ./...; then
    msg "${TEST_OS} tests passed!"
  else
    msg "${TEST_OS} tests failed!"
    failed=1
  fi

  # If tests fail then fail the build.
  if [ $failed == 1 ]; then
    msg "Testing has failed. Stopping build!"
    exit 1
  fi
fi

# Update the version in this file
if [ $UPDATE_VERSION -eq 1 ]; then
  msg "Updating the build script with new version numbers."
  sed -i -r 's/^VERSION_MAJOR=[0-9]+$/VERSION_MAJOR='"$VERSION_MAJOR"'/' $SCRIPT_PATH \
  && sed -i -r 's/^VERSION_MINOR=[0-9]+$/VERSION_MINOR='"$VERSION_MINOR"'/' $SCRIPT_PATH \
  && sed -i -r 's/^VERSION_PATCH=[0-9]+$/VERSION_PATCH='"$VERSION_PATCH"'/' $SCRIPT_PATH \
  && sed -i -r 's/^VERSION_SPECIAL=*$/VERSION_SPECIAL='"$VERSION_SPECIAL"'/' $SCRIPT_PATH

  if [ $? != 0 ]; then
    err "Failed to update the version in build script"
    exit 1
  fi
  # We have updated the file we should push the new version
  PUSH_VERSION=1
fi

if [ $BUILD_LINUX -eq 1 ]; then
  build_bin "linux"
fi

if [ $BUILD_WINDOWS -eq 1 ]; then
  build_bin "windows"
fi

if [ $BUILD_MAC -eq 1 ]; then
  build_bin "darwin"
fi

if [ $TAR_FILES -eq 1 ]; then
  msg "Starting compression of binaries"
  cd $OUT_DIR
  for d in $(ls); do
    if [ $( echo $d | grep -c "windows") -eq 0 ]; then
      msg "starting tar and gzip on $d"
      tar -czvf $d.tar.gz $d
    else
      msg "Starting zip on $d"
      zip -r $d.zip $d
    fi
  done
  cd $STARTING_DIR
fi

if [ $PUSH_VERSION -eq 1 ]; then
  travis-branch-commit "Pushing new version $VERSION." "build.sh"
  if [ $? != 0 ]; then
    exit 1
  fi
fi

msg "Finished."