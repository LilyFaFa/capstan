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

iperf3 $* > test-log 2>&1
if [ $? -ne 0 ]
then
    echo "Failed run iperf3 $*"
    cat test-log
    exit 1
fi

# Resolve result
echo "Resolving result"
cat test-log

line=`grep "receiver" test-log`

bandwidth=`echo $line|cut -d " " -f7`

RESULT="BandWidth $bandwidth $PrometheusLabel"
