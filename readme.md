# SDR-GO Lab 4 Part 1

## About
### (Part 1)
The goal of this project is to implement a network of servers arranged in a graph. 
These server will share a work load and update each other on their result using undulation algorithm.

## Compiling and running

1. Download repository
2. run the "go build" command in src directory
3. A src.exe will be created, to execute, write "./src" or "src" depending on the shell

The main program will lunch servers according to config.json file. The depth has set correctly according to the topology used.
It will also lunch a client which can be used to test the servers (ui on cmd)