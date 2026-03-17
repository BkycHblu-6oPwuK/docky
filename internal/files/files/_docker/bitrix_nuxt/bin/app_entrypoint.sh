#!/bin/bash

CREATE_SIMLINK="$(dirname "$0")/create_simlink.sh"

if [ -f "$CREATE_SIMLINK" ]; then
    bash "$CREATE_SIMLINK"
fi

if [ "$XDEBUG_DISABLED" = "1" ] || [ "$XDEBUG_DISABLED" = "true" ]; then
    for f in /usr/local/etc/php/conf.d/*xdebug*.ini; do
        [ -e "$f" ] || continue
        mv "$f" "$f.disabled"
    done
fi

exec supervisord -c /etc/supervisor/conf.d/supervisord.conf
