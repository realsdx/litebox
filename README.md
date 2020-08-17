# LiteBox

Isolated execution environment for executing programs. User can set different resource limits
on CPU, Memory, Stack Size and File size etc. After execution get a report on Resource usage
by the program.

## Build

`go build .`

## Usage

`sudo ./litebox --cpu=5 --mem=7000 --nproc=25 --exec='python3 test.py'`

## Current Status
*Currently this project is development phase, do not use it*
If you found a problem, open an issue, if you aslo have the solution to it open a Pull Request :)
