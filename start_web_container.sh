#!/bin/bash

docker run -it --name hkjnweb -p 443:4430 -v /etc/ssl:/etc/ssl/:ro hkjn/armv7l-golang
