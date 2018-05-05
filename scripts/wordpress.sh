#!/bin/sh
# Copyright (c) 2018 The ZJU-SEL Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

MYSQL=/usr/bin/mysql
TSUNG=/usr/local/tsung/bin/tsung
USER=root
check_result(){
    if [ $? -ne 0 ]
    then
        exit 1
    fi
} 

echo MYSQL_HOST=$MYSQL_HOST
echo MYSQL_PASSWORD=$MYSQL_ROOT_PASSWORD

$MYSQL -h $MYSQL_HOST -P3306 -u $USER -p$MYSQL_ROOT_PASSWORD bitnami_wordpress < /root/capstan/wordpress.sql
check_result

$TSUNG start > test-log 2>&1
if [ $? -ne 0 ]
then
    echo "Failed run wrk $*"
    cat test-log
    exit 1
fi

# Resolve result
echo "Resolving result"

line=`grep "Log directory is:" test-log`

log=`echo $line|cut -d " " -f4`
echo $log
cd $log
pwd
echo -ne '\n' |perl -MCPAN -e 'install Template'
perl /usr/local/tsung/lib/tsung/bin/tsung_stats.pl
pwd
RESULT="QPS $log $PrometheusLabel"
