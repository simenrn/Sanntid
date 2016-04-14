#!/bin/bash
clear
echo "Username set to 'Student':"

echo "Type in the last byte of the IP to the elvator you want to connect to:"
read IP
echo "Connecting to 129.241.187."$IP
scp -r /home/student/ja student@129.241.187.$IP:~/jess
ssh student@129.241.187.$IP "jess"

