#!/bin/sh

set -o errexit
set -o pipefail

VENDOR=monostream.com
DRIVER=localflex

# assuming the single driver file is located at /$DRIVER inside the DaemonSet image.

driver_dir=$VENDOR${VENDOR:+"~"}${DRIVER}
if [ ! -d "/flexmnt/$driver_dir" ]; then
	mkdir "/flexmnt/$driver_dir"
fi

cp "/$DRIVER" "/flexmnt/$driver_dir/.$DRIVER"
mv -f "/flexmnt/$driver_dir/.$DRIVER" "/flexmnt/$driver_dir/$DRIVER"

# used for deployment script to know that this has finished
echo "done"

# this is a workaround to prevent the container from exiting and k8s restarting the daemonset pod
while : ; do
	sleep 3600
done