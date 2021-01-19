#!/bin/bash
set -o errexit   # abort on nonzero exitstatus
set -o pipefail  # don't hide errors within pipes

addgroup -g ${GROUPID:-9001} user
adduser -u ${USERID:-9001} -G user -h /home/user -s /bin/sh -D user

# Give the user rights explicitly because systems like concourse CI wouldn't be able
# to work properly as subdirs are created as root.
chown user:user /home/user

exec gosu user "$@"
