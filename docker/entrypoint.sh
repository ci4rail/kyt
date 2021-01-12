#!/bin/bash
set -o errexit   # abort on nonzero exitstatus
set -o pipefail  # don't hide errors within pipes

useradd --shell /bin/bash -u ${USERID:-9001} -o -c "" -m user

# Give the user rights explicitly because systems like concourse CI wouldn't be able
# to work properly as subdirs are created as root.
chown user:user /home/user/

exec gosu user "$@"
