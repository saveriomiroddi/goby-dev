#!/bin/bash

# The PPA must have the https://launchpad.net/~gophers/+archive/ubuntu/archive PPA dependencies
# set, because it needs the `golang-1.10` package.

UBUNTU_RELEASE=xenial
GOBY_VERSION=0.1.9
MAINTAINER_EMAIL=saverio.pub2@gmail.com
PPA_PATH=saveriomiroddi/goby-lang
GOBY_BINPATH=/usr/bin
GOBY_LIBPATH="/usr/lib/goby/$GOBY_VERSION"

SOURCEDIR="$GOPATH/src/github.com/goby-lang/goby"
PACKAGEDIR="/tmp/goby-lang"

rm -rf "$PACKAGEDIR"
mkdir -p "$PACKAGEDIR/src/github.com/goby-lang"
cp -R "$SOURCEDIR" "$PACKAGEDIR/src/github.com/goby-lang/"
cd "$PACKAGEDIR/src/github.com/goby-lang/"

git checkout "v$GOBY_VERSION"

ln -s src/github.com/goby-lang/goby/Makefile

dh_make --single --native --copyright mit --email "$MAINTAINER_EMAIL"

perl -i -pe "s/unstable/$UBUNTU_RELEASE/" debian/changelog

rm debian/*.ex debian/*.EX

perl -i -pe 's/^(Section:).*/$1 devel/'                                  debian/control
perl -i -pe 's/^(Homepage:).*/$1 https:\/\/goby-lang.org/'               debian/control
perl -i -pe 's/^#(Vcs-Git:).*/$1 https:\/\/github.com\/goby-lang\/goby'  debian/control
perl -i -pe 's/^#(Vcs-Browser:).*/$1 https:\/\/github.com\/goby-lang\/goby' debian/control
perl -i -pe 's/^(Description:).*/$1 Goby language/'                      debian/control
perl -i -pe 's/^ <insert long description.*/ Goby - A new language helps you develop highly concurrent web applications/' debian/control
perl -i -pe 's/^(Standards-Version:) 3.9.6/$1 3.9.7/'                    debian/control
perl -i -pe 's/^(Build-Depends:.*)/$1, dh-golang, golang-1.10/'          debian/control

awk -i inplace 'NR==1{print; print "\n\
export PATH := /usr/lib/go-1.10/bin:$(PATH)\n\
export GOBY_BINPATH := /usr/bin\n\
export GOBY_LIBPATH := /usr/share/gocode/src/github.com/goby-lang/goby/lib\n\
"} NR!=1' debian/rules

echo $'override_dh_auto_test:\n\techo "Skipping tests..."\n' >> debian/rules

debuild -S | tee /tmp/debuild.log 2>&1
dput "ppa:$PPA_PATH" "$(perl -ne 'print $1 if /dpkg-genchanges -S >(.*)/' /tmp/debuild.log)"

echo "Remove personal infos!"
